import os
import json
from contextlib import nullcontext
import torch
from model import ModelArgs, Transformer
from chatglm_tokenizer.tokenization_chatglm import ChatGLMTokenizer
import numpy as np
import streamlit as st

@st.cache_data
def load_model():
    # 在这里加载你的模型 # retain only the top_k most likely tokens, clamp others to have 0 probability
    seed = 1337
    # device = 'cuda' if torch.cuda.is_available() else 'cpu'
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
        vocab_size=64793,  # 64793,
        multiple_of=multiple_of,
        max_seq_len=max_seq_len,
        dropout=dropout,
    )  # s
    torch.manual_seed(seed)
    torch.cuda.manual_seed(seed)
    torch.backends.cuda.matmul.allow_tf32 = True  # allow tf32 on matmul
    torch.backends.cudnn.allow_tf32 = True  # allow tf32 on cudnn
    device_type = 'cuda' if 'cuda' in device else 'cpu'  # for later use in torch.autocast
    ptdtype = {'float32': torch.float32, 'bfloat16': torch.bfloat16, 'float16': torch.float16}[dtype]
    ctx = nullcontext() if device_type == 'cpu' else torch.cuda.amp.autocast()

    ckpt_path = './out/epoch_1.pth'
    state_dict = torch.load(ckpt_path, map_location=device)
    gptconf = ModelArgs(**model_args)
    model = Transformer(gptconf)
    unwanted_prefix = '_orig_mod.'
    for k, v in list(state_dict.items()):
        if k.startswith(unwanted_prefix):
            state_dict[k[len(unwanted_prefix):]] = state_dict.pop(k)
    model.load_state_dict(state_dict, strict=False)

    model.eval()
    model.to(device)

    # load the tokenizer
    tokenizer = ChatGLMTokenizer(vocab_file='./chatglm_tokenizer/tokenizer.model')
    return model,tokenizer

def answer_question(user_question):
    # 在这里调用你的模型，并返回回答
    x = tokenizer.encode(user_question, add_special_tokens=False) + [tokenizer.special_tokens['<bos>']]
    x = (torch.tensor(x, dtype=torch.long, device=device)[None, ...])

    with torch.no_grad():
        with ctx:
            y = model.generate(x, 2, 100, temperature=1.0, top_k=30)
            answer=tokenizer.decode(y[0].tolist())
            answer = answer.replace(user_question, '')
    return answer

if __name__ == "__main__":

    out_dir = 'out'
    start = ""
    num_samples = 1  # number of samples to draw
    max_new_tokens = 100  # number of tokens generated in each sample
    temperature = 1.0  # 1.0 = no change, < 1.0 = less random, > 1.0 = more random, in predictions
    top_k = 30
    device = 'cuda' if torch.cuda.is_available() else 'cpu'

    device_type = 'cuda' if 'cuda' in device else 'cpu'
    ctx = nullcontext() if device_type == 'cpu' else torch.cuda.amp.autocast()

    model,tokenizer = load_model()

    st.title("Baby-llama")

    # 获取或创建 session_state
    session_state = st.session_state
    if not hasattr(session_state, "question_history"):
        session_state.question_history = []

    # 用户输入问题的文本框
    user_question = st.text_input("请输入你的问题:")

    # 在用户点击按钮后触发的事件
    if st.button("获取回答"):
        if user_question:
            # 调用模型回答问题
            model_answer = answer_question(user_question)
            # 显示模型的回答
            st.write(f"Baby-llama: {model_answer}")

            # 保存提问记录
            session_state.question_history.append(user_question)
        else:
            st.warning("请先输入问题再获取回答。")

    # 显示以往的提问记录
    if session_state.question_history:
        st.subheader("以往的提问记录:")
        for idx, question in enumerate(session_state.question_history[::-1]):
            st.write(f"{idx + 1}. {question}")