package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var nodes map[string]int

func main() {
	flag.Parse()
	nodes = map[string]int{}

	d := flag.Arg(0)
	i, err := os.Stat(d)
	if err != nil {
		log.Fatal(err)
	}

	s, ok := i.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatalf("%q doesn't support FileInfo.Sys()", d)
	}

	dev := s.Dev
	fmt.Fprintf(os.Stderr, "Staying on device 0x%x\n", dev)

	filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if i, err := os.Stat(path); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %q: %s\n", path, err)
			return filepath.SkipDir
		} else {
			if s == nil {
				fmt.Fprintf(os.Stderr, "WARN: %q doesn't support FileInfo.Sys()\n", path)
				return filepath.SkipDir
			}

			if s, ok := i.Sys().(*syscall.Stat_t); !ok {
				fmt.Fprintf(os.Stderr, "WARN: %q doesn't support FileInfo.Sys()\n", path)
				return filepath.SkipDir
			} else {
				if dev != s.Dev {
					fmt.Fprintf(os.Stderr, "Skipping %q\n", path)
					return filepath.SkipDir
				}
			}
		}

		p := path
		for {
			if p == "" {
				break
			}
			if p != path || info.IsDir() {
				if _, ok := nodes[p]; !ok {
					nodes[p] = 0
				}
				nodes[p] += 1
			}

			n := strings.LastIndex(p, "/")
			if n == -1 || p == "/" {
				break
			}
			p = p[0:n]
		}

		return nil
	})

	for p, n := range nodes {
		fmt.Fprintf(os.Stdout, "%7d %s\n", n, p)
	}
}
