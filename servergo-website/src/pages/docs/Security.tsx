import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiLock, FiShield, FiUserCheck, FiEye } from 'react-icons/fi'

export default function Security() {
  return (
    <section className="security-section">
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiShield style={{ marginRight: '10px' }} /> 安全特性
      </h2>

      {/* 安全概述 */}
      <Card style={{ marginBottom: '24px' }}>
        <h3 style={{ marginBottom: '16px' }}>安全概述</h3>
        <p>
          由于 ServerGo 通常用于共享文件，安全性是一个重要考虑因素。
          ServerGo 提供了认证功能，帮助您保护您的文件和服务器免受未授权访问：
        </p>
        <ul style={{ paddingLeft: '20px', marginTop: '10px' }}>
          <li><strong>认证机制</strong> - 包括基本认证、令牌认证和表单登录</li>
          <li><strong>目录访问限制</strong> - 控制可以访问的目录</li>
        </ul>
      </Card>

      {/* 认证机制 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiUserCheck style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>认证机制</h3>
        </div>

        <p>
          ServerGo 支持多种认证方式，可以根据您的需求选择最合适的方式：
        </p>

        <div style={{ marginTop: '16px' }}>
          <h4>基本认证 (Basic Auth)</h4>
          <p>最简单的认证方式，通过 HTTP 基本认证协议实现：</p>
          <CodeBlock code="servergo -a basic -u admin -w secure123" />
          <p>
            <strong>优点</strong>: 简单易用，大多数浏览器都支持<br />
            <strong>限制</strong>: 凭据以 Base64 编码传输，安全性有限
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>令牌认证 (Token Auth)</h4>
          <p>通过 URL 参数或请求头部传递访问令牌：</p>
          <CodeBlock code="servergo -a token -t your-secret-token" />
          <p>访问方式：</p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>URL 方式: <code>http://localhost:8080/?token=your-secret-token</code></li>
            <li>请求头部: <code>Authorization: Bearer your-secret-token</code></li>
          </ul>
          <p>
            <strong>优点</strong>: 适合程序化访问或 API 调用<br />
            <strong>限制</strong>: 需要在 URL 或请求中包含令牌
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>表单登录 (Form Login)</h4>
          <p>提供网页登录界面，用户友好：</p>
          <CodeBlock code="servergo -a basic -u admin -w secure123 -l" />
          <p>
            <strong>优点</strong>: 提供用户友好的登录界面，支持会话保持<br />
            <strong>限制</strong>: 需要浏览器支持 Cookie
          </p>
        </div>
      </Card>

      {/* 目录访问限制 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiEye style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>目录访问限制</h3>
        </div>

        <p>
          您可以控制哪些目录对外可见，以及是否允许浏览目录内容：
        </p>

        <div style={{ marginTop: '16px' }}>
          <h4>限制可访问目录</h4>
          <p>指定一个目录作为文件服务的根目录：</p>
          <CodeBlock code="servergo -d /path/to/public/files" />
          <p>这样用户只能访问该目录及其子目录中的文件。</p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>禁用目录列表</h4>
          <p>如不希望用户浏览目录内容，可禁用目录列表功能：</p>
          <CodeBlock code="servergo -i false" />
          <p>禁用后，用户必须知道确切的文件路径才能访问文件。</p>
        </div>
      </Card>

      {/* 安全最佳实践 */}
      <Card>
        <h3 style={{ marginBottom: '16px' }}>安全最佳实践</h3>
        <p>
          以下是使用 ServerGo 时的一些安全最佳实践建议：
        </p>

        <ol style={{ paddingLeft: '20px', marginTop: '10px' }}>
          <li>
            <strong>总是启用认证</strong> - 对公开的文件服务器使用认证
            <CodeBlock code="servergo -a basic -u admin -w complex-password" />
          </li>
          <li>
            <strong>仅共享必要文件</strong> - 创建专门用于共享的目录，而不是共享整个系统
            <CodeBlock code="servergo -d /path/to/shared/files" />
          </li>
          <li>
            <strong>使用强密码</strong> - 避免使用简单或默认密码
          </li>
          <li>
            <strong>定期轮换凭据</strong> - 特别是对长期运行的服务器
          </li>
          <li>
            <strong>监控访问日志</strong> - 启用日志记录以跟踪访问
            <CodeBlock code="servergo --enable-log-persistence" />
          </li>
          <li>
            <strong>使用表单认证</strong> - 提供更好的用户体验
            <CodeBlock code="servergo -a basic -u admin -w password -l" />
          </li>
        </ol>
      </Card>
    </section>
  )
} 