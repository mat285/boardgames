'use client'

// MUI Imports
import Card from '@mui/material/Card'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'

// Components Imports
import { Grid } from '@mui/material'

import type { GemCount } from '@/types/splendor/games'


type GemBankProps = {
    gems: GemCount
}

const GemBank = (props: GemBankProps) => {
    const { gems } = props; 

  return (
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
                    <Typography>{gems.ruby}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{gems.sapphire}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{gems.emerald}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{gems.diamond}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{gems.obsidian}</Typography>
                </Grid>
                <Grid item xs={2}>
                    <Typography>{gems.wild}</Typography>
                </Grid>
            </Grid>
        </CardContent>
    </Card>  
    )
}

export default GemBank
