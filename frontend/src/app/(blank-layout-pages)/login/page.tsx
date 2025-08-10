'use client'
import themeConfig from '@/configs/themeConfig'
import { UserContext } from '@/@contexts/User'
// Component Imports
import Login from '@views/Login'
import { useRouter } from 'next/navigation'
import { useContext } from 'react'


const LoginPage = () => {
  // Vars
  const router = useRouter()

  const { user } = useContext(UserContext)

  if (user) {
    router.push('/')
    return <></>
  }

  return <Login />
}

export default LoginPage
