package template

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
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
	)

	fs := afero.NewOsFs()
	ctx := context.Background()

	var appVersion string
	{
		dir, err := gitrepo.TopLevel(ctx, ".")
		if err != nil {
			return microerror.Mask(err)
		}

		c := gitrepo.Config{
			Dir: dir,
		}

		repo, err := gitrepo.New(c)
		if err != nil {
			return microerror.Mask(err)
		}

		appVersion, err = getProjectVersion(repo, "origin/master")
		if err != nil {
			return microerror.Mask(err)
		}
	}

	log.Printf("templating helm chart\ndir: %s\nsha: %s\ntag: %s\napp-version: %s\nversion: %s\n", chartDir, sha, tag, appVersion, version)

	var s *helmtemplate.TemplateHelmChartTask
	{
		c := helmtemplate.Config{
			Fs:           fs,
			ChartDir:     chartDir,
			Branch:       branch,
			Sha:          sha,
			ChartVersion: version,
			AppVersion:   appVersion,
		}

		s, err = helmtemplate.NewTemplateHelmChartTask(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if err := s.Run(); err != nil {
		return microerror.Mask(err)
	}

	log.Println("templated helm chart")

	return nil
}

// getProjectVersion retrieves version stored in project's Go source code.
// It looks up the value of variable `version` in `pkg/project/project.go` file
// on version defined in ref.
func getProjectVersion(repo *gitrepo.Repo, ref string) (string, error) {
	filePath := "pkg/project/project.go"
	varName := "version"

	content, err := repo.GetFileContent(filePath, ref)
	if err != nil {
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
