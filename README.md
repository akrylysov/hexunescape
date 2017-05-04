hexunescape
===========

hexunescape is a small tool to unescape hex escaped text files (e.g. Nginx logs).

Installation
------------

```
go get -u github.com/akrylysov/hexunescape/cmd/hexunescape
```

Usage
-----

```
usage: hexunescape [path]

path is optional, defaults to stdin.
```

Example
-------

Stdin as an input:

```
echo "\x22abc\x22" | hexunescape
"abc"
```

File as an input:

```
echo "{\x22foo\x22:\x22bar\x22}" > test
hexunescape test
{"foo":"bar"}
```