package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/corazawaf/coraza/v3/types"
	"github.com/corazawaf/coraza/v3/types/variables"
	"github.com/stretchr/testify/require"
)

func TestCreateWAF(t *testing.T) {
	t.Run("relative directives file", func(t *testing.T) {
		_, err := createWAF("./testdata/test.conf", "")
		require.NoError(t, err)
	})

	t.Run("absolute directives file with absolute and relative include", func(t *testing.T) {
		tmpDir := t.TempDir()

		incf, err := os.CreateTemp(tmpDir, "relative-include.conf")
		require.NoError(t, err)
		defer incf.Close()

		f, err := os.CreateTemp(tmpDir, "test.conf")
		require.NoError(t, err)
		defer f.Close()

		directives, _ := os.ReadFile("./testdata/test.conf")
		_, err = f.Write(directives)
		require.NoError(t, err)
		_, err = f.WriteString(fmt.Sprintf("Include %s\n", filepath.Base(incf.Name())))
		require.NoError(t, err)
		_, err = f.WriteString(fmt.Sprintf("Include %s\n", incf.Name()))
		require.NoError(t, err)
		err = f.Sync()
		require.NoError(t, err)

		_, err = createWAF(f.Name(), "")
		require.NoError(t, err)
	})

	t.Run("relative directives file with absolute and relative include", func(t *testing.T) {
		_, testFilepath, _, _ := runtime.Caller(0)

		rincf, err := os.CreateTemp(path.Join(path.Dir(testFilepath), "./testdata"), "relative-include.conf")
		require.NoError(t, err)
		defer func() {
			rincf.Close()
			os.Remove(rincf.Name())
		}()

		aincf, err := os.CreateTemp(t.TempDir(), "relative-include.conf")
		require.NoError(t, err)
		defer aincf.Close()

		directives, _ := os.ReadFile("./testdata/test.conf")
		f, err := os.CreateTemp(path.Join(path.Dir(testFilepath), "./testdata"), "test.conf")
		require.NoError(t, err)
		_, err = f.Write(directives)
		require.NoError(t, err)
		_, err = f.WriteString(fmt.Sprintf("Include %s\n", filepath.Base(rincf.Name())))
		require.NoError(t, err)
		_, err = f.WriteString(fmt.Sprintf("Include %s\n", aincf.Name()))
		require.NoError(t, err)
		err = f.Sync()
		require.NoError(t, err)
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		_, err = createWAF(path.Join("./testdata", filepath.Base(f.Name())), "")
		require.NoError(t, err)
	})

	t.Run("absolute directives file", func(t *testing.T) {
		_, testFilepath, _, _ := runtime.Caller(0)
		_, err := createWAF(path.Join(filepath.Dir(testFilepath), "./testdata/test.conf"), "")
		require.NoError(t, err)
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		require.Equal(t, -1, getEnvInt("TEST1", -1))
	})

	t.Run("existing env", func(t *testing.T) {
		os.Setenv("TEST2", "1")
		defer os.Unsetenv("TEST2")
		require.Equal(t, 1, getEnvInt("TEST2", -11))
	})
}

func TestGetEnvString(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		require.Equal(t, "default", getEnvString("TEST1", "default"))
	})

	t.Run("existing env", func(t *testing.T) {
		os.Setenv("TEST2", "1")
		defer os.Unsetenv("TEST2")
		require.Equal(t, "1", getEnvString("TEST2", "default"))
	})
}

func TestResolveWriter(t *testing.T) {
	t.Run("stdout", func(t *testing.T) {
		w, err := resolveWriter("/dev/stdout")
		require.NoError(t, err)
		require.Equal(t, os.Stdout, w)
	})

	t.Run("stderr", func(t *testing.T) {
		w, err := resolveWriter("/dev/stderr")
		require.NoError(t, err)
		require.Equal(t, os.Stderr, w)
	})

	t.Run("file", func(t *testing.T) {
		tmpDir := t.TempDir()
		f, err := os.CreateTemp(tmpDir, "test.log")
		require.NoError(t, err)
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		w, err := resolveWriter(f.Name())
		require.NoError(t, err)
		require.Equal(t, f.Name(), w.(*os.File).Name())
	})
}

func TestWriteLog_OneEntryPerLine(t *testing.T) {
	logEntry := "error log entry"
	t.Run("file", func(t *testing.T) {
		tmpDir := t.TempDir()
		file, err := os.CreateTemp(tmpDir, "test.log")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
			os.Remove(file.Name())
		})

		w, err := resolveWriter(file.Name())
		require.NoError(t, err)

		matchedRule := &MatchedRule{
			Rule_: RuleMetadata{
				ID_: 1234,
			},
			MatchedDatas_: []types.MatchData{
				MatchData{},
			},
			ErrorLog_: logEntry,
		}
		callback := matchedRulesLogger(w)
		callback(matchedRule)
		callback(matchedRule)

		bytes, err := io.ReadAll(file)
		require.NoError(t, err)

		message := string(bytes)
		lines := strings.Split(message, "\n")
		require.Equal(t, 3, len(lines))
		require.Contains(t, lines[0], logEntry)
		require.Contains(t, lines[1], logEntry)
		require.Empty(t, lines[2])
	})
}

type MatchedRule struct {
	// Message is the macro expanded message
	Message_ string
	// Data is the macro expanded logdata
	Data_ string
	// URI is the full request uri unparsed
	URI_ string
	// TransactionID is the transaction ID
	TransactionID_ string
	// Disruptive is whether this rule will perform disruptive actions (note also pass, allow, redirect are considered disruptive actions)
	Disruptive_ bool
	// ServerIPAddress is the address of the server
	ServerIPAddress_ string
	// ClientIPAddress is the address of the client
	ClientIPAddress_ string
	// MatchedDatas is the matched variables.
	MatchedDatas_ []types.MatchData

	Rule_ types.RuleMetadata

	AuditLog_ string
	ErrorLog_ string
}

func (mr *MatchedRule) AuditLog() string {
	return mr.AuditLog_
}

func (mr *MatchedRule) Data() string {
	return mr.Data_
}

func (mr *MatchedRule) ErrorLog() string {
	return mr.ErrorLog_
}

func (mr *MatchedRule) TransactionID() string {
	return mr.TransactionID_
}

func (mr *MatchedRule) ClientIPAddress() string {
	return mr.ClientIPAddress_
}

func (mr *MatchedRule) Disruptive() bool {
	return mr.Disruptive_
}

func (mr *MatchedRule) MatchedDatas() []types.MatchData {
	return mr.MatchedDatas_
}

func (mr *MatchedRule) Message() string {
	return mr.Message_
}

func (mr *MatchedRule) Rule() types.RuleMetadata {
	return mr.Rule_
}

func (mr *MatchedRule) ServerIPAddress() string {
	return mr.ServerIPAddress_
}

func (mr *MatchedRule) URI() string {
	return mr.URI_
}

type MatchData struct {
	// Variable
	Variable_ variables.RuleVariable
	// Key of the variable, blank if no key is required
	Key_ string
	// Value of the current VARIABLE:KEY
	Value_ string
	// Message is the expanded macro message
	Message_ string
	// Data is the expanded logdata of the macro
	Data_ string
	// Chain depth of variable match
	ChainLevel_ int
}

func (md MatchData) Variable() variables.RuleVariable {
	return md.Variable_
}

func (md MatchData) Key() string {
	return md.Key_
}

func (md MatchData) Value() string {
	return md.Value_
}

func (md MatchData) Message() string {
	return md.Message_
}

func (md MatchData) Data() string {
	return md.Data_
}

func (md MatchData) ChainLevel() int {
	return md.ChainLevel_
}

type RuleMetadata struct {
	ID_       int
	File_     string
	Line_     int
	Revision_ string
	Severity_ types.RuleSeverity
	Version_  string
	Tags_     []string
	Maturity_ int
	Accuracy_ int
	Operator_ string
	Phase_    types.RulePhase
	Raw_      string
	SecMark_  string
}

func (rm RuleMetadata) ID() int {
	return rm.ID_
}

func (rm RuleMetadata) File() string {
	return rm.File_
}

func (rm RuleMetadata) Line() int {
	return rm.Line_
}

func (rm RuleMetadata) Revision() string {
	return rm.Revision_
}

func (rm RuleMetadata) Severity() types.RuleSeverity {
	return rm.Severity_
}
func (rm RuleMetadata) Version() string {
	return rm.Version_
}

func (rm RuleMetadata) Tags() []string {
	return rm.Tags_
}

func (rm RuleMetadata) Maturity() int {
	return rm.Maturity_
}

func (rm RuleMetadata) Accuracy() int {
	return rm.Accuracy_
}

func (rm RuleMetadata) Operator() string {
	return rm.Operator_
}

func (rm RuleMetadata) Phase() types.RulePhase {
	return rm.Phase_
}

func (rm RuleMetadata) Raw() string {
	return rm.Raw_
}

func (rm RuleMetadata) SecMark() string {
	return rm.SecMark_
}
