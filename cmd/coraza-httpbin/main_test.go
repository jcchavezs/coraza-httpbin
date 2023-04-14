package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateWAF(t *testing.T) {
	t.Run("relative directives file", func(t *testing.T) {
		_, err := createWAF("./testdata/test.conf")
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

		_, err = createWAF(f.Name())
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

		_, err = createWAF(path.Join("./testdata", filepath.Base(f.Name())))
		require.NoError(t, err)
	})

	t.Run("absolute directives file", func(t *testing.T) {
		_, testFilepath, _, _ := runtime.Caller(0)
		_, err := createWAF(path.Join(filepath.Dir(testFilepath), "./testdata/test.conf"))
		require.NoError(t, err)
	})
}
