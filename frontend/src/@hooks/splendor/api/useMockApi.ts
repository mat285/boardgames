import type { GameState, Game, GemCount, GamePacket, User, Card as SplendorCard } from "@/@types/splendor/games"
import type { Api, GamesApi, UserApi, WebsocketApi } from "@/@types/splendor/api"

const MockGamesApiClient: GamesApi = {
    getRegistry: async () => Promise.resolve([]), // apiGet(API_ROUTE_REGISTRY),
    getGames: async () => Promise.resolve([]), // apiGet(API_ROUTE_GAMES_BASE),
    getGame: async (id: string) => Promise.resolve({id} as Game), // apiGet(API_ROUTE_GAME_BASE + '/' + id),
    newGame: async (name: string) => Promise.resolve({id: "1", name, description: "Description", image: "Image"} as Game), // apiPost(API_ROUTE_NEW_GAME, { name }),
    joinGame: async (id: string) => Promise.resolve({id: "1", name: "Game 1", description: "Description", image: "Image"} as Game), // apiPost(API_ROUTE_JOIN_GAME, { id }),
    startGame: async (id: string) => Promise.resolve({id: "1", name: "Game 1", description: "Description", image: "Image"} as Game), // apiPost(API_ROUTE_START_GAME, { id }),
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
            players: [
                {
                    id: "1",
                    username: "Player 1",
                    hand: {
                        cards: [],
                        reserved: [],
                        bonus: [],
                        gems: cost,
                    },
                },
                {
                    id: "2",
                    username: "Player 2",
                    hand: {
                        cards: [],
                        reserved: [],
                        bonus: [],
                        gems: cost,
                    },
                },
                {
                    id: "3",
                    username: "Player 3",
                    hand: {
                        cards: [],
                        reserved: [],
                        bonus: [],
                        gems: cost,
                    },
                },
            ],
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
                bonuses: [
                    {
                        id: 1,
                        name: "Bonus 1",
                        description: "Bonus 1",
                        image: "https://via.placeholder.com/150",
                    },
                    {
                        id: 2,
                        name: "Bonus 2",
                        description: "Bonus 2",
                        image: "https://via.placeholder.com/150",
                    },
                    {
                        id: 3,
                        name: "Bonus 3",
                        description: "Bonus 3",
                        image: "https://via.placeholder.com/150",
                    },
                    {
                        id: 4,
                        name: "Bonus 4",
                        description: "Bonus 4",
                        image: "https://via.placeholder.com/150",
                    },
                ],
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
    getUserGames: async () => Promise.resolve([]), // apiGet(API_ROUTE_USER_GAMES),
    getGamePacket: async (id: string) => Promise.resolve({} as GamePacket), // apiGet(API_ROUTE_GAME_PACKET + '/' + id),
}

export const MockUserApiClient: UserApi = {
    login: async (username: string) => Promise.resolve({} as User), // apiPost(API_ROUTE_USER_LOGIN, { username, password }),
    getCurrentUser: async () => Promise.resolve({} as User), // apiGet(API_ROUTE_USER_CURRENT),
}

export const MockWebsocketApiClient: WebsocketApi = {
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