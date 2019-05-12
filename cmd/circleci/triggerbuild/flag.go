package triggerbuild

func init() {
	Cmd.Flags().String("branch", "master", "project's branch name")
	Cmd.Flags().String("job", "", "CircleCI job name. This can only be used for project which do not have a defined workflow.")
	Cmd.Flags().String("org", "", "CircleCI organization or username")
	Cmd.Flags().String("project", "", "CircleCI project name")
	Cmd.Flags().String("token", "", "CircleCI API token")
}
