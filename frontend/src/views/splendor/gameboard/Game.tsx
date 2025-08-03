'use client'
import { useEffect, useState } from 'react'

// MUI Imports
import { Grid } from '@mui/material'

import useApi from '@/hooks/splendor/api/useApi'

import type { GameState, Player } from '@/types/splendor/games'
import PlayerHand from './PlayerHand'
import Board from './Board'


const Game = () => {

    const { games }   = useApi();

    const [loading, setLoading] = useState(true);
    const [gameState, setGameState] = useState<GameState | null>(null);

    useEffect(() => {
        games.getGameState("").then((state) => {
            console.log(state);
            setGameState(state);
            setLoading(false);
        });
    }, []);


  // Vars
  return (
    <Grid container direction="column" xs={24} md={24} lg={24} spacing={6}>
        <Grid item xs={12}>
            <Grid container direction="row" xs={24} md={24} lg={24} spacing={6}>

            </Grid>
        </Grid>
        <Grid item xs={12}></Grid>
    <Grid container direction="row" xs={24} md={24} lg={24} spacing={6}>
        <Grid container xs={12} spacing={6}>
            <Grid item xs={12} md={12} lg={12}>
                <Board gameState={gameState} loading={loading} />
            </Grid>
            {gameState && (<PlayerHands gameState={gameState} />)}
        </Grid>
        </Grid>            
        </Grid>

  )
}

const PlayerHands = ({gameState}: {gameState: GameState}) => {
    const players = gameState.players;
    let user: Player | null = null;
    const otherPlayers: Player[] = [];

    players.forEach((player) => {
        if (player.id === gameState.currentPlayer) {
            user = player;
        } else {
            otherPlayers.push(player);
        }
    });

    return (
        <>
        {user && (
            <Grid item xs={12}>
                <PlayerHand player={user} orientation="horizontal"/>
            </Grid>
        )}
        {otherPlayers.map((player) => {
            return (
                <Grid item xs={12} key={player.id}>
                    <PlayerHand player={player} orientation="horizontal"/>
            </Grid>
            )
        })}
        </>
    )
}

export default Game
