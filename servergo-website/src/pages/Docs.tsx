import { Navigate } from 'react-router-dom'

export default function Docs() {
  // 重定向到文档索引页
  return <Navigate to="/docs/index" replace />
} 