import sys
path_add = r'D:\baby-llama\baby-llama2-chinese\model\model_train'
path_add2 = r'D:\baby-llama\baby-llama2-chinese\model'
sys.path.append(path_add)
sys.path.append(path_add2)
import os
import grpc
from concurrent import futures
import torch
from contextlib import nullcontext
from model import ModelArgs, Transformer
from chatglm_tokenizer.tokenization_chatglm import ChatGLMTokenizer

# 导入生成的 gRPC 代码
import llm_service_pb2
import llm_service_pb2_grpc

# 模型服务实现
class LLMServicer(llm_service_pb2_grpc.LLMServiceServicer):
    def __init__(self):
        # 设置设备
        self.device = 'cuda' if torch.cuda.is_available() else 'cpu'
        self.device_type = 'cuda' if 'cuda' in self.device else 'cpu'
        
        # 设置上下文
        self.ctx = nullcontext() if self.device_type == 'cpu' else torch.cuda.amp.autocast()
        
        # 加载模型
        self.load_model()
    
    def load_model(self):
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
        state_dict = torch.load(ckpt_path, map_location=self.device)
        
        # 创建模型
        gptconf = ModelArgs(**model_args)
        self.model = Transformer(gptconf)
        
        # 处理状态字典中的前缀
        unwanted_prefix = '_orig_mod.'
        for k, v in list(state_dict.items()):
            if k.startswith(unwanted_prefix):
                state_dict[k[len(unwanted_prefix):]] = state_dict.pop(k)
        
        # 加载状态字典
        self.model.load_state_dict(state_dict, strict=False)
        
        # 设置为评估模式并移动到设备
        self.model.eval()
        self.model.to(self.device)
        
        # 加载分词器
        self.tokenizer = ChatGLMTokenizer(vocab_file='../chatglm_tokenizer/tokenizer.model')
        
        print("gRPC服务器中模型加载成功！")
    
    def Generate(self, request, context):
        prompt = request.prompt
        temperature = request.temperature
        max_new_tokens = request.max_new_tokens
        top_k = request.top_k
        
        # 生成回答
        x = self.tokenizer.encode(prompt, add_special_tokens=False) + [self.tokenizer.special_tokens['<bos>']]
        x = (torch.tensor(x, dtype=torch.long, device=self.device)[None, ...])
        
        with torch.no_grad():
            with self.ctx:
                y = self.model.generate(x, 2, max_new_tokens, temperature=temperature, top_k=top_k)
                answer = self.tokenizer.decode(y[0].tolist())
                answer = answer.replace(prompt, '')
        
        return llm_service_pb2.GenerateResponse(response=answer)

def serve():
    # 创建 gRPC 服务器
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    llm_service_pb2_grpc.add_LLMServiceServicer_to_server(LLMServicer(), server)
    
    # 监听端口
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC服务器已启动，监听端口50051...")
    
    # 保持服务器运行
    server.wait_for_termination()

if __name__ == '__main__':
    serve() 