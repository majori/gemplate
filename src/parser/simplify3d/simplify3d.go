package simplify3d

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/majori/gemplate/src/parser"
)

func parseStaticRow(raw string) interface{} {
	// Value is probably g-code
	if strings.Contains(raw, ";") {
		return raw
	}

	// Check if value is list
	if strings.Contains(raw, ",") {
		items := strings.Split(raw, ",")
		listValues := make([]interface{}, 0)
		for _, value := range items {
			listValues = append(listValues, parseStaticRow(value))
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

func parseStaticSettings(source *bufio.Scanner) parser.Static {
	static := make(parser.Static)

	start := false
	for source.Scan() {
		row := strings.TrimSpace(source.Text())

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

		row = strings.TrimSpace(row[1:])

		if start {
			setting := strings.SplitN(row, ",", 2)
			static[setting[0]] = parseStaticRow(setting[1])
		}
	}

	return static
}

func ParseSettings(source *bufio.Scanner) parser.Settings {
	return parser.Settings{
		Static: parseStaticSettings(source),
	}
}
