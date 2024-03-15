// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.10;

interface IERC721 {
    /**
     * @dev Emitted when `tokenId` token is transferred from `from` to `to`.
     */
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);

    /**
     * @dev Emitted when `owner` enables `approved` to manage the `tokenId` token.
     */
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);

    /**
     * @dev Emitted when `owner` enables or disables (`approved`) `operator` to manage all of its assets.
     */
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);

    /**
     * @dev Returns the number of tokens in ``owner``'s account.
     */
    function balanceOf(address owner) external view returns (uint256 balance);

    /**
     * @dev Returns the owner of the `tokenId` token.
     *
     * Requirements:
     *
     * - `tokenId` must exist.
     */
    function ownerOf(uint256 tokenId) external view returns (address owner);

    /**
     * @dev Safely transfers `tokenId` token from `from` to `to`.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must exist and be owned by `from`.
     * - If the caller is not `from`, it must be approved to move this token by either {approve} or {setApprovalForAll}.
     * - If `to` refers to a smart contract, it must implement {IERC721Receiver-onERC721Received}, which is called upon
     *   a safe transfer.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external;

    /**
     * @dev Safely transfers `tokenId` token from `from` to `to`, checking first that contract recipients
     * are aware of the ERC-721 protocol to prevent tokens from being forever locked.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must exist and be owned by `from`.
     * - If the caller is not `from`, it must have been allowed to move this token by either {approve} or
     *   {setApprovalForAll}.
     * - If `to` refers to a smart contract, it must implement {IERC721Receiver-onERC721Received}, which is called upon
     *   a safe transfer.
     *
     * Emits a {Transfer} event.
     */
    function safeTransferFrom(address from, address to, uint256 tokenId) external;

    /**
     * @dev Transfers `tokenId` token from `from` to `to`.
     *
     * WARNING: Note that the caller is responsible to confirm that the recipient is capable of receiving ERC-721
     * or else they may be permanently lost. Usage of {safeTransferFrom} prevents loss, though the caller must
     * understand this adds an external call which potentially creates a reentrancy vulnerability.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `tokenId` token must be owned by `from`.
     * - If the caller is not `from`, it must be approved to move this token by either {approve} or {setApprovalForAll}.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(address from, address to, uint256 tokenId) external;

    /**
     * @dev Gives permission to `to` to transfer `tokenId` token to another account.
     * The approval is cleared when the token is transferred.
     *
     * Only a single account can be approved at a time, so approving the zero address clears previous approvals.
     *
     * Requirements:
     *
     * - The caller must own the token or be an approved operator.
     * - `tokenId` must exist.
     *
     * Emits an {Approval} event.
     */
    function approve(address to, uint256 tokenId) external;

    /**
     * @dev Approve or remove `operator` as an operator for the caller.
     * Operators can call {transferFrom} or {safeTransferFrom} for any token owned by the caller.
     *
     * Requirements:
     *
     * - The `operator` cannot be the address zero.
     *
     * Emits an {ApprovalForAll} event.
     */
    function setApprovalForAll(address operator, bool approved) external;

    /**
     * @dev Returns the account approved for `tokenId` token.
     *
     * Requirements:
     *
     * - `tokenId` must exist.
     */
    function getApproved(uint256 tokenId) external view returns (address operator);

    /**
     * @dev Returns if the `operator` is allowed to manage all of the assets of `owner`.
     *
     * See {setApprovalForAll}
     */
    function isApprovedForAll(address owner, address operator) external view returns (bool);
}

contract Marketplace {
    event NFTListed(address indexed seller, uint256 indexed nftId);
    event NFTBought(address indexed buyer, address indexed seller, uint256 indexed nftId, uint256 price);

    IERC721 public nft;
    uint256 public nftId;
    address payable public seller;
    bool public onSale;
    uint256 public price;

    constructor(address _nft, uint256 _nftId, uint256 _price) {
        nft = IERC721(_nft);
        nftId = _nftId;
        seller = payable(msg.sender);
        price = _price;
        onSale = true;
        require(nft.ownerOf(nftId) == msg.sender, "sender not the nft owner");
        emit NFTListed(seller, nftId);
    }

    // list
    // buy
    //
    //
    // list
    // buy (nft goes to seller -> smart contract, money goes buyer -> smart contract)
    // shipping (if cancelled, nft goes back to seller, money goes back to buyer)
    // received
    // complete transaction (nft goes to buyer, money goes to seller)
    //

    // list
    // enabled listing (NFT goes buyer -> smart contract)
    // users bid (money goes bidder -> smart contract)
    // listing ends (assuming bidder exists,
    //            NFT goes to highest bidder, money from highest bidder goes smart contract -> seller
    //            all other bids go from smart contract -> each bidder (aka a refund)

    function buyNFT() external {
        address buyer = msg.sender;
        require(onSale, "not on sale");
        require(seller != buyer, "buyer cannot be the seller");
        require(nft.ownerOf(nftId) == seller, "seller no longer owns nft");

        nft.transferFrom(seller, buyer, nftId);
        // contract ERC721: function transferFrom(from, to, index)
        //   msg.sender == contract address
        //   param 0 == seller address
        //   param 1 == buyer address (previously msg.sender)
        //   param 2 == the nft id within the param 0 address

        onSale = false;
        emit NFTBought(buyer, seller, nftId, price);
    }
}

