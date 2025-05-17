import type { ReactNode, CSSProperties } from 'react'

interface CardProps {
  title?: ReactNode;
  children: ReactNode;
  extra?: ReactNode;
  style?: CSSProperties;
  className?: string;
}

export default function Card({ title, children, extra, style, className }: CardProps) {
  return (
    <div className={`card ${className || ''}`} style={style}>
      {title && (
        <div className="card-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <div className="card-title">{title}</div>
          {extra && <div className="card-extra">{extra}</div>}
        </div>
      )}
      {!title && extra}
      <div className="card-content">{children}</div>
    </div>
  )
} 