package template

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/helmtemplate"
)

func runTemplateError(cmd *cobra.Command, args []string) (err error) {
	var (
		chartDir = cmd.Flag("dir").Value.String()
		branch   = cmd.Flag("branch").Value.String()
		sha      = cmd.Flag("sha").Value.String()
		tag      = cmd.Flag("tag").Value.String()
		version  = cmd.Flag("version").Value.String()
		validate bool
		tagBuild bool
	)
	{
		var err error
		validate, err = strconv.ParseBool(cmd.Flag("validate").Value.String())
		if err != nil {
			return microerror.Mask(err)
		}
		tagBuild, err = strconv.ParseBool(cmd.Flag("tag-build").Value.String())
		if err != nil {
			return microerror.Mask(err)
		}
	}

	fs := afero.NewOsFs()
	ctx := context.Background()

	var appVersion string
	skipAppVersionCheck := false
	{
		dir, err := gitrepo.TopLevel(ctx, ".")
		if err != nil {
			return microerror.Mask(err)
		}

		appVersion, err = getProjectVersion(dir)
		if err != nil {
			return microerror.Mask(err)
		}

		// for repositories without pkg/project/project.go
		if appVersion == "" {
			appVersion = version
			skipAppVersionCheck = true
		}
	}

	log.Printf("templating helm chart\ndir: %s\nsha: %s\ntag: %s\napp-version: %s\nversion: %s\n", chartDir, sha, tag, appVersion, version)

	var s *helmtemplate.TemplateHelmChartTask
	{
		c := helmtemplate.Config{
			Fs:                  fs,
			ChartDir:            chartDir,
			Branch:              branch,
			Sha:                 sha,
			Version:             version,
			AppVersion:          appVersion,
			SkipAppVersionCheck: skipAppVersionCheck,
		}

		s, err = helmtemplate.NewTemplateHelmChartTask(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if err := s.Run(validate, tagBuild); err != nil {
		return microerror.Mask(err)
	}

	log.Println("templated helm chart")

	return nil
}

// getProjectVersion retrieves version stored in project's Go source code. It
// looks up the value of variable `version` in `pkg/project/project.go`. If the
// file doesn't exist it returns an empty string.
func getProjectVersion(repoDir string) (string, error) {
	filePath := "pkg/project/project.go"
	varName := "version"

	content, err := ioutil.ReadFile(filepath.Join(repoDir, filePath))
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", microerror.Mask(err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return "", microerror.Mask(err)
	}

	version, err := parseString(node, varName)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return version, nil
}

// parseString returns the value of the variable varName found in r AbstractSyntaxTree.
func parseString(r ast.Node, varName string) (value string, err error) {
	ast.Inspect(r, func(n ast.Node) bool {
		d, ok := n.(*ast.ValueSpec)
		if ok {
			for _, id := range d.Names {
				if id.Name == varName {
					v := id.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit)
					if v.Kind == token.STRING {
						value, err = strconv.Unquote(v.Value)
						return false
					}
				}
			}
		}
		return true
	})

	if err != nil {
		return "", microerror.Mask(err)
	}

	return value, nil
}
