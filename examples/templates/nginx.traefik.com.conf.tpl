# 生成时间: {{ timestamp }}
{{ $ServiceKeys := getr `^/traefik/http/services/\w+/loadbalancer/servers/\d+/url` | keys -}}
{{ $ServiceMatched := regexps `^/traefik/http/services/(?<Name>\w+)/loadbalancer/servers/\d+/url` $ServiceKeys -}}
{{ $Services := unique $ServiceMatched.Name }}

{{- range $Service := $Services }} 
upstream {{ $Service }} {
    {{- $LoadBalancerRegex := sprintf `^/traefik/http/services/%s/loadbalancer/servers/\d+/url` $Service -}}
    {{ $LoadBalancerKeys := getr $LoadBalancerRegex | keys }}
    {{- range $LoadBalancerKey := $LoadBalancerKeys }}
    {{ $URL := get $LoadBalancerKey | urlparse  -}}
    server {{ $URL.Host -}};
    {{- end }}
}
{{ end }}

server {
    listen 80;
    server_name api.lovezsh.com;

    {{ range $Service := $Services }}
    location /api/v1/{{ $Service }}/ {
        proxy_pass http://{{ $Service }}/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $http_x_forwarded_for;
    }
    {{ end }}
}