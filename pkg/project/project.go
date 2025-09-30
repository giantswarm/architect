package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "7.2.2-dev"
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
