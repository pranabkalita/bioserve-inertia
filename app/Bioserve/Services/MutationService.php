<?php

namespace App\Bioserve\Services;

use App\Models\Article;
use App\Models\Mutation;
use App\Models\Protein;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class MutationService
{
    protected int $protein_id;

    public function __construct(int $protein_id)
    {
        $this->protein_id = $protein_id;
    }

    /**
     * Process mutations for the given XML data.
     *
     * @param array $xml
     * @return void
     */
    public function process(array $xml)
    {
        foreach ($xml as $element) {
            DB::beginTransaction(); // Start transaction for each loop

            try {
                $pmid = (int) $element['MedlineCitation']['PMID'];
                $protein = Protein::findOrFail($this->protein_id);

                // Extract the abstract
                $abstract = $this->extractAbstract($element);

                // Define the mutation pattern based on protein name
                $mutationPattern = '/(?:\b' . $protein->name . '\b(?:[\s-])?|(?<![A-Za-z\d-]))[A-Z]\d{2,5}[A-Z]\b(?![A-Za-z\d-])/';

                // Match the mutations in the abstract
                preg_match_all($mutationPattern, $abstract, $matches);
                $uniqueMutations = array_unique($matches[0]);

                // Fetch and update the article's success field
                $article = Article::where([
                    'pmid' => $pmid,
                    'protein_id' => $this->protein_id
                ])->first();

                $article->update(['success' => 1]);

                // Insert mutations if found
                if (count($uniqueMutations) > 0) {
                    $this->insertMutations($article->id, $uniqueMutations);
                }

                DB::commit(); // Commit transaction if all is good
            } catch (\Exception $e) {
                DB::rollBack(); // Rollback transaction if there's an error

                // Log the error to continue debugging, without halting the loop
                Log::error('Error processing PMID: ' . $pmid . '. Error: ' . $e->getMessage());
            }
        }
    }

    /**
     * Insert mutations into the database.
     *
     * @param int $article_id
     * @param array $mutations
     * @return void
     */
    protected function insertMutations(int $article_id, array $mutations): void
    {
        $data = [];
        foreach ($mutations as $mutation) {
            $data[] = [
                'article_id' => $article_id,
                'name' => $mutation,
                'created_at' => now()->toDateTimeString(),
                'updated_at' => now()->toDateTimeString(),
            ];
        }

        Mutation::insert($data);
    }

    /**
     * Extract abstract from a PubMed article element.
     *
     * @param array $element
     * @return string
     */
    protected function extractAbstract(array $element): string
    {
        $abstract = '';
        if (isset($element['MedlineCitation']['Article']['Abstract'])) {
            $abstractText = $element['MedlineCitation']['Article']['Abstract']['AbstractText'];
            if (is_array($abstractText)) {
                $abstract = implode(' ', $abstractText);
            } else {
                $abstract = $abstractText;
            }
        }

        return $abstract;
    }
}
