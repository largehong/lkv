server {
    listen 80;
    server_name api.lovezsh.com;

    {{ range $service := getp "/http/services"}}
    location /api/v1/{{ .service.Name }}/ {
        proxy_pass http://{{ .service.Name }}/;
    }
    {{ end }}
}