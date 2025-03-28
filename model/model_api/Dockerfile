FROM nvidia/cuda:11.8.0-devel-ubuntu22.04 as builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    git \
    && rm -rf /var/lib/apt/lists/*

# 复制必要的文件
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt

# 复制应用文件
COPY . .

# 生成 gRPC 代码
RUN python3 -m pip install grpcio-tools
RUN python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. llm_service.proto

# 暴露端口
EXPOSE 8000
EXPOSE 50051

# 启动命令
CMD ["bash", "-c", "python3 grpc_server.py & python3 app.py"] 