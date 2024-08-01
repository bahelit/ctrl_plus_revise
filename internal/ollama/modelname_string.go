// Code generated by "stringer -linecomment -type=ModelName"; DO NOT EDIT.

package ollama

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Llama3Dot1-0]
	_ = x[CodeLlama-1]
	_ = x[CodeLlama13b-2]
	_ = x[CodeGemma-3]
	_ = x[DeepSeekCoder-4]
	_ = x[DeepSeekCoderV2-5]
	_ = x[Gemma-6]
	_ = x[Gemma2b-7]
	_ = x[Gemma2-8]
	_ = x[Gemma22B-9]
	_ = x[Llama3-10]
	_ = x[Llava-11]
	_ = x[Mistral-12]
	_ = x[Phi3-13]
}

const _ModelName_name = "llama3.1:latestcodellama:latestcodellama:13bcodegemma:7bdeepseek-coder:latestdeepseek-coder-v2:latestgemma:latestgemma:2bgemma2:latestgemma2:2bllama3:latestllava:latestmistral:latestphi3:latest"

var _ModelName_index = [...]uint8{0, 15, 31, 44, 56, 77, 101, 113, 121, 134, 143, 156, 168, 182, 193}

func (i ModelName) String() string {
	if i < 0 || i >= ModelName(len(_ModelName_index)-1) {
		return "ModelName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ModelName_name[_ModelName_index[i]:_ModelName_index[i+1]]
}