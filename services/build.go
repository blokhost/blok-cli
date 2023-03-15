package services

import (
	"github.com/babilu-online/common/context"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type BuildService struct {
	context.DefaultService
	debug bool
}

const BUILD_SVC = "build_svc"

func (svc BuildService) Id() string {
	return BUILD_SVC
}

func (svc *BuildService) Start() error {
	svc.debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	return nil
}

func (svc *BuildService) Build(src, command string) error {
	log.Printf("Building site: %s %s", src, command)
	c := exec.Command("bash", "-c", command)
	if src != "." {
		c.Dir = src
	}
	if svc.debug {
		c.Stdout = log.Writer()
	}
	return c.Run()
}
