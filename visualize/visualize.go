package visualize

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-graphviz"

	"go-micrograd/value"
)

// HTMLFile renders the expression graph rooted at root and writes it to path as
// a standalone page, creating the enclosing directory when it does not exist
// yet. The SVG is inlined, so the file needs no network and no other assets.
func HTMLFile(root *value.Value, path string) error {
	svg, err := renderSVG(root)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return os.WriteFile(path, []byte(page(string(svg))), 0o644)
}

// renderSVG lays the expression graph out with Graphviz and returns the <svg>
// element on its own.
func renderSVG(root *value.Value) ([]byte, error) {
	ctx := context.Background()

	graph, err := graphviz.ParseBytes([]byte(dot(root)))
	if err != nil {
		return nil, fmt.Errorf("parse dot: %w", err)
	}
	defer graph.Close()

	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("start graphviz: %w", err)
	}
	defer g.Close()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return nil, fmt.Errorf("render svg: %w", err)
	}

	// Graphviz emits an XML prolog and a DOCTYPE ahead of the element, neither
	// of which belongs inside an HTML page.
	out := buf.Bytes()
	if i := bytes.Index(out, []byte("<svg")); i > 0 {
		out = out[i:]
	}
	return out, nil
}

// dot renders the expression graph as Graphviz DOT source. Every Value becomes
// a record node holding its label and data, and every operation becomes a small
// node sitting between its operands and their result.
func dot(root *value.Value) string {
	var (
		b       strings.Builder
		visited = map[*value.Value]bool{}
		walk    func(*value.Value)
	)

	// Expressions are graphs rather than trees, so a Value shared by two
	// operations must be drawn once: the pointer is what identifies a node.
	walk = func(v *value.Value) {
		if visited[v] {
			return
		}
		visited[v] = true

		fmt.Fprintf(&b, "\t%q [shape=record label=\"{ %s | data %d }\"]\n", nodeID(v), escape(v.Label), v.Data)
		if v.Operation != "" {
			fmt.Fprintf(&b, "\t%q [label=%q]\n", opID(v), v.Operation)
			fmt.Fprintf(&b, "\t%q -> %q\n", opID(v), nodeID(v))
		}
		for _, child := range v.Prev {
			// An operand feeds the operation that produced its parent, not the
			// parent itself, so the operator shows up in between.
			fmt.Fprintf(&b, "\t%q -> %q\n", nodeID(child), opID(v))
			walk(child)
		}
	}

	b.WriteString("digraph {\n")
	b.WriteString("\trankdir=LR\n")
	walk(root)
	b.WriteString("}\n")
	return b.String()
}

func page(svg string) string {
	return `<!doctype html>
<meta charset="utf-8">
<title>micrograd expression graph</title>
<style>
  body { margin: 0; background: #fbfbfa; color: #1a1a19;
         font: 14px ui-sans-serif, system-ui, -apple-system, sans-serif; }
  header { padding: 14px 20px; border-bottom: 1px solid #e6e5e2; }
  h1 { margin: 0; font-size: 14px; font-weight: 600; letter-spacing: .01em; }
  main { padding: 24px 20px; overflow-x: auto; }
  svg { max-width: 100%; height: auto; }
</style>
<header><h1>micrograd expression graph</h1></header>
<main>
` + svg + `</main>
`
}

// nodeID names a Value by its address, which is unique for one program run.
func nodeID(v *value.Value) string { return fmt.Sprintf("n%p", v) }

// opID names the operation node drawn just before the Value it produced.
func opID(v *value.Value) string { return nodeID(v) + "op" }

// escape protects the characters that divide a record label into cells.
func escape(s string) string {
	return strings.NewReplacer(
		`\`, `\\`, `"`, `\"`, `{`, `\{`, `}`, `\}`, `|`, `\|`, `<`, `\<`, `>`, `\>`,
	).Replace(s)
}
