package imagegen

import (
	"fmt"
	"strings"
)

func mockSVG(prompt, size string) string {
	w, hgt := 1024, 1024
	fmt.Sscanf(size, "%dx%d", &w, &hgt)

	var hash uint32 = 2166136261
	for _, r := range prompt {
		hash = (hash ^ uint32(r)) * 16777619
	}
	h1 := hash % 360
	h2 := (hash / 7) % 360

	esc := func(s string) string {
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		return s
	}
	lines := wrap(prompt, 38, 6)
	var tspans strings.Builder
	for i, ln := range lines {
		dy := "1.3em"
		if i == 0 {
			dy = "0"
		}
		tspans.WriteString(fmt.Sprintf(`<tspan x="50%%" dy="%s">%s</tspan>`, dy, esc(ln)))
	}

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" font-family="Georgia, serif">
<defs>
<linearGradient id="g" x1="0" y1="0" x2="1" y2="1">
<stop offset="0" stop-color="hsl(%d,70%%,18%%)"/>
<stop offset="0.55" stop-color="hsl(%d,65%%,12%%)"/>
<stop offset="1" stop-color="#07060d"/>
</linearGradient>
<radialGradient id="r" cx="30%%" cy="25%%" r="70%%">
<stop offset="0" stop-color="hsl(%d,85%%,55%%)" stop-opacity="0.5"/>
<stop offset="1" stop-color="hsl(%d,85%%,55%%)" stop-opacity="0"/>
</radialGradient>
</defs>
<rect width="%d" height="%d" fill="url(#g)"/>
<rect width="%d" height="%d" fill="url(#r)"/>
<text x="50%%" y="38%%" text-anchor="middle" fill="#ece9ff" font-size="34" opacity="0.92">%s</text>
<text x="50%%" y="91%%" text-anchor="middle" fill="hsl(%d,80%%,65%%)" font-size="20" letter-spacing="6" font-family="monospace">MOCK · KEIN BILD-KI-KEY</text>
</svg>`, w, hgt, w, hgt, h1, h2, h1, h2, w, hgt, w, hgt, tspans.String(), h2)
}

func wrap(s string, width, maxLines int) []string {
	words := strings.Fields(s)
	lines := []string{}
	cur := ""
	for _, w := range words {
		if len(cur)+len(w)+1 > width {
			lines = append(lines, cur)
			cur = w
			if len(lines) >= maxLines-1 {
				break
			}
		} else if cur == "" {
			cur = w
		} else {
			cur += " " + w
		}
	}
	if cur != "" && len(lines) < maxLines {
		lines = append(lines, cur)
	}
	if len(lines) == 0 {
		lines = []string{s}
	}
	return lines
}
