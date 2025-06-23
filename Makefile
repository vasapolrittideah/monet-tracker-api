.PHONY: build-image

build-image:
	docker build -t money-tracker-api-local:latest .

# === Bazel Targets ===
.PHONY: bazel-gazelle bazel-build bazel-tidy bazel-test

bazel-gazelle:
	bazel run //:gazelle

bazel-build:
	bazel build //...

bazel-tidy:
	bazel mod tidy

bazel-test:
	bazel test //...

bazel-clean:
	bazel clean

# === Combined Commands ===
.PHONY: init all

init: bazel-gazelle bazel-tidy

all: bazel-gazelle bazel-tidy bazel-build bazel-test
