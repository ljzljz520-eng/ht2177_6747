package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"storeledger/internal/config"
	"storeledger/internal/domain"
	"storeledger/internal/parser"
	"storeledger/internal/service"
	"storeledger/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	c := config.Load()
	path := flag.NewFlagSet("store-ledger", flag.ContinueOnError)
	db := path.String("db", c.DatabasePath, "database path")
	_ = path.Parse(os.Args[2:])
	st, err := store.Open(*db)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer st.Close()
	svc, err := service.New(st, service.DeterministicClock("1970-01-01T00:00:00Z"), c.Reviewer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	switch os.Args[1] {
	case "import":
		runImport(svc)
	case "query":
		runQuery(svc)
	case "health":
		fmt.Println(svc.Health()["status"])
	default:
		usage()
	}
}

func runImport(svc *service.Service) {
	doc, err := parser.ParseCSV("CLI-BATCH", "CLI import", "stdin", bytes.NewBufferString("id,store,inspector,score,findings\nR1,S1,Alice,82,clean\n"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	result, err := svc.ImportAndValidate(doc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(result.Summary)
}
func runQuery(svc *service.Service) {
	page, err := svc.QueryRecords(domain.Query{Page: 1, PageSize: 20})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, record := range page.Items {
		fmt.Printf("%s %s %d\n", record.ID, record.Status, record.Score)
	}
}
func usage() { fmt.Println("store-ledger import|query|health [-db path]") }
