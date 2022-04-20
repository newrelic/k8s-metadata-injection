{{/* vim: set filetype=mustache: */}}

{{- /*
Naming helpers
*/ -}}

{{- define "nri-metadata-injection.name.admission" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "admission") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.admission" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "admission") }}
{{- end -}}

{{- define "nri-metadata-injection.name.admission-create" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "admission-create") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.admission-create" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "admission-create") }}
{{- end -}}

{{- define "nri-metadata-injection.name.admission-patch" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "admission-patch") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.admission-patch" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "admission-patch") }}
{{- end -}}

{{- define "nri-metadata-injection.name.self-signed-issuer" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "self-signed-issuer") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.self-signed-issuer" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "self-signed-issuer") }}
{{- end -}}

{{- define "nri-metadata-injection.name.root-issuer" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "root-issuer") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.root-issuer" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "root-issuer") }}
{{- end -}}

{{- define "nri-metadata-injection.name.webhook-cert" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.name" .) "suffix" "webhook-cert") }}
{{- end -}}

{{- define "nri-metadata-injection.fullname.webhook-cert" -}}
{{ include "newrelic.common.naming.truncateToDNSWithSuffix" (dict "name" (include "newrelic.common.naming.fullname" .) "suffix" "webhook-cert") }}
{{- end -}}
