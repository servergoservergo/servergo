import React from 'react';
import { Typography, Card, Avatar, Divider, Button, Row, Col, Timeline } from 'antd';
import { GithubOutlined, HeartOutlined, StarOutlined, TeamOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

const About: React.FC = () => {
  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', textAlign: 'left' }}>
      <Title level={1} style={{ textAlign: 'center', marginBottom: '40px' }}>
        关于 ServerGo
      </Title>

      <Row gutter={[24, 24]}>
        <Col xs={24} md={16}>
          <Card>
            <Title level={2}>项目简介</Title>
            <Paragraph>
              ServerGo 是一款简单易用的文件服务器工具，专为开发者和日常用户设计。
              它允许用户快速启动一个本地文件服务器，无需复杂配置，几秒钟内即可实现文件共享。
            </Paragraph>
            
            <Paragraph>
              无论您是需要在本地网络中共享文件，还是为前端项目提供静态资源服务，
              ServerGo 都能通过简单的命令行界面满足您的需求。
            </Paragraph>

            <Title level={3}>核心功能</Title>
            <ul>
              <li>快速启动本地文件服务器</li>
              <li>支持目录浏览和文件下载</li>
              <li>多种主题可供选择</li>
              <li>基本认证保护文件安全</li>
              <li>自动选择可用端口</li>
              <li>支持跨平台（Windows, macOS, Linux）</li>
            </ul>

            <Title level={3}>技术栈</Title>
            <Paragraph>
              ServerGo 使用 Go 语言开发，采用以下技术和框架：
            </Paragraph>
            <ul>
              <li><Text strong>Gin</Text> - 高性能 Web 框架</li>
              <li><Text strong>Go 嵌入</Text> - 用于管理静态资源</li>
              <li><Text strong>HTML Templates</Text> - 目录列表渲染</li>
              <li><Text strong>CLI</Text> - 命令行界面</li>
            </ul>
          </Card>
        </Col>

        <Col xs={24} md={8}>
          <Card>
            <div style={{ textAlign: 'center', marginBottom: '20px' }}>
              <Avatar size={100} src="https://avatars.githubusercontent.com/u/5215596" />
              <Title level={4} style={{ marginTop: '15px', marginBottom: '0' }}>CC11001100</Title>
              <Text type="secondary">项目作者</Text>
              <div style={{ marginTop: '15px' }}>
                <Button type="link" icon={<GithubOutlined />} href="https://github.com/cc11001100" target="_blank">
                  GitHub
                </Button>
              </div>
            </div>
            <Divider />
            <Title level={4}>项目统计</Title>
            <p>
              <StarOutlined /> Stars: <Text strong>300+</Text>
            </p>
            <p>
              <TeamOutlined /> Contributors: <Text strong>5+</Text>
            </p>
            <p>
              <HeartOutlined /> Latest version: <Text strong>v1.0.0</Text>
            </p>
          </Card>
          
          <Card style={{ marginTop: '20px' }}>
            <Title level={4}>支持项目</Title>
            <Paragraph>
              如果您觉得 ServerGo 对您有帮助，请考虑：
            </Paragraph>
            <ul>
              <li>在 GitHub 上给项目一个星标</li>
              <li>报告问题或提交功能请求</li>
              <li>贡献代码或文档</li>
            </ul>
            <div style={{ textAlign: 'center', marginTop: '15px' }}>
              <Button type="primary" icon={<GithubOutlined />} href="https://github.com/cc11001100/servergo" target="_blank">
                访问 GitHub
              </Button>
            </div>
          </Card>
        </Col>
      </Row>

      <Card style={{ marginTop: '30px' }}>
        <Title level={2}>项目历程</Title>
        <Timeline
          items={[
            {
              color: 'green',
              children: (
                <>
                  <Title level={4}>2023年5月 - 项目发布</Title>
                  <Paragraph>
                    ServerGo 的第一个版本发布，实现了基本的文件服务功能和目录浏览。
                  </Paragraph>
                </>
              ),
            },
            {
              children: (
                <>
                  <Title level={4}>2023年8月 - 多主题支持</Title>
                  <Paragraph>
                    增加了多主题支持，包括默认和暗黑主题。
                  </Paragraph>
                </>
              ),
            },
            {
              children: (
                <>
                  <Title level={4}>2023年10月 - 添加认证机制</Title>
                  <Paragraph>
                    实现了基本认证功能，提高了文件服务的安全性。
                  </Paragraph>
                </>
              ),
            },
            {
              color: 'red',
              children: (
                <>
                  <Title level={4}>2024年5月 - 扩展主题系统</Title>
                  <Paragraph>
                    新增蓝色、绿色和复古DOS风格主题，丰富用户界面选择。
                  </Paragraph>
                </>
              ),
            },
          ]}
        />

        <Divider />

        <Title level={2}>开源许可</Title>
        <Paragraph>
          ServerGo 在 MIT 许可证下发布。这意味着您可以自由地使用、修改和分发此软件，
          无论是在个人项目还是商业项目中。请查看 GitHub 仓库中的 LICENSE 文件了解更多信息。
        </Paragraph>
      </Card>
    </div>
  );
};

export default About; 