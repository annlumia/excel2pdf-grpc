// Package office2pdf ...
package office2pdf

import (
	"log"
	"path/filepath"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Excel ...
type Excel struct {
	app       *ole.IDispatch
	workbooks *ole.VARIANT
	xls       *ole.VARIANT
}

func (el *Excel) open(inFile string) (err error) {

	ole.CoInitialize(0)

	var unknown *ole.IUnknown

	unknown, err = oleutil.CreateObject("Excel.Application")
	if err != nil {
		return
	}

	el.app, err = unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return
	}

	_, err = oleutil.PutProperty(el.app, "Visible", false)
	if err != nil {
		return
	}

	_, err = oleutil.PutProperty(el.app, "DisplayAlerts", false)
	if err != nil {
		return
	}

	el.workbooks, err = oleutil.GetProperty(el.app, "Workbooks")
	if err != nil {
		return
	}

	el.xls, err = oleutil.CallMethod(el.workbooks.ToIDispatch(), "Open", inFile)
	if err != nil {
		return
	}

	return
}

// Export ...
func (el *Excel) Export(inFile string) (outFile string, err error) {
	outDir := filepath.Dir(inFile)

	err = el.open(inFile)
	if err != nil {
		log.Printf("E! Failed to open file. %s\n", err.Error())
		return
	}
	defer el.close()

	outFile = filepath.Join(outDir, filepath.Base(inFile+".pdf"))
	_, err = oleutil.CallMethod(el.xls.ToIDispatch(), "ExportAsFixedFormat", 0, outFile)
	if err != nil {
		log.Printf("E! Failed to export as fixed format. %s\n", err.Error())
		outFile = ""
		return
	}

	return
}

func (el *Excel) close() {

	if el.xls != nil {
		oleutil.MustPutProperty(el.xls.ToIDispatch(), "Saved", true)
	}

	if el.workbooks != nil {
		oleutil.MustCallMethod(el.workbooks.ToIDispatch(), "Close")
	}

	if el.app != nil {
		oleutil.MustCallMethod(el.app, "Quit")
		el.app.Release()
	}

	ole.CoUninitialize()
}
