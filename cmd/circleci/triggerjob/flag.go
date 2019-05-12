package triggerjob

func init() {
	Cmd.Flags().String("branch", "master", "name of the repository branch")
	Cmd.Flags().String("job", "", "CircleCI job name")
	Cmd.Flags().String("org", "", "CircleCI organization or username")
	Cmd.Flags().String("repo", "", "CircleCI repository or project name")
	Cmd.Flags().String("token", "", "CircleCI API token")
}
