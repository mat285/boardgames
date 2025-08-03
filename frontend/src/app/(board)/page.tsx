// MUI Imports
import Grid from '@mui/material/Grid'

// Components Imports
import PlayerHand from '@/views/splendor/gameboard/PlayerHand'
import Board from '@/views/splendor/gameboard/Board'
import Card from '@mui/material/Card'
import { CardContent, Typography } from '@mui/material'

const GameBoard = () => {
  return (
    <Grid container xs={24} spacing={6}>
      {/* <Grid item xs={12} md={4}>
      </Grid>
      <Grid item xs={12} md={8} lg={8}>
      </Grid> */}
      {/* <Grid item xs={12} md={6} lg={4}>
      </Grid>
      <Grid item xs={12} md={6} lg={4}>
      </Grid> */}
      {/* <Grid item xs={12} md={6} lg={4}>
        <Grid container spacing={6}>
          <Grid item xs={12} sm={6}>
          </Grid>
          <Grid item xs={12} sm={6}>
          </Grid>
          <Grid item xs={12} sm={6}>
          </Grid>
          <Grid item xs={12} sm={6}>
          </Grid>
        </Grid>
      </Grid> */}
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Typography>
              Game Board
            </Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={12} md={12} lg={12}>
        <Board />
      </Grid>
      <Grid item xs={12}>
        <PlayerHand player={null} />
      </Grid>
    </Grid>
  )
}

export default GameBoard
