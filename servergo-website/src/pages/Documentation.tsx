import React from 'react';
import { Typography, Anchor, Row, Col, Card, Tabs } from 'antd';

const { Title, Paragraph, Text } = Typography;
const { Link } = Anchor;
const { TabPane } = Tabs;

const Documentation: React.FC = () => {
  return (
    <div className="documentation-container">
      <Row gutter={24}>
        {/* 侧边导航 */}
        <Col xs={24} md={6}>
          <Card style={{ position: 'sticky', top: '20px' }}>
            <Anchor affix={false}>
              <Link href="#installation" title="安装" />
              <Link href="#basic-usage" title="基本用法" />
              <Link href="#commands" title="命令参考">
                <Link href="#start-command" title="start 命令" />
                <Link href="#version-command" title="version 命令" />
              </Link>
              <Link href="#themes" title="主题" />
              <Link href="#authentication" title="认证" />
              <Link href="#advanced" title="高级功能" />
            </Anchor>
          </Card>
        </Col>

        {/* 主要内容 */}
        <Col xs={24} md={18}>
          <Title>ServerGo 文档</Title>
          <Paragraph>
            ServerGo 是一款简单易用的文件服务器工具，设计用于快速启动本地文件服务。
            本文档提供了 ServerGo 的详细使用说明和配置选项。
          </Paragraph>

          <Title level={2} id="installation">安装</Title>
          <Tabs defaultActiveKey="1">
            <TabPane tab="二进制下载" key="1">
              <Paragraph>
                从 GitHub Releases 下载预编译的二进制文件：
              </Paragraph>
              <ol>
                <li>访问 <a href="https://github.com/cc11001100/servergo/releases" target="_blank" rel="noopener noreferrer">GitHub Releases</a> 页面</li>
                <li>下载适合您操作系统的版本（Windows、macOS 或 Linux）</li>
                <li>解压文件</li>
                <li>将 servergo 可执行文件移动到您的 PATH 环境变量路径下</li>
              </ol>
            </TabPane>
            
            <TabPane tab="使用 Homebrew (macOS)" key="2">
              <Paragraph>
                如果您使用 macOS 和 Homebrew，可以通过以下命令安装：
              </Paragraph>
              <pre>
                brew tap cc11001100/servergo<br/>
                brew install servergo
              </pre>
            </TabPane>
            
            <TabPane tab="从源码编译" key="3">
              <Paragraph>
                如果您想从源码编译，需要 Go 1.16 或更高版本：
              </Paragraph>
              <pre>
                git clone https://github.com/cc11001100/servergo.git<br/>
                cd servergo<br/>
                go build -o servergo
              </pre>
            </TabPane>
          </Tabs>

          <Title level={2} id="basic-usage">基本用法</Title>
          <Paragraph>
            启动一个提供当前目录内容的文件服务器：
          </Paragraph>
          <pre>servergo start</pre>
          
          <Paragraph>
            指定端口和目录：
          </Paragraph>
          <pre>servergo start -p 8080 -d /path/to/directory</pre>
          
          <Paragraph>
            使用基本认证保护文件：
          </Paragraph>
          <pre>servergo start -a basic -u admin -w password</pre>

          <Title level={2} id="commands">命令参考</Title>
          
          <Title level={3} id="start-command">start 命令</Title>
          <Paragraph>
            <code>start</code> 命令用于启动文件服务器，支持以下参数：
          </Paragraph>
          <ul>
            <li><Text code>-p, --port</Text> - 指定服务器端口（默认：自动选择可用端口）</li>
            <li><Text code>-d, --directory</Text> - 指定要提供服务的目录（默认：当前目录）</li>
            <li><Text code>-a, --auth</Text> - 认证类型，目前支持 <Text code>basic</Text>（基本认证）</li>
            <li><Text code>-u, --username</Text> - 基本认证的用户名</li>
            <li><Text code>-w, --password</Text> - 基本认证的密码</li>
            <li><Text code>-m, --theme</Text> - 目录列表主题（default, dark, blue, green, retro）</li>
            <li><Text code>--no-open-browser</Text> - 不自动打开浏览器</li>
            <li><Text code>--no-dir-list</Text> - 禁用目录列表</li>
          </ul>

          <Title level={3} id="version-command">version 命令</Title>
          <Paragraph>
            <code>version</code> 命令用于显示 ServerGo 的版本信息：
          </Paragraph>
          <pre>servergo version</pre>

          <Title level={2} id="themes">主题</Title>
          <Paragraph>
            ServerGo 支持多种目录列表主题，可以通过 <Text code>-m</Text> 或 <Text code>--theme</Text> 参数指定：
          </Paragraph>
          <ul>
            <li><Text code>default</Text> - 默认浅色主题</li>
            <li><Text code>dark</Text> - 暗黑主题</li>
            <li><Text code>blue</Text> - 蓝色主题</li>
            <li><Text code>green</Text> - 绿色主题</li>
            <li><Text code>retro</Text> - 复古DOS风格主题</li>
          </ul>
          <pre>servergo start -m dark</pre>

          <Title level={2} id="authentication">认证</Title>
          <Paragraph>
            ServerGo 支持基本认证 (HTTP Basic Authentication) 来保护您的文件：
          </Paragraph>
          <pre>servergo start -a basic -u admin -w password</pre>
          
          <Paragraph>
            启用认证后，访问服务器的用户需要提供正确的用户名和密码才能访问文件。
          </Paragraph>

          <Title level={2} id="advanced">高级功能</Title>
          <Paragraph>
            <strong>自定义端口映射</strong>
          </Paragraph>
          <Paragraph>
            您可以设置特定的端口来启动服务：
          </Paragraph>
          <pre>servergo start -p 8080</pre>
          
          <Paragraph>
            <strong>禁用目录列表</strong>
          </Paragraph>
          <Paragraph>
            如果只想提供文件下载而不显示目录内容，可以禁用目录列表功能：
          </Paragraph>
          <pre>servergo start --no-dir-list</pre>
          
          <Paragraph>
            <strong>禁用自动打开浏览器</strong>
          </Paragraph>
          <Paragraph>
            默认情况下，ServerGo 会自动在浏览器中打开服务地址，如果不需要此功能：
          </Paragraph>
          <pre>servergo start --no-open-browser</pre>
        </Col>
      </Row>
    </div>
  );
};

export default Documentation; 