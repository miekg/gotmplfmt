%%%
title = "gotmplfmt 1"
area = "user commands"
workgroup = "Go"
%%%

# NAME

gotmplfmt - format Go HTML template files

# SYNOPSIS

**gotmplfmt** [**FILE**]...

# DESCRIPTION

**gotmplfmt** formats Go HTML template file from **FILE**. If no file is given, it reads from standard input.

The zone is formatted according to the following rules:

- tabs are used for indentation
- the structure of open HTML tags and template verbs is followed, except for html, body, head, meta,
  main or nav tag
- before a {{block}} or {{define}} an extra newline is introduced

Note: you _can_ use this on Go text templates, but as whitespace is significant there, it will lead
to "corrupt" output.

# OPTIONS

There are is one options (and a debugging one):

`-w` _WIDTH_
: use _WIDTH_ as line width

`-t`
: show the lexed tokens

# EXAMPLE

    % cat <<'EOF' | ./gotmplfmt
    {{- if .Flash}} BLAAT {{else if not .Flash}} {{hallo}} {{meer}} BLOET {{end -}}
    EOF

Returns:

    {{- if .Flash}}
             BLAAT
    {{else if not .Flash}}
            {{hallo}}{{meer}} BLOET
    {{end -}}

# AUTHOR

Miek Gieben <miek@miek.nl>.
