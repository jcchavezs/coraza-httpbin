package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	coreruleset "github.com/corazawaf/coraza-coreruleset"
	"github.com/corazawaf/coraza-coreruleset/io"
	"github.com/corazawaf/coraza/v3"
	corazahttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/mccutchen/go-httpbin/v2/httpbin"
	"github.com/yalue/merged_fs"
)

var (
	port           int
	directivesFile string
)

func init() {
	// set a and b as flag int vars
	flag.IntVar(&port, "port", getEnvInt("PORT", 8080), "Port to listen on")
	flag.StringVar(&directivesFile, "directives", getEnvString("DIRECTIVES_FILE", ""), "Directives file to use")

	// parse flags from command line
	flag.Parse()
}

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

func main() {
	app := httpbin.New()

	wafConfig := coraza.NewWAFConfig().
		WithRootFS(merged_fs.NewMergedFS(coreruleset.FS, io.OSFS)).
		WithErrorCallback(logError)

	if directivesFile != "" {
		wafConfig = wafConfig.WithDirectivesFromFile(directivesFile)
	}

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		log.Fatal(err)
	}

	// handle route using handler function
	http.Handle("/", corazahttp.WrapHandler(waf, app.Handler()))

	// listen to port

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
