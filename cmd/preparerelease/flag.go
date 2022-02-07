package preparerelease

func init() {
	Cmd.Flags().Bool("update-changelog", true, "if true, update CHANGELOG.md")
	Cmd.Flags().String("version", "", "version to be released")
}
