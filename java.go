package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//JAVA - java template holder
var JAVA = Template{
	dockerFile: `FROM insideo/jre8:latest
WORKDIR /app
COPY {{.JarFilePath}} /app
CMD java -jar {{.JarFileName}} {{.AdditionalArguments}}
				 `,
	dockerignoreFile: `/target/
.apt_generated
.classpath
.factorypath
.project
.settings
.springBeans
.sts4-cache
.vscode
.idea
*.iws
*.iml
*.ipr
/nbproject/private/
/build/
/nbbuild/
/dist/
/nbdist/
/.nb-gradle/
/.vscode/
/.mvn/
/temp-builds/
					   `,
}

//CheckForExistingJarFile - check for existing jar file in maven/gradle build path
func CheckForExistingJarFile() (string, error) {
	jarFilePath := ""
	mavenTargetPath := filepath.Join(".", "target")
	err := filepath.Walk(mavenTargetPath, func(path string, info os.FileInfo, err error) error {
		if mavenTargetPath != path && strings.HasSuffix(path, ".jar") {
			jarFilePath = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	gradleTargetPath := filepath.Join(".", "build", "lib")
	err = filepath.Walk(gradleTargetPath, func(path string, info os.FileInfo, err error) error {
		if gradleTargetPath != path && strings.HasSuffix(path, ".jar") {
			jarFilePath = path
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return jarFilePath, nil
}

//GenerateDockerFiles - generate DockerFile and .dockerignoreFile
func (data *JavaBindingTemplate) GenerateDockerFiles() error {
	file, err := os.OpenFile("Dockerfile", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	tmpl := template.New("DockerFile")
	tmpl, err = tmpl.Parse(JAVA.dockerFile)
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
	tmpl2, err = tmpl2.Parse(JAVA.dockerignoreFile)
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
