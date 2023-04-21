package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"sync"
	"time"

	"statement-service-poc/pkg/domain"
	"statement-service-poc/pkg/report"
	"statement-service-poc/pkg/router"

	"github.com/MSK998/goq"
)

// package global variables, these are only accessible from the handler package
var (
	processMutex      sync.Mutex
	statementsRunning bool
)

func init() {
	fmt.Println("Adding statement routes")
	router.RegisteredHandlers.Register(func(sm *http.ServeMux) {
		sm.HandleFunc("/statement/", StatementHandler)
	})
}

func StatementHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is for a single merchant statement
	if path.Base(r.URL.Path) != "statement" {
		SingleStatementHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		// Lock the mutex to stop other requests from interfering with data
		processMutex.Lock()
		defer processMutex.Unlock()

		// If some other request has the lock then return a conflict
		if statementsRunning {
			http.Error(w, "Statements have already been started by another request", http.StatusConflict)
			return
		}
		// Set the statements running to true
		statementsRunning = true
		// Pass off statement generation to a new thread to keep response times fast
		go StartStatementGeneration()
		// Eagerly respond with a 200
		w.WriteHeader(http.StatusOK)
	}
}

// A handler that will deal with a single merchant statement using a path param
func SingleStatementHandler(w http.ResponseWriter, r *http.Request) {
	// Gets the MID from the URL path
	mid := path.Base(r.URL.Path)

	// Parse the html templates that will eventually be used to generate reports.
	tmpl, err := template.ParseFiles("template/test.html", "template/statement_merchant_summary.html")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new report generator
	rg := report.NewReportGenerator(tmpl)

	summ := make([]domain.StatementMerchantSummary, 0)

	// Declare a data struct - this can be anything as the generatePDF function only needs an interface{}
	data := struct {
		Summaries []domain.StatementMerchantSummary
	}{
		Summaries: summ,
	}

	// Generate a single report
	// This will block the request while generating
	// Run times at the moment are about 2-3 sec
	pdf := rg.GeneratePDF(fmt.Sprintf("temp/%s_statement.pdf", mid), "statement_merchant_summary.html", data)

	w.Header().Add("Content-Type", "application/pdf")
	w.Write(pdf)

}

// A function that is called to start generating statements.
// Right now this is just generating several identical pdfs with different names.
// The real service will probably grab data from the database and generate using that.
func StartStatementGeneration() {

	// Parse the templates and return any errors
	tmpl, err := (template.ParseFiles("template/test.html", "template/statement_merchant_summary.html"))
	if err != nil {
		println(err.Error())
	}

	// create a new report generator
	reportGenerator := report.NewReportGenerator(tmpl)

	// Create some dummy data
	data := struct {
		Title string
		Items []string
	}{
		Title: "My Statement",
		Items: []string{"Hello", "I am a statement", "I can take many forms"},
	}

	numberToGenerate := 5
	start := time.Now()
	// Create an instance of GoQ, This will be the goroutine manager to limit the amount of concurrent goroutines running at once, This is set to 5 but can be any number
	goqManager := goq.New(5)
	for i := 0; i < numberToGenerate; i++ {
		// Wait for a space to be available in our 5 at a time queue. This is a blocking call.
		goqManager.Wait()
		// Once a space is opened up we create a goroutine that will generate a PDF and call done to free up some space.
		go func(num int) {
			defer goqManager.Done()
			reportGenerator.GeneratePDF(fmt.Sprintf("temp/report_%d.pdf", num), "test.html", data)
		}(i)
	}
	// Wait for all remaining goroutines to complete.
	goqManager.WaitAllDone()
	totalRunTime := time.Since(start)

	// Lock the sync lock and modify the statements running variable.

	processMutex.Lock()
	defer processMutex.Unlock()
	statementsRunning = false

	avgTime := time.Duration(totalRunTime.Nanoseconds() / int64(numberToGenerate))

	// Report back that the statement generation has comepleted successfully
	fmt.Println("Statement Generation Complete in", totalRunTime.Round(time.Millisecond), "Average time per report:", avgTime.Round(time.Millisecond))
}
