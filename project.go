package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	enry "github.com/go-enry/go-enry/v2"
)

type Project struct {
	Path      string `json:"path"`
	Workspace string `json:"workspace"`
	URL       string `json:"url"`
}

func NewProjectFromPath(fullPath string, workspace string) Project {
	url, err := RepoURL(fullPath)
	if err != nil {
		log.Println(err)
		url = ""
	}
	return Project{
		Path:      fullPath,
		Workspace: workspace,
		URL:       url,
	}
}

func (p *Project) Name() string {
	return path.Join(p.Workspace, path.Base(p.Path))
}

var projectMarkers = []string{
	// version control
	".git", ".hg", ".svn",
	// JavaScript / TypeScript
	"package.json", "yarn.lock", "package-lock.json", "pnpm-lock.yaml",
	"tsconfig.json", "jsconfig.json", ".eslintrc", ".eslintrc.js",
	".eslintrc.json", ".eslintrc.yaml", ".eslintrc.yml",
	"next.config.js", "next.config.mjs", "next.config.ts",
	"nuxt.config.js", "nuxt.config.ts",
	"vue.config.js", "vite.config.js", "vite.config.ts", "vitest.config.ts",
	"webpack.config.js", "rollup.config.js", "rollup.config.ts",
	"babel.config.js", "babel.config.json", ".babelrc",
	"jest.config.js", "jest.config.ts", "jest.config.json",
	"prettier.config.js", ".prettierrc", ".prettierrc.json",
	"tailwind.config.js", "tailwind.config.ts", "tailwind.config.mjs",
	"postcss.config.js", "postcss.config.mjs",
	"svelte.config.js", "astro.config.mjs", "astro.config.js",
	// Go
	"go.mod", "go.sum",
	// Rust
	"Cargo.toml", "Cargo.lock",
	// Python
	"pyproject.toml", "setup.py", "setup.cfg", "requirements.txt",
	"Pipfile", "Pipfile.lock", "poetry.lock", "pyrightconfig.json",
	// Java / Kotlin
	"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle",
	"settings.gradle.kts", "gradlew", "gradlew.bat",
	// Ruby
	"Gemfile", "Gemfile.lock", "Rakefile",
	// PHP
	"composer.json", "composer.lock", "artisan",
	// Elixir
	"mix.exs", "mix.lock",
	// Haskell
	"stack.yaml", "cabal.project", "Setup.hs",
	// C / C++
	"CMakeLists.txt", "Makefile", "Makefile.am", "configure", "configure.ac",
	"meson.build", "Makefile.in", "configure.in",
	// Swift
	"Package.swift", "Package.resolved",
	// .NET
	"*.sln", "*.csproj", "*.fsproj", "*.vbproj",
	// Scala
	"build.sbt",
	// Dart / Flutter
	"pubspec.yaml", "analysis_options.yaml",
	// Zig
	"build.zig", "build.zig.zon",
	// OCaml
	"dune-project", "dune", "Makefile",
	// Erlang
	"rebar.config", "erlang.mk",
	// Nim
	"*.nimble",
	// Lua
	"*.rockspec",
	// Nix
	"flake.nix", "default.nix", "shell.nix",
	// Docker
	"Dockerfile", ".dockerignore",
	// General
	".gitignore", ".gitattributes", ".editorconfig",
	".env.example", ".env.template",
	".node-version", ".nvmrc", ".ruby-version", ".python-version",
	".tool-versions", ".java-version",
	"README.md", "CHANGELOG.md", "CONTRIBUTING.md", "LICENSE",
}

func hasProjectMarkers(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		for _, marker := range projectMarkers {
			// handle wildcard patterns
			if strings.HasPrefix(marker, "*.") && strings.HasSuffix(name, marker[1:]) {
				return true
			}
			if name == marker {
				return true
			}
		}
	}
	return false
}

func ContainsVendorDirs(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if enry.IsVendor(e.Name()) {
			return true
		}
	}
	return false
}

func ContainsCodeFiles(dir string, maxDepth uint) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			if enry.IsVendor(e.Name()) {
				continue
			}
			if maxDepth == 0 {
				continue
			}
			if ContainsCodeFiles(filepath.Join(dir, e.Name()), maxDepth-1) {
				return true
			}
		} else {
			lang := enry.GetLanguage(e.Name(), nil)
			if enry.GetLanguageType(lang) == enry.Programming {
				return true
			}
		}
	}
	return false
}

func IsProject(dir string, maxDepth uint) bool {
	return hasProjectMarkers(dir) || ContainsVendorDirs(dir) || ContainsCodeFiles(dir, maxDepth)
}

func prettyPath(p string) string {
	home := os.Getenv("HOME")
	if home != "" && strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}
