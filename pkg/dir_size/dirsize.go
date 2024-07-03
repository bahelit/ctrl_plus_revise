package dirSize

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type DirectoryInfo struct {
	DirCount  int64
	FileCount int64
	TotalSize int64
}

// GetDirInfo Traverse each root of the file tree in parallel returning the
// total file count and total size of all the files.
func GetDirInfo(fSys fs.FS) (dirInfo DirectoryInfo, errs []error) {
	err := fs.WalkDir(fSys, ".", func(p string, d fs.DirEntry, err error) error {
		fileInfo, walkErr := d.Info()
		if walkErr != nil && !errors.Is(err, fs.ErrNotExist) {
			_, _ = fmt.Fprintf(os.Stderr, "error reading file info: %v error: %v\n", p, walkErr)
			errs = append(errs, walkErr)
		}

		if fileInfo.IsDir() {
			dirInfo.DirCount++
			return nil
		}

		dirInfo.FileCount++
		dirInfo.TotalSize += fileInfo.Size()
		return nil
	})
	if err != nil {
		errs = append(errs, err)
		return dirInfo, errs
	}

	return dirInfo, nil
}
