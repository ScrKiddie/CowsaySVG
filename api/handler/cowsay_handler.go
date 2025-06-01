package handler

import (
	"cowsaysvg/api/_utility"
	"fmt"
	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	text := query.Get("text")

	if text == "" {
		apiURL := os.Getenv("API_URL")
		if apiURL != "" {
			body, err := _utility.FetchPlainText(apiURL)
			if err != nil {
				http.Error(w, "failed to fetch from API: "+err.Error(), http.StatusInternalServerError)
				return
			}
			text = body
		}
	}

	cow := query.Get("cow")
	if cow == "" {
		cow = "default"
	}

	var colors []string
	if colorsParam := query.Get("colors"); colorsParam != "" {
		colors = strings.Split(colorsParam, ",")
	}

	duration := 0.0
	if durStr := query.Get("duration"); durStr != "" {
		if dur, err := strconv.ParseFloat(durStr, 64); err == nil && dur >= 0 {
			duration = dur
		}
	}

	numColors := len(colors)

	timing := query.Get("timing")
	if timing == "" {
		if numColors > 0 {
			timing = fmt.Sprintf("steps(%d, end)", numColors)
		} else {
			timing = "steps(1, end)"
		}
	}
	ballonWidth := 40
	if bw := query.Get("ballonWidth"); bw != "" {
		if parsed, err := strconv.Atoi(bw); err == nil && parsed > 0 {
			ballonWidth = parsed
		}
	}

	charWidth := 10
	if cw := query.Get("charWidth"); cw != "" {
		if parsed, err := strconv.Atoi(cw); err == nil && parsed > 0 {
			charWidth = parsed
		}
	}

	lineHeight := 20
	if lh := query.Get("lineHeight"); lh != "" {
		if parsed, err := strconv.Atoi(lh); err == nil && parsed > 0 {
			lineHeight = parsed
		}
	}

	speech, err := cowsay.Say(
		text,
		cowsay.Type(cow),
		cowsay.BallonWidth(uint(ballonWidth)),
	)
	if err != nil {
		slog.Warn("invalid cow type, fallback to default", "cow", cow, "error", err)
		speech, _ = cowsay.Say(text, cowsay.Type("default"), cowsay.BallonWidth(uint(ballonWidth)))
	}

	config := _utility.Config{
		CharWidth:  charWidth,
		LineHeight: lineHeight,
		FontSize:   int(float64(lineHeight) * 0.85),
	}

	animationParams := _utility.AnimationParams{
		Colors:         colors,
		TimingFunction: timing,
		Duration:       duration,
	}

	renderer := _utility.NewRenderer(w, config)
	w.Header().Set("Content-Type", "image/svg+xml")
	renderer.Render(speech, animationParams, duration == 0)
}
