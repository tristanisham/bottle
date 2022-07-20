package cli


func nginxService() string {

	return `
	server {

        server_name example.com www.example.com; # Replace with website's domain

        location / {
                proxy_pass http://localhost:8080; # Replace with website's port (default 8080)
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection 'upgrade';
                proxy_set_header Host $host;
                proxy_cache_bypass $http_upgrade;
    	}
	}`
}
