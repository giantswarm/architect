package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "5.3.1-dev"
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
