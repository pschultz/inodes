package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var nodes map[string]int

func main() {
	flag.Parse()
	root := flag.Arg(0)

	nodes = map[string]int{}

	dev := getDevice(root)
	fmt.Fprintf(os.Stderr, "Staying on device 0x%x\n", dev)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err := check(path, info, err, dev); err != nil {
			return err
		}

		count(path, info.IsDir(), root)

		return nil
	})

	for p, n := range nodes {
		fmt.Fprintf(os.Stdout, "%7d %s\n", n, p)
	}
}

func getDevice(path string) uint64 {
	i, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	s, ok := i.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatalf("%q doesn't support FileInfo.Sys()", path)
	}

	return s.Dev
}

func check(path string, info os.FileInfo, err error, dev uint64) error {
	var (
		s  *syscall.Stat_t
		ok bool
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %q: %s\n", path, err)
		return filepath.SkipDir
	}

	if s, ok = info.Sys().(*syscall.Stat_t); !ok {
		fmt.Fprintf(os.Stderr, "WARN: %q doesn't support FileInfo.Sys()\n", path)
		return filepath.SkipDir
	}

	if dev != s.Dev {
		fmt.Fprintf(os.Stderr, "Skipping %q\n", path)
		return filepath.SkipDir
	}

	return nil
}

func count(path string, isDir bool, root string) {
	p := path

	if isDir {
		if _, ok := nodes[p]; !ok {
			nodes[p] = 0
		}

		nodes[p] += 1
	}

	if path != root {
		count(filepath.Dir(p), true, root)
	}
}
