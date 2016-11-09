#!/usr/bin/env zsh

# programs & scripts
wffiles=(transmit OpenFavourite.js)
# icons
wffiles+=(icon.png update.png)
# data files
wffiles+=(info.plist README.md LICENCE.txt)

here="$( cd "$( dirname "$0" )"; pwd )"

log() {
    echo "$@" > /dev/stderr
}

do_dist=0
do_clean=0

usage() {
    cat <<EOS
build-workflow.zsh [options]

Usage:
    build-workflow.zsh [-d] [-c]
    build-workflow.zsh -h

Options:
    -d      Build .alfredworkflow file
    -c      Remove build directory
    -h      Show this help message and exit
EOS
}

while getopts ":cdh" opt; do
  case $opt in
    c)
      do_clean=1
      ;;
    d)
      do_dist=1
      ;;
    h)
      usage
      exit 0
      ;;
    \?)
      log "Invalid option: -$OPTARG"
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))

cleanup() {
    local p="${here}/build"
    log "Cleaning up ..."
    test -d "$p" && rm -rf ./build/
}

pushd "$here" &> /dev/null

log "Building executable(s) ..."
go build -v -o ./transmit ./cmd/
ST_BUILD=$?
if [ "$ST_BUILD" != 0 ]; then
    log "Error building executable."
    cleanup
    popd &> /dev/null
    exit $ST_BUILD
fi

chmod 755 ./transmit

log

log "Cleaning ./build ..."
rm -rvf ./build

log

log "Copying assets to ./build ..."

mkdir -vp ./build

for n in $wffiles; do
    cp -v "$n" ./build/
done

log

if [[ $do_dist -eq 1 ]]; then
    # Get the dist filename from the executable
    zipfile="$(./transmit --distname 2> /dev/null)"
    if [ "$?" -ne 0 ]; then
        log "Error getting distname from transmit."
        cleanup
        exit 1
    fi

    log

    if test -e "$zipfile"; then
        log "Removing existing .alfredworkflow file ..."
        rm -rvf "$zipfile"
        log
    fi

    pushd ./build/ &> /dev/null

    log "Building .alfredworkflow file ..."
    zip -r8n .png "../${zipfile}" ./*
    ST_ZIP=$?
    if [ "$ST_ZIP" != 0 ]; then
        log "Error creating .alfredworkflow file."
        popd &> /dev/null
        cleanup
        popd &> /dev/null
        exit $ST_ZIP
    fi
    popd &> /dev/null

    log
    log "Wrote '${zipfile}' in '$( pwd )'"
fi

[[ $do_clean -eq 1 ]] && { cleanup }

popd &> /dev/null
