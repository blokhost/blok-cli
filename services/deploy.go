package services

import (
	"errors"
	"fmt"
	"github.com/babilu-online/common/context"
	blok_cli "github.com/blokhost/blok-cli"
	"github.com/blokhost/blok-cli/providers"
	"github.com/blokhost/blok-cli/providers/shadow"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

type DeployTarget struct {
	Name     string
	Location string
}

type DeployService struct {
	context.DefaultService

	auth *AuthService
	diff *DiffService

	providers map[string]providers.Provider
}

const DEPLOY_SVC = "deploy_svc"

func (svc DeployService) Id() string {
	return DEPLOY_SVC
}

func (svc *DeployService) Start() error {
	svc.auth = svc.Service(AUTH_SVC).(*AuthService)
	svc.diff = svc.Service(DIFF_SVC).(*DiffService)

	svc.providers = map[string]providers.Provider{
		"shdw": &shadow.ShadowProvider{},
	}

	for _, p := range svc.providers {
		err := p.Start(svc.auth.PublicKey(), svc.auth.SignTransaction, svc.auth.SignMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

//SiteActions deploys the given site src to the network
//TODO Determine if we need these or just make it generic

func (svc *DeployService) Site(cfg *blok_cli.BlokConfig) error {
	if cfg.Dst == "" {
		//TODO Create
	}

	return nil
}

//Drive deploys the given drive src to the network
//TODO Determine if we need these or just make it generic
func (svc *DeployService) Drive(cfg *blok_cli.BlokConfig) error {
	log.Println("Deploying Drive")

	provider, ok := svc.providers[cfg.DstProvider]
	if !ok {
		return errors.New(fmt.Sprintf("invalid provider %s", cfg.DstProvider))
	}

	diffResult, err := svc.diff.Diff(nil, cfg.Build.Output)
	if err != nil {
		return err
	}

	if cfg.Dst == "" { //Create a drive if its doesnt exist
		log.Println("No destination specified, creating drive...")

		driveName := fmt.Sprintf("%s-%v", strings.ReplaceAll(cfg.Src, "/", "_"), time.Now().Unix())
		cfg.Dst, err = provider.Create(driveName, uint64(diffResult.ChangeSize))
		if err != nil {
			return err
		}
	}

	log.Printf("Change set size: %v", diffResult.ChangeSize)

	log.Println("Deleting old files: ", len(diffResult.Removed))
	for _, d := range diffResult.Removed {
		err := provider.DeleteFile(cfg.Dst, d.Name)
		if err != nil {
			return err
		}
	}

	log.Println("Updating modified files: ", len(diffResult.Updated))
	for _, d := range diffResult.Updated {
		data, err := ioutil.ReadFile(d.FilePath())
		if err != nil {
			return err
		}

		_, err = provider.EditFile(cfg.Dst, d.Name, data)
		if err != nil {
			return err
		}
	}

	log.Println("Uploading new files: ", len(diffResult.Added))
	errs := svc.UploadFiles(cfg, diffResult.Added)
	for _, e := range errs {
		log.Println(e)
	}
	if len(errs) > 0 {
		return errors.New("failed to upload files")
	}

	//for _, d := range diffResult.Added {
	//	data, err := ioutil.ReadFile(d.FilePath())
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = provider.UploadFile(cfg.Dst, d.Name, data)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (svc *DeployService) UploadFiles(cfg *blok_cli.BlokConfig, files []*DiffFile) []error {
	jobs := make(chan *DiffFile, 1000)
	results := make(chan error, 1000)
	var wg sync.WaitGroup

	provider, ok := svc.providers[cfg.DstProvider]
	if !ok {
		return []error{errors.New(fmt.Sprintf("invalid provider %s", cfg.DstProvider))}
	}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			svc.uploadWorker(workerID, jobs, results, provider, cfg)
		}(i)
	}

	for _, f := range files {
		jobs <- f
	}
	close(jobs)
	wg.Wait()
	close(results)

	errorList := []error{}
	for e := range results {
		log.Println(e)
		errorList = append(errorList, e)
	}

	return errorList
}

//UploadFile opens & uploads the file via the given storage provider
func (svc *DeployService) UploadFile(provider providers.Provider, cfg *blok_cli.BlokConfig, file *DiffFile) error {
	data, err := ioutil.ReadFile(file.FilePath())
	if err != nil {
		return err
	}

	_, err = provider.UploadFile(cfg.Dst, file.Name, data)
	if err != nil {
		return err
	}
	return nil
}

func (svc *DeployService) uploadWorker(id int, jobs <-chan *DiffFile, results chan<- error, provider providers.Provider, cfg *blok_cli.BlokConfig) {
	for j := range jobs {
		//log.Println("worker", id, "uploading", j.FilePath())

		err := svc.UploadFile(provider, cfg, j)
		if err != nil {
			results <- err
		}
	}
}
