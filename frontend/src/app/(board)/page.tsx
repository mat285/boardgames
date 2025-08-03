// MUI Imports
import Grid from '@mui/material/Grid'

// Components Imports
import PlayerHand from '@/views/gameboard/PlayerHand'
import Board from '@/views/gameboard/Board'

const GameBoard = () => {
  return (
    <Grid container spacing={6}>
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
      <Grid item xs={12} md={12} lg={12}>
        <Board />
      </Grid>
      <Grid item xs={12}>
        <PlayerHand />
      </Grid>
    </Grid>
  )
}

export default GameBoard
