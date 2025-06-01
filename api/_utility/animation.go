package _utility

import (
	"fmt"
	"math"
	"strings"
)

type AnimationParams struct {
	Colors           []string
	TimingFunction   string
	Duration         float64
	CascadeDirection string
}

func generateKeyframePercentages(colorCount int) []int {
	if colorCount <= 1 {
		if colorCount == 1 {
			return []int{0, 100}
		}
		return []int{}
	}

	percentages := make([]int, colorCount)
	step := 100.0 / float64(colorCount-1)

	for i := 0; i < colorCount; i++ {
		percentages[i] = int(math.Round(float64(i) * step))
	}
	if colorCount > 0 {
		percentages[colorCount-1] = 100
	}
	return percentages
}

func GenerateColorKeyframeRules(colors []string) string {
	numColors := len(colors)
	if numColors == 0 {
		return ""
	}

	var keyframeRulesJoined string
	if numColors == 1 {
		escapedColor := escapePercent(colors[0])
		keyframeRulesJoined = fmt.Sprintf("0%% { fill: %s; }\n100%% { fill: %s; }", escapedColor, escapedColor)
	} else {
		percentages := generateKeyframePercentages(numColors)
		keyframeRules := make([]string, numColors)
		for i := 0; i < numColors; i++ {
			percent := percentages[i]
			escapedColor := escapePercent(colors[i])
			keyframeRules[i] = fmt.Sprintf("%d%% { fill: %s; }", percent, escapedColor)
		}
		keyframeRulesJoined = strings.Join(keyframeRules, "\n")
	}
	return keyframeRulesJoined
}

func CalculateAnimationDelay(params AnimationParams, lineIdx, charIdx, numberOfLines, numberOfCharsInLine int) float64 {
	var delay float64 = 0.0

	if params.Duration <= 0 || numberOfCharsInLine <= 0 {
		return 0.0
	}

	const defaultCascadeSpreadFactor = 0.5

	switch params.CascadeDirection {
	case "ltr":
		if numberOfCharsInLine > 1 {
			delay = (params.Duration * defaultCascadeSpreadFactor) * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
	case "rtl":
		if numberOfCharsInLine > 1 {
			delay = -0.5 * params.Duration * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
	case "ttb":
		lineDelayFactor := 0.0
		if numberOfLines > 1 {
			lineDelayFactor = (params.Duration * defaultCascadeSpreadFactor) * (float64(lineIdx) / float64(numberOfLines-1))
		}
		charInLineDelayFactor := 0.0
		if numberOfCharsInLine > 1 {
			charInLineDelayFactor = (params.Duration * 0.1) * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
		delay = lineDelayFactor + charInLineDelayFactor
	case "btt":
		lineDelayFactor := 0.0
		if numberOfLines > 1 {
			lineDelayFactor = (params.Duration * defaultCascadeSpreadFactor) * (float64(numberOfLines-1-lineIdx) / float64(numberOfLines-1))
		}
		charInLineDelayFactor := 0.0
		if numberOfCharsInLine > 1 {
			charInLineDelayFactor = (params.Duration * 0.1) * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
		delay = lineDelayFactor + charInLineDelayFactor
	case "diag-tlbr":
		if numberOfLines > 0 && numberOfCharsInLine > 0 {
			maxSumIdx := float64((numberOfLines - 1) + (numberOfCharsInLine - 1))
			if maxSumIdx <= 0 {
				maxSumIdx = 1
			}
			currentSumIdx := float64(lineIdx + charIdx)
			delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSumIdx / maxSumIdx)
		}
	case "diag-trbl":
		if numberOfLines > 0 && numberOfCharsInLine > 0 {
			maxSumIdx := float64((numberOfLines - 1) + (numberOfCharsInLine - 1))
			if maxSumIdx <= 0 {
				maxSumIdx = 1
			}
			currentSumIdx := float64(lineIdx + (numberOfCharsInLine - 1 - charIdx))
			delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSumIdx / maxSumIdx)
		}
	case "diag-bltr":
		if numberOfLines > 0 && numberOfCharsInLine > 0 {
			maxSumIdx := float64((numberOfLines - 1) + (numberOfCharsInLine - 1))
			if maxSumIdx <= 0 {
				maxSumIdx = 1
			}
			currentSumIdx := float64((numberOfLines - 1 - lineIdx) + charIdx)
			delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSumIdx / maxSumIdx)
		}
	case "diag-brtl":
		if numberOfLines > 0 && numberOfCharsInLine > 0 {
			maxSumIdx := float64((numberOfLines - 1) + (numberOfCharsInLine - 1))
			if maxSumIdx <= 0 {
				maxSumIdx = 1
			}
			currentSumIdx := float64((numberOfLines - 1 - lineIdx) + (numberOfCharsInLine - 1 - charIdx))
			delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSumIdx / maxSumIdx)
		}
	case "center-out":
		if numberOfCharsInLine > 1 {
			centerPosition := float64(numberOfCharsInLine-1) / 2.0
			distanceFromCenter := math.Abs(float64(charIdx) - centerPosition)
			if centerPosition > 0 {
				delay = (params.Duration * 0.3) * (distanceFromCenter / centerPosition)
			}
		}
	case "edges-in":
		if numberOfCharsInLine > 1 {
			centerPosition := float64(numberOfCharsInLine-1) / 2.0
			distanceFromCenter := math.Abs(float64(charIdx) - centerPosition)
			if centerPosition > 0 {
				delay = (params.Duration * 0.3) * ((centerPosition - distanceFromCenter) / centerPosition)
			}
		}
	default:
		if numberOfCharsInLine > 1 {
			delay = -0.5 * params.Duration * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
	}
	return delay
}
