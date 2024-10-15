<?php

namespace App\Bioserve;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

use Symfony\Component\Process\Process;
use Symfony\Component\Process\Exception\ProcessFailedException;

use Illuminate\Support\Facades\File;


class BProteinTerminal
{
    protected string $protein = '';
    protected int $pmidCount = 0;

    public function __construct()
    {
        $edirect = public_path('edirect');

        if (!File::exists($edirect)) {
            $this->installEDirect();
        }
    }

    public function searchProtein(string $protein)
    {
        $command = './edirect/esearch -db pubmed -query "' . $protein . '"';

        $process = Process::fromShellCommandline($command);
        $process->run();

        if (!$process->isSuccessful()) {
            throw new ProcessFailedException($process);
        }

        $output = $process->getOutput();

        $xml = simplexml_load_string($output);

        return (int) $xml->Count;
    }

    public function fetchPmids(string $protein)
    {
        $command = './edirect/esearch -db pubmed -query "' . $protein . '" | ./edirect/efetch -format uid > _ids.txt';

        $process = Process::fromShellCommandline($command);

        $process->run();

        if (!$process->isSuccessful()) {
            return false;
            // return "Error executing command: " . $process->getErrorOutput() . " Return code: " . $process->getExitCode();
        }

        return true;
    }

    public function getPmidsFromFile()
    {
        $filePath = public_path('_ids.txt');

        $fileContents = file($filePath, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);

        if ($fileContents !== false) {
            return $fileContents;
        } else {
            return [];
        }
    }

    public function fetchPmid(int $pmid)
    {
        $command = './edirect/nquire -get https://icite.od.nih.gov api/pubs -pmids ' . $pmid . '';

        $process = Process::fromShellCommandline($command);

        $process->run();

        if (!$process->isSuccessful()) {
            return false;
            // return "Error executing command: " . $process->getErrorOutput() . " Return code: " . $process->getExitCode();
        }

        $output = $process->getOutput();
        $data = json_decode($output, true)['data'];

        dd($data);
    }

    public function fetchPmidWithAbstract(string $pmid)
    {
        $command = './edirect/efetch -db pubmed -id ' . $pmid . ' -format xml';

        $process = Process::fromShellCommandline($command);

        $process->run();
        // dd(!$process->isSuccessful());
        if (!$process->isSuccessful()) {
            return false;
        }

        $output = $process->getOutput();

        $xml = simplexml_load_string($output);

        dd($xml);
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
}
