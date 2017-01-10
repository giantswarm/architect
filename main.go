package main

import (
	"log"

	"github.com/giantswarm/architect/cmd"
)

// func main() {
// 	log.Printf("starting build\n")
//
// 	flag.Parse()
//
// 	organisation := flag.Arg(0)
// 	project := flag.Arg(1)
// 	if organisation == "" {
// 		log.Fatalln("organisation cannot be empty")
// 	}
// 	if project == "" {
// 		log.Fatalln("project cannot be empty")
// 	}
// 	// goos := "linux"
// 	// goarch := "amd64"
//
// 	if err := os.Setenv("GOROOT", "/usr/local/go"); err != nil {
// 		log.Fatalf("could not set GOROOT env var: %v\n", err)
// 	}
// 	if err := os.Setenv("GOPATH", "/go"); err != nil {
// 		log.Fatalf("could not set GOPATH env var: %v\n", err)
// 	}
// 	if err := os.Setenv("PATH", fmt.Sprintf("/usr/local/go/bin/:%v", os.Getenv("PATH"))); err != nil {
// 		log.Fatalf("could not set GOPATH env var: %v\n", err)
// 	}
//
// 	buildDir, err := setUpGoDirectory(CodeDir, organisation, project)
// 	if err != nil {
// 		log.Fatalf("could not set up go directory: %v\n", err)
// 	}
//
// 	if err := runGoTests(buildDir); err != nil {
// 		log.Fatalf("could not run go tests: %v\n", err)
// 	}
//
// 	if err := buildGoBinary(buildDir); err != nil {
// 		log.Fatalf("could not build go binary: %v\n", err)
// 	}
//
// 	imageName, err := buildImage(buildDir, Registry, organisation, project, "test")
// 	if err != nil {
// 		log.Fatalf("could not build iamge: %v\n", err)
// 	}
//
// 	if err := runContainer(imageName); err != nil {
// 		log.Fatalf("could not run container: %v\n", err)
// 	}
// }

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
