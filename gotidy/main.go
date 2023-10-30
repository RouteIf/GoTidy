package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	tidy "github.com/RouteIf/GoTidy"
)

var (
	debug *bool = flag.Bool("debug", false, "Output debugging messages")
)

func main() {
	flag.Parse()

	t := tidy.New()
	defer t.Free()

	t.OutputXml(true)
	t.AddXmlDecl(false)
	t.QuoteAmpersand(true)
	t.TidyMark(false)
	t.CharEncoding(tidy.Utf8)
	t.AsciiChars(true)
	t.NumericEntities(true)
	t.FixUri(true)
	t.DropProprietaryAttributes(true)
	t.FixBackslash(true)
	t.JoinClasses(true)
	t.JoinStyles(true)
	t.ShowBodyOnly(tidy.True)

	in, _ := io.ReadAll(os.Stdin)
	output, err := t.Tidy(string(in))
	if *debug && err != nil {
		log.Fatal(err, output)
	}
	fmt.Println(output)
}
