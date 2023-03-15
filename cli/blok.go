package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/babilu-online/common/context"
	blok_cli "github.com/blokhost/blok-cli"
	"github.com/blokhost/blok-cli/services"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Blok CLI module
// TODO Cli docs
// - build <site_id>
// - deploy <site_id>
// - site <create/show/update/delete> <site_id>
type Blok struct {
	context.DefaultService

	cfg *blok_cli.BlokConfig
	ctx *context.Context

	deploy *services.DeployService
	build  *services.BuildService
}

func (b Blok) Id() string {
	return "blok-cli"
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(ioutil.Discard)
	ctx, err := context.NewCtx(
		&services.GitService{},
		&services.AuthService{},
		&services.BuildService{},
		&services.DeployService{},
		&services.DiffService{},
		&services.SiteService{},
		&Blok{},
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.SetOutput(os.Stdout)
	err = ctx.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done!")
}

func (b *Blok) Start() error {
	flag.Parse()

	data, err := ioutil.ReadFile("./blok_config.json")
	if err != nil {
		return err
	}

	var cfg blok_cli.BlokConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	b.cfg = &cfg

	//Service bindings
	b.deploy = b.Service(services.DEPLOY_SVC).(*services.DeployService)
	b.build = b.Service(services.BUILD_SVC).(*services.BuildService)

	return b.Run()
}

func (b *Blok) Action() string {
	if len(os.Args) < 1 {
		return ""
	}

	action := os.Args[1]
	log.Printf("Action: %v", action)

	return action
}

func (b *Blok) Run() error {
	switch b.Action() {
	case "deploy":
		return b.Deploy()
	case "build":
		return b.Build()
	case "upload":
		return b.Upload()
	case "site":
		return b.SiteActions()
	}

	log.Println("Invalid action", b.Action())
	return errors.New("invalid action")
}

func (b *Blok) Deploy() error {
	//Download source
	// Git
	// TODO Blok Drive
	log.Println("Downloading source...")
	err := b.Download()
	if err != nil {
		return err
	}

	//Build source
	log.Println("Building site...")
	err = b.Build()
	if err != nil {
		log.Println("Failed to build: ", err)
		return err
	}

	//Upload to drive
	log.Println("Uploading site...")
	err = b.Upload()
	if err != nil {
		return err
	}

	log.Println("Site successfully deployed!")
	log.Printf("https://%s.blok.host", b.cfg.Site)
	return nil
}

//Build handles this service bindings
func (b *Blok) Download() error {
	tn := time.Now()
	defer func() {
		log.Printf("Download took: %s", time.Now().Sub(tn))
	}()

	if strings.Index(b.cfg.Src, "github.com") > -1 {
		gitSvc := b.Service(services.GIT_SVC).(*services.GitService)
		return gitSvc.Download(b.cfg.Src, ".")
	}

	//TODO Download from drive
	log.Println("Unable to download", b.cfg.Src)
	return errors.New("non github drives not currently supported")
}

//Build handles this service bindings
func (b *Blok) Build() error {
	tn := time.Now()
	defer func() {
		log.Printf("Build took: %s", time.Now().Sub(tn))
	}()

	return b.build.Build(b.cfg.Build.BuildPath, b.cfg.Build.Command)
}

//Upload handles this service bindings
func (b *Blok) Upload() error {
	tn := time.Now()
	defer func() {
		log.Printf("Upload took: %s", time.Now().Sub(tn))
	}()

	return b.deploy.Drive(b.cfg)
}

//SiteActions handles this service bindings
func (b *Blok) SiteActions() error {
	//TODO Switch between site options

	return b.deploy.Site(b.cfg)
}
