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
const helper = require('./helper.js');
const logger = helper.getLogger('query-chaincode');
const config = require('../../network/general-config.json');

var queryChaincode = async function (channelName, chaincodeName, args,
	functionName, userName, orgName) {
	logger.info('========= Query chaincode on org: "%s", channelName: "%s",' +
		'chaincodeName: "%s", userName: "%s", functionName: "%s", args: "%s" ========= ', orgName, channelName, chaincodeName, userName, functionName, args);

	try {
		let client = await helper.getClientForOrg(orgName);
		let user = await helper.getRegisteredUser(userName, orgName);

		let channel = client.getChannel(channelName);

		let tx_id = client.newTransactionID();

		let request = {
			chaincodeId: chaincodeName,
			txId: tx_id,
			fcn: functionName,
			args: args,
			targets: channel.getPeers()
		};

		let results = await channel.queryByChaincode(request);

		if (results) {
			for (let i = 0; i < results.length; i++) {
				logger.info(results[i].toString('utf8'));
				return results[i].toString('utf8');
			}
		} else {
			logger.error('results is null');
			return 'results is null';
		}
	} catch (err) {
		logger.error('Failed query: ' + err.stack ? err.stack : err);
		throw new Error('Failed query: ' + err.toString());
	}
};

exports.queryChaincode = queryChaincode;