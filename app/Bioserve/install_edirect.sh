#!/bin/bash

# Install EDirect
sh -c "$(curl -fsSL https://ftp.ncbi.nlm.nih.gov/entrez/entrezdirect/install-edirect.sh)"

# Output success message
echo "EDirect installed and PATH updated for the current session."
