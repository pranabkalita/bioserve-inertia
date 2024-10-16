<?php

namespace App\Console\Commands;

use App\Jobs\PerformFetchMutations;
use App\Jobs\PerformFetchMutationsBatch;
use App\Models\Article;
use App\Models\ProcessMutation;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class FetchMutations extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'app:fetch';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Fetch mutations';

    /**
     * Execute the console command.
     */
    public function handle()
    {
        // Record start time
        $startTime = microtime(true);

        Article::where('success', 0)->chunk(50, function ($articles) {
            // Collect PMIDs in chunks of 100
            $articleIds = $articles->pluck('pmid')->implode(',');

            try {
                // Create a ProcessMutation entry for this batch
                $processMutations = ProcessMutation::create(['pmids' => $articleIds]);

                // Process the batch using PerformFetchMutationsBatch
                app(PerformFetchMutationsBatch::class)->handle();

                // Delete the process mutation record after successful processing
                $processMutations->delete();
            } catch (\Exception $e) {
                // Log errors and continue to the next chunk of articles
                Log::error('Error processing batch: ' . $e->getMessage());
            }
        });


        // Record end time
        $endTime = microtime(true);

        // Calculate the time difference
        $executionTime = $endTime - $startTime;

        // Print the time taken
        dd('DONE - Time taken: ' . $executionTime . ' seconds');

        // -------------------------------------------------------


        // Record start time
        $startTime = microtime(true);

        $articleIds = Article::where('success', 0)
            ->limit(500)
            ->pluck('pmid')
            ->implode(',');

        $processMutations = ProcessMutation::create(['pmids' => $articleIds]);

        // Article::where('success', false)->take(100)->each(function ($article) {
        //     PerformFetchMutations::dispatch($article);
        // });

        // PerformFetchMutationsBatch::dispatch();

        $allArticles = Article::where('success', 0)->get();

        app(PerformFetchMutationsBatch::class)->handle();

        $processMutations->delete();

        // Record end time
        $endTime = microtime(true);

        // Calculate the time difference
        $executionTime = $endTime - $startTime;

        // Print the time taken
        dd('DONE - Time taken: ' . $executionTime . ' seconds');

        dd('DONE');

        return Command::SUCCESS;
    }
}
