'use client'

// React Imports
import type { CSSProperties } from 'react'

// Third-party Imports
import styled from '@emotion/styled'

// Component Imports
import { Grid } from '@mui/material'

import LogoIcon from '@core/svg/Logo'

type LogoTextProps = {
  color?: CSSProperties['color']
}

const LogoText = styled.span<LogoTextProps>`
  color: ${({ color }) => color ?? 'var(--mui-palette-text-primary)'};
  font-size: 1.25rem;
  line-height: 1.2;
  font-weight: 600;
  letter-spacing: 0.15px;
  text-transform: uppercase;
  margin-inline-start: 10px;
`

const Logo = ({ color }: { color?: CSSProperties['color'] }) => {
  return (
    <div className='flex items-center min-bs-[24px]'>
      <Grid container xs={6} spacing={0}>
        <Grid item xs={6}>
          <LogoIcon />
        </Grid>
        <Grid item xs={8} mt={-6}>
          <div className='flex flex-col'>
            <LogoText color={color}>Splendor</LogoText>
          </div>
        </Grid>
      </Grid>
    </div>
  )
}

export default Logo
