/*
Script to convert files contents into string contents.

One of the main things it does is replace annoying things like
backticks with ` + "`" + ` so that the golang string literal won't get
accidently escaped.

An alternative approach to this might be to copy the bytes as a
bytearray and then to convert it into a string, but that might take up
more space and be confusing to look at in the source code.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	backtick    = "`"
	quote       = `"`
	plus        = " + "
	replacement = backtick + plus + quote + backtick + quote + plus + backtick
)

var (
	infilename  string
	outfilename string
	outvarname  string
	pkgname     string
	appendfile  bool
)

func init() {
	flag.StringVar(&infilename, "src", "", "Source Filename (Required)")
	flag.StringVar(&outfilename, "dst", "", "Destination Filename (Required)")
	flag.StringVar(&outvarname, "var", "", "Output Variable Name (Required)")
	flag.StringVar(&pkgname, "pkg", "", "Package Name (Required) unless using append.")
	flag.BoolVar(&appendfile, "append", false, "Append to output file (Optional)")
}

func main() {
	flag.Parse()

	// Check for the required flags.
	if infilename == "" || outfilename == "" || outvarname == "" {
		flag.PrintDefaults()
		log.Println("src = ", infilename)
		log.Println("dst = ", outfilename)
		log.Println("src = ", outvarname)
		log.Fatal("error : all of the required flags need to be defined.")
	}

	// Check for package name if creating a file from scratch.
	if !appendfile && pkgname == "" {
		flag.PrintDefaults()
		log.Fatal("error : package name needs to be defined.")
	}

	// Copy the entire file contents into a string.
	b, err := ioutil.ReadFile(infilename)
	if err != nil {
		log.Fatal("error reading file:", err)
	}
	s := string(b)

	// Replace ` symbols to prevent escaping the string literal.
	s = strings.Replace(s, "`", replacement, -1)

	// Format the Code content into a constant declaration.
	code := fmt.Sprintf("const %s = `%s`", outvarname, s)

	// Output the source code to the destination file.
	if appendfile {

		// when appending the file, assume the front matter
		// has already been added.
		doAppendFile(outfilename, code)

	} else {

		// When overwriting a file from scratch,
		// add front-matter to the file.
		cmd := strings.Join(os.Args[1:], " ")
		code = fmt.Sprintf(frontMatter, cmd, pkgname) + code
		doOverwrite(outfilename, code)
	}
}

const frontMatter = `
// generated by "file2string.go"; DO NOT EDIT

// command: "file2string.go %s"

package %s

`

func doAppendFile(filename, code string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("error opening file to append : ", err)
	}
	defer f.Close()
	_, err = f.WriteString(code)
	if err != nil {
		log.Fatal("error when appending file : ", err)
	}
}

func doOverwrite(filename, code string) {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("error opening file to overwrite : ", err)
	}
	defer f.Close()
	_, err = f.WriteString(code)
	if err != nil {
		log.Fatal("error when overwriting file : ", err)
	}
}
