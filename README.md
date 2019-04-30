# Snyk Import API Client - Work in Progress

# currently only for Bitbucket Server

## What it does
It takes a list of relative paths to point Snyk to specific files to import in specific location for the bitbucket integration of your org

## Usage
1. List all your relative paths into a file mypaths
2. cat mypaths | ./snykimport  -token 12345678-1234-1234-1234-123456789012 -orgId 12345678-1234-1234-1234-123456789012 -intId 12345678-1234-1234-1234-123456789012 -repo MyRepoName -projectkey AN

You can pipe whatever you want into snykimport as long as it is the relative paths, respecting case sensitivity. Use `find <folder> -name <pattern>` fu or whatever else.

## Arguments
- token: Your Snyk API token (under account settings)
- orgId: Look into your Snyk org settings
- intId: Integration ID, find it under settings->integrations->Bitbucket Server
- repo: Repo name, find it in Bitbucket Server
- projectKey: project Key, find it in Bitbucket Server

## RepoSlug option
If you repo name has a space, a slug will need to be specified. I.e "My Repo" will have a "My-Repo" slug, which you will need to specify.
