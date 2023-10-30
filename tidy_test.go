package tidy

import (
	"strings"
	"testing"
)

var corruptedHtml string = "<title id='bob' class='frank'>Hello, 世界</title><p>Foo!"

func Test_Tidy(t *testing.T) {
	tdy := New()
	defer tdy.Free()

	output, _ := tdy.Tidy(corruptedHtml)

	if !strings.Contains(output.Output, "<html>") {
		t.Errorf("Unable to fix corrupted HTML")
	}
}

func Test_AddXmlDecl(t *testing.T) {
	tdy := New()
	defer tdy.Free()

	tdy.OutputXml(true)
	tdy.AddXmlDecl(true)
	output, _ := tdy.Tidy(corruptedHtml)

	if !strings.HasPrefix(output.Output, "<?xml") {
		t.Errorf("XML declaration was not added")
	}

	tdy.AddXmlDecl(false)
	output, _ = tdy.Tidy(corruptedHtml)
	if strings.HasPrefix(output.Output, "<?xml") {
		t.Errorf("XML declaration must be omitted")
	}
}

func Test_TidyMark(t *testing.T) {
	tdy := New()
	defer tdy.Free()

	tdy.TidyMark(true)
	output, _ := tdy.Tidy(corruptedHtml)

	if !strings.Contains(output.Output, "HTML Tidy for") {
		t.Errorf("Tidy mark was not added")
	}

	tdy.TidyMark(false)
	output, _ = tdy.Tidy(corruptedHtml)
	if strings.Contains(output.Output, "HTML Tidy for") {
		t.Errorf("Tidy mark must be omitted")
	}
}

func Test_Multibyte(t *testing.T) {
	tdy := New()
	defer tdy.Free()

	tdy.InputEncoding(Utf8)
	tdy.OutputEncoding(Utf8)
	output, _ := tdy.Tidy(corruptedHtml)

	if !strings.Contains(output.Output, "世界") {
		t.Errorf("The output is not in UTF-8 or unicode symbols were encoded")
	}
}
