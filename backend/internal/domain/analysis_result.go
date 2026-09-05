package domain

// AnalysisResult 是当前架构阶段的最小分析结果，后续实现 JD 业务时再补充真实字段。
type AnalysisResult struct {
	matchScore MatchScore
}

func NewAnalysisResult(matchScore int) (AnalysisResult, error) {
	score, err := NewMatchScore(matchScore)
	if err != nil {
		return AnalysisResult{}, err
	}

	return AnalysisResult{matchScore: score}, nil
}

func (result AnalysisResult) MatchScore() int {
	return result.matchScore.Value()
}
