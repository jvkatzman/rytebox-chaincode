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
const logger = helper.getLogger('invoke-chaincode');

var invokeTransaction = async function (channelName, chaincodeName, args,
	functionName, userName, orgName) {

	logger.info('========= Invoke chaincode transaction on org: "%s", channelName: "%s",' +
		'chaincodeName: "%s", userName: "%s", functionName: "%s", args: "%s" ========= ', orgName, channelName, chaincodeName, userName, functionName, args);

	let client = await helper.getClientForOrg(orgName);
	let user = await helper.getRegisteredUser(userName, orgName);
	let channel = client.getChannel(channelName);

	let tx_id = client.newTransactionID();
	let tx_id_string = tx_id.getTransactionID();
	let payload = null;

	// send proposal to endorser
	let request = {
		targets: channel.getPeers(),
		chaincodeId: chaincodeName,
		fcn: functionName,
		args: args,
		chainId: channelName,
		txId: tx_id
	};

	return channel.sendTransactionProposal(request, 120000).then((results) => {
		var proposalResponses = results[0];
		var proposal = results[1];
		var header = results[2];
		var all_good = true;

		payload = proposalResponses[0];
		for (var i in proposalResponses) {
			let one_good = false;
			if (proposalResponses && proposalResponses[i].response &&
				proposalResponses[i].response.status === 200) {
				one_good = true;
				logger.info('transaction:%s proposal was GOOD', tx_id_string);
			} else {
				logger.error('transaction:%s proposal was BAD', tx_id_string);
			}
			all_good = all_good & one_good;
		}
		if (all_good) {
			for (var i in proposalResponses) {
				logger.debug(util.format(
					'Successfully sent Proposal for transaction:%s and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s',
					tx_id_string, proposalResponses[i].response.status, proposalResponses[i].response.message,
					proposalResponses[i].response.payload, proposalResponses[i].endorsement
					.signature));
			}
			var request = {
				proposalResponses: proposalResponses,
				proposal: proposal,
				header: header
			};
			// set the transaction listener and set a timeout of 30sec
			// if the transaction did not get committed within the timeout period,
			// fail the test
			var transactionID = tx_id.getTransactionID();
			var sendPromise = channel.sendTransaction(request);

			var eventPromises = [];

			var eventhubs = helper.getEventHubsForOrg(orgName);
			eventhubs.forEach((eh) => {
				logger.debug('invokeEventPromise - setting up event');
				let invokeEventPromise = new Promise((resolve, reject) => {
					let event_timeout = setTimeout(() => {
						let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
						logger.error(message);
						//eh.disconnect();
					}, 600000);
					eh.registerTxEvent(tx_id_string, (tx, code) => {
							logger.debug('The chaincode invoke chaincode transaction has been committed on peer %s', eh.getPeerAddr());
							logger.debug('Transaction:%s has status of %s', tx, code);
							clearTimeout(event_timeout);

							if (code !== 'VALID') {
								let message = util.format('The invoke chaincode transaction:%s was invalid, code:"%s"', tx_id_string, code);
								logger.error(message);
								reject(new Error(message));
							} else {
								let message = util.format('The invoke chaincode transaction:%s was valid.', tx_id_string);
								logger.info(message);
								resolve(message);
							}
						}, (err) => {
							clearTimeout(event_timeout);
							logger.error(err);
							reject(err);
						},
						// the default for 'unregister' is true for transaction listeners
						// so no real need to set here, however for 'disconnect'
						// the default is false as most event hubs are long running
						// in this use case we are using it only once
						{
							unregister: true,
							disconnect: true
						}
					);
					if (!eh.isconnected()) {
						logger.debug("Event hub not connected yet: " + eh.getPeerAddr());
						eh.connect();
					} else {
						logger.debug("Event hub already connected: " + eh.getPeerAddr());
					}
				});
				eventPromises.push(invokeEventPromise);
			});

			return Promise.all([sendPromise].concat(eventPromises)).then((results) => {
				logger.debug('Event promise all complete and testing complete for transaction:%s', tx_id_string);
				return results[0]; // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call
			}).catch((err) => {
				logger.error(
					'Failed to send transaction:%s and get notifications within the timeout period.', tx_id_string
				);
				logger.error(err);
				return util.format('Failed to send transaction%s and get notifications within the timeout period.', tx_id_string);
			});
		} else {
			logger.error(
				'Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...'
			);
			return 'Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...';
		}
	}, (err) => {
		logger.error('Failed to send proposal for transaction:%s due to error: ', tx_id_string, err.stack ? err.stack :
			err);
		return util.format('Failed to send proposal for transaction:%s due to error: ', tx_id_string, err.stack ? err.stack :
			err);
	}).then((response) => {
		if (response.status === 'SUCCESS') {
			logger.info('Successfully sent transaction:%s to the orderer.', tx_id_string);
			return payload.response.payload;
		} else {
			logger.error('Failed to order the transaction:%s. Error code: %s', tx_id_string, response.status);
			throw new Error(payload.message);
		}
	}, (err) => {
		logger.error('Failed to send transaction:%s due to error: %s', tx_id_string, err.stack ? err
			.stack : err);
		throw new Error('Failed to send transaction:%s due to error: %s', tx_id_string, err.stack ? err.stack :
			err);
	});
};

exports.invokeTransaction = invokeTransaction;