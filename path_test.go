package fastrouter

import (
	"bytes"
	"testing"
)

func TestPath(t *testing.T) {
	buffer := new(bytes.Buffer)
	cases := []struct {
		in  string
		out string
	}{
		{in: "", out: "/"},
		{in: "/a", out: "/a"},
		{in: "/a/b", out: "/a/b"},
		{in: "/a/b/", out: "/a/b/"},
		{in: "..", out: "/"},
		{in: "../", out: "/"},
		{in: "../a", out: "/a"},
		{in: "./a", out: "/a"},
		{in: "..a", out: "/..a"},
		{in: "a.html", out: "/a.html"},
		{in: "..a/.b/", out: "/..a/.b/"},
		{in: "/a/..", out: "/"},
		{in: "/..//////././a", out: "/a"},
		{in: "/a/b/.////././/", out: "/a/b/"},
		{in: "/a/b/c/./././../", out: "/a/b/"},
		{in: "../../b/../..c/./././..../.../", out: "/"},
		{in: "/a/b/c/./././...../../../", out: "/"},
	}

	for _, c := range cases {
		if out := cleanPath(c.in, buffer); out != c.out {
			t.Fatal(c, out)
		}
	}
}
