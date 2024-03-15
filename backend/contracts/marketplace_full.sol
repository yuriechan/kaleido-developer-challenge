// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.10;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC721/extensions/ERC721Full.sol";
import "./marketplace.sol";

contract Marketplace is ERC721Full {
    event NFTNew(uint256 indexed nftId);
    event NFTListed(address indexed seller, uint256 indexed nftId);
    event NFTBought(address indexed buyer, address indexed seller, uint256 indexed nftId, uint256 price);
    event NFTShipped(address indexed seller, uint256 indexed nftId);
    event NFTReceived(address indexed buyer, uint256 indexed nftId);
    event NFTCancel(uint256 indexed nftId);

    mapping (uint256 => bool) statuses;
    mapping (uint256 => uint256) prices;
    mapping (uint256 => address payable) buyers; // buyer is payable incase of refund

    constructor() ERC721("Marketplace NFT", "MPNFT") {}

    // When minting new token, do extra setup for our marketplace nft
    function _mint(address to, uint256 tokenId) internal override(ERC721Full) {
        super._mint(to, tokenId);
        statuses.push(tokenId, "owned");
        buyers.push(tokenId, address(0));
        prices.push(tokenId, 0);
        emit NFTNew(tokenId);
    }

    function _transferFrom(address from, address to, uint256 tokenId) internal override(ERC721Full) {
        // If item is transferred reset sale status
        super._transferFrom(from, to, tokenId);
        statuses[tokenId] = "owned";
        buyers[tokenId] = address(0);
        prices[tokenId] = 0;
    }

    function list(uint256 tokenId) external {
        require(this.ownerOf(tokenId) == msg.sender, "Only NFT owner can list an item");
        require(statuses[tokenId] == "owned", "NFT must be in owned status");
        statuses[tokenId] = "listed";
        emit NFTListed(msg.sender, tokenId);
    }

    function buy(uint256 tokenId) payable external {
        require(this.ownerOf(tokenId) != msg.sender, "Owner cannot buy own token");
        require(statuses[tokenId] == "listed", "NFT must be in listed status");
        statuses[tokenId] = "bought";
        buyers[tokenId] = payable(msg.sender);
        prices[tokenId] = msg.value;
        emit NFTBought(msg.sender, this.ownerOf(tokenId), tokenId, msg.value);
    }

    function shipped(uint256 tokenId) external {
        require(this.ownerOf(tokenId) == msg.sender, "Only seller can set item as shipped");
        require(statuses[tokenId] == "bought", "NFT must be in bought status");
        statuses[tokenId] = "shipped";
        emit NFTShipped(msg.sender, tokenId);
    }

    function received(uint256 tokenId) external {
        require(buyers[tokenId] == msg.sender, "Only buyer can call received");
        require(statuses[tokenId] == "shipped", "NFT must be in bought shipped");

        address seller = this.ownerOf(tokenId);
        // Send money to token owner
        bool ok = seller.call{value: prices[tokenId]}("");
        require(ok, "Failed to pay seller");

        // Transfer owner to buyer (_transfer override resets marketplace fields)
        this._transfer(seller, buyers[tokenId], tokenId);
        emit NFTReceived(msg.sender, tokenId);
    }

    function cancel(uint256 tokenId) external {
        require(statuses[tokenId] != "shipped", "NFT cannot be shipped");
        require(statuses[tokenId] != "owned", "NFT cannot be shipped");
        statuses[tokenId] = "owned";
        buyers[tokenId] = address(0);
        price = prices[tokenId];
        if (price > 0) {
            bool ok = buyers[tokenId].call{value: price}("");
            require(ok, "Failed to refund buyer");
        }
        statuses[tokenId] = "owned";
        buyers[tokenId] = address(0);
        prices[tokenId] = 0;
        emit NFTCancel(tokenId);
    }

}

