package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type configLoadStates int

const (
	SEARCHING configLoadStates = iota
	METRIC_DEFINITION
)

var (
	identifierRE   = `[a-zA-Z_][a-zA-Z0-9_]+`
	statsdMetricRE = `[a-zA-Z_](-?[a-zA-Z0-9_])+`

	metricLineRE = regexp.MustCompile(`^(\*\.|` + statsdMetricRE + `\.)+(\*|` + statsdMetricRE + `)$`)
	labelLineRE  = regexp.MustCompile(`^(` + identifierRE + `)\s*=\s*"(.*)"$`)
	metricNameRE = regexp.MustCompile(`^` + identifierRE + `$`)
)

type metricMapper struct {
	Mappings []metricMapping `yaml:"mappings"`
}

type metricMapping struct {
	Match  string `yaml:"match"`
	Name   string `yaml:"name"`
	regex  *regexp.Regexp
	Labels map[string]string `yaml:"labels"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must have a mapping file as only argument")
		os.Exit(-1)
	}

	mappings, err := initFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}

	m := metricMapper{
		Mappings: mappings,
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		fmt.Printf("failed to marshal config: %s", err.Error())
	}

	fmt.Println(string(data))
}

func initFromFile(fileName string) ([]metricMapping, error) {
	mappingStr, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return initFromString(string(mappingStr))
}

func initFromString(fileContents string) ([]metricMapping, error) {
	lines := strings.Split(fileContents, "\n")
	numLines := len(lines)
	state := SEARCHING

	parsedMappings := []metricMapping{}
	currentMapping := metricMapping{Labels: map[string]string{}}
	for i, line := range lines {
		line := strings.TrimSpace(line)

		switch state {
		case SEARCHING:
			if line == "" {
				continue
			}
			if !metricLineRE.MatchString(line) {
				return nil, fmt.Errorf("Line %d: expected metric match line, got: %s", i, line)
			}

			// Translate the glob-style metric match line into a proper regex that we
			// can use to match metrics later on.
			metricRe := strings.Replace(line, ".", "\\.", -1)
			metricRe = strings.Replace(metricRe, "*", "([^.]+)", -1)
			currentMapping.regex = regexp.MustCompile("^" + metricRe + "$")

			currentMapping.Match = line
			state = METRIC_DEFINITION

		case METRIC_DEFINITION:
			if (i == numLines-1) && (line != "") {
				return nil, fmt.Errorf("Line %d: missing terminating newline", i)
			}
			if line == "" {
				if len(currentMapping.Labels) == 0 {
					return nil, fmt.Errorf("Line %d: metric mapping didn't set any labels", i)
				}
				name, ok := currentMapping.Labels["name"]
				if !ok || name == "" {
					return nil, fmt.Errorf("Line %d: metric mapping didn't set a metric name", i)
				}
				currentMapping.Name = name
				delete(currentMapping.Labels, "name")
				parsedMappings = append(parsedMappings, currentMapping)
				state = SEARCHING
				currentMapping = metricMapping{Labels: map[string]string{}}
				continue
			}
			if err := updateMapping(line, i, &currentMapping); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("illegal state")
		}
	}

	return parsedMappings, nil
}

func updateMapping(line string, i int, mapping *metricMapping) error {
	matches := labelLineRE.FindStringSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("Line %d: expected label mapping line, got: %s", i, line)
	}
	label, value := matches[1], matches[2]
	if label == "name" && !metricNameRE.MatchString(value) {
		return fmt.Errorf("Line %d: metric name '%s' doesn't match regex '%s'", i, value, metricNameRE)
	}

	(*mapping).Labels[label] = value
	return nil
}
