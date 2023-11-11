package main

import (
	"flag"
	"fmt"
	stdio "io"
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

func matchedRulesLogger(w stdio.Writer) func(error types.MatchedRule) {
	return func(err types.MatchedRule) {
		msg := err.ErrorLog()
		fmt.Fprintf(w, "[%s] %s\n", err.Rule().Severity(), msg)
	}
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

func createWAF(directivesFile, matchedRulesLogDst string) (coraza.WAF, error) {
	matchedRulesLogWriter, err := resolveWriter(matchedRulesLogDst)
	if err != nil {
		return nil, fmt.Errorf("resolving error log destination: %w", err)
	}

	wafConfig := coraza.NewWAFConfig().
		WithRootFS(mergefs.Merge(coreruleset.FS, io.OSFS)).
		WithErrorCallback(matchedRulesLogger(matchedRulesLogWriter))

	if directivesFile != "" {
		wafConfig = wafConfig.WithDirectivesFromFile(directivesFile)
	}

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		return nil, fmt.Errorf("creating WAF: %w", err)
	}

	return waf, nil
}

func resolveWriter(errLogTarget string) (stdio.Writer, error) {
	switch errLogTarget {
	case "/dev/stdout", "":
		return os.Stdout, nil
	case "/dev/stderr":
		return os.Stderr, nil
	default:
		return os.OpenFile(errLogTarget, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	}
}

func main() {
	port := flag.Int("port", getEnvInt("PORT", 8080), "Port to listen on")
	matchedRulesLogDestination := flag.String("matched-rules-log", getEnvString("MATCHED_RULES_LOG", "/dev/stdout"), "Destination to log matched rules to")
	directivesFile := flag.String("directives", getEnvString("DIRECTIVES_FILE", ""), "Directives file to use")

	// parse flags from command line
	flag.Parse()

	waf, err := createWAF(*directivesFile, *matchedRulesLogDestination)
	if err != nil {
		log.Fatal(err)
	}

	app := httpbin.New()

	// handle route using handler function
	http.Handle("/", corazahttp.WrapHandler(waf, app.Handler()))

	// listen to port
	log.Printf("Listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
