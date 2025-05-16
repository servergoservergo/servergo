import React from 'react';
import { Link } from 'react-router-dom';
import { Row, Col, Button, Typography, Card, Space } from 'antd';
import {
  RocketOutlined,
  ThunderboltOutlined,
  SecurityScanOutlined,
  AppstoreOutlined,
  BgColorsOutlined,
} from '@ant-design/icons';

const { Title, Paragraph } = Typography;

const Home: React.FC = () => {
  return (
    <div>
      {/* 英雄区域 */}
      <div className="home-hero">
        <Title level={1}>ServerGo</Title>
        <Paragraph style={{ fontSize: '18px', maxWidth: '800px', margin: '0 auto' }}>
          简单易用的文件服务器，让文件共享变得简单快捷
        </Paragraph>
        <div className="hero-buttons">
          <Space>
            <Button type="primary" size="large" icon={<RocketOutlined />}>
              <Link to="/getting-started">快速开始</Link>
            </Button>
            <Button size="large" icon={<ThunderboltOutlined />}>
              <a href="https://github.com/cc11001100/servergo/releases" target="_blank" rel="noopener noreferrer">
                下载最新版本
              </a>
            </Button>
          </Space>
        </div>
      </div>

      {/* 功能特点区域 */}
      <div className="feature-section">
        <Title level={2}>主要特性</Title>
        <Row gutter={[24, 24]} justify="center">
          <Col xs={24} md={12} lg={8}>
            <Card className="feature-card">
              <ThunderboltOutlined />
              <Title level={4}>快速部署</Title>
              <Paragraph>
                一键启动，零配置，快速搭建本地文件服务器，无需复杂设置。
              </Paragraph>
            </Card>
          </Col>

          <Col xs={24} md={12} lg={8}>
            <Card className="feature-card">
              <SecurityScanOutlined />
              <Title level={4}>安全可靠</Title>
              <Paragraph>
                支持基本认证(Basic Auth)，保护您的文件安全，防止未授权访问。
              </Paragraph>
            </Card>
          </Col>

          <Col xs={24} md={12} lg={8}>
            <Card className="feature-card">
              <AppstoreOutlined />
              <Title level={4}>功能丰富</Title>
              <Paragraph>
                支持目录浏览、文件下载、自定义端口等多种功能，满足各种需求。
              </Paragraph>
            </Card>
          </Col>

          <Col xs={24} md={12} lg={8}>
            <Card className="feature-card">
              <BgColorsOutlined />
              <Title level={4}>多主题支持</Title>
              <Paragraph>
                提供多种目录浏览主题，包括默认、暗黑、蓝色、绿色和复古DOS风格。
              </Paragraph>
            </Card>
          </Col>
        </Row>
      </div>

      {/* 主题展示区域 */}
      <div className="theme-showcase">
        <Title level={2}>精美主题</Title>
        <Paragraph style={{ marginBottom: '30px' }}>
          ServerGo提供多种目录浏览主题，让您的文件服务器更加美观
        </Paragraph>

        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={8}>
            <Card title="默认主题" className="theme-card">
              <div className="theme-preview" style={{ backgroundColor: '#fff', border: '1px solid #f0f0f0' }}>
                <div style={{ padding: '20px', textAlign: 'left' }}>
                  <div style={{ marginBottom: '10px', color: '#1890ff', fontWeight: 'bold' }}>目录: /example</div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px' }}>
                    <span style={{ marginRight: '5px' }}>📁</span> 
                    <span>documents/</span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px' }}>
                    <span style={{ marginRight: '5px' }}>📄</span> 
                    <span>readme.txt</span>
                  </div>
                </div>
              </div>
            </Card>
          </Col>

          <Col xs={24} sm={12} md={8}>
            <Card title="暗黑主题" className="theme-card">
              <div className="theme-preview" style={{ backgroundColor: '#141414', color: '#fff', border: '1px solid #303030' }}>
                <div style={{ padding: '20px', textAlign: 'left' }}>
                  <div style={{ marginBottom: '10px', color: '#177ddc', fontWeight: 'bold' }}>目录: /example</div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px', color: '#d9d9d9' }}>
                    <span style={{ marginRight: '5px' }}>📁</span> 
                    <span>documents/</span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px', color: '#d9d9d9' }}>
                    <span style={{ marginRight: '5px' }}>📄</span> 
                    <span>readme.txt</span>
                  </div>
                </div>
              </div>
            </Card>
          </Col>

          <Col xs={24} sm={12} md={8}>
            <Card title="复古DOS主题" className="theme-card">
              <div className="theme-preview" style={{ backgroundColor: '#000084', color: '#00ff00', border: '1px solid #00aaaa' }}>
                <div style={{ padding: '20px', textAlign: 'left', fontFamily: 'monospace' }}>
                  <div style={{ marginBottom: '10px', color: '#ffff00', fontWeight: 'bold' }}>目录: C:\\example</div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px' }}>
                    <span style={{ marginRight: '5px' }}>[DIR]</span> 
                    <span>documents/</span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px' }}>
                    <span style={{ marginRight: '5px' }}>[FILE]</span> 
                    <span>readme.txt</span>
                  </div>
                </div>
              </div>
            </Card>
          </Col>
        </Row>
      </div>
    </div>
  );
};

export default Home; 