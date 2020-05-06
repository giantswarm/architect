package prepare

func init() {
	Cmd.Flags().String("version", "", "version to be released")
}
