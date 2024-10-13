<?php

namespace App\Bioserve;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

class BProtein
{
    protected string $proteinName = '';

    protected int $pmidCount = 0;

    protected array $pmidArray = [];

    public function __construct(string $proteinName)
    {
        $this->proteinName = $proteinName;
    }


    public function getTotalArticleCount(): int
    {
        return $this->pmidCount;
    }

    public function getPmids(): array
    {
        return $this->pmidArray;
    }

    public function fetchTotalArticleCount(): void
    {
        $client = new Client();

        $searchUrl = 'https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=' . urlencode($this->proteinName) . '&retmode=xml';

        try {
            $response = $client->get($searchUrl);
            $xml = simplexml_load_string($response->getBody()->getContents());

            $this->pmidCount = (int) $xml->Count;
        } catch (RequestException $e) {
            echo "Error fetching data: " . $e->getMessage() . "\n";
        }
    }

    public function fetchArticles(): void
    {
        if (!$this->pmidCount) return;

        // if ($this->pmidCount > 9999) {
        //     // Shell Script
        // }

        $client = new Client();

        try {
            $searchUrl = 'https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=' . urlencode($this->proteinName) . '&retmax=' . $this->pmidCount . '&retmode=xml';

            $response = $client->get($searchUrl);
            $xml = simplexml_load_string($response->getBody()->getContents());

            // Get the PMIDs
            $pmidArray = (array) $xml->IdList->Id ?? [];
            $this->pmidArray = array_map(function ($pmid) {
                return [
                    'pmid' => (string) $pmid
                ];
            }, $pmidArray);
        } catch (RequestException $e) {
            echo "Error fetching data: " . $e->getMessage() . "\n";
        }
    }
}
