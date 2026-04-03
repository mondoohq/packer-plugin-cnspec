// Copyright Mondoo, Inc. 2026
// SPDX-License-Identifier: BUSL-1.1

package provisioner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertScoreToRiskThreshold(t *testing.T) {
	tests := []struct {
		name           string
		scoreThreshold int
		riskThreshold  int
		expected       int
	}{
		{
			name:           "score_threshold converts to risk_threshold",
			scoreThreshold: 80,
			riskThreshold:  0,
			expected:       20, // 100 - 80
		},
		{
			name:           "risk_threshold takes precedence over score_threshold",
			scoreThreshold: 80,
			riskThreshold:  30,
			expected:       30,
		},
		{
			name:           "score_threshold of 100 converts to 0",
			scoreThreshold: 100,
			riskThreshold:  0,
			expected:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertScoreToRiskThreshold(tt.scoreThreshold, tt.riskThreshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetermineScoreThreshold(t *testing.T) {
	tests := []struct {
		name          string
		onFailure     string
		riskThreshold int
		expected      int
	}{
		{
			name:          "default threshold is 100",
			onFailure:     "",
			riskThreshold: 0,
			expected:      100,
		},
		{
			name:          "on_failure continue sets threshold to 0",
			onFailure:     "continue",
			riskThreshold: 0,
			expected:      0,
		},
		{
			name:          "on_failure continue ignores risk_threshold",
			onFailure:     "continue",
			riskThreshold: 50,
			expected:      0,
		},
		{
			name:          "risk_threshold is used when set",
			onFailure:     "",
			riskThreshold: 80,
			expected:      80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineScoreThreshold(tt.onFailure, tt.riskThreshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScorePassesThreshold(t *testing.T) {
	tests := []struct {
		name       string
		score      uint32
		threshold  int
		shouldPass bool
	}{
		{
			name:       "score above threshold passes",
			score:      90,
			threshold:  80,
			shouldPass: true,
		},
		{
			name:       "score equal to threshold passes",
			score:      80,
			threshold:  80,
			shouldPass: true,
		},
		{
			name:       "score below threshold fails",
			score:      50,
			threshold:  80,
			shouldPass: false,
		},
		{
			name:       "perfect score passes default threshold",
			score:      100,
			threshold:  100,
			shouldPass: true,
		},
		{
			name:       "any score passes with threshold 0",
			score:      0,
			threshold:  0,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scorePassesThreshold(tt.score, tt.threshold)
			assert.Equal(t, tt.shouldPass, result)
		})
	}
}
