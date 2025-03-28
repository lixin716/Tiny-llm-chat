import sys

path_add = r'D:\baby-llama\baby-llama2-chinese\model\model_train'
path_add2 = r'D:\baby-llama\baby-llama2-chinese\model'
sys.path.append(path_add)
sys.path.append(path_add2)

import os
from contextlib import nullcontext
import torch
from flask import Flask, request, jsonify
import json
from model import ModelArgs, Transformer
from chatglm_tokenizer.tokenization_chatglm import ChatGLMTokenizer

# 创建 Flask 应用
app = Flask(__name__)

# 全局变量来存储模型和分词器
model = None
tokenizer = None
device = None
ctx = None

# 加载模型函数
def load_model():
    global model, tokenizer, device, ctx
    
    # 设置设备
    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    device_type = 'cuda' if 'cuda' in device else 'cpu'
    
    # 设置上下文
    ctx = nullcontext() if device_type == 'cpu' else torch.cuda.amp.autocast()
    
    # 模型参数
    seed = 1337
    dtype = "float32"
    max_seq_len = 512
    dim = 1024
    n_layers = 12
    n_heads = 8
    multiple_of = 32
    dropout = 0.0
    
    model_args = dict(
        dim=dim,
        n_layers=n_layers,
        n_heads=n_heads,
        n_kv_heads=n_heads,
        vocab_size=64793,
        multiple_of=multiple_of,
        max_seq_len=max_seq_len,
        dropout=dropout,
    )
    
    # 设置随机种子
    torch.manual_seed(seed)
    torch.cuda.manual_seed(seed)
    torch.backends.cuda.matmul.allow_tf32 = True
    torch.backends.cudnn.allow_tf32 = True
    
    # 加载检查点
    ckpt_path = '../model_para/baby_llama.pth'
    state_dict = torch.load(ckpt_path, map_location=device)
    
    # 创建模型
    gptconf = ModelArgs(**model_args)
    model = Transformer(gptconf)
    
    # 处理状态字典中的前缀
    unwanted_prefix = '_orig_mod.'
    for k, v in list(state_dict.items()):
        if k.startswith(unwanted_prefix):
            state_dict[k[len(unwanted_prefix):]] = state_dict.pop(k)
    
    # 加载状态字典
    model.load_state_dict(state_dict, strict=False)
    
    # 设置为评估模式并移动到设备
    model.eval()
    model.to(device)
    
    # 加载分词器
    tokenizer = ChatGLMTokenizer(vocab_file='../chatglm_tokenizer/tokenizer.model')
    
    print("模型加载成功！")

# 生成回答函数
def generate_answer(prompt, temperature=1.0, max_new_tokens=100, top_k=30):
    x = tokenizer.encode(prompt, add_special_tokens=False) + [tokenizer.special_tokens['<bos>']]
    x = (torch.tensor(x, dtype=torch.long, device=device)[None, ...])
    
    with torch.no_grad():
        with ctx:
            y = model.generate(x, 2, max_new_tokens, temperature=temperature, top_k=top_k)
            answer = tokenizer.decode(y[0].tolist())
            answer = answer.replace(prompt, '')
    return answer

# 在启动时加载模型
@app.before_first_request
def before_first_request():
    load_model()

# API 端点
@app.route("/generate", methods=["POST"])
def generate():
    data = request.get_json()
    
    prompt = data.get("prompt", "")
    temperature = data.get("temperature", 1.0)
    max_new_tokens = data.get("max_new_tokens", 100)
    top_k = data.get("top_k", 30)
    
    response = generate_answer(
        prompt, 
        temperature=temperature,
        max_new_tokens=max_new_tokens,
        top_k=top_k
    )
    
    return jsonify({"response": response})

@app.route("/health", methods=["GET"])
def health_check():
    return jsonify({"status": "healthy"})

if __name__ == "__main__":
    # 直接启动 Flask 应用
    app.run(host="0.0.0.0", port=8000, debug=False) 