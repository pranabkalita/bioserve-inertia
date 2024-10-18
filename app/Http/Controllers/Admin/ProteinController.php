<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Models\Protein;
use Illuminate\Http\Request;
use Inertia\Inertia;

class ProteinController extends Controller
{
    public function index(Request $request)
    {
        $proteins = Protein::withCount(['articles', 'articles as mutations_count' => function ($query) {
            $query->join('mutations', 'articles.id', '=', 'mutations.article_id');
        }])
            ->when($request->get('search'), function ($query) use ($request) {
                $query->where('name', 'LIKE', '%' . $request->get('search') . '%');
            })->paginate(10)->withQueryString();

        return Inertia::render('Admin/Protein/Index', [
            'proteins' => $proteins,
            'search' => $request->get('search')
        ]);
    }
}
