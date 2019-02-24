# Kubernitify
A minimalistic tool which helps you to quickly dockerize and kubernitify your application with a few prompts.

# Instructions
* Just place the executable in your project's root directory.

# Pre-requisites :
* github.com/AlecAivazis/survey

# Install dependencies
```
go get -u github.com/AlecAivazis/survey
```

# Generate Executable
`go build -o kubernitify main.go`

# Currently supported platforms
* Java
* Golang
* NodeJS

# Usage
```
kubernitify
  -java <Optional. For java platform setup>
  -go <Optional. For golang platform setup>
  -nodejs <Optional. For nodejs platform setup>
  -run <Optional. Directly run the application using existing Dockerfile and yaml file>
```