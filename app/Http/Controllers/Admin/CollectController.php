<?php

namespace App\Http\Controllers\Admin;

use App\Bioserve\BProtein;
use App\Http\Controllers\Controller;
use App\Models\Protein;
use Illuminate\Http\Request;
use Inertia\Inertia;

class CollectController extends Controller
{
    protected BProtein $bProtein;

    protected int $pmidCount = 0;

    public function index(Request $request)
    {
        $pmidCount = 0;
        if ($request->get('search')) {
            $bProtein = new BProtein($request->get('search'));
            $bProtein->fetchTotalArticleCount();
            $pmidCount = $bProtein->getTotalArticleCount();
        }

        $protrins = Protein::paginate(10);

        return Inertia::render('Admin/Collect/Index', [
            'search' => $request->get('search'),
            'pmidCount' => $pmidCount,
            'protrins' => $protrins
        ]);
    }
}
