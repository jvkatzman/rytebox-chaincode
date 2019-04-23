# Copyright IBM Corporation. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
#  Author: Sandeep Pulluru <sandeep.pulluru@ibm.com>

CHANNEL_NAME=$1
TOTAL_CHANNELS=$2
: ${CHANNEL_NAME:="defaultchannel"}
## Let us not use more than one channel at the moment
: ${TOTAL_CHANNELS:="1"}
#echo "Using CHANNEL_NAME prefix as $CHANNEL_NAME"
ROOT_DIR=$PWD
export FABRIC_CFG_PATH=$ROOT_DIR
export PATH=$PATH:$ROOT_DIR/bin
ARCH=$(echo "$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')

function generateCerts() {
	pwd
	CRYPTOGEN=cryptogen
	echo
	echo "##########################################################"
	echo "##### Generate certificates using cryptogen tool #########"
	echo "##########################################################"
	$CRYPTOGEN generate --config=$FABRIC_CFG_PATH/cryptogen.yaml
	echo
}

## docker-compose template to replace private key file names with constants
function replacePrivateKey() {
	OPTS="-i"
	if [ $(uname -s) = "Darwin" ]; then
		OPTS="-it"
	fi
	cp docker-compose-template.yaml docker-compose.yaml
	cp network-config/network-config-org1-template.json network-config/network-config-org1.json
	cp network-config/network-config-org2-template.json network-config/network-config-org2.json

	cd crypto-config/peerOrganizations/org1.example.com/ca/
	PRIV_KEY=$(ls *_sk)
	cd $ROOT_DIR
	sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
	cd crypto-config/peerOrganizations/org2.example.com/ca/
	PRIV_KEY=$(ls *_sk)
	cd $ROOT_DIR
	sed $OPTS "s/CA2_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml

	#ORG1
	cd crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore
	PRIV_KEY=$(ls *_sk)
	cd $ROOT_DIR
	sed $OPTS "s/ORG1_ADMIN_KEY/${PRIV_KEY}/g" network-config/network-config-org1.json

	#ORG2
	cd crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/keystore
	PRIV_KEY=$(ls *_sk)
	cd $ROOT_DIR
	sed $OPTS "s/ORG2_ADMIN_KEY/${PRIV_KEY}/g" network-config/network-config-org2.json
	rm -rf network-config/network-config-org[1-2].jsont
	cd $ROOT_DIR
}

## Generate orderer genesis block , channel configuration transaction and anchor peer update transactions
function generateChannelArtifacts() {
	if [ ! -d channel ]; then
		mkdir -p channel
	fi
	CONFIGTXGEN=configtxgen
	echo "##########################################################"
	echo "#########  Generating Orderer Genesis block ##############"
	echo "##########################################################"
	# Note: For some unknown reason (at least for now) the block file can't be
	# named orderer.genesis.block or the orderer will fail to launch!
	$CONFIGTXGEN -profile TwoOrgsOrdererGenesis -channelID testchannelid -outputBlock ./channel/genesis.block

	echo "#################################################################"
	echo "### Generating channel configuration transaction '$CHANNEL_NAME$i.tx' ###"
	echo "#################################################################"
	$CONFIGTXGEN -profile TwoOrgsChannel -outputCreateChannelTx ./channel/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
	##TODO: Enable this when used multiple channels
	# $CONFIGTXGEN -profile TwoOrgsChannel -outputCreateChannelTx ./channel/$CHANNEL_NAME$i.tx -channelID $CHANNEL_NAME$i
	echo
}

generateCerts
replacePrivateKey
generateChannelArtifacts
cd $ROOT_DIR
