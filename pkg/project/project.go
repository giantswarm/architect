package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "6.17.0"
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
