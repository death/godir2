package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/death/godir2/bwalk"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var (
	verbose = flag.Bool("verbose", false, "Show debugging information.")
	root    = flag.String("root", "", "The root directory ($GOPATH/src if unspecified).")
	dir     = flag.String("dir", "", "The directory to match against.")

	match   = "."
	minDist = 1000
)

func main() {
	flag.Parse()
	if *dir == "" {
		flag.Usage()
		return
	}
	if *root == "" {
		*root = os.ExpandEnv("$GOPATH/src")
	}
	if err := bwalk.Walk(*root, walkDir); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", match)
}

func walkDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return nil
	}
	if !info.IsDir() {
		return nil
	}

	if info.Name() == ".git" {
		return filepath.SkipDir
	}
	if minDist == 0 {
		return filepath.SkipDir
	}

	rel, err := filepath.Rel(*root, path)
	if err != nil {
		rel = path
	}
	parts := strings.Split(rel, string(filepath.Separator))
	for i := range parts {
		evaluate(path, filepath.Join(parts[i:]...))
		evaluate(path, filepath.Join(parts[:i]...))
	}

	return nil
}

func evaluate(path, name string) {
	if name == "" || name == "." {
		return
	}
	if dist := levenshtein.DistanceForStrings([]rune(*dir), []rune(name), levenshtein.DefaultOptions); dist < minDist {
		if *verbose {
			fmt.Printf("%4d: %-30s %s\n", dist, name, path)
		}
		match, minDist = path, dist
	}
}
