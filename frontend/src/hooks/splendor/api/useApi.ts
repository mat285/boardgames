import type { GamesApi, UserApi, WebsocketApi, Api } from "@/types/splendor/api"

import { MockApiClient } from "./useMockApi" 

const API_BASE_PATH = '/api/v1alpha1'

const API_ROUTE_REGISTRY = API_BASE_PATH + '/registry'

const API_ROUTE_USER_BASE = API_BASE_PATH + '/user'
const API_ROUTE_USER_LOGIN = API_ROUTE_USER_BASE + '/login'
const API_ROUTE_USER_GAMES = API_ROUTE_USER_BASE + '/games'

const API_ROUTE_GAME_BASE = API_BASE_PATH + '/game' 
const API_ROUTE_GAMES_BASE = API_BASE_PATH + '/games'
const API_ROUTE_NEW_GAME = API_ROUTE_GAMES_BASE + '/:name/new'
const API_ROUTE_JOIN_GAME = API_ROUTE_GAME_BASE + '/:id/join'
const API_ROUTE_START_GAME = API_ROUTE_GAME_BASE + '/:id/start'
const API_ROUTE_GAME_STATE = API_ROUTE_GAME_BASE + '/:id/state'
const API_ROUTE_GAME_PACKET = API_ROUTE_GAME_BASE + '/:id/packet'

const API_ROUTE_WEBSOCKETS = API_BASE_PATH + '/websockets'

const apiGet = async (route: string) => {
    // TODO: Add headers
    const response = await fetch(route)

    if (!response.ok) {
        throw new Error('Failed to fetch data')
    }
    
    return response.json()
}

const apiPost = async (route: string, data: any) => {
    // TODO: Add headers
    const response = await fetch(route, {
        method: 'POST',
        body: JSON.stringify(data)
    })
    
    if (!response.ok) {
        throw new Error('Failed to fetch data')
    }   

    return response.json()
}

const GamesApiClient: GamesApi = {
        getRegistry: async () => apiGet(API_ROUTE_REGISTRY),
        getGames: async () => apiGet(API_ROUTE_GAMES_BASE),
        getGame: async (id: string) => apiGet(API_ROUTE_GAME_BASE + '/' + id),
        newGame: async (name: string) => apiPost(API_ROUTE_NEW_GAME, { name }),
        joinGame: async (id: string) => apiPost(API_ROUTE_JOIN_GAME, { id }),
        startGame: async (id: string) => apiPost(API_ROUTE_START_GAME, { id }),
        getGameState: async (id: string) => apiGet(API_ROUTE_GAME_STATE + '/' + id),
        getGamePacket: async (id: string) => apiGet(API_ROUTE_GAME_PACKET + '/' + id),
    }

const UserApiClient: UserApi = {
    getUserGames: () => apiGet(API_ROUTE_USER_GAMES),
    getUserLogin: (username: string, password: string) => apiPost(API_ROUTE_USER_LOGIN, { username, password }),
}

const WebsocketApiClient: WebsocketApi = {
    subscribeWebsocket: async () => {
        const ws = new WebSocket(API_ROUTE_WEBSOCKETS)
                  
            ws.onmessage = (event) => {
                console.log(event)
            }

            ws.onopen = () => {
                console.log('WebSocket connected')
            }

            ws.onclose = () => {
                console.log('WebSocket disconnected')
            }

        return ws
    }
}

const ApiClient: Api = {
    games: GamesApiClient,
    user: UserApiClient,
    websocket: WebsocketApiClient,
}

const useApi = (): Api => {
    const isMock = true

    if (isMock) {
        return MockApiClient
    }

    return ApiClient
}


export default useApi
