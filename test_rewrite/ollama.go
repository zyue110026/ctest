package testrewrite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	"net/http"
	"os"
	"time"
)

// OllamaResponse represents the basic structure of Ollama API response
type OllamaResponse struct {
	Model       string `json:"model"`
	RemoteModel string `json:"remote_model"`
	RemoteHost  string `json:"remote_host"`
	CreatedAt   string `json:"created_at"`
	Message     struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Thinking      string `json:"thinking"`
	Done          bool   `json:"done"`
	DoneReason    string `json:"done_reason"`
	TotalDuration int64  `json:"total_duration"`
}

// CallOllama sends the prompt to Ollama API and returns the rewritten code
func CallOllama(prompt string) (string, error) {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = ctestglobals.OllamaModelDefault
	}

	url := "http://localhost:11434/api/chat"

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert Go developer rewriting Kubernetes e2e tests for dynamic configuration. Follow user instructions strictly. Output only Go code or NONE.",
			},
			{"role": "user", "content": OneShotUserExample},
			{"role": "assistant", "content": OneShotAssistantExample},
			{"role": "user", "content": OneShotUserExample2},
			{"role": "assistant", "content": OneShotAssistantExample2},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"options": map[string]interface{}{
			"temperature": 0.0,
			"num_ctx":     131072,
		},
		"stream": false,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(bodyBytes))

	// resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	// ✅ Handle non-200 HTTP codes
	if resp.StatusCode != http.StatusOK {
		errBody, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"ollama returned HTTP %d: %s",
			resp.StatusCode,
			string(errBody),
		)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("Ollama response:")
	// fmt.Println(string(respBytes))
	if err != nil {
		return "", fmt.Errorf("failed to read Ollama response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	outTrim := ollamaResp.Message.Content
	// fmt.Println("Ollama output:")
	// fmt.Println(outTrim)
	if outTrim == "" || outTrim == "NONE" {
		return "NONE", nil
	}

	return outTrim, nil
}
