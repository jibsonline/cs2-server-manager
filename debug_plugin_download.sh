#!/usr/bin/env bash

echo "======================================"
echo "  Plugin Download Debug Script"
echo "======================================"
echo

echo "[1/5] Basic environment check..."
echo "  Working directory: $(pwd)"
echo "  Shell: $SHELL"
echo "  Bash version: $BASH_VERSION"
echo "  User: $(whoami)"
echo

echo "[2/5] Checking update.sh..."
if [[ -f ./scripts/update.sh ]]; then
  echo "  ✓ Script exists: ./scripts/update.sh"
  echo "  Permissions: $(ls -l ./scripts/update.sh | awk '{print $1}')"
  echo "  First line (shebang): $(head -n1 ./scripts/update.sh)"
else
  echo "  ✗ Script not found: ./scripts/update.sh"
  exit 1
fi
echo

echo "[3/5] Testing script syntax..."
if bash -n ./scripts/update.sh 2>&1; then
  echo "  ✓ No syntax errors detected"
else
  echo "  ✗ Syntax errors found!"
  exit 1
fi
echo

echo "[4/5] Checking dependencies..."
for cmd in curl jq unzip tar rsync; do
  if command -v "$cmd" >/dev/null 2>&1; then
    version=$($cmd --version 2>&1 | head -n1 || echo "version unavailable")
    echo "  ✓ $cmd: $version"
  else
    echo "  ✗ $cmd: NOT FOUND"
  fi
done
echo

echo "[5/5] Running update.sh with debug output..."
echo "----------------------------------------"
set -x
bash -x ./scripts/update.sh plugins 2>&1 | head -100
exit_code=${PIPESTATUS[0]}
set +x
echo "----------------------------------------"
echo
echo "Exit code: $exit_code"
echo

if [[ $exit_code -eq 0 ]]; then
  echo "✓ Plugin download completed successfully!"
else
  echo "✗ Plugin download failed with exit code $exit_code"
  echo
  echo "Please share the output above for debugging."
fi

