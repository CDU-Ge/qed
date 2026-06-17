#!/usr/bin/env sh
set -eu

REPO="${QED_REPO:-CDU-Ge/qed}"
VERSION="${QED_VERSION:-latest}"
INSTALL_DIR="${QED_INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="${QED_BINARY_NAME:-qed}"
PROGRESS="${QED_PROGRESS:-auto}"

usage() {
	cat <<EOF
Install qed from GitHub Releases.

Usage:
  ./install.sh [options]

Options:
  --version <tag>   Release tag to install, for example v0.1.0.
  --dir <path>      Install directory. Defaults to ~/.local/bin.
  --repo <repo>     GitHub repository. Defaults to CDU-Ge/qed.
  --progress <mode> Progress mode: auto, bar, quiet. Defaults to auto.
  --quiet           Same as --progress quiet.
  -h, --help        Show this help.

Environment:
  QED_VERSION       Same as --version.
  QED_INSTALL_DIR   Same as --dir.
  QED_REPO          Same as --repo.
  QED_BINARY_NAME   Installed command name. Defaults to qed.
  QED_PROGRESS      Same as --progress.
EOF
}

log() {
	[ "$PROGRESS" = "quiet" ] && return 0
	printf '%s\n' "$*" >&2
}

fail() {
	printf 'error: %s\n' "$*" >&2
	exit 1
}

while [ "$#" -gt 0 ]; do
	case "$1" in
		--version)
			[ "$#" -ge 2 ] || fail "--version requires a value"
			VERSION="$2"
			shift 2
			;;
		--dir)
			[ "$#" -ge 2 ] || fail "--dir requires a value"
			INSTALL_DIR="$2"
			shift 2
			;;
		--repo)
			[ "$#" -ge 2 ] || fail "--repo requires a value"
			REPO="$2"
			shift 2
			;;
		--progress)
			[ "$#" -ge 2 ] || fail "--progress requires a value"
			PROGRESS="$2"
			shift 2
			;;
		--quiet)
			PROGRESS="quiet"
			shift
			;;
		-h | --help)
			usage
			exit 0
			;;
		*)
			fail "unknown option: $1"
			;;
	esac
done

case "$PROGRESS" in
	auto | bar | quiet) ;;
	*) fail "unsupported progress mode: $PROGRESS" ;;
esac

detect_os() {
	case "$(uname -s)" in
		Linux) printf 'linux' ;;
		Darwin) printf 'darwin' ;;
		FreeBSD) printf 'freebsd' ;;
		OpenBSD) printf 'openbsd' ;;
		NetBSD) printf 'netbsd' ;;
		DragonFly) printf 'dragonfly' ;;
		SunOS) printf 'solaris' ;;
		AIX) printf 'aix' ;;
		MINGW* | MSYS* | CYGWIN*) printf 'windows' ;;
		*) return 1 ;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
		x86_64 | amd64) printf 'amd64' ;;
		i386 | i486 | i586 | i686) printf '386' ;;
		aarch64 | arm64) printf 'arm64' ;;
		armv5* | armv6* | armv7* | arm) printf 'arm' ;;
		ppc64) printf 'ppc64' ;;
		ppc64le) printf 'ppc64le' ;;
		s390x) printf 's390x' ;;
		riscv64) printf 'riscv64' ;;
		loongarch64) printf 'loong64' ;;
		mips) printf 'mips' ;;
		mipsel) printf 'mipsle' ;;
		mips64) printf 'mips64' ;;
		mips64el) printf 'mips64le' ;;
		*) return 1 ;;
	esac
}

show_progress() {
	case "$PROGRESS" in
		bar)
			return 0
			;;
		auto)
			[ -t 2 ]
			return
			;;
		quiet)
			return 1
			;;
	esac
}

download() {
	url="$1"
	output="$2"

	if command -v curl >/dev/null 2>&1; then
		if show_progress; then
			curl -fL --progress-bar "$url" -o "$output"
		else
			curl -fsSL "$url" -o "$output"
		fi
	elif command -v wget >/dev/null 2>&1; then
		if show_progress; then
			if wget --help 2>&1 | grep -q -- '--show-progress'; then
				wget -q --show-progress -O "$output" "$url"
			else
				wget -O "$output" "$url"
			fi
		else
			wget -qO "$output" "$url"
		fi
	else
		fail "curl or wget is required"
	fi
}

verify_checksum() {
	checksums="$1"
	archive="$2"
	line_file="$3"

	awk -v name="$archive" '$2 == name { print; found = 1 } END { exit found ? 0 : 1 }' "$checksums" >"$line_file" ||
		fail "checksum for $archive not found"

	if command -v sha256sum >/dev/null 2>&1; then
		if [ "$PROGRESS" = "quiet" ]; then
			sha256sum -c "$line_file" >/dev/null
		else
			sha256sum -c "$line_file"
		fi
	elif command -v shasum >/dev/null 2>&1; then
		if [ "$PROGRESS" = "quiet" ]; then
			shasum -a 256 -c "$line_file" >/dev/null
		else
			shasum -a 256 -c "$line_file"
		fi
	else
		log "warning: sha256sum or shasum not found; skipping checksum verification"
	fi
}

install_binary() {
	source="$1"
	target="$2"

	mkdir -p "$INSTALL_DIR"
	if command -v install >/dev/null 2>&1; then
		install -m 0755 "$source" "$target"
	else
		cp "$source" "$target"
		chmod 0755 "$target"
	fi
}

OS="$(detect_os)" || fail "unsupported operating system: $(uname -s)"
ARCH="$(detect_arch)" || fail "unsupported architecture: $(uname -m)"
PLATFORM="$OS-$ARCH"
ARTIFACT="qed.$PLATFORM"

case "$OS" in
	windows)
		ARCHIVE="$ARTIFACT.zip"
		TARGET="$INSTALL_DIR/$BINARY_NAME.exe"
		;;
	*)
		ARCHIVE="$ARTIFACT.tar.gz"
		TARGET="$INSTALL_DIR/$BINARY_NAME"
		;;
esac

if [ "$VERSION" = "latest" ]; then
	DOWNLOAD_BASE="https://github.com/$REPO/releases/latest/download"
else
	DOWNLOAD_BASE="https://github.com/$REPO/releases/download/$VERSION"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

log "downloading $ARCHIVE from $REPO ($VERSION)"
download "$DOWNLOAD_BASE/$ARCHIVE" "$TMP_DIR/$ARCHIVE"

log "downloading checksums.txt"
download "$DOWNLOAD_BASE/checksums.txt" "$TMP_DIR/checksums.txt"

(
	cd "$TMP_DIR"
	verify_checksum "checksums.txt" "$ARCHIVE" "checksum.line"
)

case "$ARCHIVE" in
	*.tar.gz)
		tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"
		;;
	*.zip)
		command -v unzip >/dev/null 2>&1 || fail "unzip is required for Windows archives"
		unzip -q "$TMP_DIR/$ARCHIVE" -d "$TMP_DIR"
		;;
esac

[ -f "$TMP_DIR/$ARTIFACT" ] || fail "archive did not contain $ARTIFACT"

install_binary "$TMP_DIR/$ARTIFACT" "$TARGET"
log "installed $TARGET"

if [ "$PROGRESS" != "quiet" ] && [ -x "$TARGET" ]; then
	"$TARGET" version || true
fi
