package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	coreruleset "github.com/corazawaf/coraza-coreruleset"
	"github.com/corazawaf/coraza/v3"
	corazahttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/jcchavezs/mergefs"
	"github.com/jcchavezs/mergefs/io"
	"github.com/mccutchen/go-httpbin/v2/httpbin"
)

var (
	port           int
	directivesFile string
)

func logError(error types.MatchedRule) {
	msg := error.ErrorLog(0)
	fmt.Printf("[%s] %s", error.Rule().Severity(), msg)
}

func getEnvInt(name string, defaultValue int) int {
	if val := os.Getenv(name); val != "" {
		intVal, _ := strconv.Atoi(val)
		return intVal
	}

	return defaultValue
}

func getEnvString(name string, defaultValue string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}

	return defaultValue
}

func createWAF(directivesFile string) (coraza.WAF, error) {
	wafConfig := coraza.NewWAFConfig().
		WithRootFS(mergefs.Merge(coreruleset.FS, io.OSFS)).
		WithErrorCallback(logError)

	if directivesFile != "" {
		wafConfig = wafConfig.WithDirectivesFromFile(directivesFile)
	}

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		return nil, err
	}

	return waf, nil
}

func main() {
	flag.IntVar(&port, "port", getEnvInt("PORT", 8080), "Port to listen on")
	flag.StringVar(&directivesFile, "directives", getEnvString("DIRECTIVES_FILE", ""), "Directives file to use")

	// parse flags from command line
	flag.Parse()

	waf, err := createWAF(directivesFile)
	if err != nil {
		log.Fatal(err)
	}

	app := httpbin.New()

	// handle route using handler function
	http.Handle("/", corazahttp.WrapHandler(waf, app.Handler()))

	// listen to port
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
