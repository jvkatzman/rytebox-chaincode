#!/bin/bash -e
#
#  Copyright 2018 IBM Corporation. All Rights Reserved.
#
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an 'AS IS' BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
#  Author: Sandeep Pulluru <sandeep.pulluru@ibm.com>

## idle timeout for Ordering service availability
TIMEOUT=60

##TODO: add a flag to enable/disable debug for peer & orderer

## usage message
function usage() {
	echo "Usage: "
	echo "  bootstrap.sh [-m start|stop|restart] [-t <release-tag>] [-l capture-logs]"
	echo "  bootstrap.sh -h|--help (print this message)"
	echo "      -m <mode> - one of 'start', 'stop', 'restart' " #or 'generate'"
	echo "      - 'start' - bring up the network with docker-compose up & start the app on port 4000"
	echo "      - 'up'    - same as start"
	echo "      - 'stop'  - stop the network with docker-compose down & clear containers , crypto keys etc.,"
	echo "      - 'down'  - same as stop"
	echo "      - 'restart' -  restarts the network and start the app on port 4000 (Typically stop + start)"
	echo "     -a ALL IN ONE 1) Launch network 2) perform admin actions & 3) Start the app "
	echo "     -r re-Generate the certs and channel Artifacts, (** Not Recommended)"
	echo "     -l capture docker logs before network teardown"
	echo
	echo "Some possible options:"
	echo
	echo "	bootstrap.sh"
	echo "	bootstrap.sh -l"
	echo "	bootstrap.sh -r"
	echo "	bootstrap.sh -m stop"
	echo "	bootstrap.sh -m stop -l"
	echo
	echo "All defaults:"
	echo "	bootstrap.sh"
	echo "	Restarts the network and uses latest docker images instead of specific TAG "
	exit 1
}

# Parse commandline args
while getopts "h?m:t:lra" opt; do
	case "$opt" in
	h | \?)
		usage
		;;
	m) MODE=$OPTARG ;;

	l) ENABLE_LOGS='y' ;;

	r) REGENERATE='y' ;;

	a) ALL='y' ;;

	esac
done
echo " _                           _        _   _      _                      _"
echo "| |    __ _ _   _ _ __   ___| |__    | \ | | ___| |___      _____  _ __| | __"
echo "| |   / _\` | | | | '_ \ / __| '_ \   |  \| |/ _ \ __\ \ /\ / / _ \| '__| |/ /"
echo "| |__| (_| | |_| | | | | (__| | | |  | |\  |  __/ |_ \ V  V / (_) | |  |   <"
echo "|_____\__,_|\__,_|_| |_|\___|_| |_|  |_| \_|\___|\__| \_/\_/ \___/|_|  |_|\_\\"
echo ""
: ${MODE:="restart"}
: ${IMAGE_TAG:="1.2.0"}
: ${ENABLE_LOGS:="n"}
: ${ALL:="n"}
: ${THIRDPARTY_IMAGE_TAG:="0.4.10"}
export ARCH=$(echo "$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
export IMAGE_TAG
export THIRDPARTY_IMAGE_TAG

COMPOSE_FILE=./docker-compose.yaml
COMPOSE_FILE_WITH_COUCH=./docker-compose-couch.yaml
COMPOSE_FILE_KEYVAL=./docker-compose-keyval.yaml
function dkcl() {
	CONTAINERS=$(docker ps -a --filter network=local_default | wc -l)
	if [ "$CONTAINERS" -gt "1" ]; then
		docker rm -f $(docker ps -aq --filter network=local_default)
	else
		printf "\n========== No containers available for deletion ==========\n"
	fi
}

function dkrm() {
	DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
	echo
	if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" = " " ]; then
		echo "========== No images available for dockdeletion ==========="
	else
		docker rmi -f $DOCKER_IMAGE_IDS
	fi
	echo
}

# delete all node modules and re-install
function cleanAndInstall() {
	## Make sure cleanup the node_moudles and re-install them again
	rm -rf ./node_modules

	printf "\n============== Installing node modules =============\n"
	npm install --no-audit
}

# Install required node modules if not installed already
function checkNodeModules() {
	echo
	if [ -d node_modules ]; then
		npm ls fabric-client fabric-ca-client || cleanAndInstall
	else
		cleanAndInstall && npm ls fabric-client fabric-ca-client
	fi
	echo
}

# Check if all the dockerimages are available or not
function checkForDockerImages() {
	printf "\n\nDocker Images : hyperledger/fabric-$IMAGE:$IMAGE_TAG\n\n"

	DOCKER_IMAGES=$(docker images | grep "$IMAGE_TAG\|$THIRDPARTY_IMAGE_TAG" | grep -v "amd" | wc -l)
	if [ $DOCKER_IMAGES -ne 8 ]; then
		printf "\n############# You don't have all fabric images, Let me pull them for you ###########\n"
		printf "######## Pulling Fabric Images ... ########\n"
		for IMAGE in peer orderer ca ccenv tools; do
			docker pull hyperledger/fabric-$IMAGE:$IMAGE_TAG
		done
		printf "######## Pulling 3rdParty Images ... ########\n"
		for IMAGE in couchdb kafka zookeeper; do
			docker pull hyperledger/fabric-$IMAGE:$THIRDPARTY_IMAGE_TAG
		done
	fi
}

## This will either re-generate the artifacts for your network (or)
## Start/Restart the fabric network
## Also, checks if images are available or not and pulls images from dockerhub
## if not available
## ex: ./bootstrap.sh          --> Restarts the network
##		 ./bootstrap.sh -m up    --> Starts the network
##		 ./bootstrap.sh -m up -r --> Regenerate all artifacts & starts the network
function startNetwork() {
	printf "\n ========= FABRIC IMAGE TAG : $IMAGE_TAG ===========\n"
	checkForDockerImages
	### dynamic generation of Org certs and channel network
	if [ "$REGENERATE" = "y" ]; then
		echo "===> Downloading platform binaries"
		rm -rf bin
		curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${IMAGE_TAG}/hyperledger-fabric-${ARCH}-${IMAGE_TAG}.tar.gz | tar xz
		rm -rf ./channel/*.block ./channel/*.tx ./crypto-config hyperledger-fabric-${ARCH}-${IMAGE_TAG}.tar.gz ./config
		source generateArtifacts.sh
		printf "\n\nStarting Network ...\n\n"
	fi

	#Launch the network
	docker-compose -f $COMPOSE_FILE -f $COMPOSE_FILE_WITH_COUCH -f $COMPOSE_FILE_KEYVAL up -d
	if [ $? -ne 0 ]; then
		printf "\n\n!!!!!!!! Unable to pull images/ start the network, Check your docker-compose !!!!!\n\n"
		exit
	fi

	##Install node modules
	checkNodeModules

	CONTAINERS=$(docker ps | grep "hyperledger/fabric" | wc -l | tr -d '[:space:]')

	if [ $CONTAINERS -eq 7 ]; then
		printf "\n\n@@@@@@@@@@@@ YOUR NETWORK IS UP & READY TO USE @@@@@@@@@@@@\n\n"
	else
		printf "\n\n!!!!!!!!!! SOMETHING IS WRONG !!!!!!!!!!\n\n"
		docker ps -a | grep Exited
	fi
}

## Teardown the fabric network and also captures all the logs if flag is enabled
## ex: ./bootstrap.sh -m down      --> teardown the network
##     ./bootstrap.sh -m down -l   --> clear the network & capture all logs
function teardownNetwork() {
	printf "\n======================= TEARDOWN NETWORK ====================\n"
	if [ "$ENABLE_LOGS" = "y" ]; then
		source ./getContainerLogs.sh
	fi
	# teardown the network and clean the containers and intermediate images
	docker-compose -f $COMPOSE_FILE -f $COMPOSE_FILE_WITH_COUCH -f $COMPOSE_FILE_KEYVAL down
	dkcl
	dkrm
	echo y | docker network prune
}

function adminActivities() {
	./admin.sh
}

function bootstrap() {
	startNetwork
	if [ "$ALL" == "y" ]; then
		adminActivities
	fi
}
## Network launch modes
## up (or Start), down (or stop) , restart
case $MODE in
'start' | 'up')
	bootstrap
	;;
'stop' | 'down')
	teardownNetwork
	;;
'restart')
	teardownNetwork
	bootstrap
	;;
*)
	usage
	;;
esac
