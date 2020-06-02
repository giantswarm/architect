package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/giantswarm/microerror"
)

func validateSingleOccurence(data []byte, regexps ...*regexp.Regexp) error {
	matches := 0
	for _, re := range regexps {
		matches += len(re.FindAllIndex(data, -1))
	}

	var combined *regexp.Regexp
	if len(regexps) == 1 {
		combined = regexps[0]
	} else {
		var patterns []string

		for _, re := range regexps {
			patterns = append(patterns, re.String())
		}

		pattern := fmt.Sprintf("(?:%s)", strings.Join(patterns, ") | (?:"))

		combined = regexp.MustCompile(pattern)
	}

	if matches == 0 {
		return microerror.Maskf(executionFailedError, "no match for pattern %#q match found in data:\n---\n%s\n---", combined, data)
	}
	if matches > 1 {
		return microerror.Maskf(executionFailedError, "%d pattern %#q matches found, expected 1 in data:\n---\n%s\n---", matches, combined, data)
	}

	return nil
}
