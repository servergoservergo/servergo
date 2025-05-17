import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiTerminal, FiDownload, FiExternalLink, FiGithub, FiCommand, FiCpu } from 'react-icons/fi'

export default function GettingStarted() {
  return (
    <section style={{ marginBottom: '40px' }}>
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiTerminal style={{ marginRight: '10px' }} /> 入门指南
      </h2>

      {/* 简介 */}
      <Card style={{ marginBottom: '24px' }}>
        <h3 style={{ marginBottom: '16px' }}>ServerGo 简介</h3>
        <p>
          ServerGo 是一个高性能的静态文件服务器，使用 Go 语言开发，旨在提供快速、简单且功能丰富的文件共享解决方案。
          您可以使用它来快速创建一个 HTTP 服务来共享文件，类似于 Python 的 <code>http.server</code> 模块，但性能更好，功能更丰富。
        </p>
        <div style={{ marginTop: '16px' }}>
          <h4 style={{ fontSize: '18px', marginBottom: '8px' }}>主要特点：</h4>
          <ul style={{ paddingLeft: '20px' }}>
            <li>零配置启动，快速分享文件</li>
            <li>跨平台支持（Windows、macOS、Linux）</li>
            <li>高性能文件传输</li>
            <li>丰富的安全选项（基本认证、令牌认证、表单登录）</li>
            <li>美观的文件浏览界面</li>
            <li>支持多种主题</li>
            <li>提供多语言国际化支持</li>
          </ul>
        </div>
      </Card>

      {/* 安装部分 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          marginBottom: '16px' 
        }}>
          <FiDownload style={{ fontSize: '20px', marginRight: '8px' }} />
          <h3 style={{ margin: 0 }}>安装 ServerGo</h3>
        </div>

        <div style={{ marginBottom: '20px' }}>
          <h4>下载预编译二进制文件</h4>
          <p>从 GitHub 发布页面下载适合您操作系统的二进制文件：</p>
          <div style={{ marginTop: '8px' }}>
            <a 
              href="https://github.com/cc11001100/servergo/releases" 
              target="_blank" 
              rel="noopener noreferrer"
              className="btn"
              style={{ display: 'inline-flex', alignItems: 'center', gap: '8px' }}
            >
              <FiGithub /> 下载最新版本
            </a>
          </div>

          <div style={{ marginTop: '16px' }}>
            <p><strong>Windows 用户</strong>：</p>
            <p>下载 <code>servergo_windows_amd64.zip</code> 文件，解压后即可使用。</p>
            
            <p><strong>macOS 用户</strong>：</p>
            <CodeBlock code={`# 解压下载的文件
tar -xzf servergo_darwin_amd64.tar.gz

# 添加执行权限
chmod +x servergo

# 可选：移动到 PATH 目录
sudo mv servergo /usr/local/bin/`} />
            
            <p><strong>Linux 用户</strong>：</p>
            <CodeBlock code={`# 解压下载的文件
tar -xzf servergo_linux_amd64.tar.gz

# 添加执行权限
chmod +x servergo

# 可选：移动到 PATH 目录
sudo mv servergo /usr/local/bin/`} />
          </div>
        </div>

        <div>
          <h4>使用 Go 安装</h4>
          <p>如果您已安装 Go 环境，可以使用以下命令直接安装：</p>
          <CodeBlock code="go install github.com/cc11001100/servergo@latest" />
        </div>
      </Card>

      {/* 基本用法 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          marginBottom: '16px' 
        }}>
          <FiCommand style={{ fontSize: '20px', marginRight: '8px' }} />
          <h3 style={{ margin: 0 }}>基本用法</h3>
        </div>

        <p>ServerGo 设计为简单易用，您可以直接在命令行中启动：</p>
        
        <div style={{ marginTop: '16px' }}>
          <h4>在当前目录启动文件服务器</h4>
          <CodeBlock code="servergo" />
          <p>这将在当前目录启动一个文件服务器，默认会自动探测可用端口。</p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>指定端口</h4>
          <CodeBlock code="servergo -p 3000" />
          <p>这将在端口 3000 上启动服务器。如果指定的端口已被占用，ServerGo 将自动查找可用端口。</p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>指定目录</h4>
          <CodeBlock code="servergo -d /path/to/share" />
          <p>这将为指定目录提供文件服务。</p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>组合使用</h4>
          <CodeBlock code="servergo -p 3000 -d /path/to/share" />
          <p>指定端口和目录启动服务器。</p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>自动打开浏览器</h4>
          <CodeBlock code="servergo -o" />
          <p>启动服务器后自动在默认浏览器中打开（默认已启用，除非在配置中禁用）。</p>
        </div>
      </Card>

      {/* 特性概览 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          marginBottom: '16px' 
        }}>
          <FiCpu style={{ fontSize: '20px', marginRight: '8px' }} />
          <h3 style={{ margin: 0 }}>特性概览</h3>
        </div>

        <ul style={{ paddingLeft: '20px' }}>
          <li><strong>自动端口探测</strong> - 如果指定端口被占用，自动查找可用端口</li>
          <li><strong>美观的文件浏览界面</strong> - 支持多种主题（default, dark, light, github, material）</li>
          <li><strong>认证机制</strong> - 支持基本认证、令牌认证和表单登录</li>
          <li><strong>多语言支持</strong> - 支持英文和简体中文</li>
          <li><strong>配置文件</strong> - 可保存和使用自定义配置</li>
          <li><strong>系统集成</strong> - 可作为系统服务安装</li>
        </ul>

        <div style={{ marginTop: '20px' }}>
          <p>
            ServerGo 是为了日常使用而设计的，适合于临时文件共享、静态网站测试、内部文档托管等场景。
            了解更多高级功能，请参阅其他文档部分。
          </p>
        </div>
      </Card>

      {/* 快速示例 */}
      <Card>
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          marginBottom: '16px' 
        }}>
          <FiExternalLink style={{ fontSize: '20px', marginRight: '8px' }} />
          <h3 style={{ margin: 0 }}>快速示例场景</h3>
        </div>

        <div style={{ marginBottom: '16px' }}>
          <h4>分享项目文件给同事</h4>
          <CodeBlock code="cd /path/to/project && servergo" />
          <p>同事可以通过 <code>http://your-ip:[端口]</code> 在浏览器中查看和下载项目文件。</p>
        </div>

        <div style={{ marginBottom: '16px' }}>
          <h4>安全地共享敏感文档</h4>
          <CodeBlock code="servergo -d /path/to/docs -a basic -u admin -w secure123" />
          <p>添加基本认证，仅允许知道用户名和密码的用户访问文件。</p>
        </div>

        <div style={{ marginBottom: '16px' }}>
          <h4>使用表单登录</h4>
          <CodeBlock code="servergo -a basic -u admin -w secure123 -l" />
          <p>提供用户友好的登录界面，而不是浏览器弹出的认证对话框。</p>
        </div>

        <div style={{ marginBottom: '16px' }}>
          <h4>自定义主题</h4>
          <CodeBlock code="servergo -m dark" />
          <p>使用暗色主题启动服务器，提供更舒适的浏览体验。</p>
        </div>

        <div>
          <p style={{ marginTop: '16px' }}>
            更多使用场景和详细配置，请查看 <a href="#usage-examples" style={{ cursor: 'pointer' }}>使用示例</a> 部分。
          </p>
        </div>
      </Card>
    </section>
  )
} 