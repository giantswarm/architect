package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "3.4.1"
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
