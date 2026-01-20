package testrewrite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// OllamaResponse represents the basic structure of Ollama API response
type OllamaResponse struct {
	Generated string `json:"generated"`
}

// CallOllamaChunk sends a single chunk of prompt to Ollama API and returns the rewritten code
func CallOllamaChunk(prompt string) (string, error) {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "deepseek-coder:33b" // default model if env not set
	}

	url := "http://localhost:11434/api/chat"
	// headers := map[string]string{"Content-Type": "application/json"}

	payload := map[string]interface{}{
		"model": model,

		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"options": map[string]interface{}{
			"temperature": 0.0,
			"num_ctx":     4096,
		},
		"stream": false,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Ollama response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return ollamaResp.Generated, nil
}

// CallOllamaWithChunks splits the file content into chunks and calls Ollama for each
func CallOllamaWithChunks(filePath, content string, maxLines int) (string, error) {
	lines := strings.Split(content, "\n")
	var outputSlices []string

	for i := 0; i < len(lines); i += maxLines {
		end := i + maxLines
		if end > len(lines) {
			end = len(lines)
		}

		chunk := strings.Join(lines[i:end], "\n")
		prompt := BuildPrompt(filePath, chunk) // your BuildPrompt function

		fmt.Printf("🔹 Processing lines %d-%d for file %s\n", i+1, end, filePath)

		out, err := CallOllamaChunk(prompt)
		if err != nil {
			return "", fmt.Errorf("chunk [%d-%d] failed: %w", i+1, end, err)
		}

		outTrim := strings.TrimSpace(out)
		if strings.EqualFold(outTrim, "NONE") || outTrim == "" {
			// LLM indicates no test needs rewrite
			return "NONE", nil
		}

		outputSlices = append(outputSlices, out)
	}

	finalOutput := strings.Join(outputSlices, "\n")
	return finalOutput, nil
}
