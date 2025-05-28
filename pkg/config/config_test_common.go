package config

import "os"

func isRunningOnGithubActions() bool {
	return len(os.Getenv("GITHUB_ACTIONS")) > 0 || os.Getenv("CI") == "true"
}
