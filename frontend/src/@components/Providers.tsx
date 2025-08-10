'use client'

// Type Imports
import type { ChildrenType, Direction } from '@core/types'

// Context Imports
import { VerticalNavProvider } from '@menu/contexts/verticalNavContext'
import ThemeProvider from '@components/theme'

// Util Imports
import { UserProvider } from '@contexts/User'

type Props = ChildrenType & {
  direction: Direction
}

const Providers = (props: Props) => {
  // Props
  const { children, direction } = props

  return (
    <UserProvider>
      <VerticalNavProvider>
        <ThemeProvider direction={direction}>
          {children}
        </ThemeProvider>
      </VerticalNavProvider>
    </UserProvider>
  )
}

export default Providers
