import { Link } from 'react-router-dom'
import { FiAlertTriangle, FiHome } from 'react-icons/fi'

export default function NotFound() {
  return (
    <div style={{ 
      textAlign: 'center', 
      padding: '80px 20px',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: 'calc(100vh - 200px)'
    }}>
      <FiAlertTriangle style={{ 
        fontSize: '60px', 
        color: 'var(--primary-color)',
        marginBottom: '20px'
      }} />
      
      <h1 style={{ 
        fontSize: '36px', 
        margin: '0 0 10px',
        color: 'var(--text-color)'
      }}>
        404 - 页面未找到
      </h1>
      
      <p style={{ 
        margin: '0 0 30px',
        color: 'var(--text-light)',
        maxWidth: '500px'
      }}>
        您要访问的页面不存在或已被移除。请检查您输入的URL是否正确，或返回首页继续浏览。
      </p>
      
      <Link to="/" className="btn" style={{ display: 'inline-flex', alignItems: 'center', gap: '8px' }}>
        <FiHome />
        返回首页
      </Link>
    </div>
  )
} 