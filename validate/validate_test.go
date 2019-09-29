package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidText(t *testing.T) {
	type args struct {
		text string
	}
	for _, tt := range []struct {
		args     args
		expected bool
	}{
		{
			args:     args{text: "郡道美玲"},
			expected: true,
		},
		{
			args:     args{text: "郡道美玲なう"},
			expected: true,
		},
		{
			args:     args{text: "ツイートなう"},
			expected: true,
		},
		{
			args:     args{text: "群道美玲"},
			expected: false,
		},
		{
			args:     args{text: "群道美玲なう"},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("isValidText(%#v)", tt.args.text), func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidText(tt.args.text))
		})
	}
}
