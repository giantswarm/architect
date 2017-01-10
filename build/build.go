package build

type ProjectType string

var (
	GolangType ProjectType = "golang"
)

type Config struct {
	Type ProjectType
}

type Builder interface {
	Test() error
	Build() error
}

func New(config Config) Builder {
	switch config.Type {
	case GolangType:
		return &GolangBuilder{}
	}

	return nil
}
