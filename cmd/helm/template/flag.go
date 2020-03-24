package template

func init() {
	Cmd.Flags().String("dir", "", "helm chart directory")
	Cmd.Flags().Bool("validate", false, "enables chart validation")
	Cmd.Flags().Bool("tag-build", false, "should be set when validating a tagged build")
}
