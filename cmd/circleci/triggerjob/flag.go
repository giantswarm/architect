package triggerjob

func init() {
	Cmd.Flags().String("token", "", "CircleCI API token")
	Cmd.Flags().String("org", "", "CircleCI organization")
	Cmd.Flags().String("repo", "", "CircleCI repository")
	Cmd.Flags().String("branch", "master", "name of the repository branch")
	Cmd.Flags().String("job", "", "CircleCI job")
}
