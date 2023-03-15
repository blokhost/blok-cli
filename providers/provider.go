package providers

import shdw "github.com/blokhost/shdw-go"

type Provider interface {
	Start(owner string, txnSignerFn shdw.SignerFunc, msgSignerFn shdw.SignerFunc) error
	Create(name string, sizeBytes uint64) (string, error)
	Delete(driveID string) error
	DeleteFile(driveID, fileName string) error
	UploadFile(driveID, fileName string, data []byte) (string, error)
	EditFile(driveID, fileName string, data []byte) (string, error)
}
