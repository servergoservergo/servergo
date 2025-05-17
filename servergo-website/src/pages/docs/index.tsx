import React, { useState, useEffect, Suspense } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { FiHome, FiCommand, FiLock, FiSettings, FiServer, FiHelpCircle } from 'react-icons/fi'

// 直接导入文档组件，而不是使用懒加载
import GettingStarted from './GettingStarted'
import CommandOptions from './CommandOptions'
import UsageExamples from './UsageExamples'
import Security from './Security'
import Configuration from './Configuration'
import FAQ from './FAQ'

// 定义Loading组件
const Loading = () => (
  <div style={{ 
    display: 'flex', 
    justifyContent: 'center', 
    alignItems: 'center',
    padding: '40px',
    minHeight: '300px'
  }}>
    <div style={{ 
      display: 'inline-block',
      width: '50px',
      height: '50px',
      border: '5px solid var(--border-color)',
      borderTopColor: 'var(--primary-color)',
      borderRadius: '50%',
      animation: 'spin 1s linear infinite'
    }} />
    <style>
      {`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}
    </style>
  </div>
)

// 文档部分定义
const sections = [
  { id: 'getting-started', name: '入门指南', icon: <FiHome />, component: GettingStarted },
  { id: 'commands', name: '命令行选项', icon: <FiCommand />, component: CommandOptions },
  { id: 'usage', name: '使用示例', icon: <FiServer />, component: UsageExamples },
  { id: 'security', name: '安全设置', icon: <FiLock />, component: Security },
  { id: 'config', name: '配置文件', icon: <FiSettings />, component: Configuration },
  { id: 'faq', name: '常见问题', icon: <FiHelpCircle />, component: FAQ },
]

export default function Docs() {
  const location = useLocation()
  const navigate = useNavigate()
  const [activeSection, setActiveSection] = useState('getting-started')

  // 从URL哈希中获取当前活动部分
  useEffect(() => {
    const hash = location.hash.replace('#', '')
    if (hash && sections.some(section => section.id === hash)) {
      setActiveSection(hash)
    } else if (!hash) {
      // 如果没有哈希，默认显示第一个部分
      setActiveSection(sections[0].id)
    }
  }, [location])

  // 变更活动部分
  const handleSectionChange = (sectionId: string) => {
    setActiveSection(sectionId)
    navigate(`#${sectionId}`)
  }

  // 查找当前活动部分
  const currentSection = sections.find(section => section.id === activeSection) || sections[0]

  return (
    <div className="docs-container">
      <style>
        {`
          .docs-container {
            display: flex;
            margin-top: 20px;
          }
          
          .docs-sidebar {
            width: 250px;
            min-width: 250px;
            padding-right: 30px;
          }
          
          .docs-content {
            flex: 1;
            padding-bottom: 60px;
          }
          
          .section-link {
            display: flex;
            align-items: center;
            padding: 12px 15px;
            margin-bottom: 8px;
            border-radius: 8px;
            color: var(--text-color);
            text-decoration: none;
            font-weight: 500;
            transition: background-color 0.2s;
            cursor: pointer;
          }
          
          .section-link:hover {
            background-color: var(--hover-color);
          }
          
          .section-link.active {
            background-color: var(--primary-light);
            color: var(--primary-color);
          }
          
          .section-link svg {
            margin-right: 10px;
            font-size: 18px;
          }
          
          @media (max-width: 768px) {
            .docs-container {
              flex-direction: column;
            }
            
            .docs-sidebar {
              width: 100%;
              padding-right: 0;
              margin-bottom: 30px;
            }
            
            .docs-sidebar-inner {
              display: flex;
              flex-wrap: wrap;
              gap: 8px;
            }
            
            .section-link {
              padding: 8px 12px;
              margin-bottom: 0;
              font-size: 0.9rem;
            }
          }
        `}
      </style>

      {/* 侧边栏导航 */}
      <div className="docs-sidebar">
        <div className="docs-sidebar-inner">
          {sections.map(section => (
            <div 
              key={section.id}
              className={`section-link ${activeSection === section.id ? 'active' : ''}`}
              onClick={() => handleSectionChange(section.id)}
            >
              {section.icon}
              <span>{section.name}</span>
            </div>
          ))}
        </div>
      </div>

      {/* 内容区域 */}
      <div className="docs-content">
        <Suspense fallback={<Loading />}>
          {/* 渲染当前活动的组件 */}
          {React.createElement(currentSection.component)}
        </Suspense>
      </div>
    </div>
  )
} 