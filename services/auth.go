package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/babilu-online/common/context"
	"github.com/gagliardetto/solana-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AuthService struct {
	context.DefaultService
	baseURL string

	jwtToken    string
	githubToken string

	rustWallet *solana.Wallet

	client *http.Client
}

const AUTH_SVC = "auth_svc"

func (svc AuthService) Id() string {
	return AUTH_SVC
}

func (svc *AuthService) Start() error {
	svc.client = &http.Client{Timeout: 5 * time.Second}
	svc.jwtToken = os.Getenv("JWT_TOKEN")
	svc.githubToken = os.Getenv("GITHUB_TOKEN")
	//svc.baseURL = "https://web2.blok.host"
	svc.baseURL = "http://localhost:9094"

	return nil
}

func (svc *AuthService) GithubToken() string {
	return svc.githubToken
}

type JWTClaims struct {
	Exp           int    `json:"exp"`
	Id            string `json:"id"`
	OrigIat       int64  `json:"orig_iat"`
	WalletAddress string `json:"wallet_addr"`
}

func (svc *AuthService) PublicKey() string {
	if svc.jwtToken != "" {
		jData := strings.Split(svc.jwtToken, ".")
		bData, err := base64.RawStdEncoding.DecodeString(jData[1])
		if err != nil {
			log.Println("Unable to get web2 base64 publickey", err)
			return ""
		}

		var vals JWTClaims
		err = json.Unmarshal(bData, &vals)
		if err != nil {
			log.Println("Unable to decode web2 publickey", err)
			return ""
		}
		return vals.WalletAddress
	}
	return svc.rustWallet.PublicKey().String()
}

func (svc *AuthService) SignMessage(payload []byte) ([]byte, error) {
	if svc.jwtToken != "" {
		return svc.signMessageWeb2(payload)
	}
	return svc.signMessageWeb3(payload)
}

func (svc *AuthService) SignTransaction(payload []byte) ([]byte, error) {
	if svc.jwtToken != "" {
		return svc.signTransactionWeb2(payload)
	}
	return svc.signTransactionWeb3(payload)
}

func (svc *AuthService) signTransactionWeb2(payload []byte) ([]byte, error) {
	rPayload := TransactionRequest{
		Transaction: SignRequest{
			Data: payload,
		},
	}

	data, err := json.Marshal(&rPayload)
	if err != nil {
		return nil, err
	}
	return svc.signWeb2(fmt.Sprintf("%s/v1/actions/transactions/sign", svc.baseURL), data)
}

type SignRequest struct {
	Data []byte `json:"data"`
	Type string `json:"type"`
}

type TransactionRequest struct {
	Transaction SignRequest `json:"transaction"`
}

type SignatureRequest struct {
	Message SignRequest `json:"message"`
}

type SignatureResponse struct {
	Signed string `json:"signed"`
}

func (svc *AuthService) signMessageWeb2(payload []byte) ([]byte, error) {
	rPayload := SignatureRequest{
		Message: SignRequest{
			Data: payload,
		},
	}

	data, err := json.Marshal(&rPayload)
	if err != nil {
		return nil, err
	}
	return svc.signWeb2(fmt.Sprintf("%s/v1/actions/messages/sign", svc.baseURL), data)
}

func (svc *AuthService) signWeb2(uri string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", svc.jwtToken))

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sigResp SignatureResponse
	err = json.Unmarshal(respData, &sigResp)
	if err != nil {
		return nil, err
	}

	s, _ := solana.SignatureFromBase58(sigResp.Signed)
	return []byte(s.String()), nil
}

func (svc *AuthService) signMessageWeb3(payload []byte) ([]byte, error) {
	sig, err := svc.rustWallet.PrivateKey.Sign(payload)
	if err != nil {
		return nil, err
	}

	return []byte(sig.String()), nil
}

func (svc *AuthService) signTransactionWeb3(payload []byte) ([]byte, error) {
	sig, err := svc.rustWallet.PrivateKey.Sign(payload)
	if err != nil {
		return nil, err
	}

	return []byte(sig.String()), nil
}
