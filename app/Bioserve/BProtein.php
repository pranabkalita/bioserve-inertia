<?php

namespace App\Bioserve;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

use Symfony\Component\Process\Process;
use Symfony\Component\Process\Exception\ProcessFailedException;

use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Storage;

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

        if ($this->pmidCount > 9999) {
            $edirect = public_path('edirect');

            if (!File::exists($edirect)) {
                $this->installEDirect();
            }

            $this->exportToText();
        }

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


    private function __installEDirect()
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

    public function moveEDirect()
    {
        // Define the source and destination paths
        $source = public_path('edirect');
        $destination = app_path('Bioserve/edirect');

        // Check if the source directory exists
        if (!File::exists($source)) {
            return response()->json(['error' => 'Source directory does not exist.'], 404);
        }

        // Create the destination directory if it doesn't exist
        if (!File::exists($destination)) {
            File::makeDirectory($destination, 0755, true);
        }

        // Move the directory
        try {
            File::move($source, $destination);
            return response()->json(['success' => 'EDirect moved successfully.']);
        } catch (\Exception $e) {
            return response()->json(['error' => 'Failed to move EDirect: ' . $e->getMessage()], 500);
        }
    }

    private function installEDirect()
    {
        // Define the full script path
        $scriptPath = base_path('app/Bioserve/install_edirect.sh');

        // Step 1: Add execute permission to the script
        $chmodProcess = new Process(['chmod', '+x', $scriptPath]);
        $chmodProcess->run();

        // Check if the chmod command executed successfully
        if (!$chmodProcess->isSuccessful()) {
            throw new ProcessFailedException($chmodProcess);
        }

        // Step 2: Run the shell script using Process
        $process = new Process(['sh', $scriptPath]);
        $process->run();

        // Check if the process executed successfully
        if (!$process->isSuccessful()) {
            throw new ProcessFailedException($process);
        }

        // Output the result of the script execution
        echo $process->getOutput();
    }

    private function exportToText()
    {
        // Define the command
        $command = './edirect/esearch -db pubmed -query "' . $this->proteinName . '" | ./edirect/efetch -format uid > _ids.txt';

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
