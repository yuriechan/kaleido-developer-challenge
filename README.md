## About this Project - "Marketplace"
Build a MVP for a e-commerce platform for second hand item where digitally sold items allows permanent tracing to ownership of the item.
- Use case:
- As a user I want to be able to track my items so that I can easily re-list and sell items I no longer need
  - A Marketplace Wallet lets me track my items
- As a buyer I want to be able to track how often an item is sold so I can understand it's resale value
  - An item's prior sales shows price history and sale frequency
- As a buyer, I want a frictionless re-listing experience so that I will sell my unused items
  - An item's NFT maintains the metadata for re-using the same data (images, name, etc) with new listings
- As a service, I want to increase re-sale of items so that I can grow more revenue from the sales fees
  - A frictionless user-experience encourages re-selling items which is how we grow revenue 
- As a service, I want transparency in items and sales so that users have more trust in buying and selling
  - A blockchain is immutable and transparent for users to trust
- As a user, I don't want my physical address publicly available so that my PII stays private
  - A private database for a hybrid public/private solution allows for handling PII with transparency for public data

## 🥞 Tech Stack
- Backend written in Go 
- DB uses Mysql
- Smart Contract written in Solidity 
- Uses Firefly API to interact with ethereum

## 🎁 Feature
- HTTP endpoint to
 - mint NFT associated to your physical / digital item 
 - buy NFT that is available in the marketplace

## 🎡 Things I have considered during the development
- Using docker compose for the ease of running this project
- Applying onion architecture to keep the codebase maintainable, taking advantages of Go's interface 
- Identifying what type of data should be stored off-chain vs on-chain

## 💡 Learning Experience
- Having a marketplace means it acts as an intermediary between the seller and buyer, which requires approval from the users for the marketplace to invoke events on smart contract on behalf of the users. 
- Multiple instances of the smart contract can be deployed to different addresses, allowing each item listing to have its own dedicated NFT attached to it.
- There are many ways to interact with ethereum account to achieve MVP - transferring directly from user A to user B VS transferring from user A to smart contract, and smart contract to user B 

## 🚧 Fix needed - Thought for better service design 
- Use event listeners to listen to status update on a blockchain node - this allows asynchronously handling event update while processing user requests.
- Implement payment features where digital asset is transferred alongside the NFT transfer.
- User authentication via wallet for more secure approach to handling asset associated to the user. 
- Expand more features to thoroughly cover the user journey - users experience extends to not only buying / selling but also item shipment, item returns, and payment.

## 🚀 Project setup
1. Install [direnv](https://github.com/direnv/direnv#install) following the installation step. <br>
Once installation is complete, run:
```shell
$ cd backend 
$ cp .env_sample .env 
$ vi .env # Add value to each ENV 
$ direnv allow 
```

2. Install Firefly CLI following this [guide](https://hyperledger.github.io/firefly/gettingstarted/).
- After running its sandbox environment, you must deploy the smart contract available under `backend/contracts/marketplace.sol`
- Also, generate Firefly's autogenerated interface to interact with the deployed smart contract 

3. Start up the project using Docker compose 
```shell 
$ cd backend 
$ make docker-compose-up 
$ make docker-compose-down 
```
- To kill the server, use `make docker-compose-down` instead.