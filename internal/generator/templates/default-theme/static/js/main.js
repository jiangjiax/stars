import { mintNFT } from './nft.js';

// 将函数挂载到 window 对象
window.mintNFT = mintNFT;

// Toast 通知功能
window.showToast = function(message, duration = 5000, type = 'success') {
    const alert = document.createElement('div');
    alert.className = `alert ${type === 'error' ? 'alert-error' : 'alert-success'} 
                      fixed bottom-8 left-1/2 -translate-x-1/2 z-50
                      max-w-md w-auto px-4 py-2 rounded-xl
                      bg-stars-secondary/95 backdrop-blur-sm
                      border border-stars-accent/20
                      text-stars-accent
                      shadow-lg
                      animate-toast-in`;
    
    alert.innerHTML = `
        <div class="flex items-center gap-2">
            <span class="text-lg">
                ${type === 'error' ? 
                    '<i class="fas fa-exclamation-circle"></i>' : 
                    '<i class="fas fa-check-circle"></i>'}
            </span>
            <span>${message}</span>
        </div>
    `;
    
    document.body.appendChild(alert);
    
    setTimeout(() => {
        alert.classList.add('animate-toast-out');
        setTimeout(() => {
            document.body.removeChild(alert);
        }, 300);
    }, duration);
};

// 星空背景效果
function initStars() {
    const container = document.getElementById('stars-container');
    if (!container) return;

    // 清空容器
    container.innerHTML = '';
    console.log('Creating stars...');

    // 根据屏幕大小计算星星数量
    const width = window.innerWidth;
    let starCount;
    if (width < 768) {  // 移动端
        starCount = 50;
    } else if (width < 1024) {  // 平板
        starCount = 100;
    } else {  // 面端
        starCount = 150;
    }

    console.log(`Creating ${starCount} stars for screen width ${width}px`);

    // 创建星星
    for (let i = 0; i < starCount; i++) {
        const star = document.createElement('div');
        star.className = 'star';
        
        // 随机大小
        const sizeClass = `size-${Math.floor(Math.random() * 3) + 1}`;
        star.classList.add(sizeClass);
        
        // 随机位置
        star.style.left = `${Math.random() * 100}%`;
        star.style.top = `${Math.random() * 100}%`;
        
        // 随机动画延迟和持续时间
        const duration = 2 + Math.random() * 4;
        const delay = Math.random() * 2;
        star.style.animation = `twinkle ${duration}s ${delay}s infinite`;
        
        container.appendChild(star);
    }
}

// 个人资料相关功能
function initProfile() {
    // QR Code 生成
    function generateQRCode() {
        const qrContainer = document.getElementById('pageQRCode');
        if (!qrContainer) return;
        
        qrContainer.innerHTML = '';
        new QRCode(qrContainer, {
            text: window.location.href,
            width: 64,
            height: 64,
            colorDark: "#112240",
            colorLight: "#ffffff",
            correctLevel: QRCode.CorrectLevel.H
        });
    }

    // 初始化
    generateQRCode();
}

// 初始化所有功能
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded');
    initProfile();     // 初始化个人资料功能
    initStars();       // 初始化星空背景
});

// 处理窗口大小改变
window.addEventListener('resize', () => {
    console.log('Window resized');
    initStars();  // 重新创建适合新屏幕大小的星星
}); 