import { useState, useEffect } from 'react'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { FiGithub, FiStar } from 'react-icons/fi'

// GitHub 仓库地址
const GITHUB_REPO = 'cc11001100/servergo'

export default function Layout() {
  const location = useLocation()
  const [stars, setStars] = useState<number | null>(null)

  // 获取GitHub星星数
  useEffect(() => {
    fetch(`https://api.github.com/repos/${GITHUB_REPO}`)
      .then(response => response.json())
      .then(data => {
        if (data.stargazers_count) {
          setStars(data.stargazers_count)
        }
      })
      .catch(error => {
        console.error('获取 GitHub 星星数失败:', error)
      })
  }, [])

  return (
    <div className="layout">
      <header className="header">
        <div className="container header-container">
          <Link to="/" className="logo">
            ServerGo
          </Link>
          
          <nav className="nav">
            <Link 
              to="/" 
              className={`nav-link ${location.pathname === '/' ? 'active' : ''}`}
            >
              首页
            </Link>
            <Link 
              to="/docs" 
              className={`nav-link ${location.pathname === '/docs' ? 'active' : ''}`}
            >
              文档
            </Link>
            <Link 
              to="/examples" 
              className={`nav-link ${location.pathname === '/examples' ? 'active' : ''}`}
            >
              示例
            </Link>
            <Link 
              to="/install" 
              className={`nav-link ${location.pathname === '/install' ? 'active' : ''}`}
            >
              安装
            </Link>
            <Link 
              to="/download" 
              className={`nav-link ${location.pathname === '/download' ? 'active' : ''}`}
            >
              下载
            </Link>
            <a 
              href={`https://github.com/${GITHUB_REPO}`} 
              target="_blank" 
              rel="noopener noreferrer"
              className="nav-link"
              style={{ display: 'flex', alignItems: 'center', gap: '8px' }}
            >
              <FiGithub />
              GitHub
              {stars !== null && (
                <span style={{ 
                  backgroundColor: 'var(--bg-light)',
                  padding: '2px 8px',
                  borderRadius: '10px',
                  fontSize: '12px',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px'
                }}>
                  <FiStar style={{ fontSize: '14px' }} />
                  {stars}
                </span>
              )}
            </a>
          </nav>
        </div>
      </header>
      
      <main className="main-content">
        <div className="container">
          <Outlet />
        </div>
      </main>
      
      <footer className="footer">
        <div className="container">
          <div style={{ textAlign: 'center' }}>
            <p>
              © {new Date().getFullYear()} ServerGo. 
              <a 
                href={`https://github.com/${GITHUB_REPO}`}
                target="_blank" 
                rel="noopener noreferrer"
                style={{ marginLeft: '8px' }}
              >
                GitHub
              </a>
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
} 