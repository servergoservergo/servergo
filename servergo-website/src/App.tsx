import { useState } from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import { Layout, Menu, Button, theme } from 'antd'
import {
  HomeOutlined,
  RocketOutlined,
  ReadOutlined,
  InfoCircleOutlined,
  GithubOutlined,
} from '@ant-design/icons'

import Home from './pages/Home'
import GettingStarted from './pages/GettingStarted'
import Documentation from './pages/Documentation'
import About from './pages/About'

import './App.css'

const { Header, Content, Footer } = Layout

function App() {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken()

  return (
    <Router>
      <Layout className="main-layout">
        <Header className="header">
          <div className="logo">
            <Link to="/">ServerGo</Link>
          </div>
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['1']}
            items={[
              {
                key: '1',
                icon: <HomeOutlined />,
                label: <Link to="/">首页</Link>,
              },
              {
                key: '2',
                icon: <RocketOutlined />,
                label: <Link to="/getting-started">快速开始</Link>,
              },
              {
                key: '3',
                icon: <ReadOutlined />,
                label: <Link to="/documentation">文档</Link>,
              },
              {
                key: '4',
                icon: <InfoCircleOutlined />,
                label: <Link to="/about">关于</Link>,
              },
              {
                key: '5',
                icon: <GithubOutlined />,
                label: <a href="https://github.com/cc11001100/servergo" target="_blank" rel="noopener noreferrer">GitHub</a>,
              },
            ]}
            style={{ flex: 1 }}
          />
        </Header>
        <Content className="main-content">
          <div
            style={{
              padding: 24,
              minHeight: 380,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
            }}
          >
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/getting-started" element={<GettingStarted />} />
              <Route path="/documentation" element={<Documentation />} />
              <Route path="/about" element={<About />} />
            </Routes>
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          ServerGo ©{new Date().getFullYear()} | 简单易用的文件服务器
        </Footer>
      </Layout>
    </Router>
  )
}

export default App
