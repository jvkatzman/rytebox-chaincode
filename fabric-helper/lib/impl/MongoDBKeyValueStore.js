/*
 Copyright 2016 IBM All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the 'License');
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an 'AS IS' BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 * 
 * Author: Sandeep Pulluru <sandeep.pulluru@ibm.com>
 */

'use strict';

var api = require('fabric-client/lib/api.js');
var fs = require('fs-extra');
var path = require('path');
var util = require('util');
var utils = require('fabric-client/lib/utils');
var nano = require('nano');

var logger = utils.getLogger('MongoDBKeyValueStore.js');
const mongo = require('mongodb')
const MongoClient = mongo.MongoClient;

/**
 * This is a sample database implementation of the [KeyValueStore]{@link module:api.KeyValueStore} API.
 * It uses a local or remote MongoDB database instance to store the keys.
 *
 * @class
 * @extends module:api.KeyValueStore
 */
var MongoDBKeyValueStore = class extends api.KeyValueStore {
    /**
     * @typedef {Object} MongoDBOpts
     * @property {string} url The MongoDB instance url, in the form of http(s)://<user>:<password>@host:port
     * @property {string} name Optional. Identifies the name of the database to use. Default: <code>member_db</code>.
     */

    /**
     * constructor
     *
     * @param {MongoDBOpts} options Settings used to connect to a MongoDB instance
     */
    constructor(options) {
        logger.debug('constructor', {
            options: options
        });

        if (!options || !options.url) {
            throw new Error('Must provide the MongoDB database url to store membership data.');
        }

        // Create the keyValStore instance
        super();

        var self = this;
        // url is the database instance url
        this._url = options.url;
        // Name of the database, optional
        if (!options.name) {
            this._name = 'member_db';
        } else {
            this._name = options.name;
        }

        return new Promise(function (resolve, reject) {
            resolve(self);
        });
    }

    getValue(name) {
        logger.debug('getValue', {
            key: name
        });

        var self = this;
        return new Promise(function (resolve, reject) {
            MongoClient.connect(self._url, {
                poolSize: 10
                // other options can go here
            }, function (err, mongodb) {
                mongodb.collection(self._name).findOne({
                    "_id": name
                }).then(function (result) {
                    mongodb.close();
                    if (result) {
                        logger.debug('result.member %s.', result.member);
                        return resolve(result.member);
                    } else {
                        return resolve(null);
                    }
                }).catch(function (error) {
                    mongodb.close();
                    logger.error('getValue: %s, ERROR: [%s.get] - ', name, self._name, error.message);
                    reject(error.message);
                });
            });
        });
    }

    setValue(name, value) {
        logger.debug('setValue', {
            key: name
        });

        var self = this;

        return new Promise(function (resolve, reject) {
            MongoClient.connect(self._url, {
                poolSize: 10
                // other options can go here
            }, function (err, mongodb) {
                mongodb.collection(self._name).findOne({
                    "_id": name
                }).then(function (result) {
                    logger.debug('result %s.', result);
                    if (result) {
                        mongodb.collection(self._name).updateOne({
                            _id: name
                        }, {
                                $set: {
                                    member: value
                                }
                            }).then(function (result) {
                                mongodb.close(true);
                                resolve(true);
                            }).catch(function (error) {
                                mongodb.close(true);
                                logger.error('setValue: %s, ERROR: [%s.get] - ', name, self._name, error.message);
                                reject(error.message);
                            });
                    } else {
                        mongodb.collection(self._name).insertOne({
                            _id: name,
                            member: value
                        }).then(function (result) {
                            mongodb.close();
                            resolve(true);
                        }).catch(function (error) {
                            mongodb.close();
                            logger.error('setValue: %s, ERROR: [%s.get] - ', name, self._name, error.message);
                            reject(error.message);
                        });

                    }
                }).catch(function (error) {
                    mongodb.close();
                    logger.error('setValue: %s, ERROR: [%s.get] - ', name, self._name, error.message);
                    reject(error.message);
                });
            });
        });
    }
};


module.exports = MongoDBKeyValueStore;