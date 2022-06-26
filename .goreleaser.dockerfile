FROM busybox:stable

WORKDIR /keeper

# Creates non root user
ENV USER=keeper
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY keeper keeper

# Running as keeper
USER keeper:keeper

ENTRYPOINT ["/keeper"]