import { ethers } from 'ethers';

// 从页面元素获取文章数据
function getArticleData() {
    const dataElement = document.getElementById('articleData');
    // 解析验证数据字符串为对象
    let verification;
    try {
        verification = JSON.parse(dataElement.dataset.verification);
    } catch (error) {
        console.error('Failed to parse verification data:', error);
        verification = {};
    }

    return {
        verification: verification,
        title: dataElement.dataset.title || ''
    };
}

// 加载 ABI
async function loadContractABI() {
    try {
        const staticPath = document.querySelector('meta[name="static-path"]').content;
        const response = await fetch(`${staticPath}/abi/ArticleNFT.json`);
        const data = await response.json();
        return data.abi;
    } catch (error) {
        console.error('Failed to load ABI:', error);
        throw new Error('Failed to load contract ABI');
    }
}

export async function mintNFT() {
    const button = document.getElementById('mintButton2');
    const buttonText = document.getElementById('mintButtonText2');
    const spinner = document.getElementById('mintSpinner2');

    // 设置按钮为加载状态
    function setLoading(loading) {
        if (loading) {
            button.disabled = true;
            button.classList.add('opacity-50', 'cursor-not-allowed');
            buttonText.textContent = '铸造中...';
            spinner.classList.remove('hidden');
        } else {
            button.disabled = false;
            button.classList.remove('opacity-50', 'cursor-not-allowed');
            buttonText.textContent = '铸造 NFT';
            spinner.classList.add('hidden');
        }
    }

    // 检是否安装了 MetaMask
    if (typeof window.ethereum === 'undefined') {
        window.showToast('请先安装并登录 MetaMask 钱包', 'error');
        return;
    }

    try {
        setLoading(true);

        // 请求用户连接钱包
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        const userAddress = accounts[0];
        
        // 获取文章数据
        const articleData = getArticleData();
        const verification = articleData.verification;
        
        // 创建合约实例
        const contractAddress = verification.NftContract;
        if (!contractAddress) {
            window.showToast('未配置 NFT 合约地址', 8000, 'error');
            return;
        }

        // 检查网络
        let provider = new ethers.providers.Web3Provider(window.ethereum);
        const network = await provider.getNetwork();
        const targetChainId = verification.NFT.ChainId;

        console.log("Network:", network);
        console.log("targetChainId:", targetChainId);

        if (targetChainId === 0) {
            window.showToast('网络切换失败，请手动切换到正确的网络', 8000, 'error');
            return;
        }
        
        if (network.chainId !== targetChainId) {
            // 请求切换网络，MetaMask 会处理提示
            await window.ethereum.request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: `0x${targetChainId.toString(16)}` }],
            });
            // 等待网络切换完成
            await new Promise((resolve) => setTimeout(resolve, 1000));
            
            // 重新获取 provider 和网络
            const updatedProvider = new ethers.providers.Web3Provider(window.ethereum);
            const updatedNetwork = await updatedProvider.getNetwork();
            
            if (updatedNetwork.chainId !== targetChainId) {
                window.showToast('网络切换失败，请手动切换到正确的网络', 8000, 'error');
                return;
            }
            
            // 使用更新后的 provider
            provider = updatedProvider;
        }

        // 检查余额
        const balance = await provider.getBalance(userAddress);
        const price = ethers.utils.parseEther(verification.NFT.Price);
        if (balance.lt(price.mul(2))) { // 确保有足够余额支付 gas
            window.showToast('钱包余额不足，请确保有足够的测试币（建议至少 0.002 ETH）', 8000, 'error');
            return;
        }

        // 加载 ABI
        const contractABI = await loadContractABI();

        // 使用 ethers.js 创建合约实例
        const signer = provider.getSigner();
        const contract = new ethers.Contract(contractAddress, contractABI, signer);

        // 获取文章信息
        if (!verification) {
            window.showToast('文章未经过验证，无法铸造 NFT', 8000, 'error');
            return;
        }

        // 准备参数
        const params = {
            author: verification.Author,
            name: articleData.title,
            contentHash: verification.ContentHash,
            arweaveId: verification.ArweaveId,
            version: verification.NFT.Version,
            price: ethers.utils.parseEther(verification.NFT.Price),
            maxSupply: verification.NFT.MaxSupply,
            royaltyFee: verification.NFT.RoyaltyFee,
            onePerAddress: verification.NFT.OnePerAddress
        };

        // 调用合约
        const tx = await contract.mintArticle(
            params.author,
            params.name,
            params.contentHash,
            params.arweaveId,
            params.version,
            params.price,
            params.maxSupply,
            params.royaltyFee,
            params.onePerAddress,
            { value: ethers.utils.parseEther(verification.NFT.Price) }
        );

        window.showToast('NFT 铸造中，请等待确认', 10000);
        
        // 等待交易确认
        const receipt = await tx.wait();
        
        if (receipt.status === 1) {
            window.showToast('NFT 铸造成功！', 8000);
            // 铸造成功后禁用按钮
            button.disabled = true;
            buttonText.textContent = '已铸造';
            button.classList.add('bg-green-500/10', 'text-green-500');
            spinner.classList.add('hidden');
        } else {
            window.showToast('NFT 铸造失败，请重试', 8000, 'error');
            setLoading(false);
        }
    } catch (error) {
        console.error('Minting failed:', error);
        let message = '铸造失败，请重试';
        
        // 声明在这里，以便后面的错误处理也能访问
        let contractABI;
        let iface;
        
        if (error.code === 4001) {
            message = '用户取消了操作';
        } else if (error.code === 'UNPREDICTABLE_GAS_LIMIT') {
            try {
                // 加载 ABI
                contractABI = await loadContractABI();
                // 获取合约 interface
                iface = new ethers.utils.Interface(contractABI);
                // 获取错误数据
                const errorData = error.error?.data || error.data;
                if (errorData) {
                    // 解码错误
                    const decodedError = iface.parseError(errorData);
                    // 根据错误名称判断
                    switch (decodedError.name) {
                        case 'RoyaltyFeeTooHigh':
                            message = '版税设置过高（不能超过100%）';
                            break;
                        case 'AlreadyMinted':
                            message = '此文章已被铸造';
                            break;
                        case 'ArweaveIdEmpty':
                            message = '缺少 Arweave ID';
                            break;
                        case 'ContentHashEmpty':
                            message = '缺少内容哈希';
                            break;
                        case 'ExceedsMaxSupply':
                            message = '超过最大发行量';
                            break;
                        case 'InsufficientPayment':
                            message = '支付金额不足';
                            break;
                        case 'MaxSupplyInvalid':
                            message = '最大发行量设置无效';
                            break;
                        case 'NameEmpty':
                            message = '文章标题不能为空';
                            break;
                        default:
                            message = `合约错误: ${decodedError.name}`;
                    }
                }
            } catch (decodeError) {
                console.error('Failed to decode error:', decodeError);
                message = '合约调用失败，请检查参数';
            }
        } else if (error.message.includes('insufficient funds')) {
            message = '钱包余额不足';
        }

        // 打印详细错误信息用于调试
        // console.log('Error details:', {
        //     code: error.code,
        //     errorData: error.error?.data,
        //     decodedError: error.error?.data && iface ? iface.parseError(error.error.data) : null,
        //     message: error.message,
        //     params: error.transaction
        // });

        window.showToast(message, 8000, 'error');
        setLoading(false);
    }
} 