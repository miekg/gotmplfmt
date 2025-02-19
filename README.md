# gotmplfmt

Fmt Go HTML templates. There are no options. The indenting used is 4 spaces.

* Tokens are put on the same line if starting with a text node, this continues until a container
  template node is seen.

## TODO

* When to keep things on a single line?
* Some extra spacing to air it out a little?
* Keep newlines, from the original, but squash them?
* {{end}} is now printed, with knowing if it was {{- end -}} or the like.
* Comment subtype.
* HTML indenting (when implemented) does not know about the template indents and vice versa
    We need to parse the html too when lexing the docs
