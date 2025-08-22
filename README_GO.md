# AI论坛 - Go语言版本

基于Go语言和MySQL数据库开发的AI学习和交流论坛系统。

## 🚀 功能特性

### 用户系统
- ✅ 用户注册和登录
- ✅ JWT身份认证
- ✅ 用户等级和积分系统
- ✅ 个人资料管理

### 内容管理
- ✅ 发布帖子和回复
- ✅ 分类管理
- ✅ 标签系统
- ✅ 搜索功能

### 问答系统
- ✅ 发布问题和回答
- ✅ 全文搜索和高级搜索
- ✅ 标签筛选和分类导航
- ✅ 悬赏积分机制
- ✅ 采纳回答功能
- ✅ 点赞和统计

### 论坛功能
- ✅ 帖子浏览和回复
- ✅ 热门帖子推荐
- ✅ 用户等级显示
- ✅ 积分奖励机制

## 🛠️ 技术栈

- **后端**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 8.0+
- **认证**: JWT
- **密码加密**: bcrypt
- **前端**: HTML5 + CSS3 + JavaScript

## 📁 项目结构

```
aiforum/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── config.env              # 环境配置
├── run.sh                  # 启动脚本
├── init_db.sql             # 数据库初始化脚本
├── README_GO.md            # 项目说明文档
├── config/                 # 配置管理
│   └── config.go
├── models/                 # 数据模型
│   ├── database.go         # 数据库初始化
│   ├── user.go            # 用户模型
│   ├── post.go            # 帖子模型
│   ├── reply.go           # 回复模型
│   ├── question.go        # 问题模型
│   ├── answer.go          # 回答模型
│   ├── category.go        # 分类模型
│   └── tag.go             # 标签模型
├── handlers/              # 请求处理器
│   ├── auth.go            # 认证相关
│   ├── post.go            # 帖子相关
│   ├── qa.go              # 问答相关
│   └── api.go             # API接口
├── middleware/            # 中间件
│   └── auth.go            # 认证中间件
├── utils/                 # 工具函数
│   ├── auth.go            # JWT工具
│   └── password.go        # 密码工具
├── templates/             # HTML模板
│   ├── layout.html        # 基础布局
│   ├── index.html         # 首页模板
│   ├── qa.html            # 问答页面模板
│   └── error.html         # 错误页面模板
├── static/                # 静态文件
│   ├── styles.css         # 样式文件
│   ├── qa.css             # 问答页面样式
│   └── script.js          # JavaScript文件
└── images/                # 图片资源
    ├── logo.png           # 网站Logo
    └── user.jpg           # 默认用户头像
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Git

### 2. 安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd aiforum

# 安装Go依赖
go mod tidy
```

### 3. 数据库配置

1. 创建MySQL数据库：
```sql
CREATE DATABASE aiforum CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `config.env`：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=aiforum
JWT_SECRET=your-secret-key
SERVER_PORT=8080
```

### 4. 运行项目

```bash
# 加载环境变量并运行
source config.env && go run main.go
```

或者使用godotenv（需要先安装）：
```bash
go install github.com/joho/godotenv/cmd/godotenv@latest
godotenv -f config.env go run main.go
```

### 5. 访问网站

打开浏览器访问：http://localhost:8080

## 📋 API接口

### 认证相关

- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `GET /auth/logout` - 用户登出

### 问答相关

- `GET /qa` - 问答页面
- `GET /qa/search` - 高级搜索
- `GET /qa/ask` - 提问页面
- `POST /qa/ask` - 发布问题
- `GET /qa/:id` - 查看问题详情
- `POST /qa/:id/answer` - 回答问题
- `POST /qa/answer/:answer_id/accept` - 采纳回答
- `POST /qa/answer/:answer_id/like` - 点赞回答

### 帖子相关

- `GET /` - 首页
- `GET /post/new` - 发帖页面
- `POST /post/new` - 创建帖子
- `GET /post/:id` - 查看帖子
- `POST /post/:id/reply` - 回复帖子

### API接口

- `GET /api/posts` - 获取帖子列表
- `GET /api/posts/:id` - 获取单个帖子
- `GET /api/categories` - 获取分类列表
- `GET /api/tags` - 获取标签列表
- `GET /search` - 搜索帖子

## 🔧 数据库表结构

### users (用户表)
- id: 用户ID
- username: 用户名
- email: 邮箱
- password: 密码（加密）
- avatar: 头像
- level: 等级
- points: 积分
- created_at: 创建时间
- updated_at: 更新时间

### categories (分类表)
- id: 分类ID
- name: 分类名称
- description: 分类描述
- post_count: 帖子数量
- created_at: 创建时间

### posts (帖子表)
- id: 帖子ID
- title: 标题
- content: 内容
- category_id: 分类ID
- user_id: 用户ID
- view_count: 浏览量
- reply_count: 回复数
- like_count: 点赞数
- tags: 标签
- created_at: 创建时间
- updated_at: 更新时间

### questions (问题表)
- id: 问题ID
- title: 标题
- content: 内容
- category_id: 分类ID
- user_id: 用户ID
- view_count: 浏览量
- answer_count: 回答数
- like_count: 点赞数
- tags: 标签
- reward: 悬赏积分
- is_solved: 是否已解决
- summary: 问题摘要
- created_at: 创建时间
- updated_at: 更新时间

### answers (回答表)
- id: 回答ID
- question_id: 问题ID
- user_id: 用户ID
- content: 回答内容
- like_count: 点赞数
- is_accepted: 是否被采纳
- created_at: 创建时间

### replies (回复表)
- id: 回复ID
- post_id: 帖子ID
- user_id: 用户ID
- content: 回复内容
- created_at: 创建时间

### tags (标签表)
- id: 标签ID
- name: 标签名称
- created_at: 创建时间

## 🎨 自定义配置

### 修改主题色彩
在 `static/styles.css` 中修改CSS变量：
```css
:root {
    --primary-color: #4A90E2;
    --secondary-color: #50C878;
    --accent-color: #FF6B35;
}
```

### 修改积分规则
在 `models/user.go` 中修改等级计算：
```go
func GetUserLevel(points int) int {
    if points >= 1000 {
        return 5
    } else if points >= 500 {
        return 4
    }
    // ...
}
```

## 🔒 安全特性

- JWT身份认证
- 密码bcrypt加密
- SQL注入防护
- XSS防护
- CSRF保护

## 📈 性能优化

- 数据库连接池
- 静态文件缓存
- 响应式设计
- 图片懒加载

## 🐛 故障排除

### 数据库连接失败
1. 检查MySQL服务是否启动
2. 验证数据库配置信息
3. 确认数据库用户权限

### 端口被占用
修改 `config.env` 中的 `SERVER_PORT` 配置

### 模板渲染错误
确保 `templates/` 目录下的模板文件存在且语法正确

## 📝 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 基础用户系统
- 帖子发布和回复
- 分类和标签系统
- 完整的问答系统

## 🤝 贡献指南

欢迎提交Issue和Pull Request来改进项目！

## 📄 许可证

MIT License

---

**AI论坛** - 让知识分享更简单，让学习交流更高效！ 