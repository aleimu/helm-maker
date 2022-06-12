package chart

const defaultChartfile = `apiVersion: v2
name: %s
description: A Helm chart for Kubernetes

# A chart can be either an 'application' or a 'library' chart.
#
# Application charts are a collection of templates that can be packaged into versioned archives
# to be deployed.
#
# Library charts provide useful utilities or functions for the chart developer. They're included as
# a dependency of application charts to inject those utilities and functions into the rendering
# pipeline. Library charts do not define any templates and therefore cannot be deployed.
type: application

# This is the chart version. This version number should be incremented each time you make changes
# to the chart and its templates, including the app version.
# Versions are expected to follow Semantic Versioning (https://semver.org/)
version: 0.1.0

# This is the version number of the application being deployed. This version number should be
# incremented each time you make changes to the application. Versions are not expected to
# follow Semantic Versioning. They should reflect the version the application is using.
# It is recommended to use it with quotes.
appVersion: "1.16.0"
`

const defaultValues = `
# Default values for <APPNAME>.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

<APPNAME>:
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

const defaultIgnore = `# Patterns to ignore when building packages.
# This supports shell glob matching, relative path matching, and
# negation (prefixed with !). Only one pattern per line.
.DS_Store
# Common VCS dirs
.git/
.gitignore
.bzr/
.bzrignore
.hg/
.hgignore
.svn/
# Common backup files
*.swp
*.bak
*.tmp
*.orig
*~
# Various IDEs
.project
.idea/
*.tmproj
.vscode/
`

const defaultIngress = `{{- if .Values.ingress.enabled -}}
{{- $fullName := include "<CHARTNAME>.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
{{- if and .Values.ingress.className (not (semverCompare ">=1.18-0" .Capabilities.KubeVersion.GitVersion)) }}
  {{- if not (hasKey .Values.ingress.annotations "kubernetes.io/ingress.class") }}
  {{- $_ := set .Values.ingress.annotations "kubernetes.io/ingress.class" .Values.ingress.className}}
  {{- end }}
{{- end }}
{{- if semverCompare ">=1.19-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1
{{- else if semverCompare ">=1.14-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1beta1
{{- else -}}
apiVersion: extensions/v1beta1
{{- end }}
kind: Ingress
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if and .Values.ingress.className (semverCompare ">=1.18-0" .Capabilities.KubeVersion.GitVersion) }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            {{- if and .pathType (semverCompare ">=1.18-0" $.Capabilities.KubeVersion.GitVersion) }}
            pathType: {{ .pathType }}
            {{- end }}
            backend:
              {{- if semverCompare ">=1.19-0" $.Capabilities.KubeVersion.GitVersion }}
              service:
                name: {{ $fullName }}
                port:
                  number: {{ $svcPort }}
              {{- else }}
              serviceName: {{ $fullName }}
              servicePort: {{ $svcPort }}
              {{- end }}
          {{- end }}
    {{- end }}
{{- end }}
`

const defaultDeployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "<CHARTNAME>.fullname" . }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "<CHARTNAME>.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "<CHARTNAME>.selectorLabels" . | nindent 8 }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "<CHARTNAME>.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
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
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
`

const defaultService = `apiVersion: v1
kind: Service
metadata:
  name: {{ include "<CHARTNAME>.fullname" . }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    {{- include "<CHARTNAME>.selectorLabels" . | nindent 4 }}
`

const defaultServiceAccount = `{{- if .Values.serviceAccount.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "<CHARTNAME>.serviceAccountName" . }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
  {{- with .Values.serviceAccount.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- end }}
`

const defaultHorizontalPodAutoscaler = `{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "<CHARTNAME>.fullname" . }}
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "<CHARTNAME>.fullname" . }}
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  metrics:
    {{- if .Values.autoscaling.targetCPUUtilizationPercentage }}
    - type: Resource
      resource:
        name: cpu
        targetAverageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
    {{- end }}
    {{- if .Values.autoscaling.targetMemoryUtilizationPercentage }}
    - type: Resource
      resource:
        name: memory
        targetAverageUtilization: {{ .Values.autoscaling.targetMemoryUtilizationPercentage }}
    {{- end }}
{{- end }}
`

const defaultNotes = `1. Get the application URL by running these commands:
{{- if .Values.ingress.enabled }}
{{- range $host := .Values.ingress.hosts }}
  {{- range .paths }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ $host.host }}{{ .path }}
  {{- end }}
{{- end }}
{{- else if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "<CHARTNAME>.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "<CHARTNAME>.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "<CHARTNAME>.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "<CHARTNAME>.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace {{ .Release.Namespace }} $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:$CONTAINER_PORT
{{- end }}
`

const defaultHelpers = `{{/*
Expand the name of the chart.
*/}}
{{- define "<CHARTNAME>.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "<CHARTNAME>.fullname" -}}
{{- if .Values.<APPNAME>.fullnameOverride }}
{{- .Values.<APPNAME>.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.<APPNAME>.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "<CHARTNAME>.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "<CHARTNAME>.labels" -}}
helm.sh/chart: {{ include "<CHARTNAME>.chart" . }}
{{ include "<CHARTNAME>.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "<CHARTNAME>.selectorLabels" -}}
app.kubernetes.io/name: {{ include "<CHARTNAME>.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "<CHARTNAME>.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "<CHARTNAME>.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
`

const defaultHelpers2 = `{{/*
Expand the name of the chart.
*/}}
{{- define "<CHARTNAME>.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "<CHARTNAME>.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "<CHARTNAME>.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "<CHARTNAME>.labels" -}}
helm.sh/chart: {{ include "<CHARTNAME>.chart" . }}
{{ include "<CHARTNAME>.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "<CHARTNAME>.selectorLabels" -}}
app.kubernetes.io/name: {{ include "<CHARTNAME>.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "<CHARTNAME>.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "<CHARTNAME>.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
`

const defaultTestConnection = `apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "<CHARTNAME>.fullname" . }}-test-connection"
  labels:
    {{- include "<CHARTNAME>.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args: ['{{ include "<CHARTNAME>.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
`

// Stderr is an io.Writer to which error messages can be written
//
// In Helm 4, this will be replaced. It is needed in Helm 3 to preserve API backward
