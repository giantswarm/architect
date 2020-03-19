package project

var (
	buildTimestamp string
	gitSHA         string
	version        = "n/a"
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
