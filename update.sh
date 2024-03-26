#!/usr/bin/env bash

curl() {
  $(type -P curl) -L -q --retry 5 --retry-delay 10 --retry-max-time 60 "$@"
}

get_version() {
  # Get myclipboard release version number
  TMP_FILE="$(mktemp)"
  if ! curl -x "${PROXY}" -sS -i -H "Accept: application/vnd.github.v3+json" -o "$TMP_FILE" 'https://api.github.com/repos/Ch1cC/myclipboard/releases/latest'; then
    "rm" "$TMP_FILE"
    echo 'error: Failed to get release list, please check your network.'
    exit 1
  fi
  HTTP_STATUS_CODE=$(awk 'NR==1 {print $2}' "$TMP_FILE")
  if [[ $HTTP_STATUS_CODE -lt 200 ]] || [[ $HTTP_STATUS_CODE -gt 299 ]]; then
    "rm" "$TMP_FILE"
    echo "error: Failed to get release list, GitHub API response code: $HTTP_STATUS_CODE"
    exit 1
  fi
  RELEASE_LATEST="$(sed 'y/,/\n/' "$TMP_FILE" | grep 'tag_name' | awk -F '"' '{print $4}')"
  "rm" "$TMP_FILE"
  RELEASE_VERSION="${RELEASE_LATEST#v}"
}

download_myclipboard() {
  DOWNLOAD_LINK="https://github.com/Ch1cC/myclipboard/releases/download/$RELEASE_VERSION/myclipboard_${RELEASE_VERSION}_linux_amd64.zip"
  echo "Downloading myclipboard archive: $DOWNLOAD_LINK"
  if ! curl -x "${PROXY}" -R -H 'Cache-Control: no-cache' -o "$ZIP_FILE" "$DOWNLOAD_LINK"; then
    echo 'error: Download failed! Please check your network or try again.'
    return 1
  fi
}

decompression() {
  if ! unzip -q "$1" -d "$TMP_DIRECTORY"; then
    echo 'error: myclipboard decompression failed.'
    "rm" -r "$TMP_DIRECTORY"
    echo "removed: $TMP_DIRECTORY"
    exit 1
  fi
  echo "info: Extract the myclipboard package to $TMP_DIRECTORY and prepare it for installation."
}

install_file() {
  NAME="$1"
  install -m 755 "${TMP_DIRECTORY}/$NAME" "/usr/local/etc/myclipboard/$NAME"
  cp -r "${TMP_DIRECTORY}/static/" "/usr/local/etc/myclipboard/"
}

install_myclipboard() {
  # Install myclipboard binary to /usr/local/bin/ and $DAT_PATH
  install_file myclipboard
}

start_myclipboard() {
  if [[ -f '/etc/systemd/system/myclipboard.service' ]]; then
    if systemctl start "${myclipboard_CUSTOMIZE:-myclipboard}"; then
      echo 'info: Start the myclipboard service.'
    else
      echo 'error: Failed to start myclipboard service.'
      exit 1
    fi
  fi
}

stop_myclipboard() {
  myclipboard_CUSTOMIZE="$(systemctl list-units | grep 'myclipboard' | awk -F ' ' '{print $1}')"
  if [[ -z "$myclipboard_CUSTOMIZE" ]]; then
    local myclipboard_daemon_to_stop='myclipboard.service'
  else
    local myclipboard_daemon_to_stop="$myclipboard_CUSTOMIZE"
  fi
  if ! systemctl stop "$myclipboard_daemon_to_stop"; then
    echo 'error: Stopping the myclipboard service failed.'
    exit 1
  fi
  echo 'info: Stop the myclipboard service.'
}
install_software() {
  package_name="$1"
  file_to_detect="$2"
  type -P "$file_to_detect" >/dev/null 2>&1 && return
  if ${PACKAGE_MANAGEMENT_INSTALL} "$package_name"; then
    echo "info: $package_name is installed."
  else
    echo "error: Installation of $package_name failed, please check your network."
    exit 1
  fi
}
main() {
  # Two very important variables
  TMP_DIRECTORY="$(mktemp -d)"
  ZIP_FILE="${TMP_DIRECTORY}/myclipboard_${RELEASE_VERSION}_linux_x86_64.zip"

  # Normal way
  install_software 'curl' 'curl'
  get_version
  NUMBER="$?"
  if [[ "$NUMBER" -eq '0' ]] || [[ "$FORCE" -eq '1' ]] || [[ "$NUMBER" -eq 2 ]]; then
    echo "info: Installing myclipboard $RELEASE_VERSION for $(uname -m)"
    download_myclipboard
    if [[ "$?" -eq '1' ]]; then
      "rm" -r "$TMP_DIRECTORY"
      echo "removed: $TMP_DIRECTORY"
      exit 1
    fi
    install_software 'unzip' 'unzip'
    decompression "$ZIP_FILE"
  elif [[ "$NUMBER" -eq '1' ]]; then
    echo "info: No new version. The current version of myclipboard is $CURRENT_VERSION ."
    exit 0
  fi

  # Determine if myclipboard is running
  if systemctl list-unit-files | grep -qw 'myclipboard'; then
    if [[ -n "$(pidof myclipboard)" ]]; then
      stop_myclipboard
      myclipboard_RUNNING='1'
    fi
  fi
  install_myclipboard
  "rm" -r "$TMP_DIRECTORY"
  echo "removed: $TMP_DIRECTORY"
  echo "info: myclipboard $RELEASE_VERSION is installed."
  echo "You may need to execute a command to remove dependent software: $PACKAGE_MANAGEMENT_REMOVE curl unzip"
  if [[ "$myclipboard_RUNNING" -eq '1' ]]; then
    start_myclipboard
  else
    echo 'Please execute the command: systemctl enable myclipboard; systemctl start myclipboard'
  fi
}

main "$@"
