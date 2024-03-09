# Use a Docker image that comes with Go version 1.18.
FROM golang:1.18

# Set working directory in a container to /app
WORKDIR /app

# Applied metadata to Docker Object
LABEL version="1.0" maintainer="yba-lamor & adiane-lamor & sefaye-lamor & npouille-lamor & papasarr-lamor"

# Copy all the project in container's working directory
COPY . .

# Run this command to install dependencies
RUN go mod download

# Run this command to build an excecutable file
RUN go build -o main .

# Listen to port 8081
EXPOSE 8081

# Command for running excecutable file 
CMD [ "./main" ]