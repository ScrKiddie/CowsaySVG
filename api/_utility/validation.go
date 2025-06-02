package _utility

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func IsValidAndSafeCSSTimingFunction(userInput string) (string, bool) {
	trimmedInput := strings.TrimSpace(userInput)

	safeKeywords := map[string]bool{
		"linear":      true,
		"ease":        true,
		"ease-in":     true,
		"ease-out":    true,
		"ease-in-out": true,
		"step-start":  true,
		"step-end":    true,
	}
	if safeKeywords[trimmedInput] {
		return trimmedInput, true
	}

	if strings.HasPrefix(trimmedInput, "steps(") && strings.HasSuffix(trimmedInput, ")") {
		argsStr := trimmedInput[len("steps(") : len(trimmedInput)-1]
		parts := strings.SplitN(argsStr, ",", 2)

		nStr := strings.TrimSpace(parts[0])
		n, err := strconv.Atoi(nStr)
		if err != nil || n <= 0 {
			slog.Warn("invalid number of steps (N) in timing function", "input_n", nStr, "full_input", userInput)
			return "", false
		}

		var stepPosition string
		if len(parts) == 2 {
			posStr := strings.TrimSpace(parts[1])
			validStepPositions := map[string]bool{
				"jump-start": true, "jump-end": true, "jump-none": true, "jump-both": true,
				"start": true, "end": true,
			}
			if !validStepPositions[posStr] {
				slog.Warn("invalid step position in timing function", "position", posStr, "full_input", userInput)
				return "", false
			}
			stepPosition = posStr
			return fmt.Sprintf("steps(%d, %s)", n, stepPosition), true
		} else if len(parts) == 1 {
			slog.Warn("steps() function in timing function requires a step position or is ambiguous", "full_input", userInput)
			return "", false
		}
	}

	slog.Warn("unsupported or potentially unsafe timing function value", "input", userInput)
	return "", false
}
