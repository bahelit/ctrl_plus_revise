// Code generated by "stringer -linecomment -type=GPU"; DO NOT EDIT.

package hardware

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AMD-0]
	_ = x[NVIDIA-1]
	_ = x[noGPU-2]
}

const _GPU_name = "AMDNvidiaCPU"

var _GPU_index = [...]uint8{0, 3, 9, 12}

func (i GPU) String() string {
	if i < 0 || i >= GPU(len(_GPU_index)-1) {
		return "GPU(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GPU_name[_GPU_index[i]:_GPU_index[i+1]]
}