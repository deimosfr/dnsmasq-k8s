{{/*
Expand the name of the chart.
*/}}
{{- define "dnsmasq-k8s-ui.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dnsmasq-k8s-ui.fullname" -}}
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
{{- define "dnsmasq-k8s-ui.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "dnsmasq-k8s-ui.labels" -}}
helm.sh/chart: {{ include "dnsmasq-k8s-ui.chart" . }}
{{ include "dnsmasq-k8s-ui.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "dnsmasq-k8s-ui.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dnsmasq-k8s-ui.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "dnsmasq-k8s-ui.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "dnsmasq-k8s-ui.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Validate values
*/}}
{{- define "dnsmasq-k8s-ui.validateValues" -}}
{{- if and (gt (int .Values.replicaCount) 1) .Values.dhcp.enabled }}
{{- fail "Can't support multiple replicas because of dnsmasq implementation with concurrency issue on a single configmap and no possible diff between 2 versions" }}
{{- end }}
{{- end }}
