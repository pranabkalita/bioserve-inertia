<?php

namespace App\Http\Controllers\Admin;

use App\Bioserve\BProteinTerminal;
use App\Http\Controllers\Controller;
use App\Models\Article;
use App\Models\Protein;
use Illuminate\Http\Request;

class ArticleController extends Controller
{
    public function store(Request $request)
    {
        if (!$request->get('id')) {
            return redirect()->back();
        }

        $protein = Protein::findOrFail($request->get('id'));

        $bProteinTerminal = new BProteinTerminal();
        $bProteinTerminal->fetchPmids($protein->name);
        $pmids = $bProteinTerminal->getPmidsFromFile();

        $pmidChunks = array_chunk($pmids, 500);

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

        return redirect()->back();
    }
}
