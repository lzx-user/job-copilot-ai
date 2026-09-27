package analysis

import "testing"

func validResultParams() AnalysisResultParams {
	return AnalysisResultParams{
		MatchScore: 80, JobSummary: "岗位摘要", CoreRequirements: []string{"Go"},
		MatchedSkills: []string{"Vue"}, MissingSkills: []string{},
		ResumeSuggestions: []string{"补充量化结果"}, PreparationTopics: []string{"并发"},
		GreetingMessage: "继续加油",
	}
}

func TestAnalysisResultValidatesScoreAndRequiredFields(t *testing.T) {
	params := validResultParams()
	result, err := NewAnalysisResult(params)
	if err != nil {
		t.Fatalf("NewAnalysisResult() error = %v", err)
	}
	if result.MatchScore() != 80 {
		t.Fatalf("MatchScore() = %d, want 80", result.MatchScore())
	}

	params.MatchScore = 101
	if _, err := NewAnalysisResult(params); err == nil {
		t.Fatal("score above 100 should fail")
	}
	params = validResultParams()
	params.CoreRequirements = nil
	if _, err := NewAnalysisResult(params); err == nil {
		t.Fatal("missing core requirements should fail")
	}
}
