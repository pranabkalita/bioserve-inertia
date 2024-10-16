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

class PerformFetchMutationsBatchCopy
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
        $xml = $bProtein->fetchBatchPmidWithAbstract($pmids);

        $this->processMutations($xml);
    }

    private function processMutations($xml)
    {
        foreach ($xml as $element) {
            DB::beginTransaction(); // Start transaction for each loop

            try {
                $pmid = (int) $element->MedlineCitation->PMID;
                $abstract = (string) $element->MedlineCitation->Article->Abstract->AbstractText;

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
}
