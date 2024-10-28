package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func parseLine(line string, section string, configMap map[string]string) (string, error) {
	trimmedLine := strings.TrimSpace(line)

	if strings.HasPrefix(trimmedLine, "#") || strings.HasPrefix(trimmedLine, ";") {
		return section, nil
	}

	if strings.HasPrefix(trimmedLine, "[") {
		section = strings.TrimPrefix(trimmedLine, "[")
		section = strings.TrimSuffix(section, "]")
		section = strings.ReplaceAll(section, " ", ".")
		section = strings.ReplaceAll(section, "\"", "")
		return section, nil
	}

	parts := strings.SplitN(trimmedLine, "=", 2)
	if len(parts) != 2 {
		return section, fmt.Errorf("invalid input format")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Remove quotes if present
	re, err := regexp.Compile(`^"(.*)"$`)
	if err != nil {
		return section, err
	}

	if match := re.FindStringSubmatch(value); match != nil {
		value = match[1]
	}

	configMap[fmt.Sprintf("%s.%s", section, key)] = value

	return section, nil
}

func ParseConfigFile(reader io.Reader) (map[string]string, error) {
	var section string
	configMap := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	var err error
	for scanner.Scan() {
		section, err = parseLine(scanner.Text(), section, configMap)
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return configMap, nil
}
