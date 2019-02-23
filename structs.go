package main

//Template - holder struct for templates
type Template struct {
	dockerFile       string
	dockerignoreFile string
}

//JavaBindingTemplate - java binding template
type JavaBindingTemplate struct {
	JarFilePath         string
	JarFileName         string
	AdditionalArguments string
}

//GolangExecutableBindingTemplate - golang executable binding template
type GolangExecutableBindingTemplate struct {
	ExecutablePath      string
	ExecutableName      string
	AdditionalArguments string
}

//GolangSourceCodeBindingTemplate - golang source code binding template
type GolangSourceCodeBindingTemplate struct {
	SourceCodePath      string
	Dependencies        string
	ExecutableName      string
	AdditionalArguments string
}

//NodeJSBindingTemplate - nodejs binding template
type NodeJSBindingTemplate struct {
	ProjectPath string
}

//KubernetesBindingTemplate - kubernetes binding template
type KubernetesBindingTemplate struct {
	AppName            string
	Namespace          string
	ExternalPortNumber string
	InternalPortNumber string
	PodCount           string
	ImageName          string
}
