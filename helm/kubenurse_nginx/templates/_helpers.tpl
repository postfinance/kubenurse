{{/* Build Kubenurse standard labels */}}
{{- define "common-labels" -}}
app.kubernetes.io/name: {{ .Chart.Name | quote }}
{{- end }}

{{- define "helm-labels" -}}
{{ include "common-labels" . }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | quote }}
{{- end }}

{{/* Build wide-used variables */}}
{{ define "name" -}}
{{ printf "%s" .Release.Name }}
{{- end }}

{{ define "image" -}}
{{ printf "%s:%s" .Values.daemonset.image.repository .Values.daemonset.image.tag }}
{{- end }}
