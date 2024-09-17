package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"

	"nicholasq.xyz/ai/internal/config"
)

type OllamaClient struct {
	client  *http.Client
	baseURL string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func NewOllamaClient() *OllamaClient {

	return &OllamaClient{
		client:  &http.Client{},
		baseURL: "http://localhost:11434",
	}
}

func (o *OllamaClient) Query(ctx context.Context, query string, config config.Config) (*AIResponse, error) {

	if config.Debug {
		color.Cyan(fmt.Sprintf("config: %v", config))
	}

	if !config.ChainOfThought {
		return o.simpleQuery(ctx, query, config)
	}

	// Stage One
	stageOneResponse, err := o.chatQuery(ctx, config.Model, []Message{
		{Role: "system", Content: stageOneChainOfThought},
		{Role: "user", Content: query},
	})
	if err != nil {
		return nil, fmt.Errorf("error in stage one: %w", err)
	}

	if config.Debug {
		color.Cyan(fmt.Sprintf("stage one response: \n%s\n", stageOneResponse))
	}

	// Stage Two
	stageTwoResponse, err := o.chatQuery(ctx, config.Model, []Message{
		{Role: "system", Content: stageOneChainOfThought},
		{Role: "user", Content: query},
		{Role: "assistant", Content: stageOneResponse},
		{Role: "system", Content: stageTwoChainOfThought},
		{Role: "user", Content: fmt.Sprintf("My initial query repeated:\n<query>\n%s\n</query>\n", query)},
	})
	if err != nil {
		return nil, fmt.Errorf("error in stage two: %w", err)
	}

	if config.Debug {
		color.Cyan(fmt.Sprintf("stage two response: \n%s\n", stageTwoResponse))
	}

	// Final Stage
	finalResponse, err := o.chatQuery(ctx, config.Model, []Message{
		{Role: "system", Content: stageOneChainOfThought},
		{Role: "user", Content: query},
		{Role: "assistant", Content: stageOneResponse},
		{Role: "system", Content: stageTwoChainOfThought},
		{Role: "user", Content: fmt.Sprintf("My initial query repeated:\n<query>\n%s\n</query>\n", query)},
		{Role: "assistant", Content: stageTwoResponse},
		{Role: "system", Content: finalStageChainOfThought},
		{Role: "user", Content: fmt.Sprintf("My initial query repeated for that last time:\n<query>\n%s\n</query>\n", query)},
	})
	if err != nil {
		return nil, fmt.Errorf("error in final stage: %w", err)
	}

	return &AIResponse{Text: finalResponse}, nil
}

func (o *OllamaClient) GetCapabilities() []string {
	return []string{"text-generation"}
}

func (o *OllamaClient) SetContext(context string) error {
	// Mock implementation
	return nil
}

func (o *OllamaClient) simpleQuery(ctx context.Context, query string, config config.Config) (*AIResponse, error) {
	response, err := o.chatQuery(ctx, config.Model, []Message{
		{Role: "user", Content: query},
	})
	if err != nil {
		return nil, fmt.Errorf("error in simple query: %w", err)
	}
	return &AIResponse{Text: response}, nil
}

func (o *OllamaClient) chatQuery(ctx context.Context, model string, messages []Message) (string, error) {
	reqBody, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/chat", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return chatResp.Message.Content, nil
}

const stageOneChainOfThought = `
You will receive my question(s) below, help solve my request by generating a detailed step-by-step plan.
Please ensure that your thought process is clear and detailed, as if you're instructing yourself on how to tailor an answer. Carefully read and understand the problem or question presented. Identify all relevant details, requirements, and objectives. List the key elements, facts, and data provided. Ensure no important information is overlooked. Examine the gathered information for patterns, relationships, or underlying principles. Consider how these elements interact or influence each other. Develop a plan or approach to solve the problem based on your analysis. Think about possible methods or solutions and decide on the most effective one. Implement your chosen strategy step by step. Apply logical reasoning and problem-solving skills to work towards a solution. If any information from previous conversations is necessary to complete this task, politely ask the me to provide that information, as you cannot access it independently. If you require more information from me, ask for it before producing the <THOUGHT> block.

Please begin your response by acknowledging these instructions and then proceed with your step-by-step analysis.
Do not return an answer, just return the thought process as if it's between you and yourself. Please provide your response strictly in the following format and respect the <THOUGHT> tags:
<THOUGHT>

[step by step plan of how to answer the user's message one per line, use bullet points and line breaks]

</THOUGHT>
`

const stageTwoChainOfThought = `
Reflect deeply on your own thought process that you created earlier in the <THOUGHT></THOUGHT> tags. If you believe your previous thoughts could be improved upon, then provide your new improvements in <IMPROVEMENT> tags. You can change ideas from your previous <THOUGHT> block. You can decide to not use some ideas, or you can add new ideas, or you can modify previous ideas. Consider if your previous <THOUGHT> block had enough information. Consider if the <THOUGHT> block did not request enough information from the user. If more information must be needed, then include the desire to request more information in the <IMPROVEMENT> tags.

Please provide your response strictly in the following format and respect the <IMPROVEMENT> tags:
<IMPROVEMENT>

[a detailed list of improvements to each step in your previous thought process.]

</IMPROVEMENT>
`

const finalStageChainOfThought = `
For the last time, reflect on your own thought process to provide an answer to me.
Your thought process that you produced was within the <THOUGHT> tags and <IMPROVEMENT> tags (if they are present).

Your task:
Provide an answer to my question(s) based on your thought process.

**Important:** Do not include the thought process or mention that you reviewed it in your final answer. Just provide an answer to me.";
`
