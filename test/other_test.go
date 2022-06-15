package test

import (
	"encoding/json"
	"fmt"
	"sigs.k8s.io/yaml"
	"strings"
	"testing"
)

// 这个json解析格式有错误...
const defaultValues = `
# Default values for <APPNAME>.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1
image:
  repository: nginx
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""
podAnnotations: {}
podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  annotations: {}
	# kubernetes.io/ingress.class: nginx
	# kubernetes.io/tls-acme: "true"
  hosts:
	- host: chart-example.local
	  paths:
		- path: /
		  pathType: ImplementationSpecific
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}
`

const defaultValuesJson = `
{
  "appname": "appname",
  "version": "version",
  "value": {
    "replicaCount": 1,
    "image": {
      "repository": "nginx",
      "pullPolicy": "IfNotPresent",
      "tag": ""
    },
    "imagePullSecrets": [],
    "nameOverride": "",
    "fullnameOverride": "",
    "podAnnotations": {},
    "podSecurityContext": {},
    "securityContext": {},
    "service": {
      "type": "ClusterIP",
      "port": 80
    },
    "resources": {},
    "nodeSelector": {},
    "affinity": {},
	"secret":"",
	"configMap": "",
	"volumeMounts":[],
	"volumes":[]
  }
}
`

const temp1 = `
appname: appname
version: version
value:
  affinity: {}
  fullnameOverride: ""
  image:
    pullPolicy: IfNotPresent
    repository: nginx
    tag: ""
  imagePullSecrets: []
  nameOverride: ""
  nodeSelector: {}
  podAnnotations: {}
  podSecurityContext: {}
  replicaCount: 1
  resources: {}
  securityContext: {}
  service:
    port: 80
    type: ClusterIP
`

const defaultDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include ".Values.<APPNAME>.fullname" . }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
spec:
  {{- if not .Values.<APPNAME>.autoscaling.enabled }}
  replicas: {{ .Values.<APPNAME>.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "<CHARTNAME>.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.<APPNAME>.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "<CHARTNAME>.selectorLabels" . | nindent 8 }}
    spec:
      {{- with .Values.<APPNAME>.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      # serviceAccountName: {{ include "<CHARTNAME>.serviceAccountName" . }}
      # securityContext:
      #   {{- toYaml .Values.<APPNAME>.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Values.<APPNAME>.Name }}
          securityContext:
            {{- toYaml .Values.<APPNAME>.securityContext | nindent 12 }}
          image: "{{ .Values.<APPNAME>.image.repository }}:{{ .Values.<APPNAME>.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {{- toYaml .Values.<APPNAME>.resources | nindent 12 }}
      {{- with .Values.<APPNAME>.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.<APPNAME>.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.<APPNAME>.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
`

func TestGenValueByJson(t *testing.T) {

	var n map[string]interface{}
	err := json.Unmarshal([]byte(defaultValuesJson), &n)
	if err != nil {

		fmt.Println(err)
	} else {
		fmt.Printf("%+v", n)
	}

	yvalue, err := yaml.Marshal(n)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(yvalue))
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal([]byte(strings.ReplaceAll(defaultDeployment, "<APPNAME>", "myapp")), &m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", m)
	}

}

func TestNewDeployment(t *testing.T) {
	var m map[string]interface{}
	dep := strings.ReplaceAll(defaultDeployment, "<APPNAME>", "myapp")
	dep = strings.ReplaceAll(dep, "<CHARTNAME>", "mychart")
	if err := yaml.Unmarshal([]byte(dep), &m); err != nil {
		fmt.Println(err)
		//fmt.Println(dep)
	} else {
		fmt.Printf("%+v", m)
	}
}
