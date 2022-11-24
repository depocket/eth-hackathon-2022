package utils

import (
	"depocket.io/app/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	DepocketTransactionURL = "https://sdk-stg.depocket.io/api/v1/transactions"
	DepocketTokenURL       = "https://sdk-stg.depocket.io/api/v1/tokens"
	GeneralTimeout         = time.Duration(60) * time.Second
)

func FetchDepocketTransaction(request model.TransactionRequest) (*model.TransactionResponse, error) {
	client := &http.Client{
		Timeout: GeneralTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, DepocketTransactionURL, nil)
	if err != nil {
		return nil, err
	}
	param := url.Values{}
	if request.Limit != nil {
		param["limit"] = []string{fmt.Sprintf("%v", *request.Limit)}
	}
	if request.Chain != nil {
		param["chain"] = []string{*request.Chain}
	}
	if request.ToAddress != nil {
		param["to_address"] = []string{*request.ToAddress}
	}
	if request.Cursor != nil {
		param["cursor"] = []string{*request.Cursor}
	}
	if request.DecodedAction != nil {
		param["decoded_action"] = []string{*request.DecodedAction}
	}
	if request.Decoded != nil {
		param["decoded"] = []string{fmt.Sprintf("%v", *request.Decoded)}
	}

	req.URL.RawQuery = param.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", string(body))
	}
	txns := &model.TransactionResponse{}
	if err := json.Unmarshal(body, txns); err != nil {
		return nil, err
	}
	return txns, nil
}

func FetchDepocketToken(request model.TokenRequest) (*model.TokenResponse, error) {
	client := &http.Client{
		Timeout: GeneralTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, DepocketTokenURL, nil)
	if err != nil {
		return nil, err
	}
	param := url.Values{}
	if request.Chain != nil {
		param["chain"] = []string{*request.Chain}
	}
	if request.Addresses != nil {
		param["addresses"] = []string{*request.Addresses}
	}
	req.URL.RawQuery = param.Encode()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", string(body))
	}
	tokens := &model.TokenResponse{}
	if err := json.Unmarshal(body, tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}
