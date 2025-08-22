// 个人中心页面JavaScript功能

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeProfilePage();
    loadUserActivity();
});

// 初始化个人中心页面
function initializeProfilePage() {
    // 绑定导航菜单点击事件
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.addEventListener('click', function() {
            const section = this.getAttribute('data-section');
            switchSection(section);
        });
    });

    // 绑定过滤标签点击事件
    const filterTabs = document.querySelectorAll('.filter-tab');
    filterTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            const filter = this.getAttribute('data-filter');
            filterActivity(filter);
        });
    });

    // 绑定表单提交事件
    const profileForm = document.getElementById('profileForm');
    if (profileForm) {
        profileForm.addEventListener('submit', handleProfileSubmit);
    }
}

// 切换内容区域
function switchSection(sectionId) {
    // 移除所有导航项的active状态
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    
    // 隐藏所有内容区域
    document.querySelectorAll('.content-section').forEach(section => {
        section.classList.remove('active');
    });
    
    // 激活选中的导航项
    document.querySelector(`[data-section="${sectionId}"]`).classList.add('active');
    
    // 显示对应的内容区域
    document.getElementById(sectionId).classList.add('active');
    
    // 根据区域加载相应数据
    loadSectionData(sectionId);
}

// 加载区域数据
function loadSectionData(sectionId) {
    switch(sectionId) {
        case 'activity':
            loadUserActivity();
            break;
        case 'my-questions':
            loadUserQuestions();
            break;
        case 'my-answers':
            loadUserAnswers();
            break;
        case 'my-shares':
            loadUserShares();
            break;
        case 'my-resources':
            loadUserResources();
            break;
        case 'my-favorites':
            loadUserFavorites();
            break;
        case 'my-following':
            loadUserFollowing();
            break;
        case 'my-followers':
            loadUserFollowers();
            break;
        case 'messages':
            loadUserMessages();
            break;
    }
}

// 过滤个人动态
function filterActivity(filter) {
    // 移除所有过滤标签的active状态
    document.querySelectorAll('.filter-tab').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // 激活选中的过滤标签
    document.querySelector(`[data-filter="${filter}"]`).classList.add('active');
    
    // 重新加载动态数据
    loadUserActivity(filter);
}

// 加载用户动态
function loadUserActivity(filter = 'all') {
    const activityList = document.getElementById('activityList');
    if (!activityList) return;
    
    // 显示加载状态
    activityList.innerHTML = '<div class="loading">加载中...</div>';
    
    // 模拟API调用
    fetch(`/api/user/activity?filter=${filter}`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderActivityList(data.activities);
            } else {
                showEmptyState(activityList, '暂无动态', '您还没有任何活动记录');
            }
        })
        .catch(error => {
            console.error('加载动态失败:', error);
            showEmptyState(activityList, '加载失败', '无法加载动态数据');
        });
}

// 渲染动态列表
function renderActivityList(activities) {
    const activityList = document.getElementById('activityList');
    
    if (!activities || activities.length === 0) {
        showEmptyState(activityList, '暂无动态', '您还没有任何活动记录');
        return;
    }
    
    const html = activities.map(activity => `
        <div class="activity-item">
            <div class="item-header">
                <h4 class="item-title">${activity.title}</h4>
                <span class="activity-type ${activity.type}">${getActivityTypeText(activity.type)}</span>
            </div>
            <div class="item-content">${activity.content}</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(activity.created_at)}</span>
                <span><i class="fas fa-eye"></i> ${activity.views} 次浏览</span>
                <span><i class="fas fa-comment"></i> ${activity.comments} 条评论</span>
            </div>
        </div>
    `).join('');
    
    activityList.innerHTML = html;
}

// 获取活动类型文本
function getActivityTypeText(type) {
    const typeMap = {
        'question': '提问',
        'answer': '回答',
        'share': '分享',
        'comment': '评论',
        'resource': '资料'
    };
    return typeMap[type] || '活动';
}

// 加载用户提问
function loadUserQuestions() {
    const questionsList = document.getElementById('questionsList');
    if (!questionsList) return;
    
    questionsList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/questions')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderQuestionsList(data.questions);
            } else {
                showEmptyState(questionsList, '暂无提问', '您还没有发布过问题');
            }
        })
        .catch(error => {
            console.error('加载提问失败:', error);
            showEmptyState(questionsList, '加载失败', '无法加载提问数据');
        });
}

// 渲染提问列表
function renderQuestionsList(questions) {
    const questionsList = document.getElementById('questionsList');
    
    if (!questions || questions.length === 0) {
        showEmptyState(questionsList, '暂无提问', '您还没有发布过问题');
        return;
    }
    
    const html = questions.map(question => `
        <div class="question-item">
            <div class="item-header">
                <h4 class="item-title">${question.title}</h4>
                <span class="question-status ${question.status}">${getQuestionStatusText(question.status)}</span>
            </div>
            <div class="item-content">${question.content.substring(0, 100)}...</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(question.created_at)}</span>
                <span><i class="fas fa-comment"></i> ${question.answer_count} 个回答</span>
                <span><i class="fas fa-eye"></i> ${question.views} 次浏览</span>
            </div>
            <div class="item-actions">
                <button class="btn-edit" onclick="editQuestion(${question.id})">
                    <i class="fas fa-edit"></i> 编辑
                </button>
                <button class="btn-delete" onclick="deleteQuestion(${question.id})">
                    <i class="fas fa-trash"></i> 删除
                </button>
            </div>
        </div>
    `).join('');
    
    questionsList.innerHTML = html;
}

// 获取问题状态文本
function getQuestionStatusText(status) {
    const statusMap = {
        'open': '待解决',
        'answered': '已回答',
        'closed': '已关闭'
    };
    return statusMap[status] || '未知';
}

// 加载用户回答
function loadUserAnswers() {
    const answersList = document.getElementById('answersList');
    if (!answersList) return;
    
    answersList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/answers')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderAnswersList(data.answers);
            } else {
                showEmptyState(answersList, '暂无回答', '您还没有回答过问题');
            }
        })
        .catch(error => {
            console.error('加载回答失败:', error);
            showEmptyState(answersList, '加载失败', '无法加载回答数据');
        });
}

// 渲染回答列表
function renderAnswersList(answers) {
    const answersList = document.getElementById('answersList');
    
    if (!answers || answers.length === 0) {
        showEmptyState(answersList, '暂无回答', '您还没有回答过问题');
        return;
    }
    
    const html = answers.map(answer => `
        <div class="answer-item">
            <div class="item-header">
                <h4 class="item-title">回答：${answer.question_title}</h4>
                <span class="answer-status ${answer.is_accepted ? 'accepted' : ''}">
                    ${answer.is_accepted ? '已采纳' : '待采纳'}
                </span>
            </div>
            <div class="item-content">${answer.content.substring(0, 150)}...</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(answer.created_at)}</span>
                <span><i class="fas fa-thumbs-up"></i> ${answer.likes} 个赞</span>
                <span><i class="fas fa-comment"></i> ${answer.comments} 条评论</span>
            </div>
            <div class="item-actions">
                <button class="btn-edit" onclick="editAnswer(${answer.id})">
                    <i class="fas fa-edit"></i> 编辑
                </button>
                <button class="btn-delete" onclick="deleteAnswer(${answer.id})">
                    <i class="fas fa-trash"></i> 删除
                </button>
            </div>
        </div>
    `).join('');
    
    answersList.innerHTML = html;
}

// 加载用户分享
function loadUserShares() {
    const sharesList = document.getElementById('sharesList');
    if (!sharesList) return;
    
    sharesList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/shares')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderSharesList(data.shares);
            } else {
                showEmptyState(sharesList, '暂无分享', '您还没有发布过分享');
            }
        })
        .catch(error => {
            console.error('加载分享失败:', error);
            showEmptyState(sharesList, '加载失败', '无法加载分享数据');
        });
}

// 渲染分享列表
function renderSharesList(shares) {
    const sharesList = document.getElementById('sharesList');
    
    if (!shares || shares.length === 0) {
        showEmptyState(sharesList, '暂无分享', '您还没有发布过分享');
        return;
    }
    
    const html = shares.map(share => `
        <div class="share-item">
            <div class="item-header">
                <h4 class="item-title">${share.title}</h4>
                <span class="share-category">${share.category}</span>
            </div>
            <div class="item-content">${share.content.substring(0, 150)}...</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(share.created_at)}</span>
                <span><i class="fas fa-eye"></i> ${share.views} 次浏览</span>
                <span><i class="fas fa-heart"></i> ${share.likes} 个赞</span>
                <span><i class="fas fa-comment"></i> ${share.comments} 条评论</span>
            </div>
            <div class="item-actions">
                <button class="btn-edit" onclick="editShare(${share.id})">
                    <i class="fas fa-edit"></i> 编辑
                </button>
                <button class="btn-delete" onclick="deleteShare(${share.id})">
                    <i class="fas fa-trash"></i> 删除
                </button>
            </div>
        </div>
    `).join('');
    
    sharesList.innerHTML = html;
}

// 加载用户资料
function loadUserResources() {
    const resourcesList = document.getElementById('resourcesList');
    if (!resourcesList) return;
    
    resourcesList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/resources')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderResourcesList(data.resources);
            } else {
                showEmptyState(resourcesList, '暂无资料', '您还没有上传过资料');
            }
        })
        .catch(error => {
            console.error('加载资料失败:', error);
            showEmptyState(resourcesList, '加载失败', '无法加载资料数据');
        });
}

// 渲染资料列表
function renderResourcesList(resources) {
    const resourcesList = document.getElementById('resourcesList');
    
    if (!resources || resources.length === 0) {
        showEmptyState(resourcesList, '暂无资料', '您还没有上传过资料');
        return;
    }
    
    const html = resources.map(resource => `
        <div class="resource-item">
            <div class="item-header">
                <h4 class="item-title">${resource.title}</h4>
                <span class="resource-type">${resource.type}</span>
            </div>
            <div class="item-content">${resource.description}</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(resource.created_at)}</span>
                <span><i class="fas fa-download"></i> ${resource.downloads} 次下载</span>
                <span><i class="fas fa-eye"></i> ${resource.views} 次浏览</span>
                <span><i class="fas fa-file"></i> ${formatFileSize(resource.file_size)}</span>
            </div>
            <div class="item-actions">
                <button class="btn-edit" onclick="editResource(${resource.id})">
                    <i class="fas fa-edit"></i> 编辑
                </button>
                <button class="btn-delete" onclick="deleteResource(${resource.id})">
                    <i class="fas fa-trash"></i> 删除
                </button>
            </div>
        </div>
    `).join('');
    
    resourcesList.innerHTML = html;
}

// 加载用户收藏
function loadUserFavorites() {
    const favoritesList = document.getElementById('favoritesList');
    if (!favoritesList) return;
    
    favoritesList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/favorites')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderFavoritesList(data.favorites);
            } else {
                showEmptyState(favoritesList, '暂无收藏', '您还没有收藏任何内容');
            }
        })
        .catch(error => {
            console.error('加载收藏失败:', error);
            showEmptyState(favoritesList, '加载失败', '无法加载收藏数据');
        });
}

// 渲染收藏列表
function renderFavoritesList(favorites) {
    const favoritesList = document.getElementById('favoritesList');
    
    if (!favorites || favorites.length === 0) {
        showEmptyState(favoritesList, '暂无收藏', '您还没有收藏任何内容');
        return;
    }
    
    const html = favorites.map(favorite => `
        <div class="favorite-item">
            <div class="item-header">
                <h4 class="item-title">${favorite.title}</h4>
                <span class="favorite-type">${favorite.type}</span>
            </div>
            <div class="item-content">${favorite.content.substring(0, 100)}...</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(favorite.created_at)}</span>
                <span><i class="fas fa-user"></i> ${favorite.author}</span>
            </div>
            <div class="item-actions">
                <button class="btn-secondary" onclick="viewFavorite(${favorite.id})">
                    <i class="fas fa-eye"></i> 查看
                </button>
                <button class="btn-delete" onclick="removeFavorite(${favorite.id})">
                    <i class="fas fa-heart-broken"></i> 取消收藏
                </button>
            </div>
        </div>
    `).join('');
    
    favoritesList.innerHTML = html;
}

// 加载用户关注
function loadUserFollowing() {
    const followingList = document.getElementById('followingList');
    if (!followingList) return;
    
    followingList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/following')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderFollowingList(data.following);
            } else {
                showEmptyState(followingList, '暂无关注', '您还没有关注任何用户');
            }
        })
        .catch(error => {
            console.error('加载关注失败:', error);
            showEmptyState(followingList, '加载失败', '无法加载关注数据');
        });
}

// 渲染关注列表
function renderFollowingList(following) {
    const followingList = document.getElementById('followingList');
    
    if (!following || following.length === 0) {
        showEmptyState(followingList, '暂无关注', '您还没有关注任何用户');
        return;
    }
    
    const html = following.map(user => `
        <div class="following-item">
            <div class="item-header">
                <div class="user-info">
                    <img src="${user.avatar}" alt="${user.username}" class="user-avatar-small">
                    <div class="user-details">
                        <h4 class="item-title">${user.username}</h4>
                        <p class="user-bio">${user.bio || '这个人很懒，还没有写简介'}</p>
                    </div>
                </div>
                <span class="follow-status">已关注</span>
            </div>
            <div class="item-meta">
                <span><i class="fas fa-users"></i> ${user.followers} 粉丝</span>
                <span><i class="fas fa-question-circle"></i> ${user.questions} 提问</span>
                <span><i class="fas fa-comment"></i> ${user.answers} 回答</span>
            </div>
            <div class="item-actions">
                <button class="btn-secondary" onclick="viewUserProfile(${user.id})">
                    <i class="fas fa-user"></i> 查看资料
                </button>
                <button class="btn-delete" onclick="unfollowUser(${user.id})">
                    <i class="fas fa-user-minus"></i> 取消关注
                </button>
            </div>
        </div>
    `).join('');
    
    followingList.innerHTML = html;
}

// 加载用户粉丝
function loadUserFollowers() {
    const followersList = document.getElementById('followersList');
    if (!followersList) return;
    
    followersList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/followers')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderFollowersList(data.followers);
            } else {
                showEmptyState(followersList, '暂无粉丝', '您还没有粉丝');
            }
        })
        .catch(error => {
            console.error('加载粉丝失败:', error);
            showEmptyState(followersList, '加载失败', '无法加载粉丝数据');
        });
}

// 渲染粉丝列表
function renderFollowersList(followers) {
    const followersList = document.getElementById('followersList');
    
    if (!followers || followers.length === 0) {
        showEmptyState(followersList, '暂无粉丝', '您还没有粉丝');
        return;
    }
    
    const html = followers.map(user => `
        <div class="follower-item">
            <div class="item-header">
                <div class="user-info">
                    <img src="${user.avatar}" alt="${user.username}" class="user-avatar-small">
                    <div class="user-details">
                        <h4 class="item-title">${user.username}</h4>
                        <p class="user-bio">${user.bio || '这个人很懒，还没有写简介'}</p>
                    </div>
                </div>
                <span class="follow-status">关注了您</span>
            </div>
            <div class="item-meta">
                <span><i class="fas fa-users"></i> ${user.followers} 粉丝</span>
                <span><i class="fas fa-question-circle"></i> ${user.questions} 提问</span>
                <span><i class="fas fa-comment"></i> ${user.answers} 回答</span>
            </div>
            <div class="item-actions">
                <button class="btn-secondary" onclick="viewUserProfile(${user.id})">
                    <i class="fas fa-user"></i> 查看资料
                </button>
                <button class="btn-primary" onclick="followUser(${user.id})">
                    <i class="fas fa-user-plus"></i> 关注
                </button>
            </div>
        </div>
    `).join('');
    
    followersList.innerHTML = html;
}

// 加载用户消息
function loadUserMessages() {
    const messagesList = document.getElementById('messagesList');
    if (!messagesList) return;
    
    messagesList.innerHTML = '<div class="loading">加载中...</div>';
    
    fetch('/api/user/messages')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                renderMessagesList(data.messages);
            } else {
                showEmptyState(messagesList, '暂无消息', '您还没有收到任何消息');
            }
        })
        .catch(error => {
            console.error('加载消息失败:', error);
            showEmptyState(messagesList, '加载失败', '无法加载消息数据');
        });
}

// 渲染消息列表
function renderMessagesList(messages) {
    const messagesList = document.getElementById('messagesList');
    
    if (!messages || messages.length === 0) {
        showEmptyState(messagesList, '暂无消息', '您还没有收到任何消息');
        return;
    }
    
    const html = messages.map(message => `
        <div class="message-item ${message.is_read ? '' : 'unread'}">
            <div class="item-header">
                <h4 class="item-title">${message.title}</h4>
                <span class="message-type">${message.type}</span>
            </div>
            <div class="item-content">${message.content}</div>
            <div class="item-meta">
                <span><i class="fas fa-clock"></i> ${formatTime(message.created_at)}</span>
                <span><i class="fas fa-user"></i> ${message.sender}</span>
            </div>
            <div class="item-actions">
                <button class="btn-secondary" onclick="markMessageRead(${message.id})">
                    <i class="fas fa-check"></i> 标记已读
                </button>
                <button class="btn-delete" onclick="deleteMessage(${message.id})">
                    <i class="fas fa-trash"></i> 删除
                </button>
            </div>
        </div>
    `).join('');
    
    messagesList.innerHTML = html;
}

// 显示空状态
function showEmptyState(container, title, message) {
    container.innerHTML = `
        <div class="empty-state">
            <i class="fas fa-inbox"></i>
            <h3>${title}</h3>
            <p>${message}</p>
        </div>
    `;
}

// 头像上传功能
function changeAvatar(input) {
    if (input.files && input.files[0]) {
        const file = input.files[0];
        const formData = new FormData();
        formData.append('avatar', file);
        
        fetch('/api/user/avatar', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // 更新头像显示
                document.getElementById('userAvatar').src = data.avatar_url;
                // 更新顶部导航栏头像
                const topNavAvatar = document.querySelector('.user-avatar img');
                if (topNavAvatar) {
                    topNavAvatar.src = data.avatar_url;
                }
                showMessage('头像上传成功', 'success');
            } else {
                showMessage('头像上传失败', 'error');
            }
        })
        .catch(error => {
            console.error('头像上传失败:', error);
            showMessage('头像上传失败', 'error');
        });
    }
}

// 处理个人资料表单提交
function handleProfileSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData.entries());
    
    fetch('/api/user/profile', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('个人资料更新成功', 'success');
            // 更新页面显示的用户名
            document.querySelector('.username').textContent = data.user.username;
        } else {
            showMessage('个人资料更新失败', 'error');
        }
    })
    .catch(error => {
        console.error('更新个人资料失败:', error);
        showMessage('个人资料更新失败', 'error');
    });
}

// 重置表单
function resetForm() {
    document.getElementById('profileForm').reset();
}

// 修改密码
function changePassword() {
    const currentPassword = document.getElementById('currentPassword').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    
    if (!currentPassword || !newPassword || !confirmPassword) {
        showMessage('请填写完整的密码信息', 'error');
        return;
    }
    
    if (newPassword !== confirmPassword) {
        showMessage('新密码与确认密码不一致', 'error');
        return;
    }
    
    fetch('/api/user/password', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            current_password: currentPassword,
            new_password: newPassword
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('密码修改成功', 'success');
            // 清空密码字段
            document.getElementById('currentPassword').value = '';
            document.getElementById('newPassword').value = '';
            document.getElementById('confirmPassword').value = '';
        } else {
            showMessage(data.message || '密码修改失败', 'error');
        }
    })
    .catch(error => {
        console.error('修改密码失败:', error);
        showMessage('密码修改失败', 'error');
    });
}

// 保存通知设置
function saveNotificationSettings() {
    const checkboxes = document.querySelectorAll('.notification-options input[type="checkbox"]');
    const settings = {};
    
    checkboxes.forEach(checkbox => {
        settings[checkbox.name] = checkbox.checked;
    });
    
    fetch('/api/user/notifications', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(settings)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('通知设置保存成功', 'success');
        } else {
            showMessage('通知设置保存失败', 'error');
        }
    })
    .catch(error => {
        console.error('保存通知设置失败:', error);
        showMessage('通知设置保存失败', 'error');
    });
}

// 标记所有消息为已读
function markAllRead() {
    fetch('/api/user/messages/read-all', {
        method: 'PUT'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('已标记所有消息为已读', 'success');
            loadUserMessages(); // 重新加载消息列表
        } else {
            showMessage('操作失败', 'error');
        }
    })
    .catch(error => {
        console.error('标记已读失败:', error);
        showMessage('操作失败', 'error');
    });
}

// 工具函数
function formatTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    
    if (minutes < 1) return '刚刚';
    if (minutes < 60) return `${minutes}分钟前`;
    if (hours < 24) return `${hours}小时前`;
    if (days < 30) return `${days}天前`;
    
    return date.toLocaleDateString();
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function showMessage(message, type = 'info') {
    // 创建消息提示元素
    const messageDiv = document.createElement('div');
    messageDiv.className = `message-toast ${type}`;
    messageDiv.textContent = message;
    
    // 添加到页面
    document.body.appendChild(messageDiv);
    
    // 显示动画
    setTimeout(() => {
        messageDiv.classList.add('show');
    }, 100);
    
    // 自动隐藏
    setTimeout(() => {
        messageDiv.classList.remove('show');
        setTimeout(() => {
            document.body.removeChild(messageDiv);
        }, 300);
    }, 3000);
}

// 编辑和删除功能（需要根据实际API实现）
function editQuestion(id) {
    window.location.href = `/qa/edit/${id}`;
}

function deleteQuestion(id) {
    if (confirm('确定要删除这个问题吗？')) {
        fetch(`/api/questions/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('问题删除成功', 'success');
                loadUserQuestions();
            } else {
                showMessage('删除失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除问题失败:', error);
            showMessage('删除失败', 'error');
        });
    }
}

function editAnswer(id) {
    window.location.href = `/qa/answer/edit/${id}`;
}

function deleteAnswer(id) {
    if (confirm('确定要删除这个回答吗？')) {
        fetch(`/api/answers/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('回答删除成功', 'success');
                loadUserAnswers();
            } else {
                showMessage('删除失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除回答失败:', error);
            showMessage('删除失败', 'error');
        });
    }
}

function editShare(id) {
    window.location.href = `/category/1/edit/${id}`;
}

function deleteShare(id) {
    if (confirm('确定要删除这个分享吗？')) {
        fetch(`/api/shares/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('分享删除成功', 'success');
                loadUserShares();
            } else {
                showMessage('删除失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除分享失败:', error);
            showMessage('删除失败', 'error');
        });
    }
}

function editResource(id) {
    window.location.href = `/category/2/edit/${id}`;
}

function deleteResource(id) {
    if (confirm('确定要删除这个资料吗？')) {
        fetch(`/api/resources/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('资料删除成功', 'success');
                loadUserResources();
            } else {
                showMessage('删除失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除资料失败:', error);
            showMessage('删除失败', 'error');
        });
    }
}

function viewFavorite(id) {
    window.location.href = `/favorite/${id}`;
}

function removeFavorite(id) {
    if (confirm('确定要取消收藏吗？')) {
        fetch(`/api/favorites/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('已取消收藏', 'success');
                loadUserFavorites();
            } else {
                showMessage('操作失败', 'error');
            }
        })
        .catch(error => {
            console.error('取消收藏失败:', error);
            showMessage('操作失败', 'error');
        });
    }
}

function viewUserProfile(id) {
    window.location.href = `/user/${id}`;
}

function followUser(id) {
    fetch(`/api/users/${id}/follow`, {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('关注成功', 'success');
            loadUserFollowers();
        } else {
            showMessage('关注失败', 'error');
        }
    })
    .catch(error => {
        console.error('关注失败:', error);
        showMessage('关注失败', 'error');
    });
}

function unfollowUser(id) {
    if (confirm('确定要取消关注吗？')) {
        fetch(`/api/users/${id}/unfollow`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('已取消关注', 'success');
                loadUserFollowing();
            } else {
                showMessage('操作失败', 'error');
            }
        })
        .catch(error => {
            console.error('取消关注失败:', error);
            showMessage('操作失败', 'error');
        });
    }
}

function markMessageRead(id) {
    fetch(`/api/messages/${id}/read`, {
        method: 'PUT'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showMessage('已标记为已读', 'success');
            loadUserMessages();
        } else {
            showMessage('操作失败', 'error');
        }
    })
    .catch(error => {
        console.error('标记已读失败:', error);
        showMessage('操作失败', 'error');
    });
}

function deleteMessage(id) {
    if (confirm('确定要删除这条消息吗？')) {
        fetch(`/api/messages/${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage('消息删除成功', 'success');
                loadUserMessages();
            } else {
                showMessage('删除失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除消息失败:', error);
            showMessage('删除失败', 'error');
        });
    }
} 