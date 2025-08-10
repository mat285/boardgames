'use client'

// Third-party Imports
import classnames from 'classnames'

// Hook Imports
import { verticalLayoutClasses } from '@layouts/utils/layoutClasses'

const FooterContent = () => {

  return (
    <div 
      className={classnames(verticalLayoutClasses.footerContent, 'flex items-center justify-between flex-wrap gap-4')}
    >
    </div>
  )
}

export default FooterContent
