# Copyright IBM Corporation. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
#  Author: Sandeep Pulluru <sandeep.pulluru@ibm.com>

echo "    _       _           _                    _   _"
echo "   / \   __| |_ __ ___ (_)_ __     __ _  ___| |_(_) ___  _ __  ___"
echo "  / _ \ / _\` | '_ \` _ \| | '_ \\   / _\` |/ __| __| |/ _ \| '_ \/ __|"
echo " / ___ \ (_| | | | | | | | | | | | (_| | (__| |_| | (_) | | | \__ \\"
echo "/_/   \_\__,_|_| |_| |_|_|_| |_|  \\__,_|\\___|\\__|_|\\___/|_| |_|___/"
CONTAINERS=$(docker ps | grep "hyperledger/fabric" | wc -l | tr -d '[:space:]')

if [ $CONTAINERS -eq 7 ]; then
	printf "\n\n##### All containers are up & running, Ready to go ... ######\n\n"
else
	printf "\n\n!!!!!!! Network doesn't seem to be available !!!!!!!\n\n"
	exit
fi

cd ../../fabric-helper

CHANNEL_NAME=defaultchannel
CC_NAME="axispoint-cc"
CC_SRC_DIR="axispoint-cc"

# C R E A T E   C H A N N E L
printf "\n\n============ C R E A T E   C H A N N E L ============\n"
NODE_ENV=local node fabric-cli.js channel create --channel-name $CHANNEL_NAME --org org1
sleep 10

# J O I N  C H A N N E L -  on all peers
printf "\n\n============ J O I N   C H A N N E L ============\n"
#Join peer0 org1
NODE_ENV=local node fabric-cli.js channel join --channel-name $CHANNEL_NAME --org org1
#Join peer0 org2
NODE_ENV=local node fabric-cli.js channel join --channel-name $CHANNEL_NAME --org org2

# I N S T A L L   C H A I N C O D E -  on all peers
printf "\n\n============ I N S T A L L    C H A I N C O D E -  on all peers ============\n"
# Install chaincode org1
NODE_ENV=local node fabric-cli.js chaincode install --src-dir ${CC_SRC_DIR} --org org1 --cc-version V1 --channel $CHANNEL_NAME --cc-name $CC_NAME
#Install chaincode org2
NODE_ENV=local node fabric-cli.js chaincode install --src-dir ${CC_SRC_DIR} --org org2 --cc-version V1 --channel $CHANNEL_NAME --cc-name $CC_NAME

# I N S T A N T I A T E   C H A I N C O D E
# Instantiating chaincode on peer0 of org1
printf "\n\n============ I N S T A N T I A T E    C H A I N C O D E ============\n"
NODE_ENV=local node fabric-cli.js chaincode instantiate --org org1 --cc-version V1 --channel $CHANNEL_NAME --cc-name $CC_NAME --init-arg ''
sleep 10
