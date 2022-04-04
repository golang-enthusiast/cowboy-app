{{/*
Expand the name of the chart.
*/}}
{{- define "cowboy-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cowboy-app.fullname" -}}
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
{{- define "cowboy-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cowboy-app.labels" -}}
helm.sh/chart: {{ include "cowboy-app.chart" . }}
{{ include "cowboy-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cowboy-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cowboy-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cowboy-app.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cowboy-app.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Change how the image is assigned based on the skaffold flag.
*/}}
{{- define "cowboy-app.image" -}}
{{- if .Values.skaffold -}}
{{- .Values.skaffoldImage -}}
{{- else -}}
{{- printf "%s%s:%s" .Values.image.hostname .Values.image.repository .Values.image.tag -}}
{{- end -}}
{{- end -}}

{{/*
Change how the cron image is assigned based on the skaffold flag.
*/}}
{{- define "cowboy-app.cron-image" -}}
{{- if .Values.skaffold -}}
{{- .Values.cron.skaffoldImage -}}
{{- else -}}
{{- printf "%s%s:%s" .Values.cron.image.hostname .Values.cron.image.repository .Values.image.tag -}}
{{- end -}}
{{- end -}}

{{/*
Change how the data migration image is assigned based on the skaffold flag.
*/}}
{{- define "cowboy-app.migration-image" -}}
{{- if .Values.skaffold -}}
{{- .Values.migration.skaffoldImage -}}
{{- else -}}
{{- printf "%s%s:%s" .Values.migration.image.hostname .Values.migration.image.repository .Values.image.tag -}}
{{- end -}}
{{- end -}}
