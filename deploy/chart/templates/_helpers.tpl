{{- define "greenstar.prefix" -}}
    {{- $prefix := "" -}}
    {{- if contains $.Chart.Name $.Release.Name -}}
        {{- print $.Release.Name -}}
    {{- else -}}
        {{- printf "%s-%s" $.Release.Name $.Chart.Name -}}
    {{- end -}}
{{- end -}}

{{- define "greenstar.versionAgnosticLabels" -}}
app.kubernetes.io/name: {{ $.Chart.Name | quote }}
app.kubernetes.io/instance: {{ $.Release.Name | quote }}
app.kubernetes.io/managed-by: {{ $.Release.Service | quote }}
{{- end -}}

{{- define "greenstar.commonLabels" -}}
{{ template "greenstar.versionAgnosticLabels" . }}
app.kubernetes.io/version: {{ $.Chart.AppVersion | replace "+" "_" | quote }}
helm.sh/chart: {{ printf "%s-%s" $.Chart.Name $.Chart.Version | replace "+" "_" | quote }}
{{- end -}}
