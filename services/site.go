package services

import (
	"github.com/babilu-online/common/context"
)

type SiteService struct {
	context.DefaultService
	//
}

const SITE_SVC = "site_svc"

func (svc SiteService) Id() string {
	return SITE_SVC
}

func (svc *SiteService) Start() error {

	return nil
}

func (svc *SiteService) Index() ([]interface{}, error) {

	return nil, nil
}

func (svc *SiteService) Create() (interface{}, error) {

	return nil, nil
}

func (svc *SiteService) Delete() error {

	return nil
}
