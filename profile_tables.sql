-- 个人中心功能数据库表结构

-- 用户关注表
CREATE TABLE IF NOT EXISTS follows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(follower_id, followed_id)
);

-- 用户收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL, -- question, share, resource
    target_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户消息表
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL, -- system, question, follow, comment
    title VARCHAR(255) NOT NULL,
    content TEXT,
    sender VARCHAR(100),
    is_read BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 扩展用户表，添加个人中心相关字段
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN website VARCHAR(255);
ALTER TABLE users ADD COLUMN profile_public BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN show_email BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN show_phone BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN email_notifications BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN browser_notifications BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN question_notifications BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN follow_notifications BOOLEAN DEFAULT TRUE;

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_followed ON follows(followed_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_type_target ON favorites(type, target_id);
CREATE INDEX IF NOT EXISTS idx_messages_user ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_read ON messages(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_questions_user ON questions(user_id);
CREATE INDEX IF NOT EXISTS idx_answers_user ON answers(user_id);
CREATE INDEX IF NOT EXISTS idx_tech_articles_user ON tech_articles(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_resources_user ON learning_resources(user_id);

-- 插入一些示例数据
INSERT OR IGNORE INTO follows (follower_id, followed_id) VALUES (1, 2);
INSERT OR IGNORE INTO follows (follower_id, followed_id) VALUES (2, 1);

INSERT OR IGNORE INTO favorites (user_id, type, target_id) VALUES (1, 'question', 1);
INSERT OR IGNORE INTO favorites (user_id, type, target_id) VALUES (1, 'share', 1);

INSERT OR IGNORE INTO messages (user_id, type, title, content, sender) VALUES 
(1, 'system', '欢迎加入AI论坛', '感谢您注册AI论坛，开始您的学习之旅吧！', '系统'),
(1, 'question', '您的问题有了新回答', '您的问题"如何学习机器学习？"收到了新的回答', '用户2'),
(1, 'follow', '新粉丝关注', '用户3关注了您', '用户3'); 