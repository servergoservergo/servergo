# ServerGo 官方网站

这是 ServerGo 项目的官方网站代码库。ServerGo 是一个轻量级、高性能、易用的静态文件服务器，使用 Go 语言开发。

## 技术栈

- React 18
- TypeScript
- Vite
- React Router
- React Icons

## 开发指南

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

这将在 [http://localhost:3000](http://localhost:3000) 启动开发服务器。

### 构建生产版本

```bash
npm run build
```

构建后的文件将位于 `dist` 目录中。

### 预览构建结果

```bash
npm run preview
```

## 代码结构

- `src/assets` - 静态资源文件
- `src/components` - 可复用组件
- `src/pages` - 页面组件
- `src/styles` - 全局样式
- `public` - 公共资源文件

## 部署

构建后的网站是一个静态网站，可以部署到任何静态网站托管服务，如 GitHub Pages, Netlify, Vercel 等。

## License

[MIT](LICENSE)
