'use client'

// Type Imports
import type { ChildrenType } from '@core/types'

// Component Imports
import RootLayout from '@layouts/Root'

const Layout = ({ children }: ChildrenType) => {
  return (
    <RootLayout>
      {children}
    </RootLayout>
  )
}

export default Layout
