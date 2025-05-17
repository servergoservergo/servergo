import { Link } from 'react-router-dom'
import { FiDownload, FiCode, FiServer, FiShield, FiZap } from 'react-icons/fi'
import Card from '../components/Card'
import CodeBlock from '../components/CodeBlock'

export default function Home() {
  return (
    <div>
      {/* 英雄区域 */}
      <section style={{ 
        textAlign: 'center',
        padding: '60px 0',
        backgroundColor: 'var(--bg-light)',
        borderRadius: 'var(--border-radius)',
        margin: '0 0 60px'
      }}>
        <h1 style={{ 
          fontSize: '2.5rem',
          marginBottom: '16px',
          color: 'var(--text-dark)'
        }}>
          <FiServer style={{ marginRight: '12px' }} />
          ServerGo
        </h1>
        
        <p style={{ 
          fontSize: '1.25rem',
          maxWidth: '700px',
          margin: '0 auto 30px',
          color: 'var(--text-light)'
        }}>
          轻量级、高性能、易使用的静态文件服务器
        </p>
        
        <div style={{ display: 'flex', gap: '16px', justifyContent: 'center' }}>
          <Link to="/download" className="btn">
            <FiDownload />
            下载
          </Link>
          <Link to="/docs" className="btn btn-secondary">
            <FiCode />
            查看文档
          </Link>
        </div>
      </section>
      
      {/* 快速开始 */}
      <section style={{ marginBottom: '60px' }}>
        <h2 style={{ textAlign: 'center', marginBottom: '30px' }}>快速开始</h2>
        
        <Card>
          <div style={{ marginBottom: '24px' }}>
            <h3>基本用法</h3>
            <p>在当前目录启动一个文件服务器：</p>
            <CodeBlock code="servergo" />
          </div>
          
          <div style={{ marginBottom: '24px' }}>
            <h3>指定端口</h3>
            <p>在指定端口启动：</p>
            <CodeBlock code="servergo -p 3000" />
          </div>
          
          <div>
            <h3>指定目录</h3>
            <p>为特定目录提供服务：</p>
            <CodeBlock code="servergo -d /path/to/files" />
          </div>
        </Card>
      </section>
      
      {/* 核心特性 */}
      <section style={{ marginBottom: '60px' }}>
        <h2 style={{ textAlign: 'center', marginBottom: '30px' }}>核心特性</h2>
        
        <div style={{ 
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
          gap: '24px'
        }}>
          <Card>
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              marginBottom: '16px'
            }}>
              <div style={{ 
                width: '48px',
                height: '48px',
                backgroundColor: 'var(--primary-light)',
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                marginRight: '16px',
                color: 'white',
                fontSize: '24px'
              }}>
                <FiZap />
              </div>
              <h3 style={{ margin: 0 }}>高性能</h3>
            </div>
            <p>使用Go语言原生开发，资源占用低，传输速度快。即使在资源有限的环境中也能高效运行。</p>
          </Card>
          
          <Card>
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              marginBottom: '16px'
            }}>
              <div style={{ 
                width: '48px',
                height: '48px',
                backgroundColor: 'var(--primary-light)',
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                marginRight: '16px',
                color: 'white',
                fontSize: '24px'
              }}>
                <FiCode />
              </div>
              <h3 style={{ margin: 0 }}>简单易用</h3>
            </div>
            <p>零配置启动，简洁的命令行参数，直观的Web界面。无需复杂设置即可快速共享文件。</p>
          </Card>
          
          <Card>
            <div style={{ 
              display: 'flex',
              alignItems: 'center',
              marginBottom: '16px'
            }}>
              <div style={{ 
                width: '48px',
                height: '48px',
                backgroundColor: 'var(--primary-light)',
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                marginRight: '16px',
                color: 'white',
                fontSize: '24px'
              }}>
                <FiShield />
              </div>
              <h3 style={{ margin: 0 }}>安全可靠</h3>
            </div>
            <p>支持HTTPS加密、基本认证以及其他安全特性。可以安全地共享敏感文件和内部文档。</p>
          </Card>
        </div>
      </section>
      
      {/* 行动号召 */}
      <section style={{ 
        textAlign: 'center',
        backgroundColor: 'var(--primary-color)',
        padding: '40px',
        borderRadius: 'var(--border-radius)',
        color: 'white'
      }}>
        <h2 style={{ marginBottom: '16px' }}>准备好了吗？立即开始使用 ServerGo！</h2>
        <p style={{ marginBottom: '24px', opacity: 0.8 }}>免费、开源，适用于所有主流操作系统</p>
        <Link to="/download" className="btn" style={{ 
          backgroundColor: 'white',
          color: 'var(--primary-color)',
          padding: '12px 24px',
          fontSize: '1.1rem'
        }}>
          <FiDownload />
          立即下载
        </Link>
      </section>
    </div>
  )
} 