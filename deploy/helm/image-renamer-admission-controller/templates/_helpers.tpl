{{- define "image-renamer-admission-plugin.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "image-renamer-admission-plugin.selectorLabels" -}}
{{- include "image-renamer-admission-plugin.labels" . -}}
{{- end }}

{{- define "image-renamer-admission-plugin.labels" -}}
helm.sh/chart: {{ include "image-renamer-admission-plugin.chart" . }}
{{ include "image-renamer-admission-plugin.selectorLabels" . }}
{{- end }}

{{- define "image-renamer-admission-plugin.chart" -}}
{{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}
