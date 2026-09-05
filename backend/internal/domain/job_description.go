package domain

import (
	"errors"
	"strings"
)

var ErrEmptyJobDescription = errors.New("job description cannot be empty")

type JobDescription struct {
	content string
}

// NewJobDescription 保证进入核心流程的 JD 不是空白内容。
func NewJobDescription(content string) (JobDescription, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return JobDescription{}, ErrEmptyJobDescription
	}

	return JobDescription{content: content}, nil
}

func (description JobDescription) Content() string {
	return description.content
}
