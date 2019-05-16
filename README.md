# Snyk Import API Client - Work in Progress

# currently only for Bitbucket Server

## What it does
It takes a list of relative paths to point Snyk to specific files to import in specific location for the bitbucket integration of your org

## Usage
1. Clone your repo
2. ./snykimport-[os]-[version] -token 12345678-1234-1234-1234-123456789012 -orgId 12345678-1234-1234-1234-123456789012 -intId 12345678-1234-1234-1234-123456789012 -repo MyRepoName -projectkey AN -path /path/to/your/repo -excludeFile /path/to/excludelist

Can add -d for debug output, seeing what gets ignored etc.

## Usage using linux pipe
1. List all your relative paths into a file mypaths
2. cat mypaths | ./snykimport  -token 12345678-1234-1234-1234-123456789012 -orgId 12345678-1234-1234-1234-123456789012 -intId 12345678-1234-1234-1234-123456789012 -repo MyRepoName -projectkey AN

You can pipe whatever you want into snykimport as long as it is the relative paths, respecting case sensitivity. Use `find <folder> -name <pattern>` fu or whatever else.

## Arguments
- token: Your Snyk API token (under account settings)
- orgId: Look into your Snyk org settings
- intId: Integration ID, find it under settings->integrations->Bitbucket Server
- repo: Repo name, find it in Bitbucket Server
- projectKey: project Key, find it in Bitbucket Server
- path: Path to your repo root folder
- excludeFile: fullpath to your exclude file
- d: debug output

## Exclude file
Simple text file list strings that if detected in path will result in file being ignore.
example "node_modules" for Node apps.

## RepoSlug option
If you repo name has a space, a slug will need to be specified. I.e "My Repo" will have a "My-Repo" slug, which you will need to specify.
