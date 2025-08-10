// React Imports
import { useContext } from 'react'

// Context Imports
import { SettingsContext } from '@core/contexts/settingsContext'
import themeConfig from '@/configs/themeConfig'

export const useSettings = () => {
  // Hooks
  const context = useContext(SettingsContext)

  if (!context) {
    return {
      mode: themeConfig.mode,
      settings: {
        mode: themeConfig.mode,
      },
    }
  }

  return context
}
