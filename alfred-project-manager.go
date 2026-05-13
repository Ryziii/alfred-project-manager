package main

import (
	"io/ioutil"
	"path"
	"strings"

	aw "github.com/deanishe/awgo"
	"go.deanishe.net/fuzzy"
)

var (
	wf           *aw.Workflow
	fuzzyOptions []fuzzy.Option
)

func init() {
	fuzzyOptions = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(aw.MaxResults(maxResults), aw.SortOptions(fuzzyOptions...))
}

func scanDirAtPath(basePath string, workspace string, maxDepth uint, detectCode bool) []Project {
	projects := []Project{}
	fullPath := path.Join(basePath, workspace)
	files, _ := ioutil.ReadDir(fullPath)
	for _, file := range files {
		if !file.IsDir() || file.Name()[0:1] == "." {
			continue
		}

		p := path.Join(fullPath, file.Name())

		if detectCode {
			if IsProject(p, maxDepth) {
				projects = append(projects, NewProjectFromPath(p, workspace))
				continue
			}
			if maxDepth > 0 {
				projects = append(projects, scanDirAtPath(basePath, path.Join(workspace, file.Name()), maxDepth-1, detectCode)...)
			}
		} else {
			projects = append(projects, NewProjectFromPath(p, workspace))
			if maxDepth > 0 {
				projects = append(projects, scanDirAtPath(basePath, path.Join(workspace, file.Name()), maxDepth-1, detectCode)...)
			}
		}
	}
	return projects
}

func scanProjects(params *Params) []Project {
	projects := []Project{}
	seen := map[string]bool{}
	for _, dir := range params.ProjectsPaths() {
		for _, p := range scanDirAtPath(dir, "", params.MaxProjectDepth, params.DetectCodeProjects) {
			if seen[p.Path] {
				continue
			}
			seen[p.Path] = true
			projects = append(projects, p)
		}
	}
	return projects
}

func projects(params *Params) []Project {
	if cached := TryCache(params); len(cached) > 0 {
		return cached
	}

	projects := scanProjects(params)

	SaveCache(params, projects)
	return projects
}

func run() {
	query := strings.TrimSpace(wf.Args()[1])
	params, err := NewParamsFromEnv()
	if err != nil {
		wf.Fatalf(err.Error())
	}

	projs := projects(params)

	label := ""
	for _, project := range projs {
		label = prettyPath(project.Path)
		item := wf.NewFileItem(project.Path).
			Title(project.Name()).
			Subtitle(label).
			Match(prettyPath(project.Path)).
			Arg(project.Path).
			Var("url", project.URL).
			UID(project.Path).
			Valid(true)
		item.Cmd().Subtitle(label + " · Open in VSCode")
		item.Opt().Subtitle(label + " · Open in Kitty")
		item.Ctrl().Subtitle(label + " · Open repo URL")
		item.Shift().Subtitle(label + " · Open in ClaudeCode")
	}

	if query != "" {
		wf.Filter(query)
	}
	wf.WarnEmpty("No matching projects found", "Try a different query")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
