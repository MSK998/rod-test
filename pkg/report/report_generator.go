package report

import (
	"bytes"
	"html/template"

	"github.com/go-rod/rod"
)

type ReportGenerator struct {
	Templates *template.Template
}

// Returns a new instance of the report generator.
func NewReportGenerator(t *template.Template) *ReportGenerator {
	return &ReportGenerator{t}
}

// Generates a single PDF with the templates and data.
func (rg *ReportGenerator) GeneratePDF(name string, templateName string, data interface{}) []byte {
	var buffer bytes.Buffer

	rg.Templates.ExecuteTemplate(&buffer, templateName, data)

	page := rod.New().MustConnect().MustPage()
	page.SetDocumentContent(buffer.String())
	pdfData := page.MustPDF(name)
	page.Browser().Close() // It's really important that we close the browser, as the process gets left dangling otherwise
	return pdfData
}
