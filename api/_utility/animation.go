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

func CalculateAnimationDelay(params AnimationParams, lineIdx, charIdx, numberOfLines, numberOfCharsInLine int, maxLineLengthGlobal int) float64 {
	if params.Duration <= 0 {
		return 0.0
	}

	delay := 0.0
	defaultCascadeSpreadFactor := 0.5
	withinLineStaggerFactor := 0.1
	centerEdgeStaggerFactor := 0.3

	maxIdxGlobal := float64(maxLineLengthGlobal - 1)
	if maxIdxGlobal < 0 {
		maxIdxGlobal = 0
	}

	normNumLines := float64(numberOfLines - 1)
	if normNumLines <= 0 {
		normNumLines = 1
	}

	normMaxIdxGlobal := maxIdxGlobal
	if normMaxIdxGlobal <= 0 {
		normMaxIdxGlobal = 1
	}

	switch params.CascadeDirection {
	case "ltr":
		delay = (params.Duration * defaultCascadeSpreadFactor) * (float64(charIdx) / normMaxIdxGlobal)
	case "rtl":
		delay = (params.Duration * defaultCascadeSpreadFactor) * ((maxIdxGlobal - float64(charIdx)) / normMaxIdxGlobal)
	case "ttb":
		lineDelayFactor := (params.Duration * defaultCascadeSpreadFactor) * (float64(lineIdx) / normNumLines)
		charInLineDelayFactor := (params.Duration * withinLineStaggerFactor) * (float64(charIdx) / normMaxIdxGlobal)
		delay = lineDelayFactor + charInLineDelayFactor
	case "btt":
		lineDelayFactor := (params.Duration * defaultCascadeSpreadFactor) * (float64(numberOfLines-1-lineIdx) / normNumLines)
		charInLineDelayFactor := (params.Duration * withinLineStaggerFactor) * (float64(charIdx) / normMaxIdxGlobal)
		delay = lineDelayFactor + charInLineDelayFactor
	case "diag-tlbr":
		maxSumForNorm := float64(numberOfLines-1) + maxIdxGlobal
		if maxSumForNorm <= 0 {
			maxSumForNorm = 1
		}
		currentSum := float64(lineIdx + charIdx)
		delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSum / maxSumForNorm)
	case "diag-trbl":
		maxSumForNorm := float64(numberOfLines-1) + maxIdxGlobal
		if maxSumForNorm <= 0 {
			maxSumForNorm = 1
		}
		currentSum := float64(lineIdx) + (maxIdxGlobal - float64(charIdx))
		delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSum / maxSumForNorm)
	case "diag-bltr":
		maxSumForNorm := float64(numberOfLines-1) + maxIdxGlobal
		if maxSumForNorm <= 0 {
			maxSumForNorm = 1
		}
		currentSum := float64(numberOfLines-1-lineIdx) + float64(charIdx)
		delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSum / maxSumForNorm)
	case "diag-brtl":
		maxSumForNorm := float64(numberOfLines-1) + maxIdxGlobal
		if maxSumForNorm <= 0 {
			maxSumForNorm = 1
		}
		currentSum := float64(numberOfLines-1-lineIdx) + (maxIdxGlobal - float64(charIdx))
		delay = (params.Duration * defaultCascadeSpreadFactor) * (currentSum / maxSumForNorm)
	case "center-out":
		centerGlobalGrid := maxIdxGlobal / 2.0
		if maxIdxGlobal > 0 {
			distanceFromGlobalCenter := math.Abs(float64(charIdx) - centerGlobalGrid)
			delay = (params.Duration * centerEdgeStaggerFactor) * (distanceFromGlobalCenter / centerGlobalGrid)
		}
	case "edges-in":
		centerGlobalGrid := maxIdxGlobal / 2.0
		if maxIdxGlobal > 0 {
			distanceFromGlobalCenter := math.Abs(float64(charIdx) - centerGlobalGrid)
			delay = (params.Duration * centerEdgeStaggerFactor) * ((centerGlobalGrid - distanceFromGlobalCenter) / centerGlobalGrid)
		}
	case "ttb-linesync":
		delay = (params.Duration * defaultCascadeSpreadFactor) * (float64(lineIdx) / normNumLines)
	case "btt-linesync":
		delay = (params.Duration * defaultCascadeSpreadFactor) * (float64(numberOfLines-1-lineIdx) / normNumLines)
	case "full-sync":
		delay = 0.0
	default:
		if numberOfCharsInLine > 1 {
			delay = -0.5 * params.Duration * (float64(charIdx) / float64(numberOfCharsInLine-1))
		}
	}
	return delay
}
