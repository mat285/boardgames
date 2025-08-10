"use client";

import type { GamesApi, UserApi, WebsocketApi, Api } from '@/@types/splendor/api'

import { MockApiClient } from './useMockApi'

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

// const API_ROUTE_GAME_STATE = API_ROUTE_GAME_BASE + '/:id/state'
// const API_ROUTE_GAME_PACKET = API_ROUTE_GAME_BASE + '/:id/packet'

const API_ROUTE_WEBSOCKETS = API_BASE_PATH + '/websockets'

const API_ROUTE_USER_CURRENT = API_ROUTE_USER_BASE

const apiGet = async (route: string, username: string, userID: string) => {
    // TODO: Add headers
    const response = await fetch(route, {
        headers: {
            'Content-Type': 'application/json',
            'X-Username': username,
            'X-UserID': userID
        }
    })

    if (!response.ok) {
        throw new Error('Failed to fetch data')
    }

    return response.json()
}

const apiPost = async (route: string, params: any, data: any, username: string, userID: string) => {
    for (const [key, value] of Object.entries(params)) {
        route = route.replace(`:${key}`, value as string)
    }

    // TODO: Add headers

    console.log('apiPost', route, data, username, userID)

    const response = await fetch(route, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type': 'application/json',
            'X-Username': username,
            'X-UserID': userID
        }
    })

    if (!response.ok) {
        throw new Error('Failed to fetch data')
    }

    return response.json()
}

const useGamesApi = (username: string, userID: string, isMock: boolean): GamesApi => {
    if (isMock) {
        return MockApiClient.games
    }

    return {
        getRegistry: async () => apiGet(API_ROUTE_REGISTRY, username, userID),
        getGames: async () => apiGet(API_ROUTE_GAMES_BASE, username, userID),
        getGame: MockApiClient.games.getGame, // async (id: string) => apiGet(API_ROUTE_GAME_BASE + '/' + id, username, userID),
        newGame: async (name: string) => apiPost(API_ROUTE_NEW_GAME, { name }, {}, username, userID),
        joinGame: async (id: string) => apiPost(API_ROUTE_JOIN_GAME, { id }, {}, username, userID),
        startGame: async (id: string) => apiPost(API_ROUTE_START_GAME, { id }, {}, username, userID),
        getGameState: MockApiClient.games.getGameState, // async (id: string) => apiGet(API_ROUTE_GAME_STATE + '/' + id, username, userID),
        getGamePacket: MockApiClient.games.getGamePacket, // async (id: string) => apiGet(API_ROUTE_GAME_PACKET + '/' + id, username, userID),
        getUserGames: async () => {
            const response = await apiGet(API_ROUTE_USER_GAMES, username, userID)

            console.log(response)

            return response
        }
    }
}

const useUserApi = (username: string, userID: string, isMock: boolean): UserApi => {
    if (isMock) {
        return MockApiClient.user
    }

    return {
        login: async (username: string) => {
            const response = await apiPost(API_ROUTE_USER_LOGIN, {}, { username }, username, userID)

            console.log(response)

            return response
        },
        getCurrentUser: async () => {
            const response = await apiGet(API_ROUTE_USER_CURRENT, username, userID)

            console.log(response)

            return response
        }
    }
}

const useWebsocketApi = (isMock: boolean): WebsocketApi => {
    if (isMock) {
        return MockApiClient.websocket
    }

    return {
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
}

const useApi = (username?: string, userID?: string): Api => {
    const isMock = false

    return {
        games: useGamesApi(username ?? '', userID ?? '', isMock),
        user: useUserApi(username ?? '', userID ?? '', isMock),
        websocket: useWebsocketApi(isMock),
    }
}


export default useApi
