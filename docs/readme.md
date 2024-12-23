# 智能合约介绍
## 1.目录结构
```
python-test/                 # python测试代码
fabric-contract/             # 智能合约
├── comment-contract/           # 评论相关合约
│   └── comment_contract.go     # 评论合约实现
├── cti-contract/              # CTI(网络威胁情报)相关合约
│   └── cti_contract.go        # CTI合约实现
├── data-contract/             # 数据统计相关合约
│   └── data_contract.go       # 数据统计合约实现
├── incentive-contract/        # 激励机制相关合约
│   └── incentive_contract.go  # 激励合约实现
├── model-contract/            # 模型相关合约
│   └── model_contract.go      # 模型合约实现
├── msgstruct/                 # 消息结构定义
│   └── model.go              # 消息数据结构定义
├── reputation-contract/       # 信誉相关合约
│   └── reputation_contract.go # 信誉合约实现
├── typestruct/               # 类型结构定义
│   └── model.go              # 基础数据类型结构定义
├── user-contract/            # 用户相关合约
│   └── user_contract.go      # 用户合约实现
├── user-point-contract/      # 用户积分相关合约
│   └── user_point_contract.go # 用户积分合约实现
├── utils/                    # 工具函数
│   ├── ecc_encrypt.go        # ECC加密相关工具
│   └── signature.go          # 签名验证相关工具
├── vendor/                   # 第三方依赖
├── go.mod                    # Go模块定义
├── go.sum                    # Go依赖版本锁定
└── main_contract.go          # 主合约入口文件
```
