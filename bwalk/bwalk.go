package bwalk

import (
	"os"
	"path/filepath"
	"sort"
)

// Walk is like filepath.Walk but walks the file tree breadth-first.
func Walk(root string, walkFn filepath.WalkFunc) error {
	queue := make([]string, 1)
	queue[0] = root

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		info, err := os.Lstat(path)
		if err != nil {
			err = walkFn(path, nil, err)
			if err != nil && err != filepath.SkipDir {
				return err
			}
			continue
		}

		err = walkFn(path, info, nil)
		if err != nil {
			if info.IsDir() && err == filepath.SkipDir {
				continue
			}
			return err
		}

		if !info.IsDir() {
			continue
		}

		names, err := readDirNames(path)
		if err != nil {
			continue
		}

		for _, name := range names {
			filename := filepath.Join(path, name)
			queue = append(queue, filename)
		}
	}
	return nil
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
