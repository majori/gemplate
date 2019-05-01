package simplify3d

import (
	"fmt"
	"bufio"
	"strconv"
	"strings"
)

func parseRow(raw string) interface{} {
	// Value is probably g-code
	if strings.Contains(raw, ";") {
		return raw
	}

	// Check if value is list
	if strings.Contains(raw, ",") {
		items := strings.Split(raw, ",")
		listValues := make([]interface{}, 0)
		for _, value := range items {
			listValues = append(listValues, parseRow(value))
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

func ParseSettings(scanner *bufio.Scanner) map[string]interface{} {
	settings := make(map[string]interface{})
	// Find start of the settings
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Settings Summary") {
			break
		}
	}


	for scanner.Scan() {
		row := strings.Trim(scanner.Text(), " ")
		if !strings.HasPrefix(row, ";") {
			break
		}

		row = strings.Trim(row[1:], " ")
		setting := strings.SplitN(row, ",", 2)
		settings[setting[0]] = parseRow(setting[1])
	}

	return settings
}
