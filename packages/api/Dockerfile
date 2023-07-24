
FROM golang:1.19

# Installing OS dependencies here.
# Avoid using `apt` command outside this RUN block, will bloat image size.
RUN apt update -y && \
  apt install -y \
  make \
  postgresql-client && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

ENV ROOT /packages/api

WORKDIR $ROOT

# Install air for live reload
RUN go install github.com/cosmtrek/air@latest

# Install libs: database to struct

RUN go install gorm.io/gen/tools/gentool@latest

# Install libs: swagger to golang code

RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 3000

# Change to entrypoint after ECS config finish
CMD ["make", "run"]
