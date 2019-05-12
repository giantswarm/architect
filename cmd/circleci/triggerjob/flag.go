package triggerjob

func init() {
	Cmd.LocalFlags().String("branch", "master", "project's branch name")
	Cmd.LocalFlags().String("job", "", "CircleCI job name. This can only be used for project which do not have a defined workflow.")
	Cmd.LocalFlags().String("org", "", "CircleCI organization or username")
	Cmd.LocalFlags().String("project", "", "CircleCI project name")
	Cmd.LocalFlags().String("token", "", "CircleCI API token")
}
