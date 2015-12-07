# Identities API
Retrieve and cache social network identities.

This API provides two endpoints:

    GET     /ids/:network/:username     # Retrieve ID given a social network and a user name
    GET     /usernames/:network/:id     # Retrieve user name given social network and user ID

# Dependencies
- In order to run this software you need a mongoDB server running with standard settings. I recommend using [Docker](https://hub.docker.com/_/mongo/).
- Not exactly a dependency but a requirement, you need an entry in your `/etc/hosts` with `mongodb-server` pointing to the server where mongoDB is running on.
- This program uses a set of tokens to connect to the Twitter API, these values are read from the file specified by the environment variable `TWITTER_KEYS_FILE`, e.g.:

        export TWITTER_KEYS_FILE="/home/user/.twitter_keys

One set of keys per line in the following way:

    consumer_key consumer_secret access_token access_token_secret

For example:

    5xY0xIxB11pOy8uEoIxgxMzxJ x91XxjIxKfxfIxo0xnqxd4xrmxSfxx9uxxv6x8txHGxipxvyxf xxxxxx5313-XfxUxmJx76xcHxjnxp1x3Px7RxhqxIqxb5x5JxO xrxmxu1EVnZ4xIhmYzyDPqzw4yYkYUrK9nlxRxvx7xxxY
    5yY0yIy2D9pOy8uEoIygyM7yJ y72yyjIyKfyfIya0ynqyd4yrmySfyy9uyyv6y8tyHGyipyvyyf yyyyyy5333-yfyUymJy76ycHyjnyp1y3Py7RyhqyIqyb5y31yO yrym3u3EVn74yIhmY7yDPq7w4yYkYUrK91ly9yvy8yyyY



# Development
## Setup

    > bin/setup

## Run

    > bin/run

## Test

    > bin/test
