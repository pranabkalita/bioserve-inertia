<?php

namespace App\Jobs;

use App\Bioserve\BProteinTerminal;
use App\Models\Article;
use App\Models\Mutation;
use App\Models\ProcessMutation;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class PerformFetchMutationsBatch
{
    // use Queueable;

    // public string $pmids = '';

    /**
     * Create a new job instance.
     */
    public function __construct()
    {
        // $this->pmids = implode(',', $pmids);
    }

    public function uniqueId()
    {
        return 'article_' . rand(1, 99999);
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        $pmids = ProcessMutation::first()->pmids;

        $bProtein = new BProteinTerminal();

        try {
            // Fetch the batch of articles using the PMIDs
            $xml = $bProtein->fetchBatchPmidWithAbstract($pmids);
            $citationsAndData = $this->prepareXMLForProcessing($xml);

            // Process the fetched mutations
            $this->processMutations($citationsAndData);
        } catch (\Exception $e) {
            // Log the error and continue with the next batch
            Log::error('Error in fetching or processing mutations: ' . $e->getMessage());
        }
    }

    private function processMutations($xml)
    {
        foreach ($xml as $element) {
            DB::beginTransaction(); // Start transaction for each loop

            try {
                $pmid = (int) $element['MedlineCitation']['PMID'];

                $abstract = '';

                if (isset($element['MedlineCitation']['Article']['Abstract'])) {
                    $abstractText = $element['MedlineCitation']['Article']['Abstract']['AbstractText'];
                    if (is_array($abstractText)) {
                        $abstract = implode(' ', $abstractText);
                    } else {
                        $abstract = $abstractText;
                    }
                }

                $mutationPattern = '/\b[A-Z]\d{2,5}[A-Z]\b/';
                $matches = [];
                preg_match_all($mutationPattern, $abstract, $matches);

                $uniqueMutations = array_unique($matches[0]);

                // Update PMID: status = 1
                $article = Article::where('pmid', $pmid)->first();

                // Update the article's success field
                $article->update(['success' => 1]);

                // If uniqueMutations length is more than 0, insert into mutation
                if (count($uniqueMutations) > 0) {
                    $data = [];

                    foreach ($uniqueMutations as $mutation) {
                        $data[] = [
                            'article_id' => $article->id,
                            'name' => $mutation,
                            'created_at' => now()->toDateTimeString(),
                            'updated_at' => now()->toDateTimeString(),
                        ];
                    }

                    // Insert mutations
                    Mutation::insert($data);
                }

                DB::commit(); // Commit transaction if all is good
            } catch (\Exception $e) {
                DB::rollBack(); // Rollback transaction if there's an error

                // Log the error to continue debugging, without halting the loop
                Log::error('Error processing PMID: ' . $pmid . '. Error: ' . $e->getMessage());
            }
        }
    }

    private function prepareXMLForProcessing($data)
    {
        $results = [];

        // Check if "PubmedArticle" is set in the array
        if (isset($data['PubmedArticle'])) {
            $pubmedArticles = $data['PubmedArticle'];

            // Check if it's a single PubmedArticle (not an indexed array)
            if (isset($pubmedArticles['MedlineCitation']) && isset($pubmedArticles['PubmedData'])) {
                // Single PubmedArticle case
                $pubmedArticles = [$pubmedArticles];
            }

            // Loop through each PubmedArticle and extract MedlineCitation with PubmedData
            foreach ($pubmedArticles as $article) {
                if (isset($article['MedlineCitation']) && isset($article['PubmedData'])) {
                    $results[] = [
                        'MedlineCitation' => $article['MedlineCitation'],
                        'PubmedData'      => $article['PubmedData'],
                    ];
                }
            }
        }

        return $results;
    }
}
