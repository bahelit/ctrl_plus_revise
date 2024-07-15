package example_test

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/bahelit/ctrl_plus_revise/pkg/dir_size"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ExampleGetDirInfo() {
	tmpDir, err := os.MkdirTemp("", "example_test")
	check(err)
	tmpDir2, err := os.MkdirTemp(tmpDir, "example_test2")
	check(err)

	test123, err := os.CreateTemp(tmpDir, "test_123")
	check(err)
	testABC, err := os.CreateTemp(tmpDir2, "test_abc")
	check(err)
	defer func() {
		_ = test123.Close()
		_ = testABC.Close()
	}()

	dirInfo, _ := dirSize.GetDirInfo(os.DirFS("."))

	fmt.Printf("%d files  size: %1.f\n", dirInfo.FileCount, float64(dirInfo.TotalSize))

	// Output:
	// 2 files  size: 0
}

func BenchmarkFilesInMemory(b *testing.B) {
	fSys := fstest.MapFS{
		"file.go":                {},
		"sub-folder/example.go":  {},
		"sub-folder2/another.go": {},
		"sub-folder2/file.go":    {},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dirSize.GetDirInfo(fSys)
	}
}
