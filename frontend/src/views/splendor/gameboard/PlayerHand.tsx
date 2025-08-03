'use client'

// MUI Imports
import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'

// Components Imports
import { Grid } from '@mui/material'

import type { Player } from '@/types/splendor/games'

export type PlayerHandProps = {
  player: Player | null
}


const PlayerHand = (props: PlayerHandProps) => {
  let { player } = props

  if (!player) {

    player = {
      id: "0",
      username: 'Player',
      hand: {
        gems: {
          ruby: 0,
          sapphire: 0,
          emerald: 0,
          diamond: 0,
          obsidian: 0,
          wild: 0,
        },
        cards: [],
        reserved: [],
        bonus: [],
      },
    }
  }

  return (
    <Card>
      <CardHeader
        title={player.username}
      />
      <CardContent>
        <Grid container spacing={2}>
          <Grid item xs={8} md={4}>
            <Card>
              <CardHeader title='Gems' />
              <CardContent>
              <Grid container spacing={8}>
                <Grid item xs={2}>
                    <Typography>Rubys</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>Sapphire</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>Emeralds</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>Diamonds</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>Obsidian</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>Wild</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.ruby}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.sapphire}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.emerald}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.diamond}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.obsidian}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{player.hand.gems.wild}</Typography>
                </Grid>
            </Grid>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={8} md={4}>
            <Card>
              <CardHeader title='Cards' />
              <CardContent>
                <Typography variant='h4'>100</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={8} md={4}>
            <Card>
              <CardHeader title='Reserved' />
              <CardContent>
                <Typography variant='h4'>100</Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default PlayerHand
