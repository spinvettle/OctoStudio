#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-up}"
VERSION="${2:-}"

MIGRATIONS_DIR="${MIGRATIONS_DIR:-./migrations}"
DB_DSN="${DB_DSN:-postgres://baishan:baishan@localhost:5432/learn?sslmode=disable}"
USE_DOCKER="${USE_DOCKER:-}"
DOCKER_NETWORK="${DOCKER_NETWORK:-host}"

run_migrate() {
  if [ -n "$USE_DOCKER" ]; then
    docker run --rm \
      --network "$DOCKER_NETWORK" \
      -v "$(pwd)/${MIGRATIONS_DIR#./}:/migrations" \
      migrate/migrate \
      -path /migrations \
      -database "$DB_DSN" \
      "$@"
    return
  fi

  migrate -path "$MIGRATIONS_DIR" -database "$DB_DSN" "$@"
}

case "$ACTION" in
  up)
    run_migrate up
    ;;

  down)
    run_migrate down "${VERSION:-1}"
    ;;

  force)
    if [ -z "$VERSION" ]; then
      echo "usage: $0 force <version>"
      exit 1
    fi
    run_migrate force "$VERSION"
    ;;

  version)
    run_migrate version
    ;;

  create)
    if [ -z "$VERSION" ]; then
      echo "usage: $0 create <name>"
      exit 1
    fi
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$VERSION"
    ;;

  *)
    echo "usage: $0 {up|down|force|version|create} [arg]"
    exit 1
    ;;
esac
