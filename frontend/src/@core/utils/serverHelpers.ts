'use client'

// Next Imports
import { Cookies } from 'react-cookie'

// Type Imports
import type { Settings } from '@core/contexts/settingsContext'
import type { SystemMode } from '@core/types'

// Config Imports
import themeConfig from '@configs/themeConfig'

export const getSettingsFromCookie = (): Settings => {
  const cookieStore = new Cookies();

  const cookieName = themeConfig.settingsCookieName

  return JSON.parse(cookieStore.get(cookieName)?.value || '{}')
}

export const getMode = () => {
  return themeConfig.mode
}

export const getSystemMode = (): SystemMode => {
  const mode = getMode()

  return mode
}

export const getServerMode = () => {
  const mode = getMode()

  return mode
}
