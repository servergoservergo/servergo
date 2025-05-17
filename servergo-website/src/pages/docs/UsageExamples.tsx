import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiCode, FiBookOpen, FiMonitor, FiLock, FiServer, FiGlobe, FiTerminal } from 'react-icons/fi'

export default function UsageExamples() {
  return (
    <section className="usage-examples-section">
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiBookOpen style={{ marginRight: '10px' }} /> 使用场景示例
      </h2>

      {/* 示例简介 */}
      <Card style={{ marginBottom: '24px' }}>
        <h3 style={{ marginBottom: '16px' }}>常见使用场景</h3>
        <p>
          ServerGo 是一个灵活的文件服务器，可以应用于多种场景。以下示例将帮助您了解如何在不同情境中使用 ServerGo：
        </p>
      </Card>

      {/* 基本文件分享 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiServer style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>快速文件分享</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景一：本地分享项目文件</h4>
          <p>当您需要与同事或客户快速分享项目文件时：</p>
          <CodeBlock code="# 在项目目录中启动服务器
cd /path/to/project
servergo" />
          <p>
            服务器将在自动探测的端口上启动，您可以将链接 <code>http://您的IP:端口</code> 分享给同事。
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景二：指定自定义端口</h4>
          <p>如果您希望使用特定端口：</p>
          <CodeBlock code="servergo -p 3000" />
          <p>
            服务器将在端口 3000 上启动。如果该端口不可用，系统会自动选择其他可用端口。
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景三：自动打开浏览器</h4>
          <p>启动服务器并自动在浏览器中打开（默认已启用）：</p>
          <CodeBlock code="servergo -o" />
          <p>
            服务器启动后会自动在默认浏览器中打开服务页面。
          </p>
        </div>
      </Card>

      {/* 前端开发 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiMonitor style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>前端开发环境</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景一：静态网站开发与测试</h4>
          <p>在开发静态网站时，您可以使用 ServerGo 快速预览：</p>
          <CodeBlock code="# 在包含 HTML/CSS/JS 文件的目录中启动
cd /path/to/website
servergo -o" />
          <p>
            每次修改代码后刷新浏览器即可查看最新效果。
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景二：与前端构建工具配合</h4>
          <p>如果您使用 Webpack、Vite 等构建工具，可以在构建输出目录启动 ServerGo：</p>
          <CodeBlock code="# 构建前端项目
npm run build

# 在构建输出目录启动服务器
cd dist
servergo" />
        </div>
      </Card>

      {/* 安全文件共享 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiLock style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>安全文件共享</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景一：共享敏感文档，需要密码保护</h4>
          <p>当您需要分享敏感或机密文档时：</p>
          <CodeBlock code="servergo -d /path/to/sensitive/docs -a basic -u client -w secure@123" />
          <p>
            接收者需要输入用户名和密码才能访问文件。
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景二：使用表单登录</h4>
          <p>使用更用户友好的登录表单：</p>
          <CodeBlock code="servergo -a basic -u admin -w password123 -l" />
          <p>
            这将提供一个网页登录表单，而不是浏览器原生的基本认证对话框。
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景三：使用令牌认证</h4>
          <p>使用令牌方式进行认证：</p>
          <CodeBlock code="servergo -a token -t your-secret-token" />
          <p>
            用户可以通过在URL中添加token参数访问：<code>http://localhost:8080/?token=your-secret-token</code>
          </p>
        </div>
      </Card>

      {/* 多语言支持 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiGlobe style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>多语言环境</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景一：查看并切换语言</h4>
          <p>查看可用的语言选项并切换：</p>
          <CodeBlock code="# 查看可用语言
servergo config set language

# 设置为中文界面
servergo config set language zh-CN

# 启动服务器
servergo" />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景二：设置英文界面</h4>
          <p>如果您的团队成员来自不同国家：</p>
          <CodeBlock code="servergo config set language en" />
          <p>
            这将使用英文作为界面语言，对国际团队更友好。
          </p>
        </div>
      </Card>

      {/* 系统集成 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiTerminal style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>系统集成</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景一：安装为系统服务</h4>
          <p>如果您需要 ServerGo 作为后台服务持续运行：</p>
          <CodeBlock code="# 安装为系统服务
servergo install

# 卸载服务
servergo uninstall" />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>场景二：在脚本中使用</h4>
          <p>在自动化脚本中使用 ServerGo：</p>
          <CodeBlock code={`#!/bin/bash
# 示例：构建并分享发布包

# 构建应用
echo "构建应用..."
make build

# 启动文件服务器分享构建结果
echo "启动文件服务器..."
servergo -d ./build -p 9000 -a basic -u release -w build@123 &
SERVER_PID=$!

# 获取本机 IP
IP=$(hostname -I | awk '{print $1}')
echo "文件可通过以下链接访问: http://$IP:9000"
echo "用户名: release"
echo "密码: build@123"

# 等待分享完成
read -p "按任意键停止服务器..." -n1 -s
kill $SERVER_PID`} />
        </div>
      </Card>

      {/* 高级配置示例 */}
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiCode style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>高级配置示例</h3>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>示例一：服务器配置组合</h4>
          <p>组合多种功能的完整示例：</p>
          <CodeBlock code={`servergo \\
  -d /path/to/share \\
  -p 8080 \\
  -a basic -u admin -w complex-password \\
  -m dark`} />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>示例二：内部文档站点</h4>
          <p>配置一个适合团队内部使用的文档站点：</p>
          <CodeBlock code={`# 首先设置持久化配置
servergo config set theme dark
servergo config set auto-open false
servergo config set enable-log-persistence true

# 启动文档服务器
servergo \\
  -d /path/to/docs \\
  -p 3000 \\
  -a form -u docs -w internal@docs \\
  -l`} />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>示例三：使用表单登录配合自定义主题</h4>
          <p>结合安全认证和美观界面：</p>
          <CodeBlock code={`# 创建临时目录
mkdir -p /tmp/project-files

# 启动服务器
servergo \\
  -d /tmp/project-files \\
  -p 7000 \\
  -a basic -u team -w project2023 \\
  -l \\
  -m blue`} />
        </div>
      </Card>
    </section>
  )
} 