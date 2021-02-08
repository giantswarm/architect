package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "3.2.3-dev"
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
