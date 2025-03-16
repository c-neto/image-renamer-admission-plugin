{{- define "image-renamer-admission-plugin.fullname" -}}
{{ .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "image-renamer-admission-plugin.selectorLabels" -}}
{{- include "image-renamer-admission-plugin.labels" . -}}
{{- end }}

{{- define "image-renamer-admission-plugin.labels" -}}
helm.sh/chart: {{ include "image-renamer-admission-plugin.chart" . | quote }}
app.kubernetes.io/name: {{ include "image-renamer-admission-plugin.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end }}

{{- define "image-renamer-admission-plugin.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "image-renamer-admission-plugin.chart" -}}
{{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}
