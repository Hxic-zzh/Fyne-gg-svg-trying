import pandas as pd
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt
import seaborn as sns
from pathlib import Path
import textwrap

# ==================== 设置 ====================
# 1. 确保CSV文件在同目录，文件名为：scientific_ToggleSwitch_vs_FyneCheckbox_20251221_123455.csv
csv_filename = "scientific_ToggleSwitch_vs_FyneCheckbox_20251221_123455.csv"

# 2. 检查文件是否存在
csv_path = Path(csv_filename)
if not csv_path.exists():
    print(f"❌ 错误: 找不到文件: {csv_filename}")
    print("请确保:")
    print("1. CSV文件在当前目录")
    print(f"2. 文件名正确: {csv_filename}")
    print("当前目录:", Path.cwd())
    exit()

# ==================== 第1步：读取数据 ====================
print(f"📁 正在读取数据: {csv_filename}")
try:
    # 读取CSV文件
    df = pd.read_csv(csv_filename)
    print(f"✅ 数据读取成功!")
    print(f"   数据形状: {df.shape[0]} 行 × {df.shape[1]} 列")
    print(f"   列名: {list(df.columns)}")
except Exception as e:
    print(f"❌ 读取文件失败: {e}")
    exit()

# 显示前几行数据
print("\n📊 数据预览:")
print(df.head())

# 检查组件名称
print(f"\n🔍 数据中的组件:")
print(df['component_name'].value_counts())

# ==================== 第2步：分离数据 ====================
# 检查列名是否有空格问题
print(f"\n🔧 数据处理...")
if 'component_name' not in df.columns:
    # 尝试修复可能的列名问题
    df.columns = df.columns.str.strip()
    print(f"   修复后的列名: {list(df.columns)}")

# 分离两组数据
toggle_data = df[df['component_name'] == 'ToggleSwitch'].copy()
native_data = df[df['component_name'] == 'FyneCheckbox'].copy()

print(f"   ToggleSwitch样本数: {len(toggle_data)}")
print(f"   FyneCheckbox样本数: {len(native_data)}")

# 检查数据是否为空
if len(toggle_data) == 0 or len(native_data) == 0:
    print("❌ 错误: 某个组件的数据为空!")
    print(f"   ToggleSwitch数据: {len(toggle_data)} 行")
    print(f"   FyneCheckbox数据: {len(native_data)} 行")
    exit()

# ==================== 第3步：小提琴图可视化 ====================
print(f"\n🎨 生成小提琴图...")

# 设置图表风格为灰度
plt.style.use('seaborn-v0_8-darkgrid')
sns.set_palette(["#FFFFFF", "#B0B0B0"])  # 白色和灰色

# 全局字体放大
plt.rcParams.update({'font.size': 18, 'axes.titlesize': 22, 'axes.labelsize': 20, 'legend.fontsize': 18, 'xtick.labelsize': 18, 'ytick.labelsize': 18})

# 创建三个子图：FPS、内存、CPU
fig, axes = plt.subplots(1, 3, figsize=(15, 5))  # 加宽图片
# fig.suptitle('Performance Distribution: TS vs FC', fontsize=24, fontweight='bold')

# 准备数据用于绘图
plot_data = []
for idx, component in enumerate(['ToggleSwitch', 'FyneCheckbox']):
    component_df = df[df['component_name'] == component]
    short = 'TS' if component == 'ToggleSwitch' else 'FC'
    for _, row in component_df.iterrows():
        plot_data.append({
            'Group': short,  # 改为Group，后续只显示TS/FC
            'FPS': row['fps'],
            'Memory (MB)': row['memory_usage_mb'],
            'CPU (%)': row['cpu_percent']
        })

plot_df = pd.DataFrame(plot_data)

# 1. FPS分布
sns.violinplot(x='Group', y='FPS', data=plot_df, ax=axes[0], inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[0].set_title('FPS Distribution', fontweight='bold', fontsize=22)
axes[0].set_xticklabels(['TS', 'FC'], fontsize=18)
axes[0].set_ylabel('Frames Per Second', fontsize=20)
axes[0].grid(True, alpha=0.3, linestyle=':')
axes[0].tick_params(axis='both', which='major', labelsize=18)

fps_means = plot_df.groupby('Group')['FPS'].mean()
for i, short in enumerate(['TS', 'FC']):
    axes[0].axhline(y=fps_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[0].text(i+0.4, fps_means[short], f'μ={fps_means[short]:.2f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 2. 内存分布
sns.violinplot(x='Group', y='Memory (MB)', data=plot_df, ax=axes[1], inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[1].set_title('Memory Usage Distribution', fontweight='bold', fontsize=22)
axes[1].set_xticklabels(['TS', 'FC'], fontsize=18)
axes[1].set_ylabel('Memory Usage (MB)', fontsize=20)
axes[1].grid(True, alpha=0.3, linestyle=':')
axes[1].tick_params(axis='both', which='major', labelsize=18)

mem_means = plot_df.groupby('Group')['Memory (MB)'].mean()
for i, short in enumerate(['TS', 'FC']):
    axes[1].axhline(y=mem_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[1].text(i+0.4, mem_means[short]+0.05, f'μ={mem_means[short]:.2f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 3. CPU分布
sns.violinplot(x='Group', y='CPU (%)', data=plot_df, ax=axes[2], inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[2].set_title('CPU Usage Distribution', fontweight='bold', fontsize=22)
axes[2].set_xticklabels(['TS', 'FC'], fontsize=18)
axes[2].set_ylabel('CPU Usage (%)', fontsize=20)
axes[2].grid(True, alpha=0.3, linestyle=':')
axes[2].tick_params(axis='both', which='major', labelsize=18)

cpu_means = plot_df.groupby('Group')['CPU (%)'].mean()
for i, short in enumerate(['TS', 'FC']):
    axes[2].axhline(y=cpu_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[2].text(i+0.4, cpu_means[short]+0.001, f'μ={cpu_means[short]:.4f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 在图下方统一加注说明
fig.text(0.5, -0.08, 'Abbreviations: TS=ToggleSwitch, FC=FyneCheckbox', ha='center', fontsize=16)

plt.tight_layout()

output_image = csv_path.stem + '_violin_plots.png'
plt.savefig(output_image, dpi=300, bbox_inches='tight')
print(f"✅ 小提琴图已保存: {output_image}")

plt.show()

# ==================== 第4步：统计显著性分析 ====================
print(f"\n📊 进行统计显著性分析...")

def check_normality(data, name="data"):
    """检查数据是否正态分布"""
    if len(data) < 3:
        return True  # 小样本假设正态
    
    try:
        stat, p = stats.shapiro(data)
        is_normal = p > 0.05
        print(f"   {name}: Shapiro-Wilk p={p:.4f}, {'normal' if is_normal else 'non-normal'}")
        return is_normal
    except Exception as e:
        print(f"   {name}: 正态性检验失败 - {e}")
        return True  # 默认假设正态

def perform_statistical_test(data1, data2, metric_name, data1_name="ToggleSwitch", data2_name="FyneCheckbox"):
    """执行完整的统计检验"""
    
    # 移除NaN值
    data1_clean = np.array(data1)[~np.isnan(data1)]
    data2_clean = np.array(data2)[~np.isnan(data2)]
    
    print(f"\n   [{metric_name}]")
    print(f"   {data1_name}: n={len(data1_clean)}, mean={np.mean(data1_clean):.4f}, std={np.std(data1_clean):.4f}")
    print(f"   {data2_name}: n={len(data2_clean)}, mean={np.mean(data2_clean):.4f}, std={np.std(data2_clean):.4f}")
    
    # 1. 正态性检验
    data1_normal = check_normality(data1_clean, f"{data1_name}_{metric_name}")
    data2_normal = check_normality(data2_clean, f"{data2_name}_{metric_name}")
    
    # 2. 选择检验方法
    if data1_normal and data2_normal:
        # 参数检验：独立样本t检验
        test_type = "Independent t-test"
        t_stat, p_value = stats.ttest_ind(data1_clean, data2_clean)
        
        # 效应量：Cohen's d
        n1, n2 = len(data1_clean), len(data2_clean)
        pooled_std = np.sqrt(((n1-1)*np.var(data1_clean) + (n2-1)*np.var(data2_clean)) / (n1+n2-2))
        mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
        cohens_d = mean_diff / pooled_std if pooled_std != 0 else 0
        
        effect_size = cohens_d
        effect_size_label = "Cohen's d"
        
    else:
        # 非参数检验：Mann-Whitney U检验
        test_type = "Mann-Whitney U test"
        u_stat, p_value = stats.mannwhitneyu(data1_clean, data2_clean)
        
        # 效应量：Cliff's delta (近似)
        mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
        pooled_std = np.std(np.concatenate([data1_clean, data2_clean]))
        effect_size = mean_diff / pooled_std if pooled_std != 0 else 0
        effect_size_label = "Standardized Mean Difference"
    
    # 3. 计算置信区间
    mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
    se_diff = np.sqrt(np.var(data1_clean)/len(data1_clean) + np.var(data2_clean)/len(data2_clean))
    ci_lower = mean_diff - 1.96 * se_diff
    ci_upper = mean_diff + 1.96 * se_diff
    
    # 4. 效应量解释
    if abs(effect_size) < 0.2:
        size_desc = "negligible"
    elif abs(effect_size) < 0.5:
        size_desc = "small"
    elif abs(effect_size) < 0.8:
        size_desc = "medium"
    else:
        size_desc = "large"
    
    # 5. 显著性判断
    significant = p_value < 0.05
    significance_desc = "SIGNIFICANT" if significant else "NOT SIGNIFICANT"
    
    return {
        'metric': metric_name,
        'test_type': test_type,
        'p_value': p_value,
        'mean_diff': mean_diff,
        'ci_95': (ci_lower, ci_upper),
        'effect_size': effect_size,
        'effect_size_label': effect_size_label,
        'effect_size_interpretation': size_desc,
        'significant': significant,
        'significance_desc': significance_desc,
        'data1_mean': np.mean(data1_clean),
        'data2_mean': np.mean(data2_clean),
        'data1_std': np.std(data1_clean),
        'data2_std': np.std(data2_clean),
        'n1': len(data1_clean),
        'n2': len(data2_clean)
    }

# 执行三个指标的检验
metrics_to_test = [
    ('fps', 'FPS'),
    ('memory_usage_mb', 'Memory Usage (MB)'),
    ('cpu_percent', 'CPU Usage (%)')
]

results = []
for col_name, display_name in metrics_to_test:
    if col_name in df.columns:
        result = perform_statistical_test(
            toggle_data[col_name].values,
            native_data[col_name].values,
            display_name
        )
        results.append(result)
    else:
        print(f"⚠️  警告: 列 '{col_name}' 不存在，跳过")

# ==================== 第5步：生成统计报告 ====================
print(f"\n" + "="*80)
print(" " * 20 + "STATISTICAL ANALYSIS REPORT")
print("="*80)

report_lines = []
for result in results:
    report_lines.append(f"\n{'='*60}")
    report_lines.append(f"METRIC: {result['metric']}")
    report_lines.append(f"{'='*60}")
    report_lines.append(f"Test Method: {result['test_type']}")
    report_lines.append(f"TS: mean={result['data1_mean']:.4f}, std={result['data1_std']:.4f}, n={result['n1']}")
    report_lines.append(f"FC: mean={result['data2_mean']:.4f}, std={result['data2_std']:.4f}, n={result['n2']}")
    report_lines.append(f"Mean Difference: {result['mean_diff']:.6f}")
    report_lines.append(f"95% Confidence Interval: [{result['ci_95'][0]:.6f}, {result['ci_95'][1]:.6f}]")
    report_lines.append(f"p-value: {result['p_value']:.10f} ({result['significance_desc']})")
    report_lines.append(f"Effect Size ({result['effect_size_label']}): {result['effect_size']:.4f} ({result['effect_size_interpretation']})")
    
    # 解释结果
    if result['significant']:
        direction = "lower" if result['mean_diff'] < 0 else "higher"
        percent_diff = abs(result['mean_diff'] / result['data2_mean'] * 100)
        report_lines.append(f"CONCLUSION: Statistically significant difference ({direction} by {percent_diff:.2f}%)")
    else:
        report_lines.append(f"CONCLUSION: No statistically significant difference")

# 打印报告
report_text = "\n".join(report_lines)
print(report_text)

# 保存报告到文件
report_filename = csv_path.stem + '_statistical_report.txt'
with open(report_filename, 'w', encoding='utf-8') as f:
    f.write(report_text)
print(f"\n✅ 统计报告已保存: {report_filename}")

# ==================== 第6步：生成汇总表格 ====================
print(f"\n📋 性能指标汇总:")
print("-" * 90)
print(f"{'Metric':<20} {'TS':<15} {'FC':<15} {'Difference':<15} {'p-value':<12} {'Significant':<10}")
print("-" * 90)

for result in results:
    toggle_val = f"{result['data1_mean']:.4f}"
    native_val = f"{result['data2_mean']:.4f}"
    diff_val = f"{result['mean_diff']:.4f}"
    p_val = f"{result['p_value']:.6f}"
    sig = "✓" if result['significant'] else "✗"
    
    print(f"{result['metric']:<20} {toggle_val:<15} {native_val:<15} {diff_val:<15} {p_val:<12} {sig:<10}")

print("-" * 90)

# ==================== 完成 ====================
print(f"\n🎉 分析完成!")
print(f"   1. 小提琴图: {output_image}")
print(f"   2. 统计报告: {report_filename}")
print(f"   3. 原始数据: {csv_filename}")