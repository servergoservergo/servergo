import { useState } from 'react'
import { FiCopy, FiCheck } from 'react-icons/fi'

interface CodeBlockProps {
  code: string
  language?: string
}

export default function CodeBlock({ code, language = 'bash' }: CodeBlockProps) {
  const [copied, setCopied] = useState(false)
  
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      
      // 3秒后重置复制状态
      setTimeout(() => {
        setCopied(false)
      }, 3000)
    } catch (err) {
      console.error('无法复制文本: ', err)
    }
  }
  
  return (
    <div className="code-block">
      <div className="code-block-header">
        <button 
          className="copy-button" 
          onClick={handleCopy}
          title={copied ? '已复制!' : '复制到剪贴板'}
        >
          {copied ? <FiCheck /> : <FiCopy />}
          <span style={{ marginLeft: '5px' }}>
            {copied ? '已复制' : '复制'}
          </span>
        </button>
      </div>
      <pre>
        <code className={`language-${language}`}>{code}</code>
      </pre>
    </div>
  )
} 