import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiHelpCircle, FiCpu, FiSettings, FiLock, FiGlobe, FiServer } from 'react-icons/fi'

export default function FAQ() {
  return (
    <section className="faq-section">
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiHelpCircle style={{ marginRight: '10px' }} /> 常见问题解答
      </h2>

      {/* 基本问题 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiServer style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>基本使用问题</h3>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>ServerGo 与其他文件服务器如 Nginx, Apache 有什么区别？</h4>
          <p>
            ServerGo 专为简单、快速的文件共享设计，区别在于：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>无需配置文件，可直接启动使用</li>
            <li>内置美观的文件浏览界面</li>
            <li>跨平台支持 (Windows, macOS, Linux)</li>
            <li>单二进制文件，不需要安装依赖</li>
          </ul>
          <p>
            Nginx 和 Apache 更适合作为生产环境的 Web 服务器，而 ServerGo 更适合临时文件共享或开发环境。
          </p>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>如何查看 ServerGo 的版本？</h4>
          <p>
            运行以下命令查看当前安装的 ServerGo 版本：
          </p>
          <CodeBlock code="servergo version" />
        </div>

        <div className="faq-item">
          <h4>我可以同时启动多个 ServerGo 实例吗？</h4>
          <p>
            可以，只需确保每个实例使用不同的端口：
          </p>
          <CodeBlock code="# 第一个实例，自动探测可用端口
servergo -d /path/to/dir1

# 第二个实例，使用端口 8081
servergo -d /path/to/dir2 -p 8081" />
        </div>
      </Card>

      {/* 性能问题 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiCpu style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>性能与资源</h3>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>ServerGo 的资源消耗如何？</h4>
          <p>
            ServerGo 设计为轻量级应用，资源消耗很低：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>空闲时内存占用通常低于 20MB</li>
            <li>CPU 使用率仅在处理请求时短暂增加</li>
            <li>二进制文件大小约 15-20MB</li>
          </ul>
          <p>
            即使在资源受限的系统（如树莓派）上，ServerGo 也能高效运行。
          </p>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>ServerGo 可以处理多少并发连接？</h4>
          <p>
            ServerGo 使用 Go 语言的高效并发模型构建，可以处理数百甚至数千个并发连接，具体取决于您的系统资源和网络带宽。
            对于普通的文件共享使用场景，性能完全足够。
          </p>
        </div>

        <div className="faq-item">
          <h4>如何提高 ServerGo 的性能？</h4>
          <p>
            虽然 ServerGo 已经相当高效，但以下是一些提高性能的建议：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>在 SSD 上运行 ServerGo，特别是分享大量小文件时</li>
            <li>使用有线网络连接而非无线连接</li>
            <li>关闭不必要的日志记录：<code>servergo --log-level error</code></li>
          </ul>
        </div>
      </Card>

      {/* 安全问题 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiLock style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>安全相关问题</h3>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>ServerGo 的认证方式安全吗？</h4>
          <p>
            ServerGo 提供的基本认证适合临时使用，但有以下注意事项：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>基本认证（Basic Auth）的凭据以明文方式传输，安全性有限</li>
            <li>对于敏感数据，建议使用在可信任的内部网络中使用</li>
            <li>尽量使用复杂的密码增加安全性</li>
          </ul>
        </div>

        <div className="faq-item">
          <h4>我如何限制只能访问特定文件或目录？</h4>
          <p>
            ServerGo 允许您指定服务的根目录，用户只能访问该目录及其子目录：
          </p>
          <CodeBlock code="# 创建专门用于共享的目录
mkdir -p ~/share
cp /path/to/files/to/share ~/share/

# 启动服务器并仅提供这个目录
servergo -d ~/share" />
        </div>
      </Card>

      {/* 故障排除 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiSettings style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>故障排除</h3>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>端口已被占用错误怎么解决？</h4>
          <p>
            如果启动 ServerGo 时出现 "端口已被占用" 错误，可以：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>指定另一个可用端口：<code>servergo -p 9000</code></li>
            <li>查找并关闭占用端口的进程：
              <CodeBlock code="# 查找占用 8080 端口的进程
lsof -i :8080
# 或在 Windows 上
netstat -ano | findstr :8080

# 然后终止相应的进程" />
            </li>
            <li>使用 <code>-p 0</code> 让 ServerGo 自动选择可用端口</li>
          </ul>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>为什么我无法在其他设备上访问我的 ServerGo 服务器？</h4>
          <p>
            这可能有几个原因：
          </p>
          <ol style={{ paddingLeft: '20px' }}>
            <li>
              <strong>防火墙阻止</strong> - 确保防火墙允许指定端口：
              <ul style={{ paddingLeft: '20px' }}>
                <li>Windows: 检查 Windows 防火墙设置</li>
                <li>macOS: 在系统偏好设置中检查防火墙</li>
                <li>Linux: 检查 iptables 或 ufw 规则</li>
              </ul>
            </li>
            <li>
              <strong>网络路由</strong> - 确保您在同一网络，或有正确的端口转发规则</li>
          </ol>
        </div>

        <div className="faq-item">
          <h4>如何查看 ServerGo 的日志和错误信息？</h4>
          <p>
            您可以通过以下方式获取更多日志信息：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>增加日志级别获取详细信息：<code>servergo --log-level debug</code></li>
            <li>启用日志持久化：<code>servergo --enable-log-persistence</code></li>
            <li>检查日志文件位置（启用持久化后）：
              <ul style={{ paddingLeft: '20px' }}>
                <li>Windows: <code>%USERPROFILE%\.servergo\logs\</code></li>
                <li>macOS/Linux: <code>~/.servergo/logs/</code></li>
              </ul>
            </li>
          </ul>
        </div>
      </Card>

      {/* 高级问题 */}
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiGlobe style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>高级使用问题</h3>
        </div>

        <div className="faq-item" style={{ marginBottom: '20px' }}>
          <h4>如何在后台运行 ServerGo？</h4>
          <p>
            取决于您的操作系统：
          </p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>
              <strong>Windows</strong>:
              <ul style={{ paddingLeft: '20px' }}>
                <li>安装为系统服务：<code>servergo install</code></li>
                <li>或使用 <code>start /b servergo</code> 作为后台进程运行</li>
              </ul>
            </li>
            <li>
              <strong>Linux/macOS</strong>:
              <ul style={{ paddingLeft: '20px' }}>
                <li>使用 nohup：<code>nohup servergo &</code></li>
                <li>或使用 screen/tmux：<code>screen -S servergo</code> 然后运行 ServerGo，按 Ctrl+A D 分离</li>
                <li>或使用 systemd 创建服务（需要额外配置）</li>
              </ul>
            </li>
          </ul>
        </div>

        <div className="faq-item">
          <h4>如何更改默认配置？</h4>
          <p>
            您可以使用配置命令来永久更改默认设置：
          </p>
          <CodeBlock code="# 查看当前配置
servergo config list

# 设置默认主题为深色
servergo config set theme dark

# 默认禁用自动打开浏览器
servergo config set auto-open false" />
          <p>
            这些设置会保存在配置文件中，成为以后启动 ServerGo 时的默认值。
          </p>
        </div>
      </Card>
    </section>
  )
} 