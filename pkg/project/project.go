package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "2.1.4-dev"
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
