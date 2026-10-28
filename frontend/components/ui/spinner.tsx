import { Loader2Icon } from 'lucide-react'

import { cn } from '@/lib/utils'

type SpinnerSize = 'sm' | 'lg'

interface SpinnerProps extends React.ComponentProps<'svg'> {
  size?: SpinnerSize
}

function Spinner({ className, size, ...props }: SpinnerProps) {
  const sizeClass =
    size === 'sm'
      ? 'h-4 w-4'
      : size === 'lg'
      ? 'h-8 w-8'
      : undefined

  return (
    <Loader2Icon
      role="status"
      aria-label="Loading"
      className={cn('animate-spin', sizeClass, className)}
      {...props}
    />
  )
}

export { Spinner }
