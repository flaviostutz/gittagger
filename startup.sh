#!/bin/bash
set -e
set -x

echo "Gittagger..."
gittagger \
    --loglevel=$LOG_LEVEL \
    --git-repo-url=$GIT_REPO_URL \
    --git-username="$GIT_USERNAME" \
    --git-email="$GIT_EMAIL" 
