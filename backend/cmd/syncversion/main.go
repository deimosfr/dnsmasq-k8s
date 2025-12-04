package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
	versionFile = "../../src/api/version.go"
	chartFile   = "../../../chart/Chart.yaml"
	readmeFile  = "../../../README.md"
	valuesFile  = "../../../chart/values.yaml"
)

func main() {
	// 1. Get the version from version.go
	version, err := getVersion()
	if err != nil {
		fmt.Printf("Error getting version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found version: %s\n", version)

	// 2. Update Chart.yaml
	if err := updateChartYaml(version); err != nil {
		fmt.Printf("Error updating Chart.yaml: %v\n", err)
		os.Exit(1)
	}

	// 3. Update README.md
	if err := updateReadme(version); err != nil {
		fmt.Printf("Error updating README.md: %v\n", err)
		os.Exit(1)
	}

	// 4. Update values.yaml
	if err := updateValuesYaml(version); err != nil {
		fmt.Printf("Error updating values.yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated all files.")
}

func getVersion() (string, error) {
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`const Version = "([^"]+)"`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "", fmt.Errorf("version constant not found in %s", versionFile)
	}
	return matches[1], nil
}

func updateChartYaml(version string) error {
	content, err := os.ReadFile(chartFile)
	if err != nil {
		return err
	}

	strContent := string(content)

	// Update version: ...
	reVersion := regexp.MustCompile(`version: .*`)
	strContent = reVersion.ReplaceAllString(strContent, fmt.Sprintf("version: %s", version))

	// Update appVersion: "..."
	reAppVersion := regexp.MustCompile(`appVersion: ".*"`)
	strContent = reAppVersion.ReplaceAllString(strContent, fmt.Sprintf("appVersion: \"%s\"", version))

	return os.WriteFile(chartFile, []byte(strContent), 0644)
}

func updateReadme(version string) error {
	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return err
	}

	// Update tag: "..."
	reTag := regexp.MustCompile(`tag: ".*"`)
	strContent := reTag.ReplaceAllString(string(content), fmt.Sprintf("tag: \"%s\"", version))

	// Update helm install --version ...
	reHelm := regexp.MustCompile(`--version \d+\.\d+\.\d+`)
	strContent = reHelm.ReplaceAllString(strContent, fmt.Sprintf("--version %s", version))

	return os.WriteFile(readmeFile, []byte(strContent), 0644)
}

func updateValuesYaml(version string) error {
	content, err := os.ReadFile(valuesFile)
	if err != nil {
		return err
	}

	// Update tag: "..."
	re := regexp.MustCompile(`tag: ".*"`)
	strContent := re.ReplaceAllString(string(content), fmt.Sprintf("tag: \"%s\"", version))

	return os.WriteFile(valuesFile, []byte(strContent), 0644)
}
