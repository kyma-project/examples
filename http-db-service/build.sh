#!/usr/bin/env bash

# docker image name.
PACKAGE_NAME=http-db-service

set -e
set -u
set -o pipefail

# Reset in case getopts.
OPTIND=1

VERBOSE=true
HELP=false
APP_VERSION=latest

declare DOCKER_SYSTEM_NAME
declare IMAGE_NAME_EXAMPLE=${PACKAGE_NAME}:${APP_VERSION}

GREEN='\033[0;32m'
NC='\033[0m' # No Color

function usage() {
    echo "Build your example/http-db-service docker image."
    echo "Usage:"
    echo "   $0 [-h] [-q] [-t]"
    echo
    echo "Flags:"
    echo "   -h        help, print usage"
    echo "   -q        quiet, to turn off VERBOSE"
    echo "   -t        tag, tag the image e.g:registry-name/image-name:tag-version, default:http-db-service:latest"
    echo "Example:"
    echo "   ./build.sh"
    echo "   ./build.sh -t registry:5000/example/http-db-service:0.0.1"
    echo "   ./build.sh -h"
    echo "   ./build.sh -q"
}

function debug() {
    echo "Config:"
    echo "Flags:"
    echo "- verbose                     : ${VERBOSE}"
    echo "result:"
    echo "- EXAMPLE IMAGE_NAME          : ${IMAGE_NAME_EXAMPLE}"
    echo "- DOCKER_SYSTEM_NAME          : ${DOCKER_SYSTEM_NAME}"
}


function setParameter() {
    if [[ $# -ne 1 ]]; then
        echo "No parameters are passed to getParameter"
        exit 1
    else
        IMAGE_NAME_EXAMPLE=$1
    fi
}

function onExit {
    local EXIT_CODE=$?
    if [[ ${VERBOSE} == "true" ]] && [[ ${HELP} == "false" ]]; then
        echo -e "${GREEN}Debug on exit${NC}"
        debug
        echo "exit code: ${EXIT_CODE}"
    fi
}
trap onExit EXIT


OPTIND=1
while getopts "ht:q" opt
 do
    case ${opt} in
        h)
            HELP=true
            usage
            exit 0
        ;;
        q)
            set +e
            VERBOSE=false
            set -e
        ;;
        t)
            set +e
            setParameter "$OPTARG"
            set -e
        ;;
        \?)
            echo "Unknown option: -$OPTARG" >&2;
            exit 1
        ;;
    esac
done
shift $((OPTIND-1))

## inspect configured docker env
DOCKER_SYSTEM_NAME=$(docker info -f "{{ .Name }}")

## perform image build
if [[ ${VERBOSE} == "true" ]]; then echo -e "${GREEN}Build Example binary locally${NC}"; fi

if [[ ${VERBOSE} == "true" ]]; then
    echo -e "${GREEN}Build docker image ${IMAGE_NAME_EXAMPLE} ${NC}"
    docker build -t "${IMAGE_NAME_EXAMPLE}" -f Dockerfile .
    echo -e "${GREEN}Built Docker image successfully...${NC}"
else
    docker build -t "${IMAGE_NAME_EXAMPLE}" -f Dockerfile . > /dev/null
fi

