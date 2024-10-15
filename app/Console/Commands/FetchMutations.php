<?php

namespace App\Console\Commands;

use App\Jobs\PerformFetchMutations;
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
        Article::where('success', false)->take(60)->each(function ($article) {
            PerformFetchMutations::dispatch($article);
        });

        return Command::SUCCESS;
    }
}
