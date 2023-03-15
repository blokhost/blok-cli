package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

type ManifestDummy struct {
	ID          string `json:"id"`
	ManifestDst string `json:"manifest_dst"`
	Config      struct {
		Build struct {
			BuildPath string `json:"build_path"`
			Output    string `json:"output"`
			Command   string `json:"command"`
		} `json:"build"`
		Deployment struct {
			Src string `json:"src"`
			Dst string `json:"dst"`
		} `json:"deployment"`
	} `json:"config"`
	Manifest struct {
		Root    string `json:"root"`
		Folders map[string]struct {
			Name    string        `json:"name"`
			Folders []interface{} `json:"folders"`
			Files   []string      `json:"files"`
		} `json:"folders"`
		Files map[string]File `json:"files"`
	} `json:"manifest"`
}

func loadDummy() (*ManifestDummy, error) {
	var md ManifestDummy

	data, err := ioutil.ReadFile("../dummy_manifest.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &md)
	if err != nil {
		return nil, err
	}

	return &md, nil
}

func TestDiff_Diff(t *testing.T) {
	d, err := loadDummy()
	if err != nil {
		t.Fatal(err)
	}
	dsvc := DiffService{}

	diff, err := dsvc.Diff(d.Manifest.Files, "../test/dist")
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Added: ", len(diff.Added))
	log.Printf("%+v\n", diff.Added)
	log.Println("Updated: ", len(diff.Updated))
	log.Printf("%+v\n", diff.Updated)
	log.Println("Removed: ", len(diff.Removed))
	log.Printf("%+v\n", diff.Removed)
}
