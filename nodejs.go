package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//NODEJS - nodejs template holder
var NODEJS = Template{
	dockerFile: `FROM node:alpine
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY {{.ProjectPath}} .
CMD [ "npm", "start" ]
				 `,
	dockerignoreFile: `node_modules
npm-debug.log
					   `,
}

//CheckForNodejsProjectPath - check for nodejs project path
func CheckForNodejsProjectPath() (string, error) {
	nodeJSTargetPath := filepath.Join(".")
	truePath := ""
	err := filepath.Walk(nodeJSTargetPath, func(path string, info os.FileInfo, err error) error {
		if nodeJSTargetPath != path && strings.HasSuffix(path, "package.json") {
			truePath = nodeJSTargetPath
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return truePath, nil
}

//GenerateDockerFiles - generate DockerFile and .dockerignoreFile
func (data *NodeJSBindingTemplate) GenerateDockerFiles() error {
	file, err := os.OpenFile("Dockerfile", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl := template.New("DockerFile")
	tmpl, err = tmpl.Parse(NODEJS.dockerFile)
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
	tmpl2, err = tmpl2.Parse(NODEJS.dockerignoreFile)
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
