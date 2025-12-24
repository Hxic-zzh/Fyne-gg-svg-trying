
import matplotlib.pyplot as plt
import numpy as np
import matplotlib

# 设置中文字体（如果需要）
try:
    matplotlib.rcParams['font.sans-serif'] = ['Microsoft YaHei', 'SimHei']
    matplotlib.rcParams['axes.unicode_minus'] = False
    print("✓ 中文字体设置成功")
except:
    print("⚠ 使用默认字体")

# 数据准备
categories = ['FPS (帧率)', '内存 (MB)\n(越低越好)', 'CPU (%)\n(越低越好)']
N = len(categories)

# 原始数据
ts_values = [61.94, 125.69, 5.4]    # ToggleSwitch
fc_values = [61.9, 127.99, 5.48]     # FyneCheckbox

# 数据标准化处理（为了让雷达图显示更合理）
# FPS：越高越好，所以直接标准化到0-1
# 内存和CPU：越低越好，所以用1-值/最大值

# 计算最大值用于标准化
fps_max = max(ts_values[0], fc_values[0])
mem_max = max(ts_values[1], fc_values[1]) * 1.1  # 留10%空间
cpu_max = max(ts_values[2], fc_values[2]) * 1.1

# 标准化后的数据（0-1之间，1表示最好）
ts_normalized = [
    ts_values[0] / fps_max,                    # FPS越大越好
    1 - (ts_values[1] / mem_max),              # 内存越小越好
    1 - (ts_values[2] / cpu_max)               # CPU越小越好
]

fc_normalized = [
    fc_values[0] / fps_max,                    # FPS越大越好
    1 - (fc_values[1] / mem_max),              # 内存越小越好
    1 - (fc_values[2] / cpu_max)               # CPU越小越好
]

# 为了使雷达图闭合，需要重复第一个点
angles = np.linspace(0, 2 * np.pi, N, endpoint=False).tolist()
ts_normalized += ts_normalized[:1]
fc_normalized += fc_normalized[:1]
angles += angles[:1]

# 创建雷达图
fig = plt.figure(figsize=(10, 8))
ax = fig.add_subplot(111, polar=True)

# 绘制多边形
ax.plot(angles, ts_normalized, 'o-', linewidth=2, label='ToggleSwitch', 
        color='#2E86AB', markersize=8, markerfacecolor='white', markeredgewidth=2)
ax.fill(angles, ts_normalized, alpha=0.25, color='#2E86AB')

ax.plot(angles, fc_normalized, 'o-', linewidth=2, label='FyneCheckbox', 
        color='#A23B72', markersize=8, markerfacecolor='white', markeredgewidth=2)
ax.fill(angles, fc_normalized, alpha=0.25, color='#A23B72')

# 设置角度标签
ax.set_xticks(angles[:-1])
ax.set_xticklabels(categories, fontsize=12, fontweight='bold')

# 设置半径标签
ax.set_ylim(0, 1)
ax.set_yticks([0.2, 0.4, 0.6, 0.8, 1.0])
ax.set_yticklabels(['20%', '40%', '60%', '80%', '100%'], fontsize=10, color='gray')
ax.set_theta_offset(np.pi / 2)  # 从顶部开始
ax.set_theta_direction(-1)      # 顺时针方向

# 添加网格线
ax.grid(True, linestyle='--', alpha=0.5)

# 添加标题
plt.title('ToggleSwitch vs FyneCheckbox 性能雷达图\n（标准化对比）', 
          fontsize=14, fontweight='bold', pad=30)

# 添加图例
plt.legend(loc='upper right', bbox_to_anchor=(1.3, 1.1), fontsize=11)

# 在顶点处添加原始数值
def add_original_values(angles, values, normalized_values, color, offset):
    for i, (angle, val, norm_val) in enumerate(zip(angles[:-1], values, normalized_values[:-1])):
        # 计算标签位置
        x = norm_val * np.cos(angle)
        y = norm_val * np.sin(angle)
        
        # 调整标签位置避免重叠
        text_offset = 0.05
        ha = 'center'
        va = 'center'
        
        if angle < np.pi/6 or angle > 5*np.pi/6:
            x += text_offset if angle < np.pi/2 else -text_offset
        elif angle > np.pi/6 and angle < 5*np.pi/6:
            y += text_offset if angle < np.pi else -text_offset
        
        # 添加标签
        ax.text(x, y, f'{val:.2f}', 
                ha=ha, va=va, fontsize=9, fontweight='bold',
                bbox=dict(boxstyle="round,pad=0.3", facecolor=color, alpha=0.7, edgecolor='none'),
                color='white')

# 添加原始数值标签
add_original_values(angles[:-1], ts_values, ts_normalized, '#2E86AB', 0.03)
add_original_values(angles[:-1], fc_values, fc_normalized, '#A23B72', -0.03)

# 添加性能说明
plt.figtext(0.5, 0.02, 
            '说明：雷达图越接近外圈表示性能越好。FPS值越高越好，内存和CPU占用越低越好。',
            ha='center', fontsize=10, style='italic', color='gray')

plt.tight_layout()
plt.savefig('radar_chart_3metrics.png', dpi=300, bbox_inches='tight', facecolor='white')
plt.savefig('radar_chart_3metrics.pdf', bbox_inches='tight', facecolor='white')

print("="*50)
print("雷达图已生成：radar_chart_3metrics.png 和 .pdf")
print("="*50)
print("原始数据：")
print(f"ToggleSwitch: FPS={ts_values[0]}, 内存={ts_values[1]}MB, CPU={ts_values[2]}%")
print(f"FyneCheckbox: FPS={fc_values[0]}, 内存={fc_values[1]}MB, CPU={fc_values[2]}%")
print("="*50)

plt.show()
