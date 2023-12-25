package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/go-errors/errors"
	yaml "github.com/goccy/go-yaml"
	"github.com/google/go-github/v57/github"
	"github.com/hugomods/releases-builder/internal/config"
)

var cfg *config.Config
var cfgFile string

func init() {
	cfg = config.NewConfig()
}

func main() {
	flag.StringVar(&cfgFile, "c", ".releases-builder.yaml", "Configuration file.")
	flag.Parse()

	err := parseConfig()
	if err != nil {
		panic(err)
	}

	var wait sync.WaitGroup

	ctx := context.Background()
	for _, repo := range cfg.Repositories {
		wait.Add(1)
		go build(ctx, repo, &wait)
	}

	wait.Wait()
}

func parseConfig() error {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, cfg); err != nil {
		return err
	}

	if len(cfg.Repositories) == 0 {
		return errors.New("no repositories specified.")
	}

	return nil
}

func build(ctx context.Context, repo string, wg *sync.WaitGroup) {
	defer wg.Done()

	paths := strings.Split(repo, "/")
	if len(paths) != 3 {
		panic("invalid repository: " + repo)
	}
	if paths[0] != "github.com" {
		panic("unsupported platform: " + paths[0])
	}

	client := github.NewClient(nil)
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		client = client.WithAuthToken(token)
	}

	releases, _, err := client.Repositories.ListReleases(ctx, paths[1], paths[2], &github.ListOptions{
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		panic(err)
	}

	for _, release := range releases {
		if err := generate(repo, release); err != nil {
			panic(err)
		}
	}
}

var tplContent = template.Must(template.New("content").Parse(`---
title: "{{ .repo }}'s {{ .release.Name }}"
date: {{ .release.CreatedAt }}
publishDate: {{ .release.PublishedAt }}
draft: {{ .release.Draft }}
prerelease: {{ .release.Prerelease }}
name: "{{ .release.Name }}"
tag_name: "{{ .release.TagName }}"
release_url: "{{ .release.HTMLURL }}"
---

{{ .release.Body }}
`))

func generate(repo string, release *github.RepositoryRelease) error {
	dir := filepath.Join(filepath.FromSlash(cfg.ContentDir), strings.Replace(repo, "/", "-", -1))
	if err := os.MkdirAll(dir, 0744); err != nil {
		return err
	}

	var buff bytes.Buffer
	data := map[string]interface{}{
		"repo":    repo,
		"release": release,
	}
	if err := tplContent.Execute(&buff, data); err != nil {
		return err
	}

	if len(cfg.Languages) > 0 {
		for _, lang := range cfg.Languages {
			if err := writeFile(filepath.Join(dir, fmt.Sprintf("%s.%s.md", *release.Name, lang.Code)), buff.Bytes()); err != nil {
				return err
			}
		}

	} else {
		return writeFile(filepath.Join(dir, *release.Name+".md"), buff.Bytes())
	}

	return nil
}

func writeFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}
