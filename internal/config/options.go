package config

import (
	"errors"
	"strconv"
	"strings"
)

func ParseArgs(args []string, base Config) (Config, error) {
	result := base
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--db":
			index++
			if index >= len(args) {
				return result, errors.New("--db needs a value")
			}
			result.DatabasePath = args[index]
		case "--address":
			index++
			if index >= len(args) {
				return result, errors.New("--address needs a value")
			}
			result.HTTPAddress = args[index]
		case "--reviewer":
			index++
			if index >= len(args) {
				return result, errors.New("--reviewer needs a value")
			}
			result.Reviewer = args[index]
		case "--minimum":
			index++
			if index >= len(args) {
				return result, errors.New("--minimum needs a value")
			}
			minimum, err := strconv.Atoi(args[index])
			if err != nil {
				return result, err
			}
			result.PolicyMinimum = minimum
		default:
			if strings.HasPrefix(args[index], "-") {
				return result, errors.New("unknown option " + args[index])
			}
		}
	}
	return result, result.Validate()
}

func Usage() string {
	return "store-ledger [import|query|health] --db path --address host:port --reviewer name --minimum score"
}

func Merge(primary, fallback Config) Config {
	result := primary
	if result.DatabasePath == "" {
		result.DatabasePath = fallback.DatabasePath
	}
	if result.HTTPAddress == "" {
		result.HTTPAddress = fallback.HTTPAddress
	}
	if result.Reviewer == "" {
		result.Reviewer = fallback.Reviewer
	}
	if result.PolicyMinimum == 0 {
		result.PolicyMinimum = fallback.PolicyMinimum
	}
	return result
}
