import Card from '../../components/Card'
import CodeBlock from '../../components/CodeBlock'
import { FiSettings, FiCommand, FiLock, FiSliders, FiList, FiInfo } from 'react-icons/fi'

// 命令行选项分类
const optionCategories = [
  {
    id: 'basic',
    title: '基本选项',
    icon: <FiCommand />,
    options: [
      { option: '-p, --port <端口>', desc: '指定HTTP服务器端口', defaultValue: '0 (自动探测)', example: 'servergo -p 3000' },
      { option: '-d, --dir <目录>', desc: '指定要提供服务的目录', defaultValue: '当前工作目录', example: 'servergo -d /path/to/files' },
      { option: '-o, --open', desc: '启动后自动打开浏览器', defaultValue: '配置文件中的值', example: 'servergo -o' }
    ]
  },
  {
    id: 'auth',
    title: '认证选项',
    icon: <FiLock />,
    options: [
      { option: '-a, --auth <类型>', desc: '认证类型 (none, basic, token, form)', defaultValue: 'none', example: 'servergo -a basic' },
      { option: '-u, --username <用户名>', desc: '认证用户名', defaultValue: '配置文件中的值', example: 'servergo -a basic -u admin' },
      { option: '-w, --password <密码>', desc: '认证密码', defaultValue: '配置文件中的值', example: 'servergo -a basic -u admin -w secure123' },
      { option: '-t, --token <令牌>', desc: '令牌认证字符串', defaultValue: '无', example: 'servergo -a token -t your-secret-token' },
      { option: '-l, --login-page', desc: '启用网页登录界面', defaultValue: 'false', example: 'servergo -a basic -l' }
    ]
  },
  {
    id: 'list',
    title: '目录列表选项',
    icon: <FiList />,
    options: [
      { option: '-i, --dir-list', desc: '启用目录列表功能', defaultValue: '配置文件中的值', example: 'servergo -i' },
      { option: '-m, --theme <主题>', desc: '目录列表主题 (default, dark, blue, green, retro, json, table)', defaultValue: '配置文件中的值', example: 'servergo -m dark' },
    ]
  },
  {
    id: 'log',
    title: '日志选项',
    icon: <FiInfo />,
    options: [
      { option: '--log-level <级别>', desc: '日志级别(debug, info, warn, error)', defaultValue: 'info', example: 'servergo --log-level debug' },
      { option: '--enable-log-persistence', desc: '启用日志持久化', defaultValue: '配置文件中的值', example: 'servergo --enable-log-persistence' }
    ]
  },
  {
    id: 'config',
    title: '配置管理',
    icon: <FiSliders />,
    options: [
      { option: 'config list', desc: '列出当前所有配置', defaultValue: '无', example: 'servergo config list' },
      { option: 'config get <键>', desc: '获取指定配置项的值', defaultValue: '无', example: 'servergo config get theme' },
      { option: 'config set <键> <值>', desc: '设置配置项的值', defaultValue: '无', example: 'servergo config set theme dark' }
    ]
  }
];

export default function CommandOptions() {
  return (
    <section style={{ marginBottom: '40px' }}>
      <h2 style={{ display: 'flex', alignItems: 'center', marginBottom: '20px' }}>
        <FiSettings style={{ marginRight: '10px' }} /> 命令行选项
      </h2>

      {/* 基本概述 */}
      <Card style={{ marginBottom: '24px' }}>
        <h3 style={{ marginBottom: '16px' }}>命令概述</h3>
        <p>
          ServerGo 提供丰富的命令行选项，让您可以根据需求自定义文件服务器。
          以下是所有可用命令和选项的详细说明，按功能分类。
        </p>

        <div style={{ marginTop: '16px' }}>
          <h4>命令格式</h4>
          <CodeBlock code="servergo [命令] [选项]" />
          <p style={{ marginTop: '8px' }}>
            <strong>注意</strong>：直接运行 <code>servergo</code> 不带任何参数时，等同于运行 <code>servergo start</code>
          </p>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>可用命令</h4>
          <ul style={{ paddingLeft: '20px' }}>
            <li><code>start</code> - 启动文件服务器（默认命令）</li>
            <li><code>config</code> - 管理配置</li>
            <li><code>version</code> - 显示版本信息</li>
            <li><code>install</code> - 安装为系统服务</li>
            <li><code>uninstall</code> - 卸载系统服务</li>
          </ul>
        </div>

        <div style={{ marginTop: '16px' }}>
          <h4>获取帮助</h4>
          <p>您可以随时使用 <code>-h</code> 或 <code>--help</code> 选项获取帮助：</p>
          <CodeBlock code="servergo -h               # 显示一般帮助信息
servergo start -h        # 显示启动命令帮助
servergo config -h       # 显示配置命令帮助" />
        </div>
      </Card>

      {/* 各类命令选项 */}
      {optionCategories.map(category => (
        <Card key={category.id} style={{ marginBottom: '24px' }}>
          <div id={category.id} style={{ 
            display: 'flex', 
            alignItems: 'center', 
            marginBottom: '16px'
          }}>
            <span style={{ 
              fontSize: '22px', 
              color: 'var(--primary-color)', 
              marginRight: '10px' 
            }}>
              {category.icon}
            </span>
            <h3 style={{ margin: 0 }}>{category.title}</h3>
          </div>

          <div className="table-responsive" style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr>
                  <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>选项</th>
                  <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>描述</th>
                  <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>默认值</th>
                  <th style={{ textAlign: 'left', padding: '10px', borderBottom: '1px solid var(--border-color)' }}>示例</th>
                </tr>
              </thead>
              <tbody>
                {category.options.map((opt, index) => (
                  <tr key={index}>
                    <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>{opt.option}</code></td>
                    <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>{opt.desc}</td>
                    <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}>{opt.defaultValue}</td>
                    <td style={{ padding: '10px', borderBottom: '1px solid var(--border-color)' }}><code>{opt.example}</code></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      ))}

      {/* 命令组合示例 */}
      <Card>
        <h3 style={{ marginBottom: '16px' }}>常用命令组合示例</h3>
        
        <div style={{ marginBottom: '16px' }}>
          <h4>启动带基本认证的服务器</h4>
          <CodeBlock code="servergo -a basic -u admin -w secure123" />
        </div>
        
        <div style={{ marginBottom: '16px' }}>
          <h4>自定义主题和端口</h4>
          <CodeBlock code="servergo -m dark -p 3000" />
        </div>
        
        <div style={{ marginBottom: '16px' }}>
          <h4>为特定目录提供服务</h4>
          <CodeBlock code="servergo -p 3000 -d ./build -a basic -u admin -w password" />
        </div>
        
        <div style={{ marginBottom: '16px' }}>
          <h4>使用配置文件</h4>
          <p>首先设置配置：</p>
          <CodeBlock code="servergo config set theme dark
servergo config set auto-open true" />
          <p>然后启动服务器（将使用配置中的设置）：</p>
          <CodeBlock code="servergo" />
        </div>
      </Card>
    </section>
  )
} 