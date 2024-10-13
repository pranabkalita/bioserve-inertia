<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Article extends Model
{
    use HasFactory;

    protected $fillable = [
        'pmid',
        'published_on',
        'success'
    ];

    // Methods

    // Relations
    public function protein()
    {
        return $this->belongsTo(Protein::class);
    }

    public function mutations()
    {
        return $this->hasMany(Mutation::class);
    }
}
