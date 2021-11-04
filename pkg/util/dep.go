package util

import "regexp"

func ParseDependency(dependencyName string) string {
	re := regexp.MustCompile(`(?m)(([?^/]).*)`)
	if !re.Match([]byte(dependencyName)) {
		return "github.com/HashLoad/" + dependencyName
	}
	re = regexp.MustCompile(`(?m)([?^/].*)(([?^/]).*)`)
	if !re.Match([]byte(dependencyName)) {
		return "github.com/" + dependencyName
	}
	return dependencyName
}
