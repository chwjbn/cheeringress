chmod +x ./out/cheer_ingress_linux
docker build -t chwjbn/cheer-ingress:prod .
docker push chwjbn/cheer-ingress:prod
