package s3d

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	p "github.com/majori/goco/pkg/parser"
)

func parseSettingsRow(raw string) interface{} {
	// Value is probably g-code
	if strings.Contains(raw, ";") {
		return raw
	}

	// Check if value is list
	if strings.Contains(raw, ",") {
		items := strings.Split(raw, ",")
		listValues := make([]interface{}, 0)
		for _, value := range items {
			listValues = append(listValues, parseSettingsRow(value))
		}
		return listValues
	}

	// Check if value is type of float
	if strings.Contains(raw, ".") {
		value, err := strconv.ParseFloat(raw, 32)
		if err == nil {
			return value
		}
	}

	// Check if value is type of int
	value, err := strconv.Atoi(raw)
	if err == nil {
		return value
	}

	// Fallback to string value
	return raw
}

func parseSettings(source *string) *p.Settings {
	static := make(p.Settings)
	scanner := bufio.NewScanner(strings.NewReader(*source))
	start := false

	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())

		if row == "" && !start {
			continue
		}

		// Find start of the settings
		if strings.Contains(row, "Settings Summary") {
			start = true
			continue
		}

		// Stop iterating when first gcode encountered
		if !strings.HasPrefix(row, ";") {
			break
		}

		// Remove ";" from beginning
		row = strings.TrimSpace(row[1:])

		if start {
			setting := strings.SplitN(row, ",", 2)
			if setting[1] == "" {
				continue
			}
			static[setting[0]] = parseSettingsRow(setting[1])
		}
	}

	return &static
}

func parseState(source *string) *p.States {
	states := p.States{}
	layer := 0
	var z float32 = 0.0

	scanner := bufio.NewScanner(strings.NewReader(*source))
	for scanner.Scan() {
		row := scanner.Text()
		state := make(map[string]interface{})

		findSubmatch := func(exp string) ([]string, bool) {
			regex := regexp.MustCompile(exp)
			submatches := regex.FindStringSubmatch(row)
			if len(submatches) > 0 {
				return submatches[1:], true
			} else {
				return nil, false
			}
		}

		if match, ok := findSubmatch("^; layer (\\d+)"); ok {
			layer, _ = strconv.Atoi(match[0])
		}

		if match, ok := findSubmatch("^;.*Z = (\\d*\\.?\\d*)"); ok {
			parsed, _ := strconv.ParseFloat(match[0], 32)
			z = float32(parsed)
		}

		state["layer"] = layer
		state["z"] = z
		states = append(states, state)
	}

	return &states
}

func Parse(source *string) (*p.Settings, *p.States) {
	return parseSettings(source), parseState(source)
}
