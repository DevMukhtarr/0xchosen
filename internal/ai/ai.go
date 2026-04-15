package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/devmukhtarr/0xchosen/internal/extractor"
	"github.com/devmukhtarr/0xchosen/internal/relationship"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"
const model = "moonshotai/kimi-k2-instruct"

type groqRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateRecon(facts *extractor.Facts, summary *relationship.Summary, apiKey string) (string, error) {
	prompt := buildPrompt(facts, summary)

	reqBody := groqRequest{
		Model: model,
		Messages: []message{
			{
				Role: "system",
				Content: `You are an expert smart contract security researcher. 
Your job is to generate detailed recon notes from structured Solidity contract facts.
Be precise, technical, and focus on what matters for security review.`,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", groqURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(body))
	}

	var groqResp groqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return groqResp.Choices[0].Message.Content, nil
}

func buildPrompt(facts *extractor.Facts, summary *relationship.Summary) string {
	var sb strings.Builder

	sb.WriteString("## Solidity Codebase Structured Facts\n\n")

	// Files in scope
	sb.WriteString("### Files in Scope\n")
	for _, f := range facts.Files {
		sb.WriteString(fmt.Sprintf("- %s\n", f.FilePath))
	}

	// Inheritance tree
	sb.WriteString("\n### Inheritance Tree\n")
	for _, inh := range summary.InheritanceTree {
		sb.WriteString(fmt.Sprintf("- %s\n", inh))
	}

	// Entry points
	sb.WriteString("\n### Public/External Entry Points\n")
	for _, ep := range summary.CallSurface {
		sb.WriteString(fmt.Sprintf("- %s\n", ep))
	}

	// Access control
	sb.WriteString("\n### Access Controlled Functions\n")
	for _, ac := range summary.AdminFunctions {
		sb.WriteString(fmt.Sprintf("- %s\n", ac))
	}

	// Key state variables
	sb.WriteString("\n### Key State Variables\n")
	for _, sv := range summary.KeyStateVars {
		sb.WriteString(fmt.Sprintf("- %s\n", sv))
	}

	// Detector warnings
	sb.WriteString("\n### Slither Detector Warnings\n")
	for _, w := range summary.DetectorWarnings {
		sb.WriteString(fmt.Sprintf("- %s\n", w))
	}

	// Dependency graph
	sb.WriteString("\n### Import Dependencies\n")
	for file, imports := range facts.DependencyGraph.Imports {
		if len(imports) > 0 {
			sb.WriteString(fmt.Sprintf("- %s imports: %s\n",
				file, strings.Join(imports, ", ")))
		}
	}

	sb.WriteString(`
## Task
Based on the structured facts above, generate detailed recon notes with these sections:

1. **Protocol Overview** — What does this protocol do?
2. **Core Modules & Their Roles** — Break down each contract and its responsibility
3. **Trust Assumptions & Roles** — Who has admin/owner access? What can they do?

Be specific and technical. Reference actual contract and function names that exist, Do not introduce
what you are not sure about.
.`)

	return sb.String()
}
