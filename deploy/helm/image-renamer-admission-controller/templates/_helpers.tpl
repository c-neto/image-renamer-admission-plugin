{{- define "image-renamer-admission-controller.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "image-renamer-admission-controller.selectorLabels" -}}
{{- include "image-renamer-admission-controller.labels" . -}}
{{- end }}

{{- define "image-renamer-admission-controller.labels" -}}
helm.sh/chart: {{ include "image-renamer-admission-controller.chart" . }}
{{ include "image-renamer-admission-controller.selectorLabels" . }}
{{- end }}

{{- define "image-renamer-admission-controller.chart" -}}
{{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}
