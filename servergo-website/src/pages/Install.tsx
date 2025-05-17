import { useState } from 'react'
import '../styles/Install.css'

// 安装选项卡类型
type InstallTab = 'windows' | 'macos' | 'linux' | 'docker'

// 安装方式类型
type InstallMethod = {
  id: string
  name: string
  description: string
  content: React.ReactNode
}

export default function Install() {
  const [activeTab, setActiveTab] = useState<InstallTab>('windows')
  const [activeMethod, setActiveMethod] = useState<string>('windows-installer')

  // Windows 安装方式
  const windowsMethods: InstallMethod[] = [
    {
      id: 'windows-installer',
      name: '安装程序',
      description: '使用安装向导快速安装',
      content: (
        <div className="install-instructions">
          <h3>使用安装向导安装 ServerGo</h3>
          <ol>
            <li>下载最新的 <a href="#" className="link">ServerGo 安装程序 (.exe)</a></li>
            <li>运行下载的安装程序并按照向导提示进行安装</li>
            <li>安装完成后，ServerGo 会自动添加到您的系统环境变量中</li>
            <li>打开命令提示符或 PowerShell 窗口，输入 <code>servergo --version</code> 验证安装</li>
          </ol>
        </div>
      )
    },
    {
      id: 'windows-portable',
      name: '便携版',
      description: '无需安装，直接解压使用',
      content: (
        <div className="install-instructions">
          <h3>使用便携版</h3>
          <ol>
            <li>下载最新的 <a href="#" className="link">ServerGo 便携版 (.zip)</a></li>
            <li>将下载的文件解压到您选择的文件夹中</li>
            <li>将解压目录添加到系统环境变量 PATH 中，或直接在解压目录下使用</li>
            <li>打开命令提示符或 PowerShell 窗口，输入 <code>servergo --version</code> 验证安装</li>
          </ol>
        </div>
      )
    },
    {
      id: 'windows-chocolatey',
      name: 'Chocolatey',
      description: '使用包管理器安装',
      content: (
        <div className="install-instructions">
          <h3>使用 Chocolatey 安装</h3>
          <p>如果您已经安装了 <a href="https://chocolatey.org/" target="_blank" rel="noopener noreferrer" className="link">Chocolatey</a>，可以使用以下命令安装 ServerGo：</p>
          <pre><code>choco install servergo</code></pre>
          <p>安装完成后，打开新的命令提示符或 PowerShell 窗口，输入以下命令验证安装：</p>
          <pre><code>servergo --version</code></pre>
        </div>
      )
    }
  ]

  // macOS 安装方式
  const macosMethods: InstallMethod[] = [
    {
      id: 'macos-homebrew',
      name: 'Homebrew',
      description: '使用 Homebrew 包管理器安装',
      content: (
        <div className="install-instructions">
          <h3>使用 Homebrew 安装</h3>
          <p>如果您已经安装了 <a href="https://brew.sh/" target="_blank" rel="noopener noreferrer" className="link">Homebrew</a>，可以使用以下命令安装 ServerGo：</p>
          <pre><code>brew install servergo</code></pre>
          <p>安装完成后，打开终端窗口，输入以下命令验证安装：</p>
          <pre><code>servergo --version</code></pre>
        </div>
      )
    },
    {
      id: 'macos-binary',
      name: '二进制安装',
      description: '下载二进制文件手动安装',
      content: (
        <div className="install-instructions">
          <h3>使用二进制文件安装</h3>
          <ol>
            <li>下载最新的 <a href="#" className="link">ServerGo macOS 二进制包 (.tar.gz)</a></li>
            <li>解压下载的文件：<code>tar -xzf servergo-macos-x64.tar.gz</code></li>
            <li>将二进制文件移动到合适的位置：<code>sudo mv servergo /usr/local/bin/</code></li>
            <li>确保文件有执行权限：<code>sudo chmod +x /usr/local/bin/servergo</code></li>
            <li>在终端中输入 <code>servergo --version</code> 验证安装</li>
          </ol>
        </div>
      )
    },
    {
      id: 'macos-installer',
      name: '安装程序',
      description: '使用安装包快速安装',
      content: (
        <div className="install-instructions">
          <h3>使用安装包安装</h3>
          <ol>
            <li>下载最新的 <a href="#" className="link">ServerGo macOS 安装包 (.pkg)</a></li>
            <li>双击打开下载的安装包并按照提示进行安装</li>
            <li>安装完成后，打开终端输入 <code>servergo --version</code> 验证安装</li>
          </ol>
        </div>
      )
    }
  ]

  // Linux 安装方式
  const linuxMethods: InstallMethod[] = [
    {
      id: 'linux-apt',
      name: 'APT (Debian/Ubuntu)',
      description: '使用 apt 包管理器安装',
      content: (
        <div className="install-instructions">
          <h3>在 Debian/Ubuntu 上使用 APT 安装</h3>
          <p>添加 ServerGo 存储库并安装：</p>
          <pre><code>curl -fsSL https://packages.servergo.dev/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://packages.servergo.dev/apt stable main"
sudo apt update
sudo apt install servergo</code></pre>
          <p>安装完成后，在终端中输入以下命令验证安装：</p>
          <pre><code>servergo --version</code></pre>
        </div>
      )
    },
    {
      id: 'linux-yum',
      name: 'YUM (RHEL/CentOS)',
      description: '使用 yum 包管理器安装',
      content: (
        <div className="install-instructions">
          <h3>在 RHEL/CentOS 上使用 YUM 安装</h3>
          <p>添加 ServerGo 存储库并安装：</p>
          <pre><code>sudo rpm --import https://packages.servergo.dev/gpg
sudo yum-config-manager --add-repo https://packages.servergo.dev/yum/servergo.repo
sudo yum install servergo</code></pre>
          <p>安装完成后，在终端中输入以下命令验证安装：</p>
          <pre><code>servergo --version</code></pre>
        </div>
      )
    },
    {
      id: 'linux-binary',
      name: '二进制安装',
      description: '下载二进制文件手动安装',
      content: (
        <div className="install-instructions">
          <h3>使用二进制文件安装</h3>
          <ol>
            <li>下载最新的 <a href="#" className="link">ServerGo Linux 二进制包 (.tar.gz)</a></li>
            <li>解压下载的文件：<code>tar -xzf servergo-linux-x64.tar.gz</code></li>
            <li>将二进制文件移动到合适的位置：<code>sudo mv servergo /usr/local/bin/</code></li>
            <li>确保文件有执行权限：<code>sudo chmod +x /usr/local/bin/servergo</code></li>
            <li>在终端中输入 <code>servergo --version</code> 验证安装</li>
          </ol>
        </div>
      )
    }
  ]

  // Docker 安装方式
  const dockerMethods: InstallMethod[] = [
    {
      id: 'docker-hub',
      name: 'Docker Hub',
      description: '从 Docker Hub 拉取镜像',
      content: (
        <div className="install-instructions">
          <h3>使用 Docker Hub 安装</h3>
          <p>从 Docker Hub 拉取最新的 ServerGo 镜像：</p>
          <pre><code>docker pull servergo/servergo:latest</code></pre>
          <p>运行 ServerGo 容器：</p>
          <pre><code>docker run -d -p 8080:8080 --name servergo servergo/servergo:latest</code></pre>
          <p>验证容器是否正常运行：</p>
          <pre><code>docker ps</code></pre>
        </div>
      )
    },
    {
      id: 'docker-compose',
      name: 'Docker Compose',
      description: '使用 Docker Compose 运行',
      content: (
        <div className="install-instructions">
          <h3>使用 Docker Compose 安装</h3>
          <p>创建 docker-compose.yml 文件：</p>
          <pre><code>version: '3'
services:
  servergo:
    image: servergo/servergo:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config:/app/config
    restart: unless-stopped</code></pre>
          <p>启动服务：</p>
          <pre><code>docker-compose up -d</code></pre>
          <p>验证服务是否正常运行：</p>
          <pre><code>docker-compose ps</code></pre>
        </div>
      )
    }
  ]

  // 根据当前选中的选项卡获取相应的安装方式
  const getActiveMethods = (): InstallMethod[] => {
    switch (activeTab) {
      case 'windows':
        return windowsMethods
      case 'macos':
        return macosMethods
      case 'linux':
        return linuxMethods
      case 'docker':
        return dockerMethods
      default:
        return windowsMethods
    }
  }

  // 切换选项卡时的处理函数
  const handleTabChange = (tab: InstallTab) => {
    setActiveTab(tab)
    // 自动选择新选项卡的第一个安装方式
    const methods = getActiveMethods()
    if (methods.length > 0) {
      setActiveMethod(methods[0].id)
    }
  }

  const methods = getActiveMethods()

  return (
    <div className="install-page">
      <h1>安装 ServerGo</h1>
      <p className="page-description">
        选择您的操作系统并按照相应的指南进行安装。ServerGo 支持多种操作系统和安装方式。
      </p>

      <div className="install-container">
        {/* 操作系统选项卡 */}
        <div className="os-tabs">
          <button 
            className={`os-tab-button ${activeTab === 'windows' ? 'active' : ''}`}
            onClick={() => handleTabChange('windows')}
          >
            Windows
          </button>
          <button 
            className={`os-tab-button ${activeTab === 'macos' ? 'active' : ''}`}
            onClick={() => handleTabChange('macos')}
          >
            macOS
          </button>
          <button 
            className={`os-tab-button ${activeTab === 'linux' ? 'active' : ''}`}
            onClick={() => handleTabChange('linux')}
          >
            Linux
          </button>
          <button 
            className={`os-tab-button ${activeTab === 'docker' ? 'active' : ''}`}
            onClick={() => handleTabChange('docker')}
          >
            Docker
          </button>
        </div>

        <div className="install-content">
          {/* 安装方式侧边栏 */}
          <div className="install-methods-sidebar">
            <h3>安装方式</h3>
            {methods.map(method => (
              <button
                key={method.id}
                className={`install-method-button ${activeMethod === method.id ? 'active' : ''}`}
                onClick={() => setActiveMethod(method.id)}
              >
                <div className="method-name">{method.name}</div>
                <div className="method-description">{method.description}</div>
              </button>
            ))}
          </div>

          {/* 安装步骤内容 */}
          <div className="install-method-content">
            {methods.find(m => m.id === activeMethod)?.content}
          </div>
        </div>
      </div>
    </div>
  )
} 