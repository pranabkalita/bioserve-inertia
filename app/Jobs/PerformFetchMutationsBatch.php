<?php

namespace App\Jobs;

use App\Bioserve\BProteinTerminal;
use App\Models\Article;
use App\Models\Mutation;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;

class PerformFetchMutationsBatch implements ShouldQueue, ShouldBeUnique
{
    use Queueable;

    // public string $pmids = '';

    /**
     * Create a new job instance.
     */
    public function __construct(public Article $article)
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
        $bProtein = new BProteinTerminal();
        $xml = $bProtein->fetchBatchPmidWithAbstract($this->article->pmid);

        $this->processMutations($xml);
    }

    private function processMutations($xml)
    {
        foreach ($xml as $element) {
            $pmid = (int) $element->MedlineCitation->PMID;
            $abstract = (string) $element->MedlineCitation->Article->Abstract->AbstractText;

            $mutationPattern = '/\b[A-Z]\d{2,5}[A-Z]\b/';
            $matches = [];
            preg_match_all($mutationPattern, $abstract, $matches);

            $uniqueMutations = array_unique($matches[0]);

            // Update PMID: status = 1
            Article::where('pmid', $pmid)->update(['success' => 1]);
            $article = Article::where('pmid', $pmid)->first();

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

                Mutation::insert($data);
            }
        }
    }
}
