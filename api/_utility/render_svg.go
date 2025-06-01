package _utility

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/ajstarks/svgo"
)

type Config struct {
	CharWidth  int
	LineHeight int
	FontSize   int
}

type AnimationParams struct {
	Colors         []string
	TimingFunction string
	Duration       float64
}

type Renderer struct {
	canvas *svg.SVG
	config Config
}

func NewRenderer(w http.ResponseWriter, config Config) *Renderer {
	return &Renderer{
		canvas: svg.New(w),
		config: config,
	}
}

func escapePercent(colorStr string) string {
	return strings.ReplaceAll(colorStr, "%", "%%")
}

func (r *Renderer) Render(text string, params AnimationParams, isStatic bool) {
	lines := strings.Split(text, "\n")
	width, height := r.calculateDimensions(lines)

	r.Start(width, height)
	defer r.End()

	if isStatic {
		r.renderStaticTextWithGradient(lines, params.Colors)
		return
	}

	numColors := len(params.Colors)

	if numColors == 0 {
		r.renderStaticTextWithGradient(lines, []string{})
		return
	}

	escapedInitialFill := escapePercent(params.Colors[0])

	r.canvas.Def()
	var keyframeRulesJoined string

	if numColors == 1 {
		escapedColor := escapePercent(params.Colors[0])
		keyframeRulesJoined = fmt.Sprintf("0%% { fill: %s; }\n100%% { fill: %s; }", escapedColor, escapedColor)
	} else {
		keyframePercentages := r.generateKeyframePercentages(numColors)
		keyframeRules := make([]string, numColors)

		for i := 0; i < numColors; i++ {
			percent := keyframePercentages[i]
			escapedColor := escapePercent(params.Colors[i])
			keyframeRules[i] = fmt.Sprintf("%d%% { fill: %s; }", percent, escapedColor)
		}
		keyframeRulesJoined = strings.Join(keyframeRules, "\n")
	}

	r.canvas.Style("text/css", fmt.Sprintf(`
        @keyframes custom_anim {
            %s
        }
        .anim-char { 
            animation: custom_anim %.2fs %s infinite;
        }`,
		keyframeRulesJoined,
		params.Duration,
		params.TimingFunction))
	r.canvas.DefEnd()

	for lineIdx, line := range lines {
		y := r.config.LineHeight * (lineIdx + 1)
		for charIdx, char := range line {
			x := r.config.CharWidth * charIdx
			var delay float64
			if len(line) > 1 {
				delay = -0.5 * params.Duration * (float64(charIdx) / float64(len(line)-1))
			}

			r.canvas.Text(x, y, string(char), fmt.Sprintf(
				`class="anim-char" fill="%s" font-family="monospace" font-size="%dpx" 
                 style="animation-delay:%.3fs"`,
				escapedInitialFill,
				r.config.FontSize,
				delay,
			))
		}
	}
}

func (r *Renderer) generateKeyframePercentages(colorCount int) []int {
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
	percentages[colorCount-1] = 100

	return percentages
}

func (r *Renderer) renderStaticTextWithGradient(lines []string, colors []string) {
	if len(colors) == 0 {
		for lineIdx, line := range lines {
			y := r.config.LineHeight * (lineIdx + 1)
			for charIdx, char := range line {
				x := r.config.CharWidth * charIdx
				r.canvas.Text(x, y, string(char),
					fmt.Sprintf(`fill="black" font-family="monospace" font-size="%dpx"`,
						r.config.FontSize))
			}
		}
		return
	}

	for lineIdx, line := range lines {
		y := r.config.LineHeight * (lineIdx + 1)
		if len(line) == 0 {
			continue
		}

		for charIdx, char := range line {
			x := r.config.CharWidth * charIdx
			var fillColor string
			if len(colors) == 1 {
				fillColor = colors[0]
			} else {
				colorIndex := (charIdx * len(colors)) / len(line)
				if colorIndex >= len(colors) {
					colorIndex = len(colors) - 1
				}
				fillColor = colors[colorIndex]
			}

			escapedFillColor := escapePercent(fillColor)
			r.canvas.Text(x, y, string(char),
				fmt.Sprintf(`fill="%s" font-family="monospace" font-size="%dpx"`,
					escapedFillColor, r.config.FontSize))
		}
	}
}

func (r *Renderer) calculateDimensions(lines []string) (width, height int) {
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	calculatedHeight := len(lines) * r.config.LineHeight
	if len(lines) == 0 {
		return maxLen * r.config.CharWidth, 0
	}

	paddingBottom := int(float64(r.config.FontSize) * 0.35)
	if paddingBottom == 0 && r.config.FontSize > 0 {
		paddingBottom = 2
	}

	return maxLen * r.config.CharWidth, calculatedHeight + paddingBottom
}

func (r *Renderer) Start(width, height int) {
	r.canvas.Start(width, height)
	r.canvas.Rect(0, 0, width, height, `fill="none"`)
}

func (r *Renderer) End() {
	r.canvas.End()
}
