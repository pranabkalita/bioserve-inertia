<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('articles', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('protein_id')->unsigned()->index();
            $table->string('pmid')->nullable()->index();
            $table->string('published_on')->nullable();
            $table->boolean('success')->default(false);
            $table->timestamps();

            $table->foreign('protein_id')->references('id')->on('proteins')->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('articles');
    }
};
