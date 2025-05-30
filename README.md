# gotmplfmt

Fmt Go HTML templates. There is one option. The indenting used is 1 tab - this allow your terminal's tab
width setting to do its work. The formatter is rather simple, there is no AST creation, it just iterates over
a list of tokens.

## Usage

1. `go build`
2. `./gotmplfmt < template.go.tmpl`

# Before

![Before fmt](Before.png)

# After

![After fmt](After.png)
