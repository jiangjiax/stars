# Stars - Web3 Personal Website Generator

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Release](https://img.shields.io/github/v/release/jiangjiax/stars)](https://github.com/jiangjiax/stars/releases)

English | [简体中文](./README.md)

## Overview

Stars is a decentralized personal website generator developed in Go, with Web3 functionality. It helps you quickly create a modern personal website with Web3 integration. This project aims to provide creators with a decentralized content publishing platform.

### Smart Contracts

Stars supports multi-chain NFT contract deployment. Users can use the officially deployed NFT contract addresses in their article metadata, or deploy their own.

- **Ethereum Sepolia**
  - Address: `0x5c83f2287833F567b1D80D7B981084eb5CaeF445`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [Etherscan](https://sepolia.etherscan.io/address/0x5c83f2287833F567b1D80D7B981084eb5CaeF445)
  - Chain ID: 11155111

- **Telos Testnet**
  - Address: `0x903e48Ca585dBF4dFeb74f2864501feB6f0dF369`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [TelosScan](https://testnet.teloscan.io/address/0x903e48Ca585dBF4dFeb74f2864501feB6f0dF369)
  - Chain ID: 41

- **EDU Chain Testnet**
  - Address: `0xcA3Dbe8eF976e606B8c96052aaC22763aDeAEE0A`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [TelosScan](https://edu-chain-testnet.blockscout.com/address/0xcA3Dbe8eF976e606B8c96052aaC22763aDeAEE0A)
  - Chain ID: 656476

### Key Features

🚀 **High-Performance Static Site Generator**

🎨 **Modern Theme System**

📱 **Mobile Responsive**

🔗 **Web3 Integration**

🛠 **Developer Friendly**

### Contributing

We welcome all forms of contributions, whether it's new features, documentation improvements, or bug reports. Please check our [Contributing Guide](./CONTRIBUTING.md) for more information.

### Contributors

Thanks to all contributors:

<a href="https://github.com/jiangjiax/stars/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=jiangjiax/stars" />
</a>

## Quick Start

[Getting Started with Stars](./internal/generator/templates/example-posts/stars/welcome.md)

## FAQ

### Q: What is Stars?
A: Stars is a decentralized personal website generator that enables creators to publish and manage content in a decentralized way while supporting Web3 functionality integration.

### Q: How is article content stored on the blockchain?
A: Similar to NFT artworks, Stars stores content in decentralized storage systems, with only the decentralized storage addresses being stored on-chain. This mechanism applies to all forms of digital works, whether they are images, articles, or other formats.

### Q: Does Stars review or censor content?
A: Stars is a fully decentralized platform where creators have complete control over their content publishing. The platform neither can nor will censor any content.

### Q: Can published content be modified?
A: Due to the immutable nature of blockchain, content that has been published on-chain cannot be directly modified. However, Stars supports version management, allowing creators to publish new versions of content (e.g., upgrading from v1.0.0 to v1.0.1). Creators can choose to display or hide specific versions of their content.

### Q: Which blockchain networks does Stars support?
A: Stars currently supports multiple blockchain networks, including Ethereum Sepolia, Telos Testnet, and EDU Chain Testnet. Users can either use the officially deployed smart contracts or deploy their own on supported networks.

### Q: How is content permanence ensured?
A: Stars recommends users to store content using decentralized storage systems and include the decentralized storage file addresses in the article metadata, ensuring content permanence and accessibility. Even if a user's Stars personal website ceases to operate, the content remains accessible through decentralized networks.

### Q: How does Stars protect creators' rights?
A: Through blockchain technology, each article generates a unique on-chain record, providing proof of content ownership. Creators can tokenize and monetize their works in the form of NFTs.

### Q: Is there a fee for using Stars?
A: Stars itself is open-source and free to use. However, publishing content to the blockchain requires paying Gas fees on the respective network. These fees are determined by the chosen blockchain network and are independent of Stars. 