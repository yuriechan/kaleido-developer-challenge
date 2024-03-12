// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

interface IERC721 {
    function safeTransferFrom(address from, address to, uint256 tokenId) external;
    function transferFrom(address, address, uint256) external;
}

contract Marketplace {
    event NFTCreated(address indexed seller, uint256 indexed nftId);
    event NFTBought(address indexed buyer, address indexed seller, uint256 indexed nftId, uint256 price);

    IERC721 public nft;
    uint256 public nftId;

    address payable public seller;
    bool public onSale;

    address public buyer;
    uint256 public price;
    mapping(address => uint256) public bids;

    constructor(address _nft, uint256 _nftId, uint256 _price) {
        nft = IERC721(_nft);
        nftId = _nftId;

        seller = payable(msg.sender);
        price = _price;
    }

    function createNFT() external {
        require(!onSale, "already on sale");
        require(msg.sender == seller, "not seller");

        nft.transferFrom(msg.sender, address(this), nftId);
        onSale = true;

        emit NFTCreated(seller, nftId);
    }

    function buyNFT() external {
        require(onSale, "not on sale");

        if (buyer == address(0)) {
            buyer = msg.sender; 
        }

        if (buyer != address(0)) {
            nft.safeTransferFrom(address(this), buyer, nftId);
            seller.transfer(price);
        } else {
            nft.safeTransferFrom(address(this), seller, nftId);
        }
        onSale = false;

       emit NFTBought(buyer, seller, nftId, price);
    }
}

