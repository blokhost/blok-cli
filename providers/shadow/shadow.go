package shadow

import (
	"errors"
	"fmt"
	shdw "github.com/blokhost/shdw-go"
	"github.com/gagliardetto/solana-go"
	http "github.com/valyala/fasthttp"
	"log"
)

type ShadowProvider struct {
	client *http.Client

	shdw *shdw.ShadowDrive
}

func (svc *ShadowProvider) Start(owner string, txnSignerFn shdw.SignerFunc, msgSignerFn shdw.SignerFunc) error {
	svc.shdw = &shdw.ShadowDrive{}
	if err := svc.shdw.Start(); err != nil {
		return err
	}

	//SetSigner must be called after start
	return svc.shdw.SetSigner(owner, txnSignerFn, msgSignerFn)
}

func (svc *ShadowProvider) Index(driveID string) ([]*shdw.DriveFile, error) {
	return svc.shdw.DriveFiles(driveID)
}

func (svc *ShadowProvider) UploadFile(driveID string, fileName string, data []byte) (string, error) {
	pk, err := solana.PublicKeyFromBase58(driveID)
	if err != nil {
		return "", err
	}

	//log.Printf("%s: Uploading %s", pk, fileName)
	resp, err := svc.shdw.UploadFile(pk, fileName, data)
	if err != nil {
		log.Printf("%s: Upload %s err: %s", pk, fileName, err)
		return "", err
	}

	//Return first upload error
	if len(resp.UploadErrors) > 0 {
		combiErr := ""
		first := true
		for _, e := range resp.UploadErrors {
			if first {
				combiErr = e.Error
				first = false
			}
			combiErr = fmt.Sprintf("%s,%s", combiErr, e.Error)
		}
		return "", errors.New(combiErr)
	}

	return fileName, nil
}

func (svc *ShadowProvider) EditFile(driveID string, fileName string, data []byte) (string, error) {
	pk, err := solana.PublicKeyFromBase58(driveID)
	if err != nil {
		return "", err
	}

	_, err = svc.shdw.EditFile(pk, fileName, data)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (svc *ShadowProvider) DeleteFile(driveID string, fileName string) error {
	pk, err := solana.PublicKeyFromBase58(driveID)
	if err != nil {
		return err
	}

	_, err = svc.shdw.DeleteFile(pk, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ShadowProvider) Create(name string, sizeBytes uint64) (string, error) {
	resp, err := svc.shdw.Create(name, sizeBytes)
	if err != nil {
		return "", err
	}

	return resp.ShdwBucket, nil
}

func (svc *ShadowProvider) Delete(driveID string) error {
	//pk, err := solana.PublicKeyFromBase58(driveID)
	//if err != nil {
	//	return err
	//}
	//
	//resp, err := svc.shdw.DeleteDrive(pk, owner, signerFn)
	//if err != nil {
	//	return err
	//}

	return nil
}
