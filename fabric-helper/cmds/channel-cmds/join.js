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
 */

'use strict';

const channelLib = require('../../lib/join-channel.js')

exports.command = 'join';
exports.desc = 'Join Channel';
exports.builder = function (yargs) {
    return yargs.option('channel-name', {
        demandOption: true,
        describe: 'Name for the channel to join',
        requiresArg: true,
        type: 'string'
    });
};

exports.handler = function (argv) {
    return channelLib.joinChannel(argv['channel-name'], argv['org']);
};