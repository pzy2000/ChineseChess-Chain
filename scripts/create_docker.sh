CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")

docker run \
 -itd  \
 -p 7999:7999  \
 -e TZ=Asia/Shanghai  \
 -v $PROJECT_PATH:/app/explorer-backend/chainmaker-explorer-backend   \
 --restart=always  \
 --privileged=true  \
 --name=explorer-backend-opennet  \
 --shm-size=4G  \
 --log-opt max-size=1000m  \
 --log-opt max-file=3  \
 ubuntu:20.04  \
 bash -c "cd /app/explorer-backend/chainmaker-explorer-backend/scripts/ && ./chainmaker-browser.bin -config ../configs/"
