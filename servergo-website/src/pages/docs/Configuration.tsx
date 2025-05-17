import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiSettings, FiSliders, FiToggleRight, FiGlobe, FiFile, FiEye } from 'react-icons/fi'

export default function Configuration() {
  return (
    <section className="config-section">
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiSettings style={{ marginRight: '10px' }} /> 配置文件与系统
      </h2>

      {/* 配置系统概述 */}
      <Card style={{ marginBottom: '24px' }}>
        <h3 style={{ marginBottom: '16px' }}>配置系统概述</h3>
        <p>
          ServerGo 提供灵活的配置管理方式，支持通过命令行参数、配置文件或环境变量来定制服务器行为。
          配置项的优先级从高到低为：
        </p>
        <ol style={{ paddingLeft: '20px', marginTop: '10px' }}>
          <li><strong>命令行参数</strong> - 每次运行时直接指定</li>
          <li><strong>环境变量</strong> - 格式为 <code>SERVERGO_配置名</code></li>
          <li><strong>配置文件</strong> - 保存永久配置</li>
        </ol>
      </Card>

      {/* 配置文件 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiFile style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>配置文件</h3>
        </div>

        <p>
          ServerGo 的配置文件默认存储在以下位置（按平台区分）：
        </p>

        <div className="table-responsive" style={{ overflowX: 'auto', marginTop: '16px', marginBottom: '16px' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>平台</th>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>配置文件路径</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>Windows</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>%USERPROFILE%\.servergo\config.yaml</code></td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>macOS</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>$HOME/.servergo/config.yaml</code></td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>Linux</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>$HOME/.servergo/config.yaml</code></td>
              </tr>
            </tbody>
          </table>
        </div>

        <p>配置文件使用 YAML 格式，示例如下：</p>
        <CodeBlock code="# ServerGo 配置文件
auto-open: true
theme: default
language: zh-CN
enable-dir-listing: true
enable-log-persistence: true
username: admin
password: " />
      </Card>

      {/* 配置项管理 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiSliders style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>配置项管理</h3>
        </div>

        <p>
          ServerGo 提供了专门的 <code>config</code> 命令来管理配置项，让您可以方便地查看、获取和设置配置。
        </p>

        <div style={{ marginTop: '16px' }}>
          <h4>列出所有配置</h4>
          <p>要查看当前所有配置项的值：</p>
          <CodeBlock code="servergo config list" />
          <p>输出示例：</p>
          <CodeBlock code={`+-----------------------+---------------+----------------------------------+
| 配置项                | 当前值          | 描述                             |
+-----------------------+---------------+----------------------------------+
| auto-open             | 启用           | 启动服务器时自动打开浏览器         |
| enable-dir-listing    | 启用           | 启用目录列表功能                  |
| theme                 | default       | 目录列表主题                      |
| language              | 简体中文       | 界面显示语言                      |
| enable-log-persistence| 启用           | 是否将日志持久化保存到文件         |
+-----------------------+---------------+----------------------------------+`} />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>获取单个配置值</h4>
          <p>要获取特定配置项的值：</p>
          <CodeBlock code="servergo config get theme" />
          <p>输出示例：</p>
          <CodeBlock code="default" />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>修改配置值</h4>
          <p>要修改配置项的值：</p>
          <CodeBlock code="servergo config set theme dark" />
          <p>输出示例：</p>
          <CodeBlock code={`配置项 "theme" 已设置为 "dark"`} />
        </div>
      </Card>

      {/* 支持的配置项 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiToggleRight style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>支持的配置项</h3>
        </div>

        <div className="table-responsive" style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>配置项</th>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>类型</th>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>默认值</th>
                <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>说明</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>auto-open</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>布尔值</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>true</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>启动服务器时是否自动打开浏览器</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>enable-dir-listing</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>布尔值</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>true</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>是否启用目录列表功能</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>theme</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>字符串</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>default</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>目录列表主题 (default, dark, light, github, material)</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>language</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>字符串</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>系统语言</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>界面显示语言 (en, zh-CN)</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>enable-log-persistence</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>布尔值</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>true</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>是否将日志持久化保存到文件</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>username</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>字符串</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>admin</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>认证的默认用户名</td>
              </tr>
              <tr>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>password</code></td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>字符串</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>空</td>
                <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>认证的默认密码</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>布尔值配置的特殊说明</h4>
          <p>对于布尔类型的配置项，ServerGo 支持多种输入形式：</p>
          <ul style={{ paddingLeft: '20px' }}>
            <li>对于 <strong>true</strong>: <code>true</code>, <code>yes</code>, <code>y</code>, <code>1</code>, <code>on</code></li>
            <li>对于 <strong>false</strong>: <code>false</code>, <code>no</code>, <code>n</code>, <code>0</code>, <code>off</code></li>
          </ul>
          <p>例如，以下命令效果相同：</p>
          <CodeBlock code="servergo config set auto-open true
servergo config set auto-open yes
servergo config set auto-open on" />
        </div>
      </Card>

      {/* 语言设置 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiGlobe style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>多语言设置</h3>
        </div>

        <p>
          ServerGo 支持多种语言界面，可以通过配置 <code>language</code> 项来更改显示语言。
        </p>

        <div style={{ marginTop: '16px' }}>
          <h4>查看支持的语言</h4>
          <p>查看支持的语言列表：</p>
          <CodeBlock code="servergo config set language" />
          <p>输出示例：</p>
          <CodeBlock code={`可用的语言选项:
  - en (English)
  - zh-CN (简体中文)`} />
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>设置语言</h4>
          <CodeBlock code="servergo config set language zh-CN" />
          <p>输出示例：</p>
          <CodeBlock code={`界面语言已更改为 "简体中文"`} />
        </div>
      </Card>

      {/* 配置文件示例 */}
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
          <FiEye style={{ fontSize: '22px', color: 'var(--primary-color)', marginRight: '10px' }} />
          <h3 style={{ margin: 0 }}>完整配置文件示例</h3>
        </div>

        <p>
          以下是一个包含所有支持配置项的完整配置文件示例：
        </p>

        <CodeBlock code="# ServerGo 配置文件
# 通用设置
auto-open: true              # 启动时自动打开浏览器
enable-dir-listing: true     # 启用目录列表

# 外观设置
theme: dark                  # 使用暗色主题
language: zh-CN              # 使用简体中文

# 日志设置
enable-log-persistence: true # 启用日志持久化到文件

# 认证设置
username: admin              # 默认用户名
password: secure123          # 默认密码

# 注意：以下为命令行专用选项，不会保存到配置文件
# 端口通常会自动选择可用的端口
# auth: basic                # 认证类型" />
      </Card>
    </section>
  )
} 