package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "3.0.6-dev"
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
