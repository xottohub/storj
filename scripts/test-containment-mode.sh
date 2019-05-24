#!/bin/bash
set -ueo pipefail
set +x

STORJ_SIM_CMD="$1"
STORJ_SIM_CMD_OPTS="-x --host 127.0.0.1 --storage-nodes 10"
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export STORJ_NETWORK_DIR="$SCRIPTDIR/../test-storj-sim-configs"

scenario_01_upload() {
  uplink --config-dir "$STORJ_NETWORK_DIR/uplink/0" mb sj://bucket01
  uplink --config-dir "$STORJ_NETWORK_DIR/uplink/0" cp "$SCRIPTDIR/../go.mod" sj://bucket01
}

scenario_01_download() {
  uplink --config-dir "$STORJ_NETWORK_DIR/uplink/0" cp  sj://bucket01/go.mod "$1"
}

scenario_01_kill_nodes() {
  procs=$(ps -A | grep storagenode | awk '{print $1}')
  procs=($procs)

  num_nodes=3
  if [ $# > 1 ]; then
    num_nodes=$1
  fi

  kill -s kill "${procs[@]:2:$1}"
}

if [ ! -d "$STORJ_NETWORK_DIR" ]; then
  storj-sim $STORJ_SIM_CMD_OPTS network setup
fi

case $STORJ_SIM_CMD in
  sim-setup)
    # the network is setup above if the directory doesn't exist
    ;;
  sim-run)
    storj-sim $STORJ_SIM_CMD_OPTS network run
    ;;
  sim-exec)
    storj-sim $STORJ_SIM_CMD_OPTS "${@:2}"
    ;;
  uplink-setup)
    if [ -d "$STORJ_NETWORK_DIR/uplink/0" ]; then
      rm -rf "$STORJ_NETWORK_DIR/uplink/0"
    fi

    mkdir -p "$STORJ_NETWORK_DIR/uplink/0"
    echo "test" > "$STORJ_NETWORK_DIR/uplink/0/.encryption.key"
    uplink --config-dir "$STORJ_NETWORK_DIR/uplink/0" setup --non-interactive \
      --api-key=$(storj-sim $STORJ_SIM_CMD_OPTS network env GATEWAY_0_API_KEY) \
      --enc.key-filepath "$STORJ_NETWORK_DIR/uplink/0/.encryption.key"
    ;;
  uplink-exec)
    uplink --config-dir "$STORJ_NETWORK_DIR/uplink/0" "${@:2}"
    ;;
  scenario-01-upload)
    scenario_01_upload
    ;;
  scenario-01-kill-nodes)
    scenario_01_kill_nodes "${@:2}"
    ;;
  scenario-01-download)
    scenario_01_download "$2"
    ;;
  clean)
    storj-sim $STORJ_SIM_CMD_OPTS network destroy
    ;;
  *)
    echo "Usage: $0 {sim-setup|sim-run|sim-exec|uplink-setup|uplink-exec|clean|scenario....}"
    exit 1
    ;;
esac
