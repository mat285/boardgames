'use client'
// Component Imports
import NotFound from '@views/NotFound'

import themeConfig from '@/configs/themeConfig'

const Error = () => {
  // Vars
  return <NotFound mode={themeConfig.mode} />
}

export default Error
