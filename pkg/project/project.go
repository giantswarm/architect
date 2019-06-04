package project

var (
	buildTimestamp string
	gitSHA         string
)

func BuildTimestamp() string {
	return buildTimestamp
}

func GitSHA() string {
	return gitSHA
}
