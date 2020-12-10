# hello-world-service

1. `docker build -t linuxdynasty/hello-world-service:latest .`
2. `docker push linuxdynasty/hello-world-service`
3. To test that it is running. `docker run -p "80:80" linuxdynasty/hello-world-service:latest`

# Environment Variables
* `STATSD_EXPORTER` defaults to `127.0.0.1:9125`

# Dockerhub
*   https://hub.docker.com/repository/docker/linuxdynasty/hello-world-service