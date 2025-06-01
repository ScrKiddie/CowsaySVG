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
				http.Error(w, "couldn't fetch from API: "+err.Error(), http.StatusInternalServerError)
				return
			}
			text = body
		}
	}

	userCow := query.Get("cow")
	cowToUse := "default"

	if userCow != "" {
		if _utility.IsValidCowName(userCow) {
			cowToUse = userCow
		} else {
			slog.Info("invalid cow requested, using default", "requested", userCow)
		}
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

	cascadeDirection := query.Get("cascadeDirection")
	if cascadeDirection == "" {
		cascadeDirection = "rtl"
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

	opts := []cowsay.Option{
		cowsay.BallonWidth(uint(ballonWidth)),
	}

	if eyesStr := query.Get("eyes"); eyesStr != "" {
		opts = append(opts, cowsay.Eyes(eyesStr))
	}

	if tongueStr := query.Get("tongue"); tongueStr != "" {
		opts = append(opts, cowsay.Tongue(tongueStr))
	}

	if thinkParam := query.Get("think"); thinkParam == "true" || thinkParam == "1" {
		opts = append(opts, cowsay.Thinking())
	}

	if thoughtsCharStr := query.Get("thoughtsChar"); thoughtsCharStr != "" {
		runes := []rune(thoughtsCharStr)
		if len(runes) > 0 {
			opts = append(opts, cowsay.Thoughts(runes[0]))
		}
	}

	if noWrapParam := query.Get("noWrap"); noWrapParam == "true" || noWrapParam == "1" {
		opts = append(opts, cowsay.DisableWordWrap())
	}

	tryOpts := make([]cowsay.Option, len(opts))
	copy(tryOpts, opts)
	cowName := cowToUse
	randomMode := false

	if randomParam := query.Get("randomCow"); randomParam == "true" || randomParam == "1" {
		randomMode = true
		cowName = "random"
		tryOpts = append(tryOpts, cowsay.Random())
	} else {
		tryOpts = append(tryOpts, cowsay.Type(cowToUse))
	}

	speech, err := cowsay.Say(text, tryOpts...)
	if err != nil {
		slog.Warn("couldn't make cow talk, trying default",
			"attempted", cowName,
			"random", randomMode,
			"error", err)

		fallbackOpts := make([]cowsay.Option, len(opts))
		copy(fallbackOpts, opts)
		fallbackOpts = append(fallbackOpts, cowsay.Type("default"))

		speech, err = cowsay.Say(text, fallbackOpts...)
		if err != nil {
			slog.Error("still can't make cow talk", "error", err)
			http.Error(w, "something went wrong with the cow", http.StatusInternalServerError)
			return
		}
	}

	cfg := _utility.Config{
		CharWidth:  charWidth,
		LineHeight: lineHeight,
		FontSize:   int(float64(lineHeight) * 0.85),
	}

	anim := _utility.AnimationParams{
		Colors:           colors,
		TimingFunction:   timing,
		Duration:         duration,
		CascadeDirection: cascadeDirection,
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	_utility.NewRenderer(w, cfg).Render(speech, anim, duration == 0)
}
