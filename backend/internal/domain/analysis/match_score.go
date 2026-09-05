package analysis

import "fmt"

const (
	minMatchScore = 0
	maxMatchScore = 100
)

type MatchScore struct {
	value int
}

// NewMatchScore 通过统一入口维护分数范围，避免外部代码绕过业务校验。
func NewMatchScore(value int) (MatchScore, error) {
	if value < minMatchScore || value > maxMatchScore {
		return MatchScore{}, fmt.Errorf(
			"match score must be between %d and %d: %d",
			minMatchScore,
			maxMatchScore,
			value,
		)
	}

	return MatchScore{value: value}, nil
}

// Value 返回经过 Domain 校验的原始整数值。
func (score MatchScore) Value() int {
	return score.value
}
