package main

import (
	"html/template"
	"os"
)

//K8S - template kubernetes yaml file
const K8S = `#IMAGE={{.ImageName}}
apiVersion: v1
kind: Service
metadata:
  name: {{.AppName}}
  namespace: {{.Namespace}}
spec:
  type: LoadBalancer
  selector:
    app: {{.AppName}}
  ports:
    - protocol: TCP
      port: {{.ExternalPortNumber}}
      nodePort:  {{.ExternalPortNumber}}
      targetPort: {{.InternalPortNumber}}
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: {{.AppName}}
  namespace: {{.Namespace}}
spec:
  replicas: {{.PodCount}}
  selector:
    matchLabels:
      app: {{.AppName}}
  template:
    metadata:
      labels:
        app: {{.AppName}}
    spec:
      containers:
        - name: {{.AppName}}
          image: '{{.ImageName}}'
          ports:
            - containerPort: {{.InternalPortNumber}}
`

//GenerateYamlFile -generate k8s yaml file
func (data *KubernetesBindingTemplate) GenerateYamlFile() error {
	tmpl := template.New("YamlFile")
	tmpl, err := tmpl.Parse(K8S)
	if err != nil {
		return err
	}
	file, err := os.OpenFile("k8s.yaml", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
