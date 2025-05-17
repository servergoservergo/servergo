import Card from '../components/Card'
import CodeBlock from '../components/CodeBlock'
import { FiDownload, FiTerminal, FiInfo } from 'react-icons/fi'

// 假设的下载链接和版本信息
const VERSION = 'v1.0.0'
const DOWNLOAD_BASE_URL = 'https://github.com/cc11001100/servergo/releases/download/'

// 定义各平台的下载项
const downloadItems = [
  {
    id: 'windows',
    name: 'Windows',
    arch: [
      { name: 'Windows x64', filename: `servergo_${VERSION}_windows_amd64.zip` },
      { name: 'Windows x86', filename: `servergo_${VERSION}_windows_386.zip` }
    ],
    icon: '🪟',
    installCommand: `# 解压下载的zip文件
# 双击servergo.exe运行，或在命令提示符中运行
servergo.exe`
  },
  {
    id: 'macos',
    name: 'macOS',
    arch: [
      { name: 'macOS Intel', filename: `servergo_${VERSION}_darwin_amd64.tar.gz` },
      { name: 'macOS Apple Silicon', filename: `servergo_${VERSION}_darwin_arm64.tar.gz` }
    ],
    icon: '🍎',
    installCommand: `# 解压下载的tar.gz文件
tar -xzf servergo_${VERSION}_darwin_*.tar.gz

# 添加执行权限
chmod +x servergo

# 运行
./servergo`
  },
  {
    id: 'linux',
    name: 'Linux',
    arch: [
      { name: 'Linux x64', filename: `servergo_${VERSION}_linux_amd64.tar.gz` },
      { name: 'Linux x86', filename: `servergo_${VERSION}_linux_386.tar.gz` },
      { name: 'Linux ARM64', filename: `servergo_${VERSION}_linux_arm64.tar.gz` },
      { name: 'Linux ARMv7', filename: `servergo_${VERSION}_linux_armv7.tar.gz` }
    ],
    icon: '🐧',
    installCommand: `# 解压下载的tar.gz文件
tar -xzf servergo_${VERSION}_linux_*.tar.gz

# 添加执行权限
chmod +x servergo

# 运行
./servergo`
  }
]

export default function Download() {
  return (
    <div>
      <h1 style={{ textAlign: 'center', marginBottom: '40px' }}>下载 ServerGo</h1>
      
      <section style={{ marginBottom: '40px' }}>
        <Card 
          title={<><FiInfo style={{ marginRight: '8px' }} /> 当前版本</>}
        >
          <p>ServerGo {VERSION} - 发布于 2023年12月1日</p>
          <div style={{ marginTop: '16px' }}>
            <a 
              href="https://github.com/cc11001100/servergo/releases" 
              target="_blank"
              rel="noopener noreferrer"
              className="btn"
              style={{ marginRight: '16px' }}
            >
              查看所有版本
            </a>
            <a 
              href="https://github.com/cc11001100/servergo" 
              target="_blank"
              rel="noopener noreferrer"
              className="btn"
              style={{ background: '#fff', color: 'var(--primary-color)', border: '1px solid var(--border-color)' }}
            >
              源代码
            </a>
          </div>
        </Card>
      </section>
      
      <section>
        {downloadItems.map(platform => (
          <Card 
            key={platform.id}
            title={
              <div id={platform.id}>
                <span style={{ fontSize: '24px', marginRight: '10px' }}>{platform.icon}</span>
                {platform.name}
              </div>
            }
            style={{ marginBottom: '30px' }}
          >
            <div style={{ marginBottom: '20px' }}>
              <h3 style={{ marginBottom: '16px' }}>下载</h3>
              <div style={{ 
                display: 'grid', 
                gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', 
                gap: '16px'
              }}>
                {platform.arch.map(item => (
                  <a 
                    key={item.filename}
                    href={`${DOWNLOAD_BASE_URL}${VERSION}/${item.filename}`}
                    className="btn"
                    style={{ justifyContent: 'flex-start' }}
                  >
                    <FiDownload />
                    {item.name}
                  </a>
                ))}
              </div>
            </div>
            
            <div>
              <h3 style={{ marginBottom: '16px', display: 'flex', alignItems: 'center' }}>
                <FiTerminal style={{ marginRight: '8px' }} /> 
                安装与使用
              </h3>
              <CodeBlock code={platform.installCommand} />
            </div>
          </Card>
        ))}
      </section>
      
      <section style={{ marginTop: '40px' }}>
        <Card 
          title={<><FiInfo style={{ marginRight: '8px' }} /> 验证完整性</>}
        >
          <p style={{ marginBottom: '16px' }}>
            建议下载后验证文件的SHA256哈希值以确保文件完整性。每个版本的发布页面都提供了对应的哈希值。
          </p>
          <CodeBlock 
            code={`# 在Windows上验证
CertUtil -hashfile servergo_${VERSION}_windows_amd64.zip SHA256

# 在macOS或Linux上验证
shasum -a 256 servergo_${VERSION}_darwin_amd64.tar.gz`} 
          />
        </Card>
      </section>
    </div>
  )
} 