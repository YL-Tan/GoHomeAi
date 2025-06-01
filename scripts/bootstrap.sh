#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# Bootstrap script for the **ai-smarthome** Python development environment.
# - Recreates / updates the micromamba env from a SINGLE authoritative file
#   (environment.yml in the repo root)
# - Installs git pre-commit hooks
# -----------------------------------------------------------------------------
set -euo pipefail

###############################################################################
# Editable settings
###############################################################################
ENV_NAME="ai-smarthome"                       # micromamba env name
ROOT_PREFIX="${MAMBA_ROOT_PREFIX:-${HOME}/.local/share/mamba}"
MAMBA_BIN="${ROOT_PREFIX}/bin/micromamba"
ENV_FILE="$(git rev-parse --show-toplevel)/environment.yml"  # absolute path
###############################################################################

# ensure the directory exists _and_ is empty if brand-new
mkdir -p "${ROOT_PREFIX}"

# ---------- helpers ----------------------------------------------------------
log()   { printf "\e[1;32m[+]\e[0m %s\n" "$*"; }
error() { printf "\e[1;31m[!]\e[0m %s\n" "$*" >&2; }

mm() { "${MAMBA_BIN}" --root-prefix "${ROOT_PREFIX}" "$@"; }

# -------- 0. Ensure micromamba exists ---------------------------------------
if [ ! -x "${MAMBA_BIN}" ]; then
  log "micromamba not found → downloading the latest release …"

  case "$(uname -s)-$(uname -m)" in
    Linux-x86_64)   PLATFORM="linux-64"      ;;
    Linux-aarch64)  PLATFORM="linux-aarch64" ;;
    Darwin-arm64)   PLATFORM="osx-arm64"     ;;
    Darwin-x86_64)  PLATFORM="osx-64"        ;;
    *) error "Unsupported platform"; exit 1 ;;
  esac

  mkdir -p "${ROOT_PREFIX}/bin"
  # download & extract micromamba → $ROOT_PREFIX/bin
  TMP_TAR=$(mktemp)
  trap 'rm -f "${TMP_TAR:-}"' EXIT
  curl -sL "https://micro.mamba.pm/api/micromamba/${PLATFORM}/latest" -o "${TMP_TAR}"
  tar --extract --file="${TMP_TAR}" --bzip2 \
      --directory="${ROOT_PREFIX}/bin" --strip-components=1 "bin/micromamba"
  rm -f "${TMP_TAR}"
  chmod +x "${MAMBA_BIN}"
fi

export PATH="${ROOT_PREFIX}/bin:${PATH}"

# Make the micromamba binary discoverable by later CI steps
if [ -n "${GITHUB_PATH:-}" ]; then
  echo "${ROOT_PREFIX}/bin" >> "${GITHUB_PATH}"
fi

# -------- 1. Create / update the environment --------------------------------
if mm env list | grep -q " ${ENV_NAME} "; then
  log "Environment ${ENV_NAME} exists → updating from environment.yml …"
  mm env update -n "${ENV_NAME}" -f "${ENV_FILE}"
else
  log "Creating environment ${ENV_NAME} from environment.yml …"
  mm env create -n "${ENV_NAME}" -f "${ENV_FILE}"
fi

# -------- 2. Install pre-commit hooks ---------------------------------------
if [ -d .git ] && [ -f .pre-commit-config.yaml ]; then
  log "Ensuring pre-commit hooks are installed"
  mm run -n "${ENV_NAME}" pre-commit install
else
  log "No .pre-commit-config.yaml found - skipping hook install"
fi

# -------- 3. Export deterministic lock file & clean caches ------------------
mm env export -n "${ENV_NAME}" --no-builds > environment.lock.yml
mm clean -a -y >/dev/null

log "✓ Environment ${ENV_NAME} is ready."
printf "   Activate with:\n      micromamba activate %s\n" "${ENV_NAME}"