package walker

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultSize = 512
)

type Walker interface {
	Walk(path string, info os.FileInfo, err error) error
	Paths() []string
}

type DebugWalker struct {
	pathS []string
}

func NewDebugWalker() Walker {
	var w = &DebugWalker{}
	w.pathS = make([]string, 0, defaultSize)
	return w
}

func (w *DebugWalker) Walk(path string, info os.FileInfo, err error) error {
	w.pathS = append(w.pathS, path)
	return nil
}

func (w *DebugWalker) Paths() []string {
	return w.pathS
}

type FilterWalker struct {
	pattern string
	pathS   []string
}

func NewFilterWalker(pattern string) Walker {
	var w = &FilterWalker{}
	w.pattern = pattern
	w.pathS = make([]string, 0, defaultSize)
	return w
}

func (w *FilterWalker) Walk(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, w.pattern) {
		w.pathS = append(w.pathS, path)
	}
	return nil
}

func (w *FilterWalker) Paths() []string {
	return w.pathS
}

type FilterPathWalker struct {
	pattern string
	pathS   []string
}

func NewFilterPathWalker(pattern string) Walker {
	var w = &FilterPathWalker{}
	w.pattern = pattern
	w.pathS = make([]string, 0, defaultSize)
	return w
}

func (w *FilterPathWalker) Walk(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, w.pattern) {
		var p, _ = filepath.Split(path)
		w.pathS = append(w.pathS, p)
	}
	return nil
}

func (w *FilterPathWalker) Paths() []string {
	return w.pathS
}
