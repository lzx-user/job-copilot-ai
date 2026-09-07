package analysis

// AnalysisResult 是 8.21 阶段的最小分析结果；8.22 再扩展完整结构化字段。
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
