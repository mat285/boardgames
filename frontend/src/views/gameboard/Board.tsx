'use client'

// Next Imports
import dynamic from 'next/dynamic'

// MUI Imports
import Card from '@mui/material/Card'
import { useTheme } from '@mui/material/styles'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'

// Third Party Imports
import type { ApexOptions } from 'apexcharts'

// Components Imports
import { colors, Grid } from '@mui/material'

const Board = () => {
  // Hooks
  const theme = useTheme()

  // Vars
  return (
    <Card sx={{ width: '100%' }}>
      <CardContent>
        <Grid container spacing={8}>
        <Grid item xs={12}>
        <Grid container spacing={8}>
        <Grid item xs={3}/>
          <Grid item xs={6}>
            <Card>
                <CardContent>
                    <Grid container spacing={2}>
                        <Grid item xs={3}>
                            <Card>
                                <CardContent>
                                    <Typography>Bonus 1</Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item xs={3}>
                            <Card>
                                <CardContent>
                                    <Typography>Bonus 2</Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item xs={3}>
                            <Card>
                                <CardContent>
                                    <Typography>Bonus 3</Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                        <Grid item xs={3}>
                            <Card>
                                <CardContent>
                                    <Typography>Bonus 4</Typography>
                                </CardContent>
                            </Card>
                            </Grid>
                    </Grid>
                </CardContent>
            </Card>
          </Grid>
          <Grid item xs={3}></Grid>
          <Grid item xs={12}>
        <Grid container spacing={8}>
          <Grid item xs={12}>
          <Card>
              <CardContent>
                <Grid container xs={12} spacing={6}>
                    <Grid item xs={1} md={2}></Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 1</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 2</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 3</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 4</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={1} md={2}></Grid>
                </Grid>
              </CardContent>
            </Card>
            </Grid>
            </Grid>
          </Grid>
        </Grid>
        <Grid container spacing={6}>
        <Grid item xs={12}></Grid>
          <Grid item xs={12}>
          <Card>
              <CardContent>
                <Grid container xs={12} spacing={6}>
                    <Grid item xs={1} md={2}></Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 1</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 2</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 3</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 4</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={1} md={2}></Grid>
                </Grid>
              </CardContent>
            </Card>
            </Grid>
            </Grid>
            <Grid container spacing={6}>
                <Grid item xs={12}></Grid>
                <Grid item xs={12}>
                    <Card>
              <CardContent>
                <Grid container xs={12} spacing={6}>
                    <Grid item xs={1} md={2}></Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                            <Typography variant='h4'>Card 1</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 2</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 3</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={2} md={2}>
                        <Card>
                            <CardContent>
                                <Typography variant='h4'>Card 4</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                    <Grid item xs={1} md={2}></Grid>
                </Grid>
              </CardContent>
            </Card>
            </Grid>
            </Grid>
            <Grid item xs={12}></Grid>
            <Grid container spacing={8}>
                <Grid item xs={12}></Grid>
                <Grid item xs={2}></Grid>
                <Grid item xs={8}>
                    <Card>
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
                                    <Typography>5</Typography>
                                </Grid>
                                <Grid item xs={2}>
                                    <Typography>4</Typography>
                                </Grid>
                                <Grid item xs={2}>
                                    <Typography>6</Typography>
                                </Grid>
                                <Grid item xs={2}>
                                    <Typography>3</Typography>
                                </Grid>
                                <Grid item xs={2}>
                                    <Typography>2</Typography>
                                </Grid>
                                <Grid item xs={2}>
                                    <Typography>1</Typography>
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={2}></Grid>
            </Grid>
          </Grid>
          <Grid item xs={8} md={4}>
            
          </Grid>
          <Grid item xs={8} md={4}>
           
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default Board
