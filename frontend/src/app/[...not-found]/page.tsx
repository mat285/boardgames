'use client'
// Component Imports
import themeConfig from '@/configs/themeConfig'
import Providers from '@components/Providers'
import BlankLayout from '@layouts/BlankLayout'
import NotFound from '@views/NotFound'

// Util Imports

const NotFoundPage = () => {
  // Vars
  const direction = 'ltr'

  return (
    <Providers direction={direction}>
      <BlankLayout>
        <NotFound mode={themeConfig.mode} />
      </BlankLayout>
    </Providers>
  )
}

export default NotFoundPage
