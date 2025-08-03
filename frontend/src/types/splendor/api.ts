import type { Game, GamePacket, GameState } from "./games"
import type { User } from './games';

export interface Api {
    games:  GamesApi
    user: UserApi
    websocket: WebsocketApi
}

export interface GamesApi {
    getRegistry: () => Promise<Game[]>
    getGames: () => Promise<Game[]>
    getGame: (id: string) => Promise<Game>
    newGame: (name: string) => Promise<Game>
    joinGame: (id: string) => Promise<Game>
    startGame: (id: string) => Promise<Game>
    getGameState: (id: string) => Promise<GameState>
    getGamePacket: (id: string) => Promise<GamePacket>
    getUserGames: () => Promise<Game[]>
}

export interface UserApi {
    getUserLogin: (username: string, password: string) => Promise<User>
}

export interface WebsocketApi {
    subscribeWebsocket: () => Promise<WebSocket>
}