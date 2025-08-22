// 等待DOM加载完成
document.addEventListener('DOMContentLoaded', function() {
    // 初始化所有功能
    initCarousel();
    initSearch();
    initLoginForm();
    initTagCloud();
    initNavigation();
    initProgressBar();
});

// 轮播图功能
function initCarousel() {
    const carouselItems = document.querySelectorAll('.carousel-item');
    let currentIndex = 0;

    function showSlide(index) {
        carouselItems.forEach((item, i) => {
            item.classList.remove('active');
            if (i === index) {
                item.classList.add('active');
            }
        });
    }

    function nextSlide() {
        currentIndex = (currentIndex + 1) % carouselItems.length;
        showSlide(currentIndex);
    }

    // 自动轮播
    if (carouselItems.length > 1) {
        setInterval(nextSlide, 5000);
    }
}

// 搜索功能
function initSearch() {
    const searchInput = document.querySelector('.search-box input');
    const searchIcon = document.querySelector('.search-box i');

    if (searchInput) {
        searchInput.addEventListener('focus', function() {
            this.parentElement.style.boxShadow = '0 0 0 2px rgba(74, 144, 226, 0.2)';
        });

        searchInput.addEventListener('blur', function() {
            this.parentElement.style.boxShadow = 'none';
        });

        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                performSearch(this.value);
            }
        });

        searchIcon.addEventListener('click', function() {
            performSearch(searchInput.value);
        });
    }
}

function performSearch(query) {
    if (query.trim()) {
        // 这里可以添加实际的搜索逻辑
        console.log('搜索:', query);
        alert(`正在搜索: ${query}`);
    }
}

// 登录表单处理
function initLoginForm() {
    const loginForm = document.querySelector('.login-form');
    const loginBtn = document.querySelector('.btn-login');

    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = this.querySelector('input[type="text"]').value;
            const password = this.querySelector('input[type="password"]').value;

            if (!username || !password) {
                showNotification('请填写完整的登录信息', 'error');
                return;
            }

            // 模拟登录过程
            loginBtn.textContent = '登录中...';
            loginBtn.disabled = true;

            setTimeout(() => {
                // 模拟登录成功
                showNotification('登录成功！', 'success');
                loginBtn.textContent = '登录';
                loginBtn.disabled = false;
                
                // 更新用户信息显示
                updateUserDisplay(username);
            }, 2000);
        });
    }
}

function updateUserDisplay(username) {
    const userInfo = document.querySelector('.user-info');
    if (userInfo) {
        const userDetails = userInfo.querySelector('.user-details h4');
        if (userDetails) {
            userDetails.textContent = username;
        }
    }
}

// 标签云功能
function initTagCloud() {
    const tags = document.querySelectorAll('.tag-cloud .tag');
    
    tags.forEach(tag => {
        tag.addEventListener('click', function() {
            const tagText = this.textContent;
            console.log('点击标签:', tagText);
            
            // 添加点击效果
            this.style.transform = 'scale(0.95)';
            setTimeout(() => {
                this.style.transform = '';
            }, 150);
            
            // 这里可以添加标签筛选逻辑
            showNotification(`正在筛选标签: ${tagText}`, 'info');
        });
    });
}

// 导航功能
function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            // 移除所有活动状态
            navItems.forEach(nav => nav.classList.remove('active'));
            
            // 添加当前活动状态
            this.classList.add('active');
            
            // 这里可以添加页面导航逻辑
            const navText = this.querySelector('span').textContent;
            console.log('导航到:', navText);
        });
    });
}

// 进度条动画
function initProgressBar() {
    const progressFill = document.querySelector('.progress-fill');
    
    if (progressFill) {
        // 添加进度条动画
        setTimeout(() => {
            progressFill.style.transition = 'width 1s ease-in-out';
            progressFill.style.width = '65%';
        }, 500);
    }
}

// 问答项交互
function initQAInteraction() {
    const qaItems = document.querySelectorAll('.qa-item');
    
    qaItems.forEach(item => {
        item.addEventListener('click', function() {
            const title = this.querySelector('.qa-title').textContent;
            console.log('查看问答:', title);
            
            // 添加点击效果
            this.style.transform = 'translateY(-2px)';
            this.style.boxShadow = '0 8px 25px rgba(74, 144, 226, 0.15)';
            
            setTimeout(() => {
                this.style.transform = '';
                this.style.boxShadow = '';
            }, 200);
        });
    });
}

// 技术分享交互
function initTechInteraction() {
    const techItems = document.querySelectorAll('.tech-item');
    
    techItems.forEach(item => {
        item.addEventListener('click', function() {
            const title = this.querySelector('.tech-title').textContent;
            console.log('查看技术分享:', title);
            
            // 添加点击效果
            this.style.transform = 'scale(1.02)';
            setTimeout(() => {
                this.style.transform = '';
            }, 200);
        });
    });
}

// 学习资料下载
function initLearningDownload() {
    const downloadLinks = document.querySelectorAll('.download-link');
    
    downloadLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            const title = this.closest('.learning-item').querySelector('.learning-title').textContent;
            console.log('下载资料:', title);
            
            // 模拟下载过程
            this.textContent = '下载中...';
            this.style.pointerEvents = 'none';
            
            setTimeout(() => {
                this.textContent = '下载';
                this.style.pointerEvents = '';
                showNotification('下载完成！', 'success');
            }, 2000);
        });
    });
}

// 社群推荐交互
function initCommunityInteraction() {
    const communityItems = document.querySelectorAll('.community-item');
    
    communityItems.forEach(item => {
        item.addEventListener('click', function() {
            const name = this.querySelector('h4').textContent;
            console.log('加入社群:', name);
            
            showNotification(`正在加入社群: ${name}`, 'info');
        });
    });
}

// 专家推荐交互
function initExpertInteraction() {
    const expertItems = document.querySelectorAll('.expert-item');
    
    expertItems.forEach(item => {
        item.addEventListener('click', function() {
            const name = this.querySelector('h4').textContent;
            const field = this.querySelector('p').textContent;
            console.log('查看专家:', name, field);
            
            showNotification(`正在查看专家: ${name} (${field})`, 'info');
        });
    });
}

// 通知系统
function showNotification(message, type = 'info') {
    // 创建通知元素
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <i class="fas fa-${getNotificationIcon(type)}"></i>
            <span>${message}</span>
        </div>
    `;
    
    // 添加样式
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: ${getNotificationColor(type)};
        color: white;
        padding: 15px 20px;
        border-radius: 8px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
        z-index: 10000;
        transform: translateX(100%);
        transition: transform 0.3s ease;
        max-width: 300px;
    `;
    
    // 添加到页面
    document.body.appendChild(notification);
    
    // 显示动画
    setTimeout(() => {
        notification.style.transform = 'translateX(0)';
    }, 100);
    
    // 自动隐藏
    setTimeout(() => {
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

function getNotificationIcon(type) {
    const icons = {
        success: 'check-circle',
        error: 'exclamation-circle',
        warning: 'exclamation-triangle',
        info: 'info-circle'
    };
    return icons[type] || 'info-circle';
}

function getNotificationColor(type) {
    const colors = {
        success: '#50C878',
        error: '#e74c3c',
        warning: '#f39c12',
        info: '#4A90E2'
    };
    return colors[type] || '#4A90E2';
}

// 滚动效果
function initScrollEffects() {
    window.addEventListener('scroll', function() {
        const scrolled = window.pageYOffset;
        const parallax = document.querySelector('.carousel-section');
        
        if (parallax) {
            const speed = scrolled * 0.5;
            parallax.style.transform = `translateY(${speed}px)`;
        }
    });
}

// 响应式菜单
function initResponsiveMenu() {
    const menuToggle = document.createElement('button');
    menuToggle.className = 'menu-toggle';
    menuToggle.innerHTML = '<i class="fas fa-bars"></i>';
    menuToggle.style.cssText = `
        display: none;
        background: none;
        border: none;
        font-size: 20px;
        color: #333;
        cursor: pointer;
        padding: 10px;
    `;
    
    const sidebar = document.querySelector('.sidebar');
    if (sidebar) {
        sidebar.parentNode.insertBefore(menuToggle, sidebar);
        
        menuToggle.addEventListener('click', function() {
            sidebar.classList.toggle('mobile-open');
        });
        
        // 在小屏幕上显示菜单按钮
        function checkScreenSize() {
            if (window.innerWidth <= 992) {
                menuToggle.style.display = 'block';
                sidebar.style.display = 'none';
            } else {
                menuToggle.style.display = 'none';
                sidebar.style.display = 'block';
            }
        }
        
        checkScreenSize();
        window.addEventListener('resize', checkScreenSize);
    }
}

// 初始化所有交互功能
document.addEventListener('DOMContentLoaded', function() {
    initQAInteraction();
    initTechInteraction();
    initLearningDownload();
    initCommunityInteraction();
    initExpertInteraction();
    initScrollEffects();
    initResponsiveMenu();
});

// 添加一些CSS样式到页面
const additionalStyles = `
    .notification-content {
        display: flex;
        align-items: center;
        gap: 10px;
    }
    
    .notification-content i {
        font-size: 16px;
    }
    
    .sidebar.mobile-open {
        display: block !important;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1000;
        background: white;
        overflow-y: auto;
    }
    
    .menu-toggle {
        position: fixed;
        top: 20px;
        left: 20px;
        z-index: 1001;
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }
`;

// 将样式添加到页面
const styleSheet = document.createElement('style');
styleSheet.textContent = additionalStyles;
document.head.appendChild(styleSheet); 