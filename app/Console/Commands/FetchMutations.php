<?php

namespace App\Console\Commands;

use App\Jobs\PerformFetchMutations;
use App\Jobs\PerformFetchMutationsBatch;
use App\Models\Article;
use Illuminate\Console\Command;

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
        // $article = Article::where('success', 0)->first();

        Article::where('success', false)->take(10)->each(function ($article) {
            PerformFetchMutationsBatch::dispatch($article);
        });

        return Command::SUCCESS;
    }
}
