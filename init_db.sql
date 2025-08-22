-- AI论坛数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS aiforum CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE aiforum;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255) DEFAULT '/images/user.jpg',
    level INT DEFAULT 1,
    points INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    post_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 帖子表
CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category_id INT NOT NULL,
    user_id INT NOT NULL,
    view_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    tags VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 回复表
CREATE TABLE IF NOT EXISTS replies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 问题表
CREATE TABLE IF NOT EXISTS questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category_id INT NOT NULL,
    user_id INT NOT NULL,
    view_count INT DEFAULT 0,
    answer_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    tags VARCHAR(500),
    reward INT DEFAULT 0,
    is_solved BOOLEAN DEFAULT FALSE,
    summary TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 回答表
CREATE TABLE IF NOT EXISTS answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    is_accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 回答点赞表
CREATE TABLE IF NOT EXISTS answer_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    answer_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_like (answer_id, user_id),
    FOREIGN KEY (answer_id) REFERENCES answers(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 问题收藏表
CREATE TABLE IF NOT EXISTS question_favorites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_favorite (question_id, user_id),
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 问题举报表
CREATE TABLE IF NOT EXISTS question_reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_id INT NOT NULL,
    user_id INT NOT NULL,
    reason TEXT NOT NULL,
    status ENUM('pending', 'reviewed', 'resolved') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 技术文章表
CREATE TABLE IF NOT EXISTS tech_articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    category VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    cover_image VARCHAR(255),
    tags VARCHAR(500),
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    topic_slug VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 技术文章点赞表
CREATE TABLE IF NOT EXISTS tech_article_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_like (article_id, user_id),
    FOREIGN KEY (article_id) REFERENCES tech_articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 专题表
CREATE TABLE IF NOT EXISTS topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户关注表
CREATE TABLE IF NOT EXISTS user_follows (
    id INT AUTO_INCREMENT PRIMARY KEY,
    following_id INT NOT NULL,
    follower_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_follow (following_id, follower_id),
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章收藏表
CREATE TABLE IF NOT EXISTS article_favorites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_favorite (article_id, user_id),
    FOREIGN KEY (article_id) REFERENCES tech_articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章评论表
CREATE TABLE IF NOT EXISTS article_comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    parent_id INT NULL,
    like_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES tech_articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES article_comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 评论点赞表
CREATE TABLE IF NOT EXISTS comment_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    comment_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_comment_like (comment_id, user_id),
    FOREIGN KEY (comment_id) REFERENCES article_comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章举报表
CREATE TABLE IF NOT EXISTS article_reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    user_id INT NOT NULL,
    reason VARCHAR(50) NOT NULL,
    description TEXT,
    status ENUM('pending', 'resolved', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES tech_articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 学习资料表
CREATE TABLE IF NOT EXISTS learning_resources (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type ENUM('ebook', 'video', 'slides', 'dataset', 'code', 'paper') NOT NULL,
    level ENUM('beginner', 'intermediate', 'advanced') NOT NULL,
    category VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    cover_image VARCHAR(255),
    file_paths TEXT NOT NULL,
    total_size BIGINT DEFAULT 0,
    tags VARCHAR(500),
    rating DECIMAL(3,2) DEFAULT 0.00,
    download_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    download_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 资料分类表
CREATE TABLE IF NOT EXISTS resource_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 资料下载记录表
CREATE TABLE IF NOT EXISTS resource_downloads (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 资料评分表
CREATE TABLE IF NOT EXISTS resource_ratings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    user_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_rating (resource_id, user_id),
    FOREIGN KEY (resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 资料评论表
CREATE TABLE IF NOT EXISTS resource_comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认分类
INSERT IGNORE INTO categories (name, description) VALUES 
('知识问答', 'AI相关的技术问答'),
('技术分享', 'AI技术分享和经验交流'),
('学习资料', 'AI学习资源和教程'),
('社群交流', '社区活动和交流');

-- 插入默认标签
INSERT IGNORE INTO tags (name) VALUES 
('AI'),
('机器学习'),
('深度学习'),
('Python'),
('算法'),
('数据科学'),
('神经网络'),
('计算机视觉'),
('自然语言处理'),
('强化学习');

-- 创建索引
CREATE INDEX idx_posts_category ON posts(category_id);
CREATE INDEX idx_posts_user ON posts(user_id);
CREATE INDEX idx_posts_created ON posts(created_at);
CREATE INDEX idx_replies_post ON replies(post_id);
CREATE INDEX idx_replies_user ON replies(user_id);

-- 插入测试用户（密码: 123456）
INSERT IGNORE INTO users (username, email, password) VALUES 
('admin', 'admin@aiforum.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('testuser', 'test@aiforum.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi');

-- 插入测试帖子
INSERT IGNORE INTO posts (title, content, category_id, user_id, tags) VALUES 
('如何开始学习机器学习？', '作为一个初学者，我想了解如何系统地学习机器学习。请推荐一些学习路径和资源。', 1, 1, '机器学习,入门,学习路径'),
('深度学习在图像识别中的应用', '分享一些深度学习在计算机视觉领域的最新应用和研究成果。', 2, 2, '深度学习,计算机视觉,图像识别'),
('Python机器学习库对比', '对比scikit-learn、TensorFlow、PyTorch等主流机器学习库的优缺点。', 2, 1, 'Python,机器学习,库对比');

-- 插入测试回复
INSERT IGNORE INTO replies (post_id, user_id, content) VALUES 
(1, 2, '建议从吴恩达的机器学习课程开始，然后学习Python基础。'),
(1, 1, '感谢分享！我会按照这个路径学习。'),
(2, 1, '这篇文章很有帮助，特别是关于CNN的部分。');

-- 插入测试问题
INSERT IGNORE INTO questions (title, content, category_id, user_id, tags, reward, summary, is_solved) VALUES 
('如何开始学习机器学习？', '作为一个初学者，我想了解如何系统地学习机器学习。请推荐一些学习路径和资源。', 1, 1, '机器学习,入门,学习路径', 50, '作为一个初学者，我想了解如何系统地学习机器学习。请推荐一些学习路径和资源。', 1),
('深度学习在图像识别中的应用', '分享一些深度学习在计算机视觉领域的最新应用和研究成果。', 1, 2, '深度学习,计算机视觉,图像识别', 100, '分享一些深度学习在计算机视觉领域的最新应用和研究成果。', 0),
('Python机器学习库对比', '对比scikit-learn、TensorFlow、PyTorch等主流机器学习库的优缺点。', 1, 1, 'Python,机器学习,库对比', 30, '对比scikit-learn、TensorFlow、PyTorch等主流机器学习库的优缺点。', 0),
('神经网络训练技巧', '有哪些实用的神经网络训练技巧可以提高模型性能？', 1, 2, '神经网络,训练技巧,模型优化', 80, '有哪些实用的神经网络训练技巧可以提高模型性能？', 0),
('数据预处理的最佳实践', '在机器学习项目中，数据预处理有哪些最佳实践？', 1, 1, '数据预处理,机器学习,最佳实践', 40, '在机器学习项目中，数据预处理有哪些最佳实践？', 0);

-- 插入测试回答
INSERT IGNORE INTO answers (question_id, user_id, content, is_accepted) VALUES 
(1, 2, '建议从吴恩达的机器学习课程开始，然后学习Python基础。', 1),
(1, 1, '感谢分享！我会按照这个路径学习。', 0),
(2, 1, '这篇文章很有帮助，特别是关于CNN的部分。', 0),
(3, 2, 'scikit-learn适合传统机器学习，TensorFlow和PyTorch适合深度学习。', 0),
(4, 1, '使用批量归一化、dropout等技术可以提高训练效果。', 0);

-- 插入专题数据
INSERT IGNORE INTO topics (name, slug, description, icon) VALUES 
('机器学习入门', 'machine-learning', '机器学习基础知识和入门教程', 'fas fa-brain'),
('机器人控制技术', 'robotics', '机器人控制算法和技术实现', 'fas fa-robot'),
('深度学习', 'deep-learning', '深度学习理论和实践应用', 'fas fa-network-wired'),
('计算机视觉', 'computer-vision', '计算机视觉算法和应用', 'fas fa-eye'),
('自然语言处理', 'nlp', '自然语言处理技术', 'fas fa-language'),
('强化学习', 'reinforcement-learning', '强化学习理论和实践', 'fas fa-gamepad'),
('AI伦理与安全', 'ai-ethics', '人工智能伦理和安全问题', 'fas fa-balance-scale'),
('边缘AI', 'edge-ai', '边缘计算和AI应用', 'fas fa-microchip');

-- 插入资料分类数据
INSERT IGNORE INTO resource_categories (name, slug, description, icon) VALUES 
('机器学习', 'machine-learning', '机器学习相关学习资料', 'fas fa-brain'),
('深度学习', 'deep-learning', '深度学习相关学习资料', 'fas fa-network-wired'),
('机器人学', 'robotics', '机器人学相关学习资料', 'fas fa-robot'),
('计算机视觉', 'computer-vision', '计算机视觉相关学习资料', 'fas fa-eye'),
('自然语言处理', 'nlp', '自然语言处理相关学习资料', 'fas fa-language'),
('强化学习', 'reinforcement-learning', '强化学习相关学习资料', 'fas fa-gamepad'),
('数据科学', 'data-science', '数据科学相关学习资料', 'fas fa-chart-bar'),
('AI伦理', 'ai-ethics', 'AI伦理相关学习资料', 'fas fa-balance-scale');

-- 插入测试学习资料
INSERT IGNORE INTO learning_resources (title, description, type, level, category, user_id, cover_image, file_paths, total_size, tags, rating, download_count, comment_count) VALUES 
('机器学习实战：从理论到实践', '一本全面的机器学习入门书籍，涵盖监督学习、无监督学习、深度学习等核心概念，配有大量实战案例。', 'ebook', 'beginner', 'machine-learning', 1, '/images/ebook1.jpg', '/uploads/resources/ml_book.pdf', 15728640, '机器学习,入门,实战', 4.5, 156, 23),
('深度学习视频教程：神经网络详解', '深入浅出的深度学习视频教程，从基础概念到高级应用，包含完整的代码实现和项目实战。', 'video', 'intermediate', 'deep-learning', 2, '/images/video1.jpg', '/uploads/resources/dl_course.mp4', 524288000, '深度学习,神经网络,视频教程', 4.8, 89, 15),
('计算机视觉算法课件', '计算机视觉核心算法的详细课件，包含图像处理、特征提取、目标检测等内容。', 'slides', 'advanced', 'computer-vision', 1, '/images/slides1.jpg', '/uploads/resources/cv_slides.pptx', 20971520, '计算机视觉,算法,课件', 4.2, 67, 8),
('MNIST手写数字数据集', '经典的MNIST手写数字识别数据集，包含60000张训练图片和10000张测试图片。', 'dataset', 'beginner', 'machine-learning', 2, '/images/dataset1.jpg', '/uploads/resources/mnist.zip', 10485760, '数据集,MNIST,手写数字', 4.6, 234, 12),
('Python机器学习代码示例', '完整的Python机器学习代码示例，包含数据预处理、模型训练、评估等完整流程。', 'code', 'intermediate', 'machine-learning', 1, '/images/code1.jpg', '/uploads/resources/ml_code.zip', 5242880, 'Python,机器学习,代码', 4.4, 178, 19),
('Transformer论文解析', 'Transformer模型论文的详细解析，包含注意力机制、编码器-解码器架构等核心概念。', 'paper', 'advanced', 'nlp', 2, '/images/paper1.jpg', '/uploads/resources/transformer_paper.pdf', 2097152, 'Transformer,NLP,论文', 4.7, 95, 7),
('强化学习入门指南', '强化学习基础概念和算法的入门指南，包含Q-learning、策略梯度等经典算法。', 'ebook', 'beginner', 'reinforcement-learning', 1, '/images/ebook2.jpg', '/uploads/resources/rl_guide.pdf', 8388608, '强化学习,入门,Q-learning', 4.3, 123, 14),
('机器人控制算法实现', '机器人运动控制算法的完整实现，包含PID控制、路径规划、避障算法等。', 'code', 'advanced', 'robotics', 2, '/images/code2.jpg', '/uploads/resources/robot_control.zip', 15728640, '机器人,控制算法,PID', 4.5, 76, 11);

-- 更新帖子统计
UPDATE posts SET reply_count = (SELECT COUNT(*) FROM replies WHERE post_id = posts.id);
UPDATE categories SET post_count = (SELECT COUNT(*) FROM posts WHERE category_id = categories.id);

SELECT '数据库初始化完成！' as message; 