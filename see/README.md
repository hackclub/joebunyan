# See the hackclub slack!

There is a small quick [demo](https://slack-files.com/T0266FRGM-F02ARE0RBA8-69fdab8fff).

## How it works

Using a slack bot called Joe Bunyan created by the famous [Caleb Denio](https://github.com/cjdenio) this program uses a WebSocket API to get live messages and reactions right from the slack. Messages are in white and reactions are in green. Each channel also gets its coordinate on the matrix display.

## Install

Do you have your own raspberry pi and want to run this on it? Just follow these steps to get all setup:

1. Make sure you have golang 1.16 or later installed on your machine
2. Change the SCP hostname (`pi@joebunyan.local`) to the hostname of your pi in the [Makefile](./Makefile).
3. Run `make deploy`
4. SSH to the pi
5. Run `./see`
