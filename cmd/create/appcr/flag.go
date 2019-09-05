package appcr

func init() {
	Cmd.Flags().String("catalog", "", "app catalog name")
	Cmd.Flags().String("name", "", "cr name")
	Cmd.Flags().String("app-name", "", "app name")
	Cmd.Flags().StringP("output", "o", "yaml", "output format. allowed: yaml,json")
	Cmd.Flags().String("app-version", "", "app version")
}
