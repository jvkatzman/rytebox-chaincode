/**
 * Copyright 2018 IBM Corporation. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an 'AS IS' BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * 
 * Author: Sandeep Pulluru <sandeep.pulluru@ibm.com>
 */

'use strict';
const path = require('path');
const fs = require('fs');
const util = require('util');
const helper = require('./helper.js');
const config = require('../../network/general-config.json');
const logger = helper.getLogger('install-chaincode');

var instantiateChaincode = async function (channelName,
	chaincodeName, chaincodeVersion, args,
	orgName, upgrade) {
	logger.debug('\n============ Instantiate chaincode on organization ' + orgName + ' ============\n');

	try {
		let client = await helper.getAdminClientForOrg(orgName);
		let channel = client.getChannel(channelName);
		let tx_id = client.newTransactionID(true);
		let results = null;

		let request = {
			chaincodeId: chaincodeName,
			chaincodeVersion: chaincodeVersion,
			args: args,
			txId: tx_id
		};
		if (upgrade) {
			results = await channel.sendUpgradeProposal(request);
		} else {
			results = await channel.sendInstantiateProposal(request);
		}

		helper.inspectProposalResult(results);

		let eventhubs = client.getEventHubsForOrg(client.getMspid());
		let deployId = tx_id.getTransactionID();
		let eventPromises = [];

		eventhubs.forEach((eh) => {
			eh.connect();
			let txPromise = new Promise((resolve, reject) => {
				let handle = setTimeout(function () {
					eh.disconnect();
					reject();
				}, parseInt(config.eventWaitTime));
				eh.registerTxEvent(deployId, function (tx, code) {
					logger.info(
						'The transaction has been committed on peer ' +
						eh._ep._endpoint.addr);
					clearTimeout(handle);
					eh.unregisterTxEvent(deployId);
					eh.disconnect();

					if (code !== 'VALID') {
						logger.error('The transaction was invalid, code = ' + code);
						reject();
					} else {
						logger.info('The chaincode instantiate transaction was valid.');
						resolve();
					}
				});
			});
			eventPromises.push(txPromise);
		});

		request = {
			txId: tx_id,
			proposalResponses: results[0],
			proposal: results[1]
		};
		let sendPromise = channel.sendTransaction(request);
		let response = await Promise.all([sendPromise].concat(eventPromises));

		if (response[0].status === 'SUCCESS') {
			let message = util.format('Successfully sent transaction to the orderer.  Chaincode Instantiation is SUCCESS');
			logger.info(message);
			let response = {
				success: true,
				message: message
			};
		} else {
			logger.error('Failed to order the transaction. Error code: ' + response[0].status);
			throw new Error('Failed to order the transaction. Error code: ' + response[0].status);
		}
	} catch (err) {
		logger.error('Failed to instantiate chaincode on the channel: ' + err.stack ? err.stack : err);
		throw new Error('Failed to instantiate chaincode on the channel: ' + err.toString());
	}
};

exports.instantiateChaincode = instantiateChaincode;