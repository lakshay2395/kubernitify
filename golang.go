package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//GOLANGEXECUTABLE - golang executable template holder
var GOLANGEXECUTABLE = Template{
	dockerFile: `FROM alpine:3.8
WORKDIR /app
COPY {{.ExecutablePath}} /app
CMD ./{{.ExecutableName}} {{.AdditionalArguments}}
				 `,
	dockerignoreFile: ``,
}

//GOLANGSOURCECODE - golang source code template holder
var GOLANGSOURCECODE = Template{
	dockerFile: `FROM golang:1.11-alpine
WORKDIR /app
COPY {{.SourceCodePath}} /app
{{.Dependencies}}
CMD go run ./{{.ExecutableName}} {{.AdditionalArguments}}
				 `,
	dockerignoreFile: ``,
}

//CheckForExistingExecutables - check for existing linux-386 executable
func CheckForExistingExecutables() (string, error) {
	executableFilePath := ""
	targetPath := filepath.Join(".")
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if targetPath != path && strings.HasSuffix(path, "-linux-386") {
			executableFilePath = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return executableFilePath, nil
}

//CheckIfAnyMainGoFileExists - check if any main.go file exists in folder
func CheckIfAnyMainGoFileExists() (string, error) {
	mainGoFilePath := ""
	targetPath := filepath.Join(".")
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if targetPath != path && strings.HasSuffix(path, "main.go") {
			mainGoFilePath = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return mainGoFilePath, nil
}

//GenerateDockerFiles - generate DockerFile and .dockerignoreFile
func (data *GolangExecutableBindingTemplate) GenerateDockerFiles() error {
	file, err := os.OpenFile("Dockerfile", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl := template.New("DockerFile")
	tmpl, err = tmpl.Parse(GOLANGEXECUTABLE.dockerFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	file.Close()

	file, err = os.OpenFile(".dockerignore", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl2 := template.New(".dockerignore")
	tmpl2, err = tmpl2.Parse(GOLANGEXECUTABLE.dockerignoreFile)
	if err != nil {
		return err
	}
	err = tmpl2.Execute(file, nil)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

//GenerateDockerFiles - generate DockerFile and .dockerignoreFile
func (data *GolangSourceCodeBindingTemplate) GenerateDockerFiles() error {
	file, err := os.OpenFile("Dockerfile", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl := template.New("DockerFile")
	tmpl, err = tmpl.Parse(GOLANGSOURCECODE.dockerFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	file.Close()

	file, err = os.OpenFile(".dockerignore", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl2 := template.New(".dockerignore")
	tmpl2, err = tmpl2.Parse(GOLANGSOURCECODE.dockerignoreFile)
	if err != nil {
		return err
	}
	err = tmpl2.Execute(file, nil)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
