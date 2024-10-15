<?php

namespace App\Http\Controllers\Admin;

use App\Bioserve\BProtein;
use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\ProteinStoreRequest;
use App\Http\Requests\Admin\SearchProteinRequest;
use App\Models\Article;
use App\Models\Protein;
use Illuminate\Http\Request;
use Inertia\Inertia;

class CollectController extends Controller
{
    public function index(SearchProteinRequest $request)
    {
        $pmidCount = 0;

        if ($request->get('search')) {
            $bProtein = new BProtein($request->get('search'));
            $bProtein->fetchTotalArticleCount();
            $pmidCount = $bProtein->getTotalArticleCount();
        }

        $proteins = Protein::paginate(10);
        // dd($proteins);

        return Inertia::render('Admin/Collect/Index', [
            'search' => $request->get('search'),
            'pmidCount' => $pmidCount,
            'proteins' => $proteins
        ]);
    }

    public function store(ProteinStoreRequest $request)
    {
        if ($request->get('search')) {
            $bProtein = new BProtein($request->get('search'));
            $bProtein->fetchTotalArticleCount();
            $bProtein->fetchArticles();

            $pmidChunks = array_chunk($bProtein->getPmids(), 500);

            $protein = Protein::create(['name' => $request->get('search')]);
            foreach ($pmidChunks as $chunk) {
                $data = array_map(function ($item) use ($protein) {
                    $item['protein_id'] = $protein->id;
                    $item['created_at'] = now()->toDateTimeString();
                    $item['updated_at'] = now()->toDateTimeString();
                    return $item;
                }, $chunk);

                Article::insert($data);
            }
        }

        return redirect()->route('collect.index')->with('Protein data collected.');
    }
}
