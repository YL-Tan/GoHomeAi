#!/usr/bin/env bash
# -----------------------------------------------------------------------------
# Bootstrap script for the **ai-smarthome** Python development environment.
#   • Creates or updates a reproducible micromamba env under $MAMBA_ROOT_PREFIX
#   • Installs core tooling (formatters, linters, Jupyter) and ML libraries
#   • Installs git pre-commit hooks
# -----------------------------------------------------------------------------
set -euo pipefail

###############################################################################
#  Editable settings                                                           #
###############################################################################
ENV_NAME="ai-smarthome"               # micromamba environment name
PY_VERSION="3.11"                    # pinned Python major.minor
ROOT_PREFIX="${HOME}/.local/share/mamba"  # where micromamba stores envs
MAMBA_BIN="${ROOT_PREFIX}/bin/micromamba" # expected micromamba path
###############################################################################

# ---------- helpers ----------------------------------------------------------
log()   { printf "\e[1;32m[+]\e[0m %s\n" "$*"; }
error() { printf "\e[1;31m[!]\e[0m %s\n" "$*" >&2; }

mm() { "${MAMBA_BIN}" --root-prefix "${ROOT_PREFIX}" "$@"; }

# -------- 0. Ensure micromamba exists ---------------------------------------
if ! command -v "${MAMBA_BIN}" >/dev/null 2>&1; then
  log "micromamba not found → downloading the latest release …"

  # infer platform string used by the official binaries
  case "$(uname -s)-$(uname -m)" in
    Linux-x86_64)   PLATFORM="linux-64"   ;;
    Linux-aarch64)  PLATFORM="linux-aarch64" ;;
    Darwin-arm64)   PLATFORM="osx-arm64" ;;
    Darwin-x86_64)  PLATFORM="osx-64"     ;;
    *) error "Unsupported platform"; exit 1 ;;
  esac

  mkdir -p "${ROOT_PREFIX}/bin"
  pushd "${ROOT_PREFIX}/bin" >/dev/null
    curl -sL "https://micro.mamba.pm/api/micromamba/${PLATFORM}/latest" \
      | tar -xvjf - micromamba
    chmod +x micromamba
  popd >/dev/null
fi

# -------- 1. build a deterministic spec (temp YAML) -------------------------------
SPEC=$(mktemp)
cat > "${SPEC}" <<EOF
name: ${ENV_NAME}
channels:
  - pytorch
  - conda-forge
dependencies:
  - python=${PY_VERSION}
  - pip
  - ipykernel
  - jupyterlab
  - black
  - isort
  - flake8
  - pre-commit
  - numpy
  - pandas
  - scipy
  - scikit-learn
  - matplotlib
  - seaborn
  - pytorch
  - torchvision
  - torchaudio
  - cpuonly
EOF

# -------- 2. create or update the environment ---------------------------------
if mm env list | grep -q " ${ENV_NAME} "; then
  log "Environment ${ENV_NAME} exists → updating to match spec …"
  mm env update -n "${ENV_NAME}" -f "${SPEC}"
else
  log "Creating environment ${ENV_NAME} …"
  mm env create -n "${ENV_NAME}" -f "${SPEC}"
fi
rm -f "${SPEC}"

# ---------- 3. initialise git hooks -----------------------------------------
if [ -d .git ]; then
  log "Installing pre-commit hooks"
  mm run -n "${ENV_NAME}" pre-commit install
fi

# ---------- 4. export lockfile & clean caches --------------------------------
mm env export -n "${ENV_NAME}" --no-builds > environment.lock.yml
mm clean -a -y >/dev/null  # free disk space

log "✓ Environment ${ENV_NAME} is ready."
printf "   To activate it:\n      micromamba activate %s\n" "${ENV_NAME}"
