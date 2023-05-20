package translator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type azureTranslator struct {
	apiKey string
}

func (t azureTranslator) Translate(contentTH string) (string, error) {
	url := "https://api.cognitive.microsofttranslator.com/translate"
	method := "POST"
	body := []struct {
		Text string `json:"text"`
	}{{Text: contentTH}}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", t.apiKey)
	req.Header.Set("Ocp-Apim-Subscription-Region", "southeastasia")
	q := req.URL.Query()
	q.Add("api-version", "3.0")
	q.Add("from", "th")
	q.Add("to", "en")
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var response []struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if len(response) > 0 && len(response[0].Translations) > 0 {
		return response[0].Translations[0].Text, nil
	}

	return "", nil
}

func NewAzureTranslator(apiKey string) Translator {
	return &azureTranslator{apiKey: apiKey}
}
