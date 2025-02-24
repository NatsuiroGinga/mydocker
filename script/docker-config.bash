sudo tee /etc/docker/daemon.json <<-'EOF'
{
    "registry-mirrors": [
    	"dockerpull.cn",
        "docker.tbedu.top",
        "https://docker.1panel.live"
    ]
}
EOF

systemctl daemon-reload && sudo systemctl restart docker
