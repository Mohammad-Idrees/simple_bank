#### Docker commands
docker ps
docker ps -a
docker images
docker pull <image>:<tag>
docker run --name <container_name> -p <host_ports:container_ports> -e <environment_variable> -d <image>:<tag>
docker exec -it <container_name_or_id> <command> [args]
docker logs <container_name_or_id>
docker stop <container_name_or_id>
docker start <container_name_or_id>
docker rm <container_name_or_id>