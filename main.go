package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey"
)

var java = flag.Bool("java", false, "Intialize For Java Platform")
var golang = flag.Bool("go", false, "Initialize For Golang Platform")
var nodejs = flag.Bool("nodejs", false, "Initialize For Nodejs Platform")
var run = flag.Bool("run", false, "Dockerize the application and run on kubernetes")

var imageName = ""
var newDockerFile = true
var newYamlFile = true

func initApp() {
	fmt.Println("Welcome to kubernitify.")
	fmt.Println("This tool helps you to quickly dockerize and kubernitify your application with few prompts.")
	if !*run {
		if _, err := os.Stat("Dockerfile"); err == nil {
			confirmation := true
			confirmationPrompt := &survey.Confirm{
				Message: "Dockerfile found on current path. Do you wish to continue using it ?",
			}
			survey.AskOne(confirmationPrompt, &confirmation, nil)
			if confirmation {
				newDockerFile = false
			}
		}
		if _, err := os.Stat("k8s.yaml"); err == nil {
			confirmation := true
			confirmationPrompt := &survey.Confirm{
				Message: "k8s.yaml found on current path. Do you wish to continue using it ?",
			}
			survey.AskOne(confirmationPrompt, &confirmation, nil)
			if confirmation {
				newYamlFile = false
			}
		}
	}
}

func main() {
	flag.Parse()
	initApp()
	if !*run {
		if newDockerFile {
			if *java {
				javaHandler()
			} else if *golang {
				golangHandler()
			} else if *nodejs {
				nodeJsHandler()
			} else {
				platform := ""
				prompt := &survey.Select{
					Message: "Select your application platform:",
					Options: []string{"java", "golang", "nodejs"},
					Default: "java",
				}
				survey.AskOne(prompt, &platform, nil)
				switch platform {
				case "java":
					javaHandler()
					break
				case "golang":
					golangHandler()
					break
				case "nodejs":
					nodeJsHandler()
					break
				}
			}
		}
		if newYamlFile {
			generateKubernetesYamlFile()
		}
	}
	image, err := extractImageName()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	imageName = image
	if imageName == "" {
		prompt := &survey.Input{
			Message: "Unable to identify image name for the application. Please specify the image name:",
		}
		survey.AskOne(prompt, &imageName, nil)
	}
	confirmation := true
	if !*run {
		confirmationPrompt := &survey.Confirm{
			Message: "Do you wish to run the application ?",
		}
		survey.AskOne(confirmationPrompt, &confirmation, nil)
	}
	if !confirmation {
		fmt.Println("Copy the following commands to run the application: ")
		fmt.Println("-> docker built -t " + imageName + " .")
		fmt.Println("-> kubectl create -f k8s.yaml")
		exit()
	}
	runApplication()
	exit()
}

func javaHandler() {
	jarFilePath, err := CheckForExistingJarFile()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	prompt := &survey.Input{
		Message: "Enter path for the jar file:",
	}
	if jarFilePath != "" {
		prompt = &survey.Input{
			Message: "Enter path for the jar file:",
			Default: jarFilePath,
		}
	}
	survey.AskOne(prompt, &jarFilePath, nil)
	additionalCLIArgs := ""
	additionalCLIArgsPrompt := &survey.Input{
		Message: "Additional CLI Args(if any):",
	}
	_, fileName := filepath.Split(jarFilePath)
	survey.AskOne(additionalCLIArgsPrompt, &additionalCLIArgs, nil)
	template := JavaBindingTemplate{
		JarFilePath:         jarFilePath,
		JarFileName:         fileName,
		AdditionalArguments: additionalCLIArgs,
	}
	err = template.GenerateDockerFiles()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
}

func golangHandler() {
	applicationType := ""
	prompt := &survey.Select{
		Message: "Select the type of golang application:",
		Options: []string{"Standalone executable", "Source code"},
		Default: "Standalone executable",
	}
	survey.AskOne(prompt, &applicationType, nil)
	switch applicationType {
	case "Standalone executable":
		fmt.Println("NOTE: Before continuing, please ensure you have *-linux-386 executable of the application.")
		confirmation := false
		confirmationPrompt := &survey.Confirm{
			Message: "Continue?",
		}
		survey.AskOne(confirmationPrompt, &confirmation, nil)
		if !confirmation {
			exit()
		}
		executablePath, err := CheckForExistingExecutables()
		if err != nil {
			fmt.Println(fmt.Sprintf("%s", err))
			os.Exit(0)
		}
		executablePrompt := &survey.Input{
			Message: "Enter path for the executable file:",
		}
		if executablePath != "" {
			executablePrompt = &survey.Input{
				Message: "Enter path for the executable file:",
				Default: executablePath,
			}
		}
		survey.AskOne(executablePrompt, &executablePath, nil)
		additionalCLIArgs := ""
		additionalCLIArgsPrompt := &survey.Input{
			Message: "Additional CLI Args(if any):",
		}
		survey.AskOne(additionalCLIArgsPrompt, &additionalCLIArgs, nil)
		_, fileName := filepath.Split(executablePath)
		template := GolangExecutableBindingTemplate{
			ExecutablePath:      executablePath,
			ExecutableName:      fileName,
			AdditionalArguments: additionalCLIArgs,
		}
		err = template.GenerateDockerFiles()
		if err != nil {
			fmt.Println(fmt.Sprintf("%s", err))
			os.Exit(0)
		}
		break
	case "Source code":
		text := ""
		dependenciesPrompt := &survey.Multiline{
			Message: "Enter the list of dependencies in each line(if any)",
		}
		survey.AskOne(dependenciesPrompt, &text, nil)
		dependencies := strings.Split(text, "\n")
		fmt.Println(dependencies)
		dependenciesString := ""
		for _, dependency := range dependencies {
			dependenciesString = dependenciesString + "RUN go get -u " + dependency + "\r\n"
		}
		mainFilePath, err := CheckIfAnyMainGoFileExists()
		if err != nil {
			fmt.Println(fmt.Sprintf("%s", err))
			os.Exit(0)
		}
		mainFilePathPrompt := &survey.Input{
			Message: "Enter path for the main.go file(relative to current directory):",
		}
		if mainFilePath != "" {
			mainFilePathPrompt = &survey.Input{
				Message: "Enter path for the main.go file(relative to current directory):",
				Default: mainFilePath,
			}
		}
		survey.AskOne(mainFilePathPrompt, &mainFilePath, nil)
		additionalCLIArgs := ""
		additionalCLIArgsPrompt := &survey.Input{
			Message: "Additional CLI Args(if any):",
		}
		survey.AskOne(additionalCLIArgsPrompt, &additionalCLIArgs, nil)
		template := GolangSourceCodeBindingTemplate{
			SourceCodePath:      mainFilePath,
			ExecutableName:      mainFilePath,
			Dependencies:        dependenciesString,
			AdditionalArguments: additionalCLIArgs,
		}
		err = template.GenerateDockerFiles()
		if err != nil {
			fmt.Println(fmt.Sprintf("%s", err))
			os.Exit(0)
		}
		break
	}
}

func nodeJsHandler() {
	projectPath, err := CheckForNodejsProjectPath()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	projectPathPrompt := &survey.Input{
		Message: "Enter path for the nodejs project:",
	}
	if projectPath != "" {
		projectPathPrompt = &survey.Input{
			Message: "Enter path for the nodejs file:",
			Default: projectPath,
		}
	}
	survey.AskOne(projectPathPrompt, &projectPath, nil)
	template := NodeJSBindingTemplate{
		ProjectPath: projectPath,
	}
	err = template.GenerateDockerFiles()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
}

func generateKubernetesYamlFile() {
	AppName := ""
	prompt := &survey.Input{
		Message: "Enter kubernetes app name:",
		Default: "myapp",
	}
	survey.AskOne(prompt, &AppName, nil)
	PodCount := ""
	prompt = &survey.Input{
		Message: "Enter kubernetes pod count:",
		Default: "1",
	}
	survey.AskOne(prompt, &PodCount, nil)
	Namespace := ""
	prompt = &survey.Input{
		Message: "Enter kubernetes namespace:",
		Default: "default",
	}
	survey.AskOne(prompt, &Namespace, nil)
	ImageName := AppName + ":latest"
	imageName = ImageName
	InternalPortNumber := ""
	prompt = &survey.Input{
		Message: "Enter port number on which application runs natively:",
		Default: "80",
	}
	survey.AskOne(prompt, &InternalPortNumber, nil)
	ExternalPortNumber := ""
	prompt = &survey.Input{
		Message: "Enter port number on which application should be accesible from kubernetes:",
		Default: "80",
	}
	survey.AskOne(prompt, &ExternalPortNumber, nil)
	template := KubernetesBindingTemplate{
		AppName:            AppName,
		Namespace:          Namespace,
		ExternalPortNumber: ExternalPortNumber,
		InternalPortNumber: InternalPortNumber,
		PodCount:           PodCount,
		ImageName:          ImageName,
	}
	err := template.GenerateYamlFile()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
}

func exit() {
	fmt.Println("Exiting Kubernitify.")
	fmt.Println("Thank you.")
	os.Exit(0)
}

func runApplication() {
	cmd := exec.Command("sh", "-c", "docker build -t "+imageName+" .")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	cmd = exec.Command("sh", "-c", "kubectl create -f k8s.yaml")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
}

func extractImageName() (string, error) {
	file, err := os.Open("k8s.yaml")
	if err != nil {
		return "", err
	}
	defer file.Close()
	image := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		if strings.Contains(data, "#IMAGE=") {
			image = strings.Split(data, "=")[1]
		}
		break
	}
	return image, nil
}
