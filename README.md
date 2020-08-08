# Algo Benchmark Discord Bot

This bot was build for the fxgears.com algo-trading discord channel group.
It is designed to allow for easily reporting current Roi for the day and using
it has an additional benchmark in backtesting.

There is two parts to this bot. The discord bot side that handles registration and roi submission
and the rest api side that handles calls to get current daily benchmark and historic daily benchmark.

Before submission users must be approved by another existing user. This is to avoid random users from flooding the benchmark

## Discord Comamands

* !register - User can request to register with the benchmark bot. They will be given a user-id, and rest-api token
* !giveInfo - User can request the bot to give them their user-id and rest-api token after registration
* !approve <uuid> - An approved user can approve another user to submit roi to the benchmark bot and have access to the rest api
* !submitRoi <roiAmount> - An approved user can submit their current daily close roi ex. !submitRoi 0.25 (For 0.25%)
* !dailyBm - An approved user can request the dailyBm be displayed in the channel it was sent in

## Rest Api Calls

Approved users also have access to three api endpoints. They must include the following two header values in their request

* ABDB-PARTICIPANT-ID - The Participant ID given to them by the bot in discord
* ABDB-REST-API-KEY - The Rest api token given to them by the bot in discord.


* /api/healthCheck - Blank health check to see if the api service is online
* /api/getDailyBm - Returns the current days benchmark has a json struct
* /api/getHistoricDailyBm - Returns all recorded daily benchmark values
