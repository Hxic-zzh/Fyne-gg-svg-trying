import pandas as pd
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt
import seaborn as sns
from pathlib import Path

# ==================== Configuration ====================
# CSV file path
csv_filename = "scientific_ParticleButton_vs_FyneButton_20251222_101855.csv"

# Check if file exists
csv_path = Path(csv_filename)
if not csv_path.exists():
    print(f"❌ Error: File not found: {csv_filename}")
    print("Please ensure:")
    print(f"1. CSV file is in current directory")
    print(f"2. File name is correct: {csv_filename}")
    print("Current directory:", Path.cwd())
    exit()

# ==================== Step 1: Load Data ====================
print(f"📁 Loading Linux test data: {csv_filename}")
try:
    df = pd.read_csv(csv_filename)
    print(f"✅ Data loaded successfully!")
    print(f"   Data shape: {df.shape[0]} rows × {df.shape[1]} columns")
    print(f"   System: Linux, CPU cores: 32")
except Exception as e:
    print(f"❌ Failed to read file: {e}")
    exit()

print("\n📊 Data preview:")
print(df.head())

print(f"\n🔍 Components in data:")
print(df['component_name'].value_counts())

# ==================== Step 2: Separate Data ====================
print(f"\n🔧 Processing data...")

# Clean column names if needed
if 'component_name' not in df.columns:
    df.columns = df.columns.str.strip()
    print(f"   Fixed column names: {list(df.columns)}")

# Separate data for two components
pb_data = df[df['component_name'] == 'ParticleButton'].copy()
fb_data = df[df['component_name'] == 'FyneButton'].copy()

print(f"   ParticleButton samples: {len(pb_data)}")
print(f"   FyneButton samples: {len(fb_data)}")

if len(pb_data) == 0 or len(fb_data) == 0:
    print("❌ Error: Empty data for one component!")
    print(f"   ParticleButton data: {len(pb_data)} rows")
    print(f"   FyneButton data: {len(fb_data)} rows")
    exit()

# ==================== Step 3: Violin Plot Visualization ====================
print(f"\n🎨 Generating violin plots...")

# Set style to grayscale
plt.style.use('seaborn-v0_8-darkgrid')
sns.set_palette(["#FFFFFF", "#B0B0B0"])  # White and gray

# Create figure with 3 subplots
fig, axes = plt.subplots(1, 3, figsize=(15, 5))
fig.suptitle('Performance Distribution: ParticleButton vs FyneButton (Linux)', 
             fontsize=16, fontweight='bold')

# Prepare data for plotting
plot_data = []
for component in ['ParticleButton', 'FyneButton']:
    component_df = df[df['component_name'] == component]
    short = 'PB' if component == 'ParticleButton' else 'FB'
    for _, row in component_df.iterrows():
        plot_data.append({
            'Group': short,  # 只显示PB/FB
            'FPS': row['fps'],
            'Memory (MB)': row['memory_usage_mb'],
            'CPU (%)': row['cpu_percent']
        })

plot_df = pd.DataFrame(plot_data)

# 1. FPS Distribution
sns.violinplot(x='Group', y='FPS', data=plot_df, ax=axes[0], 
               inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[0].set_title('FPS Distribution', fontweight='bold', fontsize=22)
axes[0].set_xticklabels(['PB', 'FB'], fontsize=18)
axes[0].set_ylabel('Frames Per Second', fontsize=20)
axes[0].grid(True, alpha=0.3, linestyle=':')
axes[0].tick_params(axis='both', which='major', labelsize=18)

fps_means = plot_df.groupby('Group')['FPS'].mean()
for i, short in enumerate(['PB', 'FB']):
    axes[0].axhline(y=fps_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[0].text(i+0.4, fps_means[short], f'μ={fps_means[short]:.2f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 2. Memory Distribution
sns.violinplot(x='Group', y='Memory (MB)', data=plot_df, ax=axes[1], 
               inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[1].set_title('Memory Usage Distribution', fontweight='bold', fontsize=22)
axes[1].set_xticklabels(['PB', 'FB'], fontsize=18)
axes[1].set_ylabel('Memory Usage (MB)', fontsize=20)
axes[1].grid(True, alpha=0.3, linestyle=':')
axes[1].tick_params(axis='both', which='major', labelsize=18)

mem_means = plot_df.groupby('Group')['Memory (MB)'].mean()
for i, short in enumerate(['PB', 'FB']):
    axes[1].axhline(y=mem_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[1].text(i+0.4, mem_means[short]+0.05, f'μ={mem_means[short]:.2f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 3. CPU Distribution
sns.violinplot(x='Group', y='CPU (%)', data=plot_df, ax=axes[2], 
               inner='quartile', palette=["#FFFFFF", "#B0B0B0"])
axes[2].set_title('CPU Usage Distribution', fontweight='bold', fontsize=22)
axes[2].set_xticklabels(['PB', 'FB'], fontsize=18)
axes[2].set_ylabel('CPU Usage (%)', fontsize=20)
axes[2].grid(True, alpha=0.3, linestyle=':')
axes[2].tick_params(axis='both', which='major', labelsize=18)

cpu_means = plot_df.groupby('Group')['CPU (%)'].mean()
for i, short in enumerate(['PB', 'FB']):
    axes[2].axhline(y=cpu_means[short], color='black', linestyle='--', 
                    alpha=0.7, linewidth=1.5)
    axes[2].text(i+0.4, cpu_means[short]+0.001, f'μ={cpu_means[short]:.4f}', 
                 fontsize=15, fontweight='bold', color='black',
                 bbox=dict(boxstyle='round,pad=0.3', facecolor='#E0E0E0', edgecolor='black', alpha=0.7))

# 在图下方统一加注说明
fig.text(0.5, -0.08, 'Abbreviations: PB=ParticleButton, FB=FyneButton', ha='center', fontsize=16)

plt.tight_layout()

# Save image
output_image = 'linux_violin_plots.png'
plt.savefig(output_image, dpi=300, bbox_inches='tight')
print(f"✅ Violin plots saved: {output_image}")

plt.show()

# ==================== Step 4: Statistical Significance Analysis ====================
print(f"\n📊 Performing statistical significance analysis...")

def check_normality(data, name="data"):
    """Check if data is normally distributed"""
    if len(data) < 3:
        return True  # Small sample assume normal
    
    try:
        stat, p = stats.shapiro(data)
        is_normal = p > 0.05
        print(f"   {name}: Shapiro-Wilk p={p:.4f}, {'normal' if is_normal else 'non-normal'}")
        return is_normal
    except Exception as e:
        print(f"   {name}: Normality test failed - {e}")
        return True  # Default assume normal

def perform_statistical_test(data1, data2, metric_name, data1_name="ParticleButton", data2_name="FyneButton"):
    """Perform complete statistical testing"""
    
    # Remove NaN values
    data1_clean = np.array(data1)[~np.isnan(data1)]
    data2_clean = np.array(data2)[~np.isnan(data2)]
    
    print(f"\n   [{metric_name}]")
    print(f"   {data1_name}: n={len(data1_clean)}, mean={np.mean(data1_clean):.4f}, std={np.std(data1_clean):.4f}")
    print(f"   {data2_name}: n={len(data2_clean)}, mean={np.mean(data2_clean):.4f}, std={np.std(data2_clean):.4f}")
    
    # 1. Normality test
    data1_normal = check_normality(data1_clean, f"{data1_name}_{metric_name}")
    data2_normal = check_normality(data2_clean, f"{data2_name}_{metric_name}")
    
    # 2. Choose test method
    if data1_normal and data2_normal:
        # Parametric test: Independent t-test
        test_type = "Independent t-test"
        t_stat, p_value = stats.ttest_ind(data1_clean, data2_clean)
        
        # Effect size: Cohen's d
        n1, n2 = len(data1_clean), len(data2_clean)
        pooled_std = np.sqrt(((n1-1)*np.var(data1_clean) + (n2-1)*np.var(data2_clean)) / (n1+n2-2))
        mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
        cohens_d = mean_diff / pooled_std if pooled_std != 0 else 0
        
        effect_size = cohens_d
        effect_size_label = "Cohen's d"
        
    else:
        # Non-parametric test: Mann-Whitney U test
        test_type = "Mann-Whitney U test"
        u_stat, p_value = stats.mannwhitneyu(data1_clean, data2_clean)
        
        # Effect size approximation
        mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
        pooled_std = np.std(np.concatenate([data1_clean, data2_clean]))
        effect_size = mean_diff / pooled_std if pooled_std != 0 else 0
        effect_size_label = "Standardized Mean Difference"
    
    # 3. Calculate confidence interval
    mean_diff = np.mean(data1_clean) - np.mean(data2_clean)
    se_diff = np.sqrt(np.var(data1_clean)/len(data1_clean) + np.var(data2_clean)/len(data2_clean))
    ci_lower = mean_diff - 1.96 * se_diff
    ci_upper = mean_diff + 1.96 * se_diff
    
    # 4. Effect size interpretation
    if abs(effect_size) < 0.2:
        size_desc = "negligible"
    elif abs(effect_size) < 0.5:
        size_desc = "small"
    elif abs(effect_size) < 0.8:
        size_desc = "medium"
    else:
        size_desc = "large"
    
    # 5. Significance judgment
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
        'pb_mean': np.mean(data1_clean),
        'fb_mean': np.mean(data2_clean),
        'pb_std': np.std(data1_clean),
        'fb_std': np.std(data2_clean),
        'n_pb': len(data1_clean),
        'n_fb': len(data2_clean)
    }

# Test three metrics
metrics_to_test = [
    ('fps', 'FPS'),
    ('memory_usage_mb', 'Memory Usage (MB)'),
    ('cpu_percent', 'CPU Usage (%)')
]

results = []
for col_name, display_name in metrics_to_test:
    if col_name in df.columns:
        result = perform_statistical_test(
            pb_data[col_name].values,
            fb_data[col_name].values,
            display_name
        )
        results.append(result)
    else:
        print(f"⚠️  Warning: Column '{col_name}' not found, skipping")

# ==================== Step 5: Generate Statistical Report ====================
print(f"\n" + "="*80)
print(" " * 20 + "STATISTICAL ANALYSIS REPORT")
print("="*80)

report_lines = []
report_lines.append("="*80)
report_lines.append("PERFORMANCE COMPARISON: ParticleButton vs FyneButton")
report_lines.append("TEST ENVIRONMENT: Linux, 32 CPU cores")
report_lines.append("="*80)

for result in results:
    report_lines.append(f"\n{'='*60}")
    report_lines.append(f"METRIC: {result['metric']}")
    report_lines.append(f"{'='*60}")
    report_lines.append(f"Test Method: {result['test_type']}")
    report_lines.append(f"ParticleButton: mean={result['pb_mean']:.4f}, std={result['pb_std']:.4f}, n={result['n_pb']}")
    report_lines.append(f"FyneButton: mean={result['fb_mean']:.4f}, std={result['fb_std']:.4f}, n={result['n_fb']}")
    report_lines.append(f"Mean Difference: {result['mean_diff']:.6f}")
    report_lines.append(f"95% Confidence Interval: [{result['ci_95'][0]:.6f}, {result['ci_95'][1]:.6f}]")
    report_lines.append(f"p-value: {result['p_value']:.10f} ({result['significance_desc']})")
    report_lines.append(f"Effect Size ({result['effect_size_label']}): {result['effect_size']:.4f} ({result['effect_size_interpretation']})")
    
    # Interpret results
    if result['significant']:
        direction = "lower" if result['mean_diff'] < 0 else "higher"
        percent_diff = abs(result['mean_diff'] / result['fb_mean'] * 100)
        report_lines.append(f"CONCLUSION: Statistically significant difference (ParticleButton is {direction} by {percent_diff:.2f}%)")
    else:
        report_lines.append(f"CONCLUSION: No statistically significant difference")

# Print report
report_text = "\n".join(report_lines)
print(report_text)

# Save report to file
report_filename = 'linux_statistical_report.txt'
with open(report_filename, 'w', encoding='utf-8') as f:
    f.write(report_text)
print(f"\n✅ Statistical report saved: {report_filename}")

# ==================== Step 6: Summary Table ====================
print(f"\n📋 Performance Metrics Summary:")
print("-" * 100)
print(f"{'Metric':<20} {'ParticleButton':<15} {'FyneButton':<15} {'Difference':<15} {'p-value':<12} {'Significant':<12}")
print("-" * 100)

for result in results:
    pb_val = f"{result['pb_mean']:.4f}"
    fb_val = f"{result['fb_mean']:.4f}"
    diff_val = f"{result['mean_diff']:.4f}"
    p_val = f"{result['p_value']:.6f}"
    sig = "YES" if result['significant'] else "NO"
    
    print(f"{result['metric']:<20} {pb_val:<15} {fb_val:<15} {diff_val:<15} {p_val:<12} {sig:<12}")

print("-" * 100)

# ==================== Completion ====================
print(f"\n🎉 Analysis completed!")
print(f"   1. Violin plots: {output_image}")
print(f"   2. Statistical report: {report_filename}")
print(f"   3. Original data: {csv_filename}")