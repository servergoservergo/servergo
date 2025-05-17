#!/bin/bash

# 创建目录结构
mkdir -p src/components
mkdir -p src/pages
mkdir -p src/styles
mkdir -p src/assets
mkdir -p public

# 创建配置文件
echo '{
  "name": "servergo-website",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-icons": "^5.0.1",
    "react-router-dom": "^6.22.3"
  },
  "devDependencies": {
    "@types/react": "^18.2.66",
    "@types/react-dom": "^18.2.22",
    "@typescript-eslint/eslint-plugin": "^7.2.0",
    "@typescript-eslint/parser": "^7.2.0",
    "@vitejs/plugin-react": "^4.2.1",
    "eslint": "^8.57.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.6",
    "typescript": "^5.2.2",
    "vite": "^5.2.0"
  }
}' > package.json

# 复制已创建的文件到正确的位置
if [ -f "src/App.tsx" ]; then
  cp src/App.tsx src/App.tsx.bak
fi

if [ -f "src/main.tsx" ]; then
  cp src/main.tsx src/main.tsx.bak
fi

if [ -f "src/styles/index.css" ]; then
  cp src/styles/index.css src/styles/index.css.bak
fi

# 创建或更新主要组件文件
echo "console.log('Setup script completed successfully!');"

echo "项目结构已创建完成！请运行以下命令来安装依赖：

npm install

然后，您可以启动开发服务器：

npm run dev
" 