import type { ReactNode, CSSProperties } from 'react'
import './Card.css'

interface CardProps {
  title?: ReactNode;
  children: ReactNode;
  extra?: ReactNode;
  style?: CSSProperties;
  className?: string;
  variant?: 'default' | 'primary' | 'info' | 'success' | 'warning';
  bordered?: boolean;
  hoverable?: boolean;
}

export default function Card({ 
  title, 
  children, 
  extra, 
  style, 
  className, 
  variant = 'default',
  bordered = true,
  hoverable = false
}: CardProps) {
  const cardClasses = [
    'card',
    `card-${variant}`,
    bordered ? 'card-bordered' : '',
    hoverable ? 'card-hoverable' : '',
    className || ''
  ].filter(Boolean).join(' ');

  return (
    <div className={cardClasses} style={style}>
      {title && (
        <div className="card-header">
          <div className="card-title">{title}</div>
          {extra && <div className="card-extra">{extra}</div>}
        </div>
      )}
      {!title && extra && <div className="card-extra-only">{extra}</div>}
      <div className="card-content">{children}</div>
    </div>
  )
} 