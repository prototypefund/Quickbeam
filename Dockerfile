FROM debian:12

RUN printf -- '--> Create build env\n' \
 && apt-get update \
 && apt-get install -y --no-install-recommends golang-go make ca-certificates\
 && apt-get clean \
 && rm -rf /var/cache/apt \
 && mkdir /src \
 && printf -- '<-- Done creating build env\n'

VOLUME /src

WORKDIR /src
CMD ["make", "build"]
