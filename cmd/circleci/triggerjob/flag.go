package triggerjob

func init() {
	Cmd.PersistentFlags().String("branch", "master", "project's branch name")
	Cmd.PersistentFlags().String("job", "", "CircleCI job name. This can only be used for project which do not have a defined workflow.")
	Cmd.PersistentFlags().String("org", "", "CircleCI organization or username")
	Cmd.PersistentFlags().String("project", "", "CircleCI project name")
	Cmd.PersistentFlags().String("token", "", "CircleCI API token")
}
