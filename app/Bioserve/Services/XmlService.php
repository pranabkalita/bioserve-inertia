<?php

namespace App\Bioserve\Services;

class XmlService
{
    /**
     * Prepare XML data for processing by extracting the necessary fields.
     *
     * @param array $data
     * @return array
     */
    public function prepareForProcessing(array $data): array
    {
        $results = [];

        if (isset($data['PubmedArticle'])) {
            $pubmedArticles = $data['PubmedArticle'];

            // Check if it's a single PubmedArticle (not an indexed array)
            if (isset($pubmedArticles['MedlineCitation']) && isset($pubmedArticles['PubmedData'])) {
                $pubmedArticles = [$pubmedArticles];
            }

            foreach ($pubmedArticles as $article) {
                if (isset($article['MedlineCitation']) && isset($article['PubmedData'])) {
                    $results[] = [
                        'MedlineCitation' => $article['MedlineCitation'],
                        'PubmedData'      => $article['PubmedData'],
                    ];
                }
            }
        }

        return $results;
    }

    /**
     * Extract abstract from a PubMed article element.
     *
     * @param array $element
     * @return string
     */
    public function extractAbstract(array $element): string
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
