%%%
title = "gotmplfmt 1"
area = "Gnu"
workgroup = "Go"
%%%

# NAME

gotmplfmt - format Go HTML template files

# SYNOPSIS

**dnsfmt** [**FILE**]...

# DESCRIPTION

**gotmplfmt** formats Go HTML template file from **FILE**. If no file is given, it reads from standard input.

The zone is formatted according to the following rules:

- rule1

No semantic checks are done, this is purely text manipulation with some basic zone file syntax
understanding.

# OPTIONS

There are two debugging options:

`-d`
: show debug output (tabs are shown as +)

`-t`
: show the lexed tokens

# EXAMPLE

    % cat <<'EOF' | ./dnsfmt
    $TTL 6H
    $ORIGIN example.org.
    @       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 7200
                    IN      NS  ns
    example.org.            IN      NS  ns-ext.nlnetlabs.nl.
    EOF

Returns:

    $TTL 6H
    $ORIGIN example.org.
    @               IN   SOA        ns miek.miek.nl. (
                                       1712997354   ; serial  Sat, 13 Apr 2024 08:35:54 UTC
                                       4H           ; refresh
                                       1H           ; retry
                                       1W           ; expire
                                       2H           ; minimum
                                       )
                    IN   NS         ns
                    IN   NS         ns-ext.nlnetlabs.nl.

# AUTHOR

Miek Gieben <miek@miek.nl>.
