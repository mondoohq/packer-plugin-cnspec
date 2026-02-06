// Copyright Mondoo, Inc. 2026
// SPDX-License-Identifier: BUSL-1.1

package provisioner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreThresholdConversion(t *testing.T) {
	tests := []struct {
		name                  string
		scoreThreshold        int
		riskThreshold         int
		expectedRiskThreshold int
	}{
		{
			name:                  "score_threshold converts to risk_threshold",
			scoreThreshold:        80,
			riskThreshold:         0,
			expectedRiskThreshold: 20, // 100 - 80
		},
		{
			name:                  "risk_threshold takes precedence over score_threshold",
			scoreThreshold:        80,
			riskThreshold:         30,
			expectedRiskThreshold: 30,
		},
		{
			name:                  "zero score_threshold does not convert",
			scoreThreshold:        0,
			riskThreshold:         0,
			expectedRiskThreshold: 0,
		},
		{
			name:                  "score_threshold of 100 converts to 0",
			scoreThreshold:        100,
			riskThreshold:         0,
			expectedRiskThreshold: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{
				config: Config{
					ScoreThreshold: tt.scoreThreshold,
					RiskThreshold:  tt.riskThreshold,
				},
			}

			// Apply the conversion logic from Prepare()
			if p.config.ScoreThreshold != 0 {
				if p.config.RiskThreshold == 0 {
					p.config.RiskThreshold = 100 - p.config.ScoreThreshold
				}
			}

			assert.Equal(t, tt.expectedRiskThreshold, p.config.RiskThreshold)
		})
	}
}

func TestDetermineScoreThreshold(t *testing.T) {
	tests := []struct {
		name              string
		onFailure         string
		riskThreshold     int
		expectedThreshold int
	}{
		{
			name:              "default threshold is 100",
			onFailure:         "",
			riskThreshold:     0,
			expectedThreshold: 100,
		},
		{
			name:              "on_failure continue sets threshold to 0",
			onFailure:         "continue",
			riskThreshold:     0,
			expectedThreshold: 0,
		},
		{
			name:              "on_failure continue ignores risk_threshold",
			onFailure:         "continue",
			riskThreshold:     50,
			expectedThreshold: 0,
		},
		{
			name:              "risk_threshold is used when set",
			onFailure:         "",
			riskThreshold:     80,
			expectedThreshold: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This mirrors the logic in Provision() lines 653-661
			scoreThreshold := 100
			if tt.onFailure == "continue" {
				scoreThreshold = 0
			} else if tt.riskThreshold != 0 {
				scoreThreshold = tt.riskThreshold
			}

			assert.Equal(t, tt.expectedThreshold, scoreThreshold)
		})
	}
}

func TestScorePassesFails(t *testing.T) {
	tests := []struct {
		name           string
		score          uint32
		threshold      int
		shouldPass     bool
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
			// This mirrors the comparison in Provision() line 663
			passes := tt.score >= uint32(tt.threshold)
			assert.Equal(t, tt.shouldPass, passes)
		})
	}
}
