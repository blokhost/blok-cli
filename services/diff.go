package services

import (
	"github.com/babilu-online/common/context"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path         string    `json:"path"`
	FileName     string    `json:"file_name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
}

type DiffService struct {
	context.DefaultService
}

const DIFF_SVC = "diff_svc"

func (svc DiffService) Id() string {
	return DIFF_SVC
}

func (svc *DiffService) Start() error {

	return nil
}

type DiffResult struct {
	checked    map[string]struct{}
	ChangeSize int64 `json:"change_size"`

	Added   []*DiffFile `json:"added"`
	Updated []*DiffFile `json:"updated"`
	Removed []*DiffFile `json:"removed"`
}

type DiffFile struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (df *DiffFile) FilePath() string {
	return df.Path
	//return fmt.Sprintf("%s/%s", df.Path, df.Name)
}

func newDiffResult() *DiffResult {
	return &DiffResult{
		checked:    map[string]struct{}{},
		ChangeSize: 0,
		Added:      []*DiffFile{},
		Updated:    []*DiffFile{},
		Removed:    []*DiffFile{},
	}
}

//Diff returns the difference between the manifest and target path
func (svc *DiffService) Diff(manifestFiles map[string]File, path string) (*DiffResult, error) {
	dr := newDiffResult()

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		dr.onFile(manifestFiles, path, info)
		return nil
	})
	if err != nil {
		return nil, err
	}

	//Check for removed files
	dr.checkRemoved(manifestFiles)

	return dr, nil
}

func (r *DiffResult) onFile(manifestFiles map[string]File, path string, file os.FileInfo) {
	if file.IsDir() {
		return
	}

	r.checked[file.Name()] = struct{}{}

	mf, ok := manifestFiles[file.Name()]
	if !ok { //Added
		r.ChangeSize += file.Size()
		r.Added = append(r.Added, &DiffFile{path, file.Name()})
		return
	}

	if file.ModTime().After(mf.LastModified) { //Updated
		r.ChangeSize += file.Size() - mf.Size
		r.Updated = append(r.Updated, &DiffFile{path, file.Name()})
		return
	}
}

func (r *DiffResult) checkRemoved(manifestFiles map[string]File) {
	for mf, f := range manifestFiles {
		if _, ok := r.checked[mf]; ok {
			continue //File still exists
		}

		//File wasnt in checked map so was removed
		r.Removed = append(r.Removed, &DiffFile{f.Path, f.FileName})
	}
}
