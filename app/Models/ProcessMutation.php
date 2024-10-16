<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ProcessMutation extends Model
{
    use HasFactory;

    protected $fillable = [
        'pmids'
    ];
}
