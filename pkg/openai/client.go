package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
}

type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type Response struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewGptClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) ChatComplete(ctx context.Context, messages []*Message) (string, error) {

	request := &Request{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	payload, _ := json.Marshal(request)
	payloadReader := bytes.NewReader(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", payloadReader)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}

	choice := response.Choices[0].Message.Content
	choice = strings.Trim(strings.TrimSpace(choice), `"`)

	return choice, nil

}
