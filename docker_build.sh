chmod +x ./out/cheer_ingress_admin_linux
docker build -t harbor.aiagain.com/ik-arch/cheer-ingress:prod .
docker push harbor.aiagain.com/ik-arch/cheer-ingress:prod