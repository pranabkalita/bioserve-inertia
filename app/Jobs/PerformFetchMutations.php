<?php

namespace App\Jobs;

use App\Models\Article;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class PerformFetchMutations implements ShouldQueue, ShouldBeUnique
{
    use Queueable;

    /**
     * Create a new job instance.
     */
    public function __construct(public Article $article)
    {
        //
    }

    public function uniqueId()
    {
        return 'article_' . $this->article->id;
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        $url = 'https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=' . $this->article->pmid . '&rettype=abstract&retmode=text';
        $articleData = $this->fetchArticleMutations($url);

        // dd($url, $articleData, $this->article->id);

        if ($articleData) {
            foreach ($articleData['mutations'] as $mutation) {
                $this->article->mutations()->create([
                    'name' => $mutation
                ]);
            }
        }

        $this->article->update(['success' => true]);
    }

    public function fetchArticleMutations($url)
    {
        // Fetch the contents of the URL
        $articleContent = file_get_contents($url);

        if ($articleContent === false) {
            return []; // If unable to fetch the content, return an empty array
        }

        // $articleContent = $articleContent . " G1043D, G14313D";

        // Match all mutations using the regex ^[A-Z]\d{2,5}[A-Z]$
        // $mutationPattern = '/[A-Z]\d{2,5}[A-Z]/';
        // $mutationPattern = '/(?<![A-Z])[A-Z]\d{2,5}[A-Z](?![A-Z])/';
        $mutationPattern = '/\b[A-Z]\d{2,5}[A-Z]\b/';
        // $mutationPattern = '/[A-Z]\d{2,5}[A-Z](?=[,\s]|$)/';
        // $mutationPattern = '/\b[A-Z]\d{2,5}[A-Z]\b/';
        $matches = [];
        preg_match_all($mutationPattern, $articleContent, $matches);

        if (!empty($matches[0])) {
            // Get the first two lines of the article
            $lines = explode("\n", trim($articleContent));
            $firstTwoLines = array_slice($lines, 0, 2);

            // Prepare the result
            return [
                'mutations' => $matches[0],    // All matched mutations
                'first_two_lines' => implode("\n", $firstTwoLines) // First two lines as a single string
            ];
        }

        // Return empty array if no mutations found
        return [];
    }
}
