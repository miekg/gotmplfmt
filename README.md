# gotmplfmt

Fmt Go HTML templates. There are no options. The indenting used is 4 spaces.
It can deal with 1 left open HTML tag in (for instance) a partial template, if there are multiple
open HTML tags, all but the last will be automatically closed. This might be bug, but OTOH you can
argue you shouldn't really do that.

## Usage

1. `go build`
2. `./gotmplfmt < template.go.tmpl`

# Before

![Before fmt](Before.png)

# After

![After fmt](After.png)
