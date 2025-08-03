import type { GameState, Game, GemCount, GamePacket, User, Card as SplendorCard } from "@/types/splendor/games"
import type { Api, GamesApi, UserApi, WebsocketApi } from "@/types/splendor/api"

const MockGamesApiClient: GamesApi = {
    getRegistry: async () => Promise.resolve([]), // apiGet(API_ROUTE_REGISTRY),
    getGames: async () => Promise.resolve([]), // apiGet(API_ROUTE_GAMES_BASE),
    getGame: async (id: string) => Promise.resolve({} as Game), // apiGet(API_ROUTE_GAME_BASE + '/' + id),
        newGame: async (name: string) => Promise.resolve({} as Game), // apiPost(API_ROUTE_NEW_GAME, { name }),
        joinGame: async (id: string) => Promise.resolve({} as Game), // apiPost(API_ROUTE_JOIN_GAME, { id }),
        startGame: async (id: string) => Promise.resolve({} as Game), // apiPost(API_ROUTE_START_GAME, { id }),
        getGameState: async (id: string) => {
            const cost: GemCount = {
                emerald: 1,
                sapphire: 1,
                ruby: 1,
                diamond: 10,
                obsidian: 1,
                wild: 1,
            }

            const mockShown: SplendorCard[] = [
                    {
                        id: 1,
                        level: 1,
                        value: 1,
                        type: "diamond",
                        cost: cost
                    },
                    {
                        id: 2,
                        level: 1,
                        value: 1,
                        type: "diamond",
                        cost: cost
                    },
                    {
                        id: 3,
                        level: 1,
                        value: 1,
                        type: "diamond",
                        cost: cost
                    },
                    {
                        id: 4,
                        level: 1,
                        value: 1,
                        type: "diamond",
                        cost: cost
                    }
            ]

            const mockState: GameState = {
                turn: 0,
                players: [],
                currentPlayer: "",
                done: false,
                config: {
                    id: "1",
                    name: "Splendor",
                    description: "A game of Splendor",
                    image: "https://via.placeholder.com/150",
                },
                board: {
                    gems: cost,
                    bonuses: [],
                    levelOne: {
                        shown: mockShown,
                        hidden: []
                    },
                    levelTwo: {
                        shown: mockShown,
                        hidden: []
                    },
                    levelThree: {
                        shown: mockShown,
                        hidden: []
                    }
                }
            }

            return Promise.resolve(mockState)
        },
    getGamePacket: async (id: string) => Promise.resolve({} as GamePacket), // apiGet(API_ROUTE_GAME_PACKET + '/' + id),
}

const MockUserApiClient: UserApi = {
    getUserGames: async () => Promise.resolve([]), // apiGet(API_ROUTE_USER_GAMES),
    getUserLogin: async (username: string, password: string) => Promise.resolve({} as User), // apiPost(API_ROUTE_USER_LOGIN, { username, password }),
}


const MockWebsocketApiClient: WebsocketApi = {
    subscribeWebsocket: async () => Promise.resolve({} as WebSocket), // apiGet(API_ROUTE_WEBSOCKETS),
}

export const MockApiClient: Api = {
    games: MockGamesApiClient,
    user: MockUserApiClient,
    websocket: MockWebsocketApiClient,
}

const useMockApi = (): Api => {
    return MockApiClient
}

export default useMockApi