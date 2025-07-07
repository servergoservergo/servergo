#!/bin/bash

# 主题列表
themes=(
    "default" "dark" "blue" "green" "retro" 
    "modern" "material" "minimal" "glass" 
    "ocean" "forest" "sunset" "autumn" "winter" "spring" "summer"
    "cyberpunk" "neon" "matrix" "terminal" "space" 
    "neon-blue" "neon-pink" "gradient" "monochrome"
    "arctic" "desert" "volcano" "galaxy" 
    "vintage" "corporate" "paper" "bootstrap" "nature" "technology" "elegant"
)

echo "🎨 开始自动切换主题进行截图"
echo "📱 请打开浏览器访问: http://localhost:8080"
echo "📸 每个主题停留10秒供您截图"
echo ""

for i in "${!themes[@]}"; do
    theme=${themes[i]}
    echo "[$((i+1))/38] 切换到主题: $theme"
    ./servergo-screenshot config set theme $theme > /dev/null 2>&1
    echo "   🔗 当前访问: http://localhost:8080"
    echo "   💾 截图保存为: data/imgs/$theme.png"
    echo "   ⏱️  倒计时: 10秒..."
    
    for j in {10..1}; do
        echo -ne "   ⏰ $j秒后切换下一个主题...\r"
        sleep 1
    done
    echo ""
    echo ""
done

echo "🎉 所有38个主题展示完成!"
echo "📁 请检查 data/imgs/ 目录中的截图文件"
