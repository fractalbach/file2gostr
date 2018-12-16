file2gostr
========================================

Converts any file to a 
[StringLiteral](https://golang.org/ref/spec#String_literals) 
in a `.go` file.


Can be useful as a `go generate` directive:

~~~go
//go:generate file2gostr -src "index.html" -dst "index.html.go" -var "data" -pkg "main"
~~~


Main use case is to convert html pages into strings, but it could be
used for any file.  The output file is in the form:

~~~go
package main

const data = `

`
~~~

The output filename, the package name, and the variable name can all
be controlled by using the command line flags.


Backtick symbols are handled automatically, and are converted like so:
~~~go
const data = `` + "`" + ``
~~~
