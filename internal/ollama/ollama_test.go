package ollama_test

import (
	"fmt"
	"testing"

	ollama "github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

// TestStringToModel this test could fail if stringer needs to be run.
func TestStringToModel(t *testing.T) {
	tests := []struct {
		input    string
		expected ollama.ModelName
	}{
		{input: "0", expected: ollama.Llama3Dot2},
		{input: "1", expected: ollama.CodeLlama},
		{input: "2", expected: ollama.CodeLlama13b},
		{input: "3", expected: ollama.CodeGemma},
		{input: "4", expected: ollama.DeepSeekCoder},
		{input: "5", expected: ollama.DeepSeekCoderV2},
		{input: "6", expected: ollama.Gemma},
		{input: "7", expected: ollama.Gemma2b},
		{input: "8", expected: ollama.Gemma2},
		{input: "9", expected: ollama.Gemma22B},
		{input: "10", expected: ollama.LlamaVision},
		{input: "11", expected: ollama.Llama3Dot1},
		{input: "12", expected: ollama.Llama3Dot21B},
		{input: "13", expected: ollama.Llama3},
		{input: "14", expected: ollama.Llava},
		{input: "15", expected: ollama.Mistral},
		{input: "16", expected: ollama.MistralNemo},
		{input: "17", expected: ollama.NemoTronMini},
		{input: "18", expected: ollama.Phi3},
	}

	// Add more test cases as needed
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %s", tt.input), func(t *testing.T) {
			actual := ollama.StringToModel(tt.input)
			if actual != tt.expected {
				t.Errorf("StringToModel(%q) = %v; want %v", tt.input, actual, tt.expected)
			}
		})
	}

}
