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
        //     $this->installEDirect();
        //     $this->exportToText();

        //     dd('done');
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


    private
    function installEDirect()
    {
        // Detect home directory and installation path
        $homeDirectory = getenv('HOME');
        if (!$homeDirectory) {
            $homeDirectory = shell_exec('echo $HOME');
            $homeDirectory = trim($homeDirectory);
        }

        // Define the install command
        $installCommand = 'sh -c "$(curl -fsSL https://ftp.ncbi.nlm.nih.gov/entrez/entrezdirect/install-edirect.sh)"';

        // Run the installation command
        echo "Installing EDirect...\n";
        $output = shell_exec($installCommand);
        echo $output;

        // Check if installation was successful
        if (strpos($output, 'successfully downloaded and installed') === false) {
            echo "EDirect installation failed.\n";
            return false;
        }

        // Command to add PATH to the shell configuration

        // Export PATH for the current terminal session
        $exportPathCommand = 'export PATH=${HOME}/edirect:${PATH}';
        shell_exec($exportPathCommand);

        echo "EDirect installed and PATH updated successfully.\n";
        return true;
    }

    private function exportToText()
    {
        // Define the command
        $command = 'esearch -db pubmed -query "' . $this->proteinName . '" | efetch -format uid > pmids.txt';

        // Execute the command
        exec($command, $output, $return_var);

        // Check if the command was successful
        if ($return_var === 0) {
            echo "PMIDs have been successfully exported to pmids.txt";
        } else {
            echo "Error executing command: return code " . $return_var;
        }
    }
}
