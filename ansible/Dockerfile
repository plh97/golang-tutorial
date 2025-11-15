# Dockerfile
# 使用官方 Python 3.11 镜像作为基础，这是 Ansible 的推荐环境
FROM python:3.11-slim

# 设置工作目录
WORKDIR /app

# 确保 pip 和 setuptools 是最新的
RUN pip install --no-cache-dir --upgrade pip setuptools

# 明确安装 Ansible 8.x 的最新版本 (例如，使用 <9 来锁定 8 系列)
# 确保安装时包含 ssh 依赖
RUN pip install --no-cache-dir ansible==7.0.0

# 也可以安装其他你在 Ansible playbook 中需要的 Python 依赖，例如 boto3 (如果与 AWS 交互)
# RUN pip install --no-cache-dir boto3

# (可选) 确保 SSH 客户端存在，用于连接 EC2
RUN apt-get update && apt-get install -y openssh-client && rm -rf /var/lib/apt/lists/*

# 设置默认启动命令，尽管 GitHub Actions 会覆盖它，但这是个好习惯
CMD ["ansible", "--version"]
