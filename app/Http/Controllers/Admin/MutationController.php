<?php

namespace App\Http\Controllers\Admin;

use App\Bioserve\BProteinTerminal;
use App\Bioserve\Services\MutationService;
use App\Bioserve\Services\XmlService;
use App\Http\Controllers\Controller;
use App\Models\Article;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class MutationController extends Controller
{
    protected int $protein_id;
    protected XmlService $xmlService;

    public function __construct(XmlService $xmlService)
    {
        $this->xmlService = $xmlService;
    }

    public function store(Request $request)
    {
        if (!$request->get('protein_id')) {
            return redirect()->back();
        }

        $this->protein_id = $request->get('protein_id');
        $bProtein = new BProteinTerminal();

        $pmids = Article::where([
            'protein_id' => $this->protein_id,
            'success' => 0
        ])->pluck('pmid')->toArray();

        $pmidChunks = array_chunk($pmids, 50);

        foreach ($pmidChunks as $chunk) {
            try {
                // Fetch the batch of articles using the PMIDs
                $xml = $bProtein->fetchBatchPmidWithAbstract(implode(",", $chunk));
                $citationsAndData = $this->xmlService->prepareForProcessing($xml);

                // Use the MutationService to process the mutations
                $mutationService = new MutationService($this->protein_id);
                $mutationService->process($citationsAndData);
            } catch (\Exception $e) {
                Log::error('Error in fetching or processing mutations: ' . $e->getMessage());
            }
        }
    }
}
