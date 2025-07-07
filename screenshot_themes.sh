#!/bin/bash
echo "开始为所有主题截图..."
echo "请手动访问 http://localhost:8080 并为每个主题截图"
echo ""

themes=("default" "dark" "blue" "green" "retro" "modern" "material" "minimal" "glass" "ocean" "forest" "sunset" "autumn" "winter" "spring" "summer" "cyberpunk" "neon" "matrix" "terminal" "space" "neon-blue" "neon-pink" "gradient" "monochrome" "arctic" "desert" "volcano" "galaxy" "vintage" "corporate" "paper" "bootstrap" "nature" "technology" "elegant")

for theme in "${themes[@]}"; do
    echo "切换到主题: $theme"
    ./servergo-screenshot config set theme $theme
    echo "请访问 http://localhost:8080 并截图保存为 data/imgs/$theme.png"
    read -p "完成后按回车继续下一个主题..."
done

echo "所有主题截图完成!"
