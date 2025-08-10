'use client'
// MUI Imports

import NotFound from '@views/NotFound'
import type { Game } from '@types/splendor/games'
import GameComponent from '@views/splendor/gameboard/Game'
import { useParams } from 'next/navigation'

import { useContext, useEffect, useState } from 'react'
import useApi from '@hooks/splendor/api/useApi'
import { redirect } from 'next/navigation'
import { UserContext } from '@contexts/User'
import themeConfig from '@configs/themeConfig'
import CircularProgress from '@mui/material/CircularProgress'

const GameBoard = () => {
  console.log("GameBoard")
  
  const { user } = useContext(UserContext);

  if (!user) {
    redirect('/login');
  }

  const { games } = useApi(user?.username, user?.id)

  const params = useParams()

  const gameID = params.game as string

  console.log('gameID', gameID)


  const [game, setGame] = useState<Game | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!gameID) {
      return;
    }

    if (loading) {
      return;
    }

    setLoading(true)
    games.getGame(gameID).then((game) => {
      console.log(game)
      setGame(game)
      setLoading(false)
    })
    .finally(() => {
      setLoading(false)
    })
  }, [gameID, games])


  if (!gameID) {
    console.log("No game ID")

    return <NotFound mode={themeConfig.mode} />
  }

  return (
    <>
      {!loading && game && <GameComponent game={game} />}
      {!loading && !game && <NotFound mode={themeConfig.mode} />}
      {loading && <CircularProgress />}
    </>
  )
}

export default GameBoard
