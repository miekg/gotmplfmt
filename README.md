# gotmplfmt

Fmt Go HTML templates. There are no options. The indenting used is 4 spaces.

## Usage

1. `go build`
2. `./gotmplfmt < template.go.tmpl`

## TODO

* Tests. As well as unit test for the various functions.
* Put code in other packages than `main`?
* Some extra spacing to air it out a little?
* {{end}} is now printed, with knowing if it was {{- end -}} or the like.
* Comment subtype? Needed to know when to add bit of extra spacing.
