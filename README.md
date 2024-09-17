my-webapp/
│
├── backend/               # Go后端代码
│   ├── cmd/               # 主应用程序入口
│   │   └── main.go        # 入口文件，启动Go服务器
│   ├── internal/          # 内部逻辑（业务逻辑、数据库交互等）
│   │   ├── handlers/      # API路由处理程序
│   │   │   └── api.go     # API处理程序示例
│   │   ├── models/        # 数据模型（与数据库交互）
│   │   │   └── user.go    # 数据库模型示例
│   │   ├── services/      # 服务层（封装复杂业务逻辑）
│   │   │   └── user.go    # 业务逻辑示例
│   │   └── database/      # 数据库交互
│   │       └── db.go      # 数据库连接和查询
│   ├── pkg/               # 公共库代码
│   │   └── utils.go       # 工具函数
│   ├── api/               # API接口定义
│   │   └── router.go      # API路由定义
│   ├── config/            # 配置文件
│   │   └── config.go      # 应用配置
│   ├── static/            # 前端编译后文件
│   └── go.mod             # Go模块依赖文件
│
├── frontend/              # React前端代码
│   ├── app/           
│   │   └── fonts    
│   ├── package.json       # 前端依赖和脚本
│   ├── package-lock.json
│   ├── webpack.config.js  # Webpack配置（如果使用）
│   └── .env               # 前端环境变量
│
├── docker-compose.yml      # Docker Compose文件，前后端容器化
├── Dockerfile.backend      # 后端Dockerfile
├── Dockerfile.frontend     # 前端Dockerfile
└── README.md               # 项目说明文档