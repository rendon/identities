# Identities API
Retrieve and cache social network identities.

This API provides two endpoints:

    GET     /ids/:network/:username     # Retrieve ID given a social network and a user name
    GET     /usernames/:network/:id     # Retrieve user name given social network and user ID

# Dependencies
In order to run this software you need a mongoDB server running with standard settings. I recommend using [Docker](https://hub.docker.com/_/mongo/).

Not exactly a dependency but a requirement, you need an entry in your `/etc/hosts` with `mongodb-server` pointing to the server where mongoDB is running on.

# Development
## Setup

    > bin/setup

## Run

    > bin/run

## Test

    > bin/test
