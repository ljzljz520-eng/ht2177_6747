# Official Go image with the complete toolchain required by the evaluator.
FROM golang:1.22

WORKDIR /app

COPY ["go.mod", "go.sum", "/app/"]
RUN GOWORK=off GOTOOLCHAIN=local go mod download

COPY . .

# The regression test is intentionally red while the injected bug is present.
# See BUG_REPRO.md for the reproducible failing workflow.
RUN GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go build ./...

CMD ["bash"]
