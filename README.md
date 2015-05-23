# godir

DWIM chdir for gopath source directories.

# Installation

```sh
go get github.com/death/godir2
```

Place the following in your `.bashrc`:

```sh
function godir() {
    cd `godir2 -dir="$1"`
}
```

# Example

```sh
[death@sneeze ~]$ godir godir
[death@sneeze godir2]$ pwd
/home/death/dev/gocode/src/github.com/death/godir2
```

# License

MIT
