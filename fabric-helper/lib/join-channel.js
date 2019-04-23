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
const util = require('util');
const config = require('../../network/general-config.json');
const helper = require('./helper.js');
const logger = helper.getLogger('Join-Channel');

var joinChannel = async function (channelName, orgName) {
	logger.debug('\n============ Join Channel ============\n')
	logger.info(util.format(
		'Calling peers in organization "%s" to join the channel', orgName));

	let client = await helper.getAdminClientForOrg(orgName);
	let channel = client.getChannel(channelName);

	let request = {
		txId: client.newTransactionID(true)
	};

	let genesis_block = await channel.getGenesisBlock(request);

	request = {
		targets: client.getPeersForOrg(client.getMspid()),
		txId: client.newTransactionID(true),
		block: genesis_block
	};

	let eventhubs = client.getEventHubsForOrg(orgName);
	var eventPromises = [];

	eventhubs.forEach((eh) => {
		eh.connect();
		let txPromise = new Promise((resolve, reject) => {
			let handle = setTimeout(reject, config.eventWaitTime);
			eh.registerBlockEvent((block) => {
				clearTimeout(handle);
				if (eh && eh.isconnected()) {
					eh.disconnect();
				}
				// in real-world situations, a peer may have more than one channels so
				// we must check that this block came from the channel we asked the peer to join
				if (block.data.data.length === 1) {
					// Config block must only contain one transaction
					var channel_header = block.data.data[0].payload.header.channel_header;
					if (channel_header.channel_id === channelName) {
						resolve();
					} else {
						reject();
					}
				}
			});
		});
		eventPromises.push(txPromise);
	});
	let sendPromise = channel.joinChannel(request);
	let results = await Promise.all([sendPromise].concat(eventPromises));

	if (results[0] && results[0][0] && results[0][0].response && results[0][0]
		.response.status == 200) {
		logger.info(util.format(
			'Successfully joined peers in organization %s to the channel \'%s\'',
			orgName, channelName));
		let response = {
			success: true,
			message: util.format(
				'Successfully joined peers in organization %s to the channel \'%s\'',
				orgName, channelName)
		};
		return response;
	} else {
		logger.error(' Failed to join channel');
		closeConnections();
		throw new Error('Failed to join channel');
	}
};

exports.joinChannel = joinChannel;