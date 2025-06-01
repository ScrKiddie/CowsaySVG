package _utility

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"unicode"

	"github.com/ajstarks/svgo"
)

type Config struct {
	CharWidth  int
	LineHeight int
	FontSize   int
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
	rawLines := strings.Split(text, "\n")

	lines := make([]string, len(rawLines))
	for i, line := range rawLines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}

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

	keyframeRulesJoined := GenerateColorKeyframeRules(params.Colors)

	r.canvas.Def()
	if keyframeRulesJoined != "" {
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
	}
	r.canvas.DefEnd()

	escapedInitialFill := escapePercent(params.Colors[0])
	numberOfLines := len(lines)

	maxLineLengthGlobal := 0
	for _, currentLine := range lines {
		if len(currentLine) > maxLineLengthGlobal {
			maxLineLengthGlobal = len(currentLine)
		}
	}

	for lineIdx, line := range lines {
		y := r.config.LineHeight * (lineIdx + 1)
		numberOfCharsInLine := len(line)

		for charIdx, char := range line {
			x := r.config.CharWidth * charIdx

			delay := CalculateAnimationDelay(params, lineIdx, charIdx, numberOfLines, numberOfCharsInLine, maxLineLengthGlobal)

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
				colorIndex := 0
				if len(line) > 1 {
					colorIndex = int(math.Round(float64(charIdx) * float64(len(colors)-1) / float64(len(line)-1)))
				}
				if colorIndex >= len(colors) {
					colorIndex = len(colors) - 1
				}
				if colorIndex < 0 {
					colorIndex = 0
				}
				fillColor = colors[colorIndex]
			}
			escapedFillColor := escapePercent(fillColor)
			r.canvas.Text(x, y, string(char), fmt.Sprintf(`fill="%s" font-family="monospace" font-size="%dpx"`, escapedFillColor, r.config.FontSize))
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
		return 0, 0
	}
	if maxLen == 0 {
		return 0, calculatedHeight
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
func (r *Renderer) End() { r.canvas.End() }
