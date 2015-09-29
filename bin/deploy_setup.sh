#!/usr/bin/env bash
sudo apt-get update
sudo apt-get install -y nodejs mongodb-clients
sudo ln -sfv /usr/bin/nodejs /usr/bin/node
echo 'export NODE_PATH="/usr/local/lib/node_modules"' | sudo tee --append  /etc/profile
sudo apt-get install -y npm
sudo npm install -g lazy 
sudo npm install -g pg 

