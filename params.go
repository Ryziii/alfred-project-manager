package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	projectDirEnvVar      = "ProjectDirectory"
	detectCodeProjectsVar = "DetectCodeProjects"
	maxProjDepthEnvVar    = "MaxProjectDepth"
	maxResults            = 0
)

type Params struct {
	ProjectsDirs     []string
	DetectCodeProjects bool
	MaxProjectDepth  uint
	MaxResults       uint
}

func NewParamsFromEnv() (*Params, error) {
	projDir := os.Getenv(projectDirEnvVar)
	if len(projDir) == 0 {
		return nil, fmt.Errorf("Please set project directories before using the workflow")
	}

	detectCode, err := strconv.ParseBool(os.Getenv(detectCodeProjectsVar))
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", detectCodeProjectsVar, err)
	}

	maxProjDepth, err := strconv.Atoi(os.Getenv(maxProjDepthEnvVar))
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %w", maxProjDepthEnvVar, err)
	}

	return &Params{
		ProjectsDirs:      splitDirs(projDir),
		DetectCodeProjects: detectCode,
		MaxProjectDepth:   uint(maxProjDepth),
		MaxResults:        maxResults,
	}, nil
}

func (p *Params) Equal(p2 Params) bool {
	if p.DetectCodeProjects != p2.DetectCodeProjects ||
		p.MaxProjectDepth != p2.MaxProjectDepth ||
		p.MaxResults != p2.MaxResults ||
		len(p.ProjectsDirs) != len(p2.ProjectsDirs) {
		return false
	}
	for i, d := range p.ProjectsDirs {
		if d != p2.ProjectsDirs[i] {
			return false
		}
	}
	return true
}

func (p *Params) ProjectsPaths() []string {
	paths := make([]string, len(p.ProjectsDirs))
	for i, d := range p.ProjectsDirs {
		if path.IsAbs(d) {
			paths[i] = d
		} else {
			paths[i] = path.Join(os.Getenv("HOME"), d)
		}
	}
	return paths
}

func splitDirs(raw string) []string {
	parts := strings.Split(raw, ":")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
