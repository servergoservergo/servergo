import React from 'react';
import { Typography, Steps, Card, Divider, Alert, Tabs } from 'antd';
import { DesktopOutlined, DownloadOutlined, CodeOutlined, RocketOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;
const { Step } = Steps;
const { TabPane } = Tabs;

const GettingStarted: React.FC = () => {
  return (
    <div style={{ textAlign: 'left', maxWidth: '900px', margin: '0 auto' }}>
      <Title level={1}>快速开始</Title>
      <Paragraph>
        ServerGo 是一款简单易用的文件服务器，无需复杂配置，几秒钟内即可启动一个本地文件服务。
        本指南将帮助您快速上手 ServerGo。
      </Paragraph>

      <Card style={{ marginTop: '30px', marginBottom: '30px' }}>
        <Steps direction="vertical" current={-1}>
          <Step 
            title="安装 ServerGo" 
            description={
              <div>
                <Paragraph>
                  您可以通过多种方式安装 ServerGo：
                </Paragraph>
                
                <Tabs defaultActiveKey="1">
                  <TabPane tab="二进制下载" key="1">
                    <Paragraph>
                      1. 前往 <a href="https://github.com/cc11001100/servergo/releases" target="_blank" rel="noopener noreferrer">GitHub Releases</a> 页面
                    </Paragraph>
                    <Paragraph>
                      2. 下载适合您操作系统的最新版本
                    </Paragraph>
                    <Paragraph>
                      3. 解压并移动到您的 PATH 路径中
                    </Paragraph>
                  </TabPane>
                  
                  <TabPane tab="Homebrew (macOS)" key="2">
                    <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace' }}>
                      brew tap cc11001100/servergo<br/>
                      brew install servergo
                    </div>
                  </TabPane>
                  
                  <TabPane tab="从源码编译" key="3">
                    <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace' }}>
                      git clone https://github.com/cc11001100/servergo.git<br/>
                      cd servergo<br/>
                      go build -o servergo
                    </div>
                  </TabPane>
                </Tabs>
              </div>
            }
            icon={<DownloadOutlined />}
          />

          <Step 
            title="基本使用" 
            description={
              <div>
                <Paragraph>
                  启动一个简单的文件服务器，浏览当前目录：
                </Paragraph>

                <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace' }}>
                  servergo start
                </div>
                
                <Paragraph style={{ marginTop: '15px' }}>
                  启动后，ServerGo 会自动在浏览器中打开服务地址。默认情况下，ServerGo 会选择一个可用的端口。
                </Paragraph>
              </div>
            }
            icon={<RocketOutlined />}
          />

          <Step 
            title="常用参数" 
            description={
              <div>
                <Paragraph>
                  ServerGo 支持多种参数来自定义服务器行为：
                </Paragraph>

                <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace', overflowX: 'auto' }}>
                  # 指定端口<br/>
                  servergo start -p 8080<br/><br/>
                  
                  # 使用基本认证<br/>
                  servergo start -a basic -u admin -w password<br/><br/>
                  
                  # 指定主题<br/>
                  servergo start -m dark<br/><br/>
                  
                  # 指定要服务的目录<br/>
                  servergo start -d /path/to/directory<br/><br/>
                </div>
              </div>
            }
            icon={<CodeOutlined />}
          />

          <Step 
            title="高级功能" 
            description={
              <div>
                <Paragraph>
                  ServerGo 提供了多种主题选择，包括：
                </Paragraph>
                <ul>
                  <li><Text code>default</Text> - 默认浅色主题</li>
                  <li><Text code>dark</Text> - 暗黑主题</li>
                  <li><Text code>blue</Text> - 蓝色主题</li>
                  <li><Text code>green</Text> - 绿色主题</li>
                  <li><Text code>retro</Text> - 复古DOS风格主题</li>
                </ul>
                
                <Paragraph>
                  使用特定主题：
                </Paragraph>
                
                <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace' }}>
                  servergo start -m retro
                </div>
              </div>
            }
            icon={<DesktopOutlined />}
          />
        </Steps>
      </Card>

      <Divider />

      <Title level={2}>帮助与支持</Title>
      <Paragraph>
        如果您需要查看更多命令和选项，可以使用：
      </Paragraph>
      <div style={{ background: '#f6f8fa', padding: '12px', borderRadius: '6px', fontFamily: 'monospace', marginBottom: '20px' }}>
        servergo --help
      </div>

      <Alert
        type="info"
        message="提示"
        description={
          <div>
            有关更多详细信息和高级用法，请参阅 <a href="/documentation">文档</a> 页面。
            如果您遇到任何问题，请在 <a href="https://github.com/cc11001100/servergo/issues" target="_blank" rel="noopener noreferrer">GitHub Issues</a> 上提交问题。
          </div>
        }
        showIcon
      />
    </div>
  );
};

export default GettingStarted; 