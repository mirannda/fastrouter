package fastrouter

import (
	"bytes"
	"net/http"
	"testing"
)

func handler(w http.ResponseWriter, req *http.Request) {}

func TestParseParam(t *testing.T) {
	type out struct {
		name    string
		value   string
		newPath string
	}

	cases := []struct {
		in  string
		out out
	}{
		{in: "{a}", out: out{name: "a", value: "", newPath: ""}},
		{in: "{a}/", out: out{name: "a", value: "", newPath: "/"}},
		{in: "{a}/b", out: out{name: "a", value: "", newPath: "/b"}},
		{in: "{a:}/b", out: out{name: "a", value: "", newPath: "/b"}},
		{in: "{a:\\d+}/b", out: out{name: "a", value: "\\d+", newPath: "/b"}},
	}

	for _, c := range cases {
		if name, value, newPath := parseParam(c.in); name != c.out.name ||
			value != c.out.value || newPath != c.out.newPath {
			t.Fatal(c)
		}
	}
}

func TestRadixTreeNode(t *testing.T) {
	// insert /aa/
	node := (*radixTreeNode)(nil).insert("/aa/{param1:}/{param2:\\d+}", handler)
	if node.chunk != "/aa/" || node.indices != "" || node.isParam ||
		node.regex != nil || node.handler != nil || len(node.children) != 1 {
		t.Fatal("case 0", node)
	}

	if child := node.children[0]; child.chunk != "param1" ||
		child.indices != "/" || !child.isParam || child.regex != nil ||
		child.handler != nil || len(child.children) != 1 {
		t.Fatal("case 1", child)
	}

	if child := node.children[0].children[0].children[0]; child.chunk != "param2" ||
		child.indices != "" || !child.isParam ||
		child.regex.String() != "\\d+" || child.handler == nil ||
		len(child.children) != 0 {
		t.Fatal("case 2", child, child.handler == nil)
	}

	// insert /ab/
	node = node.insert("/ab/", handler)
	if node.chunk != "/a" || node.indices != "ab" || node.isParam ||
		node.regex != nil || node.handler != nil {
		t.Fatal("case 3", node)
	}

	// insert /a
	node = node.insert("/a", handler)
	if node.chunk != "/a" || node.indices != "ab" || node.isParam ||
		node.regex != nil || node.handler == nil {
		t.Fatal("case 4", node)
	}

	// insert /a/
	node = node.insert("/a/", handler)
	if node.chunk != "/a" || node.indices != "/ab" || node.isParam ||
		node.regex != nil || node.handler == nil {
		t.Fatal("case 5", node)
	}

	// insert /ac/
	node = node.insert("/ac", handler)
	if node.chunk != "/a" || node.indices != "/abc" || node.isParam ||
		node.regex != nil || node.handler == nil {
		t.Fatal("case 5", node)
	}
}

func TestTree(t *testing.T) {
	tree := newRadixTree()

	tree.insert("/aa/{a}/{b:^[0-9]+$}", handler)
	tree.insert("/ab/", handler)
	tree.insert("/ac/", handler)
	tree.insert("/bb/", handler)

	cases := []struct {
		in  string
		out bool
	}{
		{in: "/aa/", out: false},
		{in: "/ab", out: false},
		{in: "/ab/", out: true},
		{in: "/aa/a/123", out: true},
		{in: "/aa/a/12aa", out: false},
	}

	for i, c := range cases {
		req, _ := http.NewRequest("GET", c.in, nil)
		buffer := new(bytes.Buffer)

		if h := tree.search(req, buffer); (h != nil) != (c.out) {
			t.Fatalf("case %d: ", i, c, h != nil)

			if i == 3 {
				req.ParseForm()
				if len(req.Form) != 2 || req.Form.Get("a") != "a" ||
					req.Form.Get("b") != "123" {
					t.Fatal(req.Form)
				}
			}
		}
	}
}
