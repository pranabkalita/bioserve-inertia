<?php

namespace App\Http\Controllers\Admin;

use App\Bioserve\BProtein;
use App\Bioserve\BProteinTerminal;
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
            $bProteinTerminal = new BProteinTerminal();
            // $bProteinTerminal->fetchPmidWithAbstract('11595708, 8872468');
            $pmidCount = $bProteinTerminal->searchProtein($request->get('search'));
        }

        $proteins = Protein::withCount('articles')->paginate(10);
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
            $bProteinTerminal = new BProteinTerminal();
            $bProteinTerminal->fetchPmids($request->get('search'));
            $pmids = $bProteinTerminal->getPmidsFromFile();

            $pmidChunks = array_chunk($pmids, 500);

            $protein = Protein::create(['name' => $request->get('search')]);

            foreach ($pmidChunks as $chunk) {
                $data = array_map(function ($pmid) use ($protein) {
                    return [
                        'pmid' => (int) $pmid,
                        'protein_id' => $protein->id,
                        'created_at' => now()->toDateTimeString(),
                        'updated_at' => now()->toDateTimeString(),
                    ];
                }, $chunk);

                Article::insert($data);
            }
        }

        return redirect()->route('collect.index')->with('Protein data collected.');
    }

    public function show(Request $request, Protein $protein)
    {
        $articles = Article::with(['mutations'])
            ->where('protein_id', $protein->id)
            ->where('success', true)
            ->has('mutations')
            ->paginate(10);

        // Map to add the mutation list as a comma-separated string
        $articles->getCollection()->transform(function ($article) {
            $article->mutationList = $article->mutations->pluck('name')->implode(', ');
            return $article;
        });

        return Inertia::render('Admin/Collect/Show', [
            'protein' => $protein,
            'articles' => $articles
        ]);
    }
}
