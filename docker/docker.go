package docker

import (
	"fmt"
	"log"
)

func BuildImage(buildDir, registry, organisation, project, tag string) (string, error) {
	log.Printf("building image\n")

	imageName := fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, tag)

	if err := command.RunCommand([]string{"docker", "build", "-t", imageName, "."}, buildDir); err != nil {
		return "", err
	}

	return imageName, nil
}

func RunContainer(imageName string) error {
	log.Printf("running container\n")

	if err := command.RunCommand([]string{"docker", "run", imageName, "version"}, ""); err != nil {
		return err
	}

	if err := command.RunCommand([]string{"docker", "run", imageName, "--help"}, ""); err != nil {
		return err
	}

	return nil
}
