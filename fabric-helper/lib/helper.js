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
 * Author: Sandeep Pulluru <sandeep.pulluru@itpeoplecorp.com>
 */

'use strict';
const path = require('path');
const util = require('util');
const fs = require('fs-extra');
const User = require('fabric-client/lib/User.js');
const hfc = require('fabric-client');
const log4js = require('log4js');
const logger = log4js.getLogger('fabric-helper');
const request = require('request');
const utils = require('fabric-client/lib/utils.js');

//Setting default environment type if not mentioned to local
if (!process.env.NODE_ENV) {
	process.env.NODE_ENV = 'local'
}

hfc.addConfigFile(path.join(__dirname, '../../network/general-config.json'));
hfc.setLogger(logger);
logger.setLevel(hfc.getConfigSetting('loglevel'));

let clients = {};
let evenHubs = {};
var connectionProfiles = {};

// Sets up the connection profile for all organisations
module.exports.getConnectionProfile = function (orgName) {
	return new Promise(function (resolve, reject) {
		if (connectionProfiles[orgName] === undefined) {
			if (process.env.NODE_ENV === 'production') {
				request.get({
					url: 'https://blockchain-starter.ng.bluemix.net/api/v1/networks/' + process.env.NETWORK_ID + '/connection_profile',
					headers: {
						'Authorization': 'Basic ' + new Buffer(orgName + ':' +
							process.env[orgName.toUpperCase() + '_SECRET']).toString('base64')
					}
				}, function (err, response, body) {
					if (err) {
						reject(err);
					} else {
						let connectionProfile = JSON.parse(body);
						connectionProfiles[orgName] = connectionProfile;
						resolve(connectionProfile);
					}
				});
			} else {
				let connectionProfile = JSON.parse(fs.readFileSync(path.join(__dirname + '../../../network/local/network-config/network-config-' + orgName + '.json'), 'utf8'));
				connectionProfiles[orgName] = connectionProfile;
				resolve(connectionProfile);
			}
		} else {
			resolve(connectionProfiles[orgName]);
		}
	});
};

module.exports.getClientForOrg = async function (org) {
	if (clients[org] == undefined) {
		let connectionProfile = await this.getConnectionProfile(org);
		//Set configuration files for Org
		hfc.setConfigSetting('request-timeout', hfc.getConfigSetting('eventWaitTime'));
		hfc.setConfigSetting(org + '-network', connectionProfile);
		hfc.setConfigSetting(org, path.join(__dirname, '../../network/' + process.env.NODE_ENV + '/network-config/' + org + '.json'));

		//Get Client for Org
		var client = hfc.loadFromConfig(hfc.getConfigSetting(org + '-network'));
		client.loadFromConfig(hfc.getConfigSetting(org));
		client.network = connectionProfile;
		//Intitalize Key Value Store
		//await client.initCredentialStores();

		//External DB KeyStore Implementation
		let clientConfig = client.getClientConfig();
		utils.setConfigSetting('crypto-keysize', clientConfig['crypto-keysize']);
		utils.setConfigSetting('key-value-store', path.join(__dirname, clientConfig['key-value-store-impl']));

		let options = {
			name: clientConfig['key-value-store-name'],
			url: clientConfig['key-value-store-db']
		};
		let kvs = await utils.newKeyValueStore(options);

		let cryptoSuite = hfc.newCryptoSuite();
		cryptoSuite.setCryptoKeyStore(hfc.newCryptoKeyStore(options));
		client.setCryptoSuite(cryptoSuite);
		client.setStateStore(kvs);

		clients[org] = client;
		evenHubs[org] = client.getEventHubsForOrg(client.getMspid())
		return client;
	} else {
		return clients[org];
	}
};

module.exports.getEventHubsForOrg = function (userOrg) {
	return evenHubs[userOrg];
}

module.exports.getRegistrarForOrg = async function (userOrg) {
	let client = await this.getClientForOrg(userOrg);
	let certificateAuthority = client.network.organizations[userOrg].certificateAuthorities[0];

	return client.network.certificateAuthorities[certificateAuthority].registrar[0];
}

module.exports.getRegisteredUser = async function (username, userOrg) {
	try {
		var client = await this.getClientForOrg(userOrg);
		let user = await client.getUserContext(username, true);

		if (user && user.isEnrolled()) {
			logger.info('Successfully loaded "%s" of org "%s" from persistence', username, userOrg);
			return user;
		} else {
			throw new Error('username or password incorrect');
		}
	} catch (err) {
		logger.error('Failed to get Registered User: ' + err.stack ? err.stack : err);
		throw new Error('Failed to get Registered User: ' + err.toString());
	}
};

module.exports.getAdminUser = async function (userOrg) {
	try {
		let client = await this.getClientForOrg(userOrg);
		let registrar = await this.getRegistrarForOrg(userOrg);
		let user = await client.getUserContext(registrar.enrollId, true);

		if (user && user.isEnrolled()) {
			logger.info('Successfully loaded "%s" of org "%s" from persistence', registrar.enrollId, userOrg);
			return user;
		} else {
			let admin = await client.setUserContext({
				username: registrar.enrollId,
				password: registrar.enrollSecret
			});
			user = await client.getUserContext(registrar.enrollId, true);

			if (user && user.isEnrolled()) {
				logger.info('Successfully loaded "%s" of org "%s" from persistence', registrar.enrollId, userOrg);
				return user;
			} else {
				throw new Error('username or password incorrect');
			}
		}
	} catch (err) {
		logger.error('Failed to get Registered User: ' + err.stack ? err.stack : err);
		throw new Error('Failed to get Registered User: ' + err.toString());
	}
};

module.exports.registerUser = async function (username, secret, userOrg, isJson) {
	try {
		let client = await this.getClientForOrg(userOrg);
		let message = null;

		logger.info('User "%s" was not enrolled, so we will need an admin user object to register', username);
		let registrar = await this.getRegistrarForOrg(userOrg);
		let adminUserObj = await client.setUserContext({
			username: registrar.enrollId,
			password: registrar.enrollSecret
		});
		let caClient = client.getCertificateAuthority();

		// TODO: Temporary fix till Starter Plan gets changed
		let affiliation = null;
		if (process.env.NODE_ENV == 'production') {
			affiliation = 'org1.department1'
		} else {
			affiliation = userOrg.toLowerCase() + '.department1'
		}

		try {
			await caClient.register({
				enrollmentID: username,
				enrollmentSecret: secret,
				affiliation: affiliation
			}, adminUserObj);
		} catch (error) {
			logger.error('Failed to get registered user: "%s" with error: "%s"', username, error.toString());
			throw new Error(error.toString());
		}

		logger.info('Successfully got the secret for user "%s"', username);

		try {
			message = await caClient.enroll({
				enrollmentID: username,
				enrollmentSecret: secret
			});
		} catch (error) {
			logger.error('Failed to get enroll user: "%s" with error: "%s"', username, error.toString());
			throw new Error(error.toString());
		}

		logger.info(username + ' enrolled successfully on ' + userOrg);

		let member = new User(username);
		member._enrollmentSecret = secret;
		member.setCryptoSuite(client.getCryptoSuite());
		message = await member.setEnrollment(message.key, message.certificate,
			client.getMspid(userOrg));

		if (member === null)
			throw new Error(message);

		await client.setUserContext(member);
		if (member && member.isEnrolled) {
			if (isJson && isJson === true) {
				var response = {
					success: true,
					username: username,
					password: member._enrollmentSecret
				};
				return response;
			}
		} else {
			throw new Error('User was not enrolled ');
		}
	} catch (error) {
		logger.error('Failed to get registered user: "%s" with error: "%s"', username, error.toString());
		throw new Error(error.toString());
	}
}

// TODO: Temporary fix till Fabric SDK gets Fixed.  Remove this method and use getClientForOrg instead
module.exports.getAdminClientForOrg = async function (userOrg) {
	var client = await this.getClientForOrg(userOrg);

	let privateKeyPEM = null;
	let signedCertPEM = null;
	let orgAdmin = client._network_config._network_config.organizations[client.getClientConfig().organization];

	if (orgAdmin.adminPrivateKey.pem) {
		privateKeyPEM = orgAdmin.adminPrivateKey.pem
	} else {
		privateKeyPEM = fs.readFileSync(path.join(__dirname, '../', orgAdmin.adminPrivateKey.path));
	}

	if (orgAdmin.signedCert.pem) {
		signedCertPEM = orgAdmin.signedCert.pem
	} else {
		signedCertPEM = fs.readFileSync(path.join(__dirname, '../', orgAdmin.signedCert.path));
	}

	await client.createUser({
		username: 'peer' + userOrg + 'Admin',
		mspid: client.getMspid(),
		cryptoContent: {
			privateKeyPEM: privateKeyPEM,
			signedCertPEM: signedCertPEM
		}
	});

	return client;
};

module.exports.setupChaincodeDeploy = function () {
	logger.info('chaincodePath' + path.join(__dirname, hfc.getConfigSetting('chaincodePath')));
	process.env.GOPATH = path.join(__dirname, hfc.getConfigSetting('chaincodePath'));
};

module.exports.inspectProposalResult = function (proposalResult) {
	let proposalResponses = proposalResult[0];
	let proposal = proposalResult[1];
	let all_good = true;
	for (var i in proposalResponses) {
		let one_good = false;
		if (proposalResponses && proposalResponses[i].response &&
			proposalResponses[i].response.status === 200) {
			one_good = true;
			logger.info('Proposal was good');
		} else {
			logger.error('Proposal was bad');
		}
		all_good = all_good & one_good;
	}
	if (all_good) {
		logger.info(util.format(
			'Successfully sent Proposal and received ProposalResponse: Status - "%s", message - ""%s""',
			proposalResponses[0].response.status, proposalResponses[0].response.message));
	} else {
		throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200.');
	}

	return all_good;
}

module.exports.getLogger = function (moduleName) {
	var logger = log4js.getLogger(moduleName);
	logger.setLevel(hfc.getConfigSetting('loglevel'));
	return logger;
};