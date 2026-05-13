package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

const (
	cacheName   = "projects.json"
	maxCacheAge = 6 * time.Hour
)

type cache struct {
	Params           Params
	MaxDirModifiedAt time.Time
	Projects         []Project
}

func TryCache(params *Params) []Project {
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("Cache does not exist or has expired -- skipping cache")
		return nil
	}
	cache := cache{}
	_ = wf.Cache.LoadJSON(cacheName, &cache)
	if !params.Equal(cache.Params) {
		log.Printf("Search params do not match cached -- skipping cache")
		return nil
	}
	if getLatestModifiedAt(params.ProjectsPaths()).After(cache.MaxDirModifiedAt) {
		log.Printf("Projects have been added/removed since last caching -- skipping cache")
		return nil
	}
	return cache.Projects
}

func SaveCache(params *Params, projects []Project) {
	err := wf.Cache.StoreJSON(cacheName, cache{*params, getLatestModifiedAt(params.ProjectsPaths()), projects})
	if err != nil {
		log.Printf("could not save cache: %s", err.Error())
	}
}

func getLatestModifiedAt(paths []string) time.Time {
	var latest time.Time
	for _, p := range paths {
		t := getModifiedAt(p)
		if t.After(latest) {
			latest = t
		}
	}
	return latest
}

func getModifiedAt(path string) time.Time {
	dir, err := os.Stat(path)
	if err != nil {
		log.Printf("could not verify modified time of projects dir: %s", err.Error())
		return time.Time{}
	}
	timespec := dir.Sys().(*syscall.Stat_t).Ctimespec
	return time.Unix(timespec.Unix())
}
