# Xsh - A pluggable interactive shell written Go

`Xsh` (or Go shell) is a fork of Gosh (vladimirvivien/gosh) framework that uses Go's plugin system to create
for building interactive console-based shell programs.  A xsh shell is
comprised of a collection of Go plugins which implement one or more commands.
When `xsh` starts, it searches directory `./plugins` for available shared object
files that implement command plugins. Compatible with Windows and Cygwin systems

## Getting started

#### Pre-requisites

* Go 1.8 or above
* Linux
* Mac OSX
* Windows

Xsh makes it easy to create shell programs.  First, download or clone this 
repository.  For a quick start, run the following:

```bash
go run shell/xsh.go
```
This will produce the following output:
```bash
                hhh      	
                hhh      	
                hhh      	
XXX  XXX SSSSSS hhhhhhh  	
XXX  XXXSSS     hhh hhh 	
  XXXX   SSSSSS hhh  hhh 	
XXX  XXX     SSShhh  hhh 	
XXX  XXX SSSSSS hhh  hhh 	

No commands found
```
After the splashscreen is displayed, `xsh` informs you that `no commands found`, as expected.  Next,
exit the `xsh` shell (`Ctrl-C`) and let us compile the example plugins that comes with the source code.

Linux:
```bash
go build -buildmode=plugin  -o plugins/sys_command.so plugins/syscmd.go
```
Windows:
```cmd
go build -buildmode=plugin  -o plugins/sys_command.so plugins/syscmd.go
-buildmode=plugin not supported on windows/amd64
```
Finding a solution

The previous command will compile `plugins/syscmd.go` and outputs shared object
`plugins/sys_command.so`, as a Go plugin file.  Verify the shared object file was created:

```
> ls -lh plugins/
total 3.2M
-rw-rw-r-- 1  4.5K Mar 19 18:23 syscmd.go
-rw-rw-r-- 1  3.2M Mar 19 19:14 sys_command.so
-rw-rw-r-- 1  1.4K Mar 19 18:23 testcmd.go
```
Now, when xsh is restarted, it will dynamically load the commands implemented in the shared object file:

```bash
> go run shell/xsh.go
...

Loaded 4 command(s)...
Type help for available commands

xsh>
```

As indicated, typing `help` lists all available commands in the shell:

```bash
xsh> help

help: prints help information for other commands.

Available commands
------------------
      prompt:	sets a new shell prompt
         sys:	sets a new shell prompt
        help:	prints help information for other commands.
        exit:	exits the interactive shell immediately

Use "help <command-name>" for detail about the specified command
```
## A command
A Xsh `Command` is represented by type `api/Command`:
```go
type Command interface {
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	Exec(context.Context, []string) (context.Context, error)
}
```

The Xsh framework searches for Go plugin files in the `./plugins` directory.  Each package plugin must 
export a variable named `Commands` which is of type  :
```go
type Commands interface {
  ...
	Registry() map[string]Command
}
```
Type `Commands` type returns a list of `Command` via the `Registry()`.  

The following shows example command file [plugins/testcmd.go](./plugins/testcmd.go). It implements
two commands via types `helloCmd` and `goodbyeCmd`. The commands are exported via type `testCmds` using
method `Registry()`:

```go
package main

import (
	"context"
	"fmt"
	"io"

	"github.com/ROMSDEV/xsh/api"
)

type helloCmd string

func (t helloCmd) Name() string      { return string(t) }
func (t helloCmd) Usage() string     { return `hello` }
func (t helloCmd) ShortDesc() string { return `prints greeting "hello there"` }
func (t helloCmd) LongDesc() string  { return t.ShortDesc() }
func (t helloCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	out := ctx.Value("xsh.stdout").(io.Writer)
	fmt.Fprintln(out, "hello there")
	return ctx, nil
}

type goodbyeCmd string

func (t goodbyeCmd) Name() string      { return string(t) }
func (t goodbyeCmd) Usage() string     { return t.Name() }
func (t goodbyeCmd) ShortDesc() string { return `prints message "bye bye"` }
func (t goodbyeCmd) LongDesc() string  { return t.ShortDesc() }
func (t goodbyeCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	out := ctx.Value("xsh.stdout").(io.Writer)
	fmt.Fprintln(out, "bye bye")
	return ctx, nil
}

// command module
type testCmds struct{}

func (t *testCmds) Init(ctx context.Context) error {
	out := ctx.Value("xsh.stdout").(io.Writer)
	fmt.Fprintln(out, "test module loaded OK")
	return nil
}

func (t *testCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"hello":   helloCmd("hello"),
		"goodbye": goodbyeCmd("goodbye"),
	}
}

var Commands testCmds
```

## License
MIT
