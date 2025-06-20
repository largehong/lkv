# 生成时间: {{ timestamp }}
# traefik docs: https://doc.traefik.io/traefik/routing/providers/kv/
{{ $RouterKeys := getr `^/traefik/http/routers/.*` | keys -}}
{{ $RouterMatched := regexps `^/traefik/http/routers/(?<Name>\w+)/.*` $RouterKeys -}}
{{ $Routers := unique $RouterMatched.Name -}}
http:
    routers:
{{- range $router := $Routers }}
        {{ $router -}}:
            service: {{ sprintf "/traefik/http/routers/%s/service" $router | get }}
            rule: {{ sprintf "/traefik/http/routers/%s/rule" $router | get }}
            entryPoints:
            {{- $EntryPointRegex := sprintf `^/traefik/http/routers/%s/entrypoints/\d+` $router -}}
            {{ $EntryPointKeys := getr $EntryPointRegex | keys }}
            {{- range $EntryPointKey := $EntryPointKeys }}
                - {{ get $EntryPointKey }}
            {{ end -}}
{{- end }}
{{- $ServiceKeys := getr `^/traefik/http/services/.*` | keys -}}
{{ $ServiceMatched := regexps `^/traefik/http/services/(?<Name>\w+)/.*` $ServiceKeys -}}
{{ $Services := unique $ServiceMatched.Name }}
    services:
{{- range $service := $Services }}
        {{ $service -}}:
            loadBalancer:
            {{- $LoadBalancerRegex := sprintf `^/traefik/http/services/%s/loadbalancer/servers/\d+/url` $service -}}
            {{ $LoadBalancerKeys := getr $LoadBalancerRegex | keys }}
            {{- range $LoadBalancerKey := $LoadBalancerKeys }}
                - {{ get $LoadBalancerKey -}}
            {{ end }}
{{ end }}
