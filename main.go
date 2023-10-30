// GoTidy is a cGo wrapper around libtidy
package tidy

/*
#include <tidy.h>
#include <tidybuffio.h>
#include <tidyplatform.h>
#include <errno.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

type TidyOutput struct {
	Output      string
	Diagnostics string
	Errors      bool
}

type Tidy struct {
	tdoc   C.TidyDoc
	errbuf C.TidyBuffer
}

func New() *Tidy {
	t := &Tidy{}
	t.tdoc = C.tidyCreate()
	C.tidyBufInit(&t.errbuf)

	return t
}

func (tidy *Tidy) Free() {
	C.tidyBufFree(&tidy.errbuf)
	C.tidyRelease(tidy.tdoc)
}

func (tidy *Tidy) Tidy(htmlSource string) (TidyOutput, error) {
	input := C.CString(htmlSource)
	defer C.free(unsafe.Pointer(input))

	var output C.TidyBuffer
	defer C.tidyBufFree(&output)

	var rc C.int = -1

	rc = C.tidySetErrorBuffer(tidy.tdoc, &tidy.errbuf) // Capture diagnostics

	if rc >= 0 {
		rc = C.tidyParseString(tidy.tdoc, (*C.tmbchar)(input)) // Parse the input
	}

	if rc >= 0 {
		rc = C.tidyCleanAndRepair(tidy.tdoc) // Tidy it up!
	}

	if rc >= 0 {
		rc = C.tidyRunDiagnostics(tidy.tdoc) // Kvetch
	}

	if rc > 1 { // If error, force output.
		if C.tidyOptSetBool(tidy.tdoc, C.TidyForceOutput, C.yes) == 0 {
			rc = -1
		}
	}

	if rc >= 0 {
		rc = C.tidySaveBuffer(tidy.tdoc, &output) // Pretty Print
	}

	if rc >= 0 {
		output := C.GoStringN((*C.char)(unsafe.Pointer(output.bp)), C.int(output.size))
		tidyOutput := TidyOutput{Output: output}
		if rc > 0 {
			tidyOutput.Diagnostics = C.GoStringN((*C.char)(unsafe.Pointer(tidy.errbuf.bp)), C.int(tidy.errbuf.size))
			tidyOutput.Errors = rc > 1
		}

		return tidyOutput, nil
	}

	return TidyOutput{}, os.NewSyscallError(fmt.Sprintf("A severe error (%d) occurred.\n", int(rc)), errors.New(string(rc)))
}
