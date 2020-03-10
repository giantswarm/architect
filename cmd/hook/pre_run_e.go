package hook

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

func PreRunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var err error

	var repo *gitrepo.Repo
	{

		dir, err := gitrepo.TopLevel(ctx, ".")
		if err != nil {
			return microerror.Mask(err)
		}

		c := gitrepo.Config{
			Dir: dir,
		}

		repo, err = gitrepo.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	defaultSha, err := repo.HeadSHA(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	defaultTag, err := repo.HeadTag(ctx)
	if gitrepo.IsReferenceNotFound(err) {
		defaultTag = ""
	} else if err != nil {
		return microerror.Mask(err)
	}

	defaultBranch, err := repo.HeadBranch(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	// Define the version we are building.
	gitVersion, err := repo.ResolveVersion(ctx, "HEAD")
	if err != nil {
		return microerror.Mask(err)
	}

	srcVersion, err := getProjectVersion(repo)
	if err != nil {
		return microerror.Mask(err)
	}

	cmd.PersistentFlags().String("branch", defaultBranch, "git branch being built")
	cmd.PersistentFlags().String("sha", defaultSha, "git SHA1 being built")
	cmd.PersistentFlags().String("tag", defaultTag, "git tag being built")
	cmd.PersistentFlags().String("version", gitVersion, "version found in git")
	cmd.PersistentFlags().String("source-version", srcVersion, "version found in source code")

	return nil
}

// getProjectVersion retrieves version stored in project's Go source code.
// It looks up the value of variable `version` in `pkg/project/project.go` file
// on master branch.
func getProjectVersion(repo *gitrepo.Repo) (string, error) {
	filePath := "pkg/project/project.go"
	varName := "version"

	content, err := repo.GetFileContent(filePath)
	if err != nil {
		return "", microerror.Mask(err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, 0)
	if err != nil {
		return "", microerror.Mask(err)
	}

	value, err := parseString(node, varName)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return value, nil
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
