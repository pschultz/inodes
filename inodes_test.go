package main

import "testing"

type testcase struct {
	root         string
	paths        map[string]bool
	expectations map[string]int
	missing      []string
}

func testcases() []testcase {
	return []testcase{
		testcase{
			root: "/",
			paths: map[string]bool{
				"/":            true,
				"/foo":         true,
				"/foo/bar":     true,
				"/foo/bar/baz": false,
			},
			expectations: map[string]int{
				"/":        4,
				"/foo":     3,
				"/foo/bar": 2,
			},
			missing: []string{
				"/foo/bar/baz",
			},
		},
		testcase{
			root: "/foo/bar",
			paths: map[string]bool{
				"/foo/bar":     true,
				"/foo/bar/baz": false,
				"/foo/bar/qux": false,
			},
			expectations: map[string]int{
				"/foo/bar": 3,
			},
			missing: []string{
				"/",
				"/foo",
				"/foo/bar/baz",
				"/foo/bar/qux",
			},
		},
	}
}

func TestCount(t *testing.T) {

	for _, tc := range testcases() {
		nodes = map[string]int{}

		root := tc.root

		for p, isDir := range tc.paths {
			count(p, isDir, root)
		}

		for p, n := range tc.expectations {
			if nodes[p] != n {
				t.Errorf("Expected %q to have %d inodes, got %d", p, n, nodes[p])
			}
		}

		for _, p := range tc.missing {
			if _, ok := nodes[p]; ok {
				t.Errorf("Didn't expected %q in nodes map", p)
			}
		}
	}
}
