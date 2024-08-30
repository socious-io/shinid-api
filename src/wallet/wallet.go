package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"shin/src/config"
	"shin/src/shortner"
	"shin/src/utils"
	"strings"
)

type H map[string]interface{}

type Connect struct {
	ID      string
	URL     string
	ShortID string
}

func CreateDID() (string, error) {

	// Document template with public keys and empty services
	documentTemplate := H{
		"publicKeys": []H{
			{
				"id":      "auth-1",
				"purpose": "authentication",
			},
			{
				"id":      "issue-1",
				"purpose": "assertionMethod",
			},
		},
		"services": []interface{}{},
	}

	// First API call to create DID
	didRes, err := makeRequest("/cloud-agent/did-registrar/dids", H{
		"documentTemplate": documentTemplate,
	})
	if err != nil {
		return "", err
	}

	var didResponse H
	err = json.Unmarshal(didRes, &didResponse)
	if err != nil {
		return "", err
	}

	longFormDid := didResponse["longFormDid"].(string)

	// Second API call to publish DID
	pubRes, err := makeRequest(fmt.Sprintf("/cloud-agent/did-registrar/dids/%s/publications", longFormDid), H{})
	if err != nil {
		return "", err
	}

	var publishResponse H
	err = json.Unmarshal(pubRes, &publishResponse)
	if err != nil {
		return "", err
	}

	didRef := publishResponse["scheduledOperation"].(H)["didRef"].(string)
	return didRef, nil
}

func CreateConnection() (*Connect, error) {
	res, err := makeRequest("/cloud-agent/connections", H{"label": "Shin connection"})
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}
	data := body["data"].(H)
	c := &Connect{
		ID: data["connectionId"].(string),
		URL: strings.ReplaceAll(
			data["invitation"].(H)["invitationUrl"].(string),
			"https://my.domain.com/path",
			config.Config.Wellet.Connect,
		),
	}
	short, err := shortner.New(c.URL)
	if err != nil {
		return nil, err
	}
	c.ShortID = short.ShortID
	return c, nil
}

func ProofRequest(connectionID string) (string, error) {
	res, err := makeRequest("/cloud-agent/present-proof/presentations", H{
		"connectionId": connectionID,
		"proofs":       []H{},
		"options": H{
			"challenge": "A challenge for the holder to sign",
			"domain":    "shinid.com",
		},
	})
	if err != nil {
		return "", err
	}
	var body H

	if err := json.Unmarshal(res, &body); err != nil {
		return "", err
	}

	return body["presentationId"].(string), nil
}

func ProofVerify(presentID string) (H, error) {
	res, err := getRequest(fmt.Sprintf("cloud-agent/present-proof/presentations/%s", presentID))
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}
	if body["status"].(string) != "PresentationVerified" {
		return nil, fmt.Errorf("presentation not verified")
	}
	_, payload, err := utils.DecodeJWT(body["data"].([]string)[0])
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode data")
	}
	var data H
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	_, payload, err = utils.DecodeJWT(data["vp"].(H)["verifiableCredential"].([]string)[0])
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode vc")
	}
	var vc H
	if err := json.Unmarshal(payload, &vc); err != nil {
		return nil, err
	}
	return vc, nil
}

func makeRequest(path string, body H) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", config.Config.Wellet.Agent, path)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wellet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}

func getRequest(path string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", config.Config.Wellet.Agent, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wellet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}
