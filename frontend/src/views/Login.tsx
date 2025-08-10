'use client'

// React Imports
import { useContext, useState } from 'react'

// Next Imports
import Link from 'next/link'
import { useRouter } from 'next/navigation'

import CircularProgress from '@mui/material/CircularProgress'

// MUI Imports
import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'

// Hook Imports
import useApi from '@/@hooks/splendor/api/useApi'
import { UserContext } from '@/@contexts/User'

const Login = () => {
  const router = useRouter()

  const [usernameField, setUsernameField] = useState('')

  // Vars
  // const darkImg = '/images/pages/auth-v1-mask-dark.png'
  // const lightImg = '/images/pages/auth-v1-mask-light.png'

  const { user, setUser } = useContext(UserContext)

  // useEffect(() => {
  //   if (user?.username) {
  //     console.log('User is logged in')
  //     router.push('/')
  //     return
  //   }
  // }, [user])

  const api = useApi(usernameField, '')

  const [loading, setLoading] = useState(false)

  const handleSubmit = (e: any) => {
    setLoading(true)
    e.preventDefault()

    if (!usernameField) {
      throw new Error('Username is required')
    }

    api.user.login(usernameField).then((user) => {
      console.log('User logged in', user)
      setUser({
        username: user.username,
        id: user.id,
        password: ''
      })
      setLoading(false)
      router.replace('/')
    }).catch((err) => {
      console.error(err)
      setLoading(false)
    })
  }

  return (
    <div className='flex flex-col justify-center items-center min-bs-[100dvh] relative p-6'>
      <Card className='flex flex-col sm:is-[450px]'>
        <CardContent className='p-6 sm:!p-12'>
          <Link href='/' className='flex justify-center items-center mbe-6'>
            {/* <Logo /> */}
          </Link>
          <div className='flex flex-col gap-5'>
            <div>
              <Typography variant='h4'>{`Welcome to Splendor!`}</Typography>
              <Typography className='mbs-1'>Please enter a username to continue</Typography>
            </div>
            {/* <form noValidate autoComplete='off' onSubmit={handleSubmit} className='flex flex-col gap-5'> */}
            <TextField id='usernameField' autoFocus fullWidth label='Username' disabled={loading} value={usernameField} onChange={(e) => setUsernameField(e.target.value)} />
            <button type='button' onClick={handleSubmit}>
              Log In
            </button>
            {/* </form> */}
            {loading && <CircularProgress />}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default Login
