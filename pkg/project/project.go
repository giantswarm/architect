package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "4.0.2-dev"
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
