package epubsvc

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func Convert(epub, root string) (string, error) {
	dir, err := os.ReadDir(root)
	if err != nil {
		return "", errors.New("convert err " + err.Error())
	}

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return "", errors.New("convert err " + err.Error())
	}
	traverseDir(dir, root, pdfg)
	err = pdfg.Create()
	if err != nil {
		return "", errors.New("convert err " + err.Error())
	}
	path := strings.Replace(epub, "epub", "pdf", -1)
	err = pdfg.WriteFile(path)
	if err != nil {
		return "", errors.New("convert err " + err.Error())
	}
	return path, nil
}

func traverseDir(dir []fs.DirEntry, root string, pdfg *wkhtml.PDFGenerator) {
	for _, entry := range dir {
		if !entry.IsDir() {
			if !isAllowedType(entry.Name()) {
				continue
			}
			page := wkhtml.NewPage(fmt.Sprintf("%v/%v", root, entry.Name()))
			page.EnableLocalFileAccess.Set(true)
			page.FooterRight.Set("[page]")
			page.FooterFontSize.Set(10)
			page.Zoom.Set(0.95)
			page.Encoding.Set("UTF-8")
			pdfg.AddPage(page)
		} else {
			newDir, _ := os.ReadDir(fmt.Sprintf("%v/%v", root, entry.Name()))
			traverseDir(newDir, fmt.Sprintf("%v/%v", root, entry.Name()), pdfg)
		}
	}
}

func isAllowedType(in string) bool {
	return strings.Contains(in, "htm")
}
