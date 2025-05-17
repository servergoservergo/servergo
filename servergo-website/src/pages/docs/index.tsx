import React, { useState, useEffect, Suspense } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { FiHome, FiCommand, FiLock, FiSettings, FiServer, FiHelpCircle, FiChevronRight } from 'react-icons/fi'
import './docs.css'

// 直接导入文档组件，而不是使用懒加载
import GettingStarted from './GettingStarted'
import CommandOptions from './CommandOptions'
import UsageExamples from './UsageExamples'
import Security from './Security'
import Configuration from './Configuration'
import FAQ from './FAQ'

// 定义Loading组件
const Loading = () => (
  <div className="loading-container">
    <div className="loading-spinner"></div>
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
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

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
    setIsMobileMenuOpen(false) // 在移动端点击后关闭菜单
  }

  // 查找当前活动部分
  const currentSection = sections.find(section => section.id === activeSection) || sections[0]

  // 切换移动端菜单
  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen)
  }

  return (
    <div className="docs-container">
      {/* 移动设备菜单切换按钮 */}
      <div className="mobile-menu-toggle" onClick={toggleMobileMenu}>
        <span className="current-section">
          {currentSection.icon}
          {currentSection.name}
        </span>
        <FiChevronRight className={isMobileMenuOpen ? 'open' : ''} />
      </div>

      {/* 侧边栏导航 */}
      <div className={`docs-sidebar ${isMobileMenuOpen ? 'mobile-open' : ''}`}>
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