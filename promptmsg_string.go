// Code generated by "stringer -linecomment -type=PromptMsg"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CorrectGrammar-0]
	_ = x[MakeItAList-1]
	_ = x[MakeItFriendly-2]
	_ = x[MakeItFriendlyRedo-3]
	_ = x[MakeItProfessional-4]
	_ = x[MakeASummary-5]
	_ = x[MakeExplanation-6]
	_ = x[MakeExpanded-7]
	_ = x[MakeHeadline-8]
}

const _PromptMsg_name = "Correct GrammarMake it a ListMake it FriendlyMake it FriendlyMake it ProfessionalMake a SummaryMake an ExplanationExpand on the textMake a Headline"

var _PromptMsg_index = [...]uint8{0, 15, 29, 45, 61, 81, 95, 114, 132, 147}

func (i PromptMsg) String() string {
	if i < 0 || i >= PromptMsg(len(_PromptMsg_index)-1) {
		return "PromptMsg(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PromptMsg_name[_PromptMsg_index[i]:_PromptMsg_index[i+1]]
}
