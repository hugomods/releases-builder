# Releases Builder

The CLI tool is used to fetch releases of the given repositories, and save the content for static site generators to render.
It's useful for showing the releases on your project sites, real use cases: [HugoMods Releases](https://hugomods.com/releases/) and [HB Framework Releases](https://hbstack.dev/releases/).

## Installation

```sh
go install github.com/hugomods/releases-builder@latest
```

## Usage

```sh
releases-builder -h
```

## Integrate With GitHub Action

1. Create a [confgiuration](#configuration) file and save it as `.releases-builder.yaml` on repo root.
2. Create a workflow `.github/workflows/releases-builder.yaml`.
```yaml
name: Releases Builder

on:
  workflow_dispatch:
  schedule:
    - cron: '0 * * * *' # run every hour, change it at will.

jobs:
  contributors:
    permissions:
      contents: write
    runs-on: ubuntu-latest
    name: Releases Builder
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          token: ${{ secrets.GITHUB_TOKEN }} # You should replace it with another token, so that it will trigger other workflows.

      - uses: actions/setup-go@v5
        with:
          go-version: 'stable'
      
      - name: Install Releases Builder CLI
        run: go install github.com/hugomods/releases-builder@latest
      
      - name: Fetch and Save Releases
        run: releases-builder
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Commit Changes
        uses: stefanzweifel/git-auto-commit-action@v5
        with:
          commit_message: "chore: update releases"
          commit_author: "github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>" 
          file_pattern: 'content/releases' # should be same as the contentDir.
```

There are some real repos using this tool: https://github.com/hbstack/site and https://github.com/hugomods/site.

## Configuration

### `repositories`

A list of repositories, e.g. `github.com/hugomods/search`, `github.com/hugomods/images`.

### `contentDir`

The location where the releases content to be put into, default to `content/releases`.

### `languages`

An array of languages, when provied, the releases content will be saved in `index.[lang].md` pattern.

### `params`

Front matter params.

#### `params.images`

An array of images.

#### `params.authors`

An array of authors.

#### `params.categories`

An array of categories.

#### `params.series`

An array of series.

#### `params.tags`

An array of tags.

### Full Configuration Example

```yaml
languages:
  - code: en
  - code: fr
  - code: zh-hans
contentDir: content/releases
repositories:
  - github.com/hugomods/images
  - github.com/hugomods/seo
params:
  authors:
    - HugoMods
  categories:
    - C1
    - C2
  series:
    - Releases
  images:
    - /images/releases.png
  tags:
    - T1
    - T2
```
