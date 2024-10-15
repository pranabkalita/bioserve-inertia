<?php

use Illuminate\Foundation\Inspiring;
use Illuminate\Support\Facades\Artisan;

Artisan::command('inspire', function () {
    $this->comment(Inspiring::quote());
})->purpose('Display an inspiring quote')->hourly();

// \Illuminate\Support\Facades\Schedule::job(new \App\Console\Commands\FetchMutations)->everyMinute();

\Illuminate\Support\Facades\Schedule::command('app:fetch')->everyMinute();

// Artisan::command('app:fetch', function () {
//     Article::where('success', false)->take(1)->each(function ($article) {
//         PerformFetchMutations::dispatch($article);
//     });

//     return Command::SUCCESS;
// })->everyMinute();
