'use client'

// MUI Imports
import Card from '@mui/material/Card'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'

// Components Imports
import { Grid } from '@mui/material'

import SplendorCard from '@components/games/spendor/cards/Card'

import type { Board, GameState } from '@types/splendor/games'
import GemBank from './GemBank'
import Bonus from '@components/games/spendor/bonus/Bonus'


export type BoardProps = {
    gameState: GameState | null
    loading: boolean
}


const Board = (props: BoardProps) => {
    const { gameState, loading } = props;
    const board = gameState?.board;

    return (
        <Card sx={{ width: '100%' }}>
            {loading && (
                <CardContent>
                    <Typography>Loading...</Typography>
                </CardContent>
            )}
            {!loading && board && (
                <CardContent>
                    <Grid container spacing={8}>
                        <Grid item xs={12}>
                            <Grid container spacing={8}>
                                <Grid item xs={3} />
                                <Grid item xs={6}>
                                    <Card>
                                        <CardContent>
                                            <Grid container spacing={2}>
                                                <Grid item xs={3}>
                                                    <Bonus bonus={board.bonuses[0]} />
                                                </Grid>
                                                <Grid item xs={3}>
                                                    <Bonus bonus={board.bonuses[1]} />
                                                </Grid>
                                                <Grid item xs={3}>
                                                    <Bonus bonus={board.bonuses[2]} />
                                                </Grid>
                                                <Grid item xs={3}>
                                                    <Bonus bonus={board.bonuses[3]} />
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
                                                            <SplendorCard card={board?.levelOne.shown[0]} />
                                                        </Grid>
                                                        <Grid item xs={2} md={2}>
                                                            <SplendorCard card={board?.levelOne.shown[1]} />
                                                        </Grid>
                                                        <Grid item xs={2} md={2}>
                                                            <SplendorCard card={board?.levelOne.shown[2]} />
                                                        </Grid>
                                                        <Grid item xs={2} md={2}>
                                                            <SplendorCard card={board?.levelOne.shown[3]} />
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
                                                    <SplendorCard card={board?.levelTwo.shown[0]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelTwo.shown[1]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelTwo.shown[2]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelTwo.shown[3]} />
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
                                                    <SplendorCard card={board?.levelThree.shown[0]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelThree.shown[1]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelThree.shown[2]} />
                                                </Grid>
                                                <Grid item xs={2} md={2}>
                                                    <SplendorCard card={board?.levelThree.shown[3]} />
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
                                    <GemBank gems={board.gems} />
                                </Grid>
                                <Grid item xs={2}></Grid>
                            </Grid>
                        </Grid>
                        <Grid item xs={8} md={4}>

                        </Grid>
                        <Grid item xs={8} md={4}>

                        </Grid>
                    </Grid>
                </CardContent>)}
        </Card>
    )
}

export default Board
