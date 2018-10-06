package gucumber

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageImportPathIsExtractedFromFilePath(t *testing.T) {
	//arrange
	packageName := "github.com/username/reponame"
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		file := filepath.Join(p, "src/github.com/username/reponame/main.go")
		//act
		importPath := assembleImportPath(file)
		//assert
		assert.Equal(t, packageName, importPath)
	}
}

func TestPackageImportUsesModulesIfPresent(t *testing.T) {
	// Check that a go.mod doesn't exist. This is to guard against the future
	// where this project may be converted to use go modules.

	// Get the absolute path so we ensure we're deleting the same file at the
	// end of the test.
	goModPath, err := filepath.Abs("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	createGoMod := func(module string) {
		if _, err := os.Stat(goModPath); !os.IsNotExist(err) {
			t.Fatalf("go.mod file exists. remove it or update this test")
		}
		f, err := os.OpenFile("go.mod", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("could not create go.mod file: %s", err)
		}
		if _, err = f.WriteString("module " + module + "\n\n"); err != nil {
			t.Fatalf("could not write to go.mod: %s", err)
		}
	}

	//arrange
	createGoMod("example.com/package")
	defer os.Remove(goModPath)

	file := "internal/features/stuff/step_definitions.go"
	//act
	importPath := assembleImportPath(file)
	//assert
	assert.Equal(t, importPath, "example.com/package/internal/features/stuff")
}
