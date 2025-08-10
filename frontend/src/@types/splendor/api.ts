import type { Game, GamePacket, GameState, User } from './games'

export interface Api {
    games: GamesApi
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
    login: (username: string) => Promise<User>
    getCurrentUser: () => Promise<User>
}

export interface WebsocketApi {
    subscribeWebsocket: () => Promise<WebSocket>
}