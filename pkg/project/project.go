package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "7.0.3-dev"
)

func BuildTimestamp() string {
	return buildTimestamp
}

func GitSHA() string {
	return gitSHA
}

func Version() string {
	return version
}
