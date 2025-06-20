{{- $keys := getr `^/http/services/\w+/servers/\d+` | keys -}}
{{ $matchdServices := regexps `^/http/services/(?<ServiceName>\w+)/servers` $keys -}}
{{ $services := unique $matchdServices.ServiceName }}
{{ sprintf "services: %s" (join $services ",") }}