import { BrowserRouter, Routes, Route } from 'react-router-dom'

// 引入布局组件
import Layout from './components/Layout'

// 引入页面组件
import Home from './pages/Home'
import Docs from './pages/Docs'
import DocsIndex from './pages/docs/index'
import Examples from './pages/Examples'
import Download from './pages/Download'
import NotFound from './pages/NotFound'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="docs" element={<Docs />} />
          <Route path="docs/index" element={<DocsIndex />} />
          <Route path="examples" element={<Examples />} />
          <Route path="download" element={<Download />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
} 