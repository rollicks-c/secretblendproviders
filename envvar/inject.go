package envvar

import (
	"fmt"
	"os"
	"regexp"
)

func injectVars(exp string) (string, error) {

	// regex to match $VAR or ${VAR}
	re := regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)

	// process all matches
	missingValues := []string{}
	result := re.ReplaceAllStringFunc(exp, func(match string) string {

		// get var name
		name := re.FindStringSubmatch(match)[1]

		// replace
		if value, exists := os.LookupEnv(name); exists {
			return value
		}

		// remember missing values
		missingValues = append(missingValues, name)
		return match
	})

	// report missing values
	if len(missingValues) > 0 {
		return "", fmt.Errorf("failed to load secret - envvars not provided: %s", missingValues)
	}

	return result, nil

}
