package report

// Will eventually be what we pass into the report generator but havent gotten there yet.
type Report struct {
	ReportName   string
	TemplateName string
	ReportData   interface{}
}
