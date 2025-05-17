import Card from '../components/Card'
import CodeBlock from '../components/CodeBlock'
import { FiCode, FiFolder, FiGlobe, FiServer, FiLock, FiUpload, FiFileText } from 'react-icons/fi'

const examples = [
  {
    id: 'local-sharing',
    title: '局域网文件共享',
    icon: <FiFolder />,
    description: '快速在本地网络共享文件，无需复杂设置',
    code: 'servergo -p 8080',
    explanation: '默认情况下，ServerGo绑定到所有网络接口，使得您可以通过局域网内的IP地址访问共享的文件。',
    notes: [
      '确保您的防火墙允许8080端口（或您选择的端口）',
      '局域网内的设备可以通过浏览器使用您电脑的IP地址访问，例如：http://192.168.1.100:8080',
      '适合临时共享大型文件或多个文件，无需使用U盘等物理媒介'
    ]
  },
  {
    id: 'development-server',
    title: '开发环境测试服务器',
    icon: <FiCode />,
    description: '为前端开发提供本地服务器，快速测试静态页面',
    code: 'servergo -p 3000 -d ./build --cors',
    explanation: '通过启用CORS选项，允许跨域请求，对于前端开发和API测试非常有用。',
    notes: [
      '指定特定目录（如构建输出目录）以专注于测试产品代码',
      '启用CORS支持允许您的前端应用从不同的源访问资源',
      '可以与前端构建工具集成，在构建后自动提供服务'
    ]
  },
  {
    id: 'protected-sharing',
    title: '受保护的文件共享',
    icon: <FiLock />,
    description: '使用基本认证和HTTPS加密保护共享文件',
    code: 'servergo --auth username:password --ssl --cert ./cert.pem --key ./key.pem',
    explanation: '添加基本的HTTP认证和SSL加密，提供安全的文件共享环境，适合敏感文档。',
    notes: [
      '使用强密码保护敏感文件',
      'HTTPS加密确保数据传输安全',
      '如果需要更高的安全性，建议使用专业的文件共享解决方案'
    ]
  },
  {
    id: 'upload-server',
    title: '临时上传服务器',
    icon: <FiUpload />,
    description: '允许他人上传文件到您的服务器',
    code: 'servergo --auth admin:secure_password',
    explanation: 'ServerGo默认启用文件上传功能，添加认证以确保只有授权用户可以上传文件。',
    notes: [
      '在Web界面中，用户可以通过"上传"按钮上传文件',
      '如果不需要上传功能，使用--no-upload选项禁用',
      '记得设置认证以防止未授权上传'
    ]
  },
  {
    id: 'documentation-hosting',
    title: '文档托管服务',
    icon: <FiFileText />,
    description: '托管团队或项目文档，方便内部访问',
    code: 'servergo -d ./docs -p 8088 --host 10.0.0.5',
    explanation: '指定文档目录和特定的IP地址，只在内部网络提供文档访问服务。',
    notes: [
      '适合托管内部文档、API文档或技术规范',
      '可以与静态网站生成器（如MkDocs、VuePress）集成',
      '指定特定的主机IP以限制访问范围'
    ]
  },
  {
    id: 'docker-integration',
    title: 'Docker容器化部署',
    icon: <FiServer />,
    description: '在Docker容器中运行ServerGo',
    code: `# Dockerfile
FROM alpine:latest
WORKDIR /app
COPY servergo /app/
EXPOSE 8080
VOLUME ["/data"]
CMD ["/app/servergo", "-d", "/data", "-p", "8080"]

# 构建并运行
docker build -t servergo .
docker run -p 8080:8080 -v $(pwd):/data servergo`,
    explanation: '将ServerGo容器化，便于在各种环境中一致部署。卷挂载允许容器访问主机文件系统。',
    notes: [
      '容器化使部署更简单、更一致',
      '可以在任何支持Docker的环境中运行',
      '通过卷挂载轻松管理共享内容'
    ]
  },
  {
    id: 'cloud-deployment',
    title: '云服务器部署',
    icon: <FiGlobe />,
    description: '在云服务器上部署ServerGo提供文件服务',
    code: `# 在云服务器上安装
curl -L https://github.com/cc11001100/servergo/releases/download/v1.0.0/servergo_v1.0.0_linux_amd64.tar.gz -o servergo.tar.gz
tar -xzf servergo.tar.gz
chmod +x servergo

# 使用systemd设置为服务
cat > /etc/systemd/system/servergo.service << EOF
[Unit]
Description=ServerGo File Server
After=network.target

[Service]
Type=simple
User=www-data
ExecStart=/path/to/servergo -p 80 -d /var/www/files --auth admin:secure_password
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable servergo
systemctl start servergo`,
    explanation: '在云服务器上部署ServerGo并设置为系统服务，实现持久运行和自动重启。',
    notes: [
      '确保云服务器的安全组/防火墙允许相应端口',
      '建议添加反向代理（如Nginx）以提供SSL和额外的安全层',
      '定期备份和更新以确保安全'
    ]
  }
]

export default function Examples() {
  return (
    <div>
      <h1 style={{ textAlign: 'center', marginBottom: '40px' }}>ServerGo 使用示例</h1>
      
      <p style={{ textAlign: 'center', maxWidth: '800px', margin: '0 auto 40px', lineHeight: 1.6 }}>
        ServerGo 是一个简单而强大的文件服务器，可以用于多种场景。以下是一些常见的使用案例和示例配置，帮助您充分利用ServerGo的功能。
      </p>
      
      <div id="examples-toc" style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', 
        gap: '20px',
        margin: '40px 0'
      }}>
        {examples.map(example => (
          <a 
            key={example.id}
            href={`#${example.id}`}
            style={{ 
              textDecoration: 'none',
              color: 'inherit'
            }}
          >
            <Card 
              style={{ height: '100%' }}
              className="card-hover"
            >
              <div style={{ 
                display: 'flex', 
                alignItems: 'center',
                marginBottom: '10px' 
              }}>
                <span style={{ 
                  marginRight: '12px',
                  fontSize: '24px',
                  color: 'var(--primary-color)'
                }}>
                  {example.icon}
                </span>
                <h3 style={{ margin: 0 }}>{example.title}</h3>
              </div>
              <p>{example.description}</p>
            </Card>
          </a>
        ))}
      </div>
      
      <div>
        {examples.map(example => (
          <section 
            key={example.id}
            id={example.id}
            style={{ 
              marginBottom: '60px',
              scrollMarginTop: '100px' // 为锚点导航留出空间
            }}
          >
            <h2 style={{ 
              display: 'flex', 
              alignItems: 'center',
              marginBottom: '20px' 
            }}>
              <span style={{ 
                marginRight: '15px',
                fontSize: '28px',
                color: 'var(--primary-color)'
              }}>
                {example.icon}
              </span>
              {example.title}
            </h2>
            
            <Card style={{ marginBottom: '30px' }}>
              <p style={{ fontSize: '18px', marginBottom: '20px' }}>{example.description}</p>
              <div style={{ marginBottom: '20px' }}>
                <h3 style={{ marginBottom: '10px' }}>示例命令</h3>
                <CodeBlock code={example.code} />
              </div>
              <p>{example.explanation}</p>
            </Card>
            
            <Card 
              title="使用说明与提示"
              style={{ backgroundColor: 'var(--bg-light)' }}
            >
              <ul style={{ paddingLeft: '20px' }}>
                {example.notes.map((note, index) => (
                  <li key={index} style={{ marginBottom: '8px' }}>{note}</li>
                ))}
              </ul>
            </Card>
          </section>
        ))}
      </div>
      
      <Card style={{ marginTop: '40px', marginBottom: '40px' }}>
        <h2 style={{ textAlign: 'center', margin: '0 0 20px' }}>创建您自己的使用场景</h2>
        <p style={{ textAlign: 'center' }}>
          ServerGo的灵活性使其适用于各种文件共享和静态服务场景。结合不同的选项，您可以创建满足特定需求的自定义配置。
          查看<a href="/docs">文档</a>了解所有可用选项，或在<a href="https://github.com/cc11001100/servergo/issues">GitHub</a>上分享您的使用案例！
        </p>
      </Card>
    </div>
  )
} 