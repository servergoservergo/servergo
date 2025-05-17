import { useState, useEffect, useRef } from 'react'
import { FiCopy, FiCheck, FiTerminal } from 'react-icons/fi'
import Prism from 'prismjs'
import 'prismjs/themes/prism-tomorrow.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-toml'
import './CodeBlock.css'

interface CodeBlockProps {
  code: string
  language?: string
  title?: string
  showLineNumbers?: boolean
}

export default function CodeBlock({ 
  code, 
  language = 'bash', 
  title,
  showLineNumbers = true
}: CodeBlockProps) {
  const [copied, setCopied] = useState(false)
  const codeRef = useRef<HTMLElement>(null)
  
  // 高亮代码
  useEffect(() => {
    if (codeRef.current) {
      Prism.highlightElement(codeRef.current)
    }
  }, [code, language])
  
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
        <div className="language-label">
          <FiTerminal />
          {title || `${language.toUpperCase()}`}
        </div>
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
      <pre className={showLineNumbers ? 'line-numbers' : ''}>
        <code ref={codeRef} className={`language-${language}`}>{code}</code>
      </pre>
    </div>
  )
} 