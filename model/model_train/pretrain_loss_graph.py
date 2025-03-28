import re
import matplotlib.pyplot as plt

# 读取log文件
log_file_path = './out/log.log'
with open(log_file_path, 'r') as file:
    log_data = file.readlines()

# 提取批次和损失值
batch_pattern = re.compile(r'Epoch:\[(\d+)/(\d+)\]\((\d+)/(\d+)\) loss:(\d+\.\d+)')
batches = []
losses = []
batch_count = 0
batch_interval = 1000  # 设置每1000个batch记录一次损失值

for line in log_data:
    match = batch_pattern.search(line)
    if match:
        epoch_num = int(match.group(1))
        current_batch = int(match.group(3))
        total_batches = int(match.group(4))
        loss = float(match.group(5))
        batch_count += 100
        if batch_count % batch_interval == 0:
            batches.append(current_batch + (epoch_num - 1) * total_batches + (batch_count // batch_interval - 1) * total_batches)
            losses.append(loss)

# 绘制图表
plt.plot(batches, losses, label='Loss')
plt.title(f'Training Loss Over Batches (Every {batch_interval} Batches)')
plt.xlabel('Batch')
plt.ylabel('Loss')

# 设置纵坐标范围
plt.ylim(2, 10)

plt.legend()
plt.show()
