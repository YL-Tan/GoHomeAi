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
ROOT_PREFIX="${HOME}/.local/share/mamba"      # where micromamba stores envs
MAMBA_BIN="${ROOT_PREFIX}/bin/micromamba"     # expected micromamba path
ENV_FILE="$(git rev-parse --show-toplevel)/environment.yml"  # absolute path
###############################################################################

# ---------- helpers ----------------------------------------------------------
log()   { printf "\e[1;32m[+]\e[0m %s\n" "$*"; }
error() { printf "\e[1;31m[!]\e[0m %s\n" "$*" >&2; }

mm() { "${MAMBA_BIN}" --root-prefix "${ROOT_PREFIX}" "$@"; }

# -------- 0. Ensure micromamba exists ---------------------------------------
if ! command -v "${MAMBA_BIN}" >/dev/null 2>&1; then
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
  curl -sL "https://micro.mamba.pm/api/micromamba/${PLATFORM}/latest" -o "${TMP_TAR}"
  tar --extract --verbose --file="${TMP_TAR}" --bzip2 \
      --directory="${ROOT_PREFIX}/bin" --strip-components=1 "bin/micromamba"
  rm -f "${TMP_TAR}"
  chmod +x "${MAMBA_BIN}"
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
if [ -d .git ]; then
  log "Ensuring pre-commit hooks are installed"
  mm run -n "${ENV_NAME}" pre-commit install
fi

# -------- 3. Export deterministic lock file & clean caches ------------------
mm env export -n "${ENV_NAME}" --no-builds > environment.lock.yml
mm clean -a -y >/dev/null

log "✓ Environment ${ENV_NAME} is ready."
printf "   Activate with:\n      micromamba activate %s\n" "${ENV_NAME}"