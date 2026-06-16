APP := qed
CMD := ./cmd/qed
DIST_DIR := dist
LOCAL_BIN := $(HOME)/.local/bin
GLOBAL_BIN := /usr/local/bin
PLATFORMS ?= \
	aix/ppc64 \
	darwin/amd64 \
	darwin/arm64 \
	dragonfly/amd64 \
	freebsd/386 \
	freebsd/amd64 \
	freebsd/arm \
	freebsd/arm64 \
	freebsd/riscv64 \
	illumos/amd64 \
	linux/386 \
	linux/amd64 \
	linux/arm \
	linux/arm64 \
	linux/loong64 \
	linux/mips \
	linux/mips64 \
	linux/mips64le \
	linux/mipsle \
	linux/ppc64 \
	linux/ppc64le \
	linux/riscv64 \
	linux/s390x \
	netbsd/386 \
	netbsd/amd64 \
	netbsd/arm \
	netbsd/arm64 \
	openbsd/386 \
	openbsd/amd64 \
	openbsd/arm \
	openbsd/arm64 \
	openbsd/ppc64 \
	openbsd/riscv64 \
	plan9/386 \
	plan9/amd64 \
	plan9/arm \
	solaris/amd64 \
	windows/386 \
	windows/amd64 \
	windows/arm64
GOFLAGS ?=
LDFLAGS ?= -s -w

.PHONY: build clean install uninstall install-global uninstall-global

build: clean
	@mkdir -p "$(DIST_DIR)"
	@set -e; \
	for platform in $(PLATFORMS); do \
		goos=$${platform%/*}; \
		goarch=$${platform#*/}; \
		target="$(DIST_DIR)/$(APP).$$(printf '%s' "$$platform" | tr / -)"; \
		echo "build $$target"; \
		CGO_ENABLED=0 GOOS=$$goos GOARCH=$$goarch go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o "$$target" "$(CMD)"; \
	done

clean:
	@rm -rf "$(DIST_DIR)"

install:
	@mkdir -p "$(LOCAL_BIN)"
	@echo "install $(LOCAL_BIN)/$(APP)"
	@CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o "$(LOCAL_BIN)/$(APP)" "$(CMD)"

uninstall:
	@echo "uninstall $(LOCAL_BIN)/$(APP)"
	@rm -f "$(LOCAL_BIN)/$(APP)"

install-global:
	@mkdir -p "$(GLOBAL_BIN)"
	@echo "install $(GLOBAL_BIN)/$(APP)"
	@CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o "$(GLOBAL_BIN)/$(APP)" "$(CMD)"

uninstall-global:
	@echo "uninstall $(GLOBAL_BIN)/$(APP)"
	@rm -f "$(GLOBAL_BIN)/$(APP)"
