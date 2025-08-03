export interface User {
    id: string
    username: string
    password: string
}


export interface GamePacket {
    id: string
    packet: string
}

export interface GameList {
    games: Game[]
}


export interface Game {
    id: string
    name: string
    description: string
    image: string
    createdAt: string
}

export interface GameState {
    config: Config
    board: Board
    turn: number
    players: Player[]
    currentPlayer: string
    done: boolean
}

export interface Config {
    id: string
    name: string
    description: string
    image: string
}

export interface Player {
    id: string
    username: string
    hand: Hand
}

export interface Hand {
    cards: Card[]   
    reserved: Card[]
    bonus: Bonus[]
    gems: GemCount
}

export interface Board {
    gems: GemCount
    levelOne: Deck
    levelTwo: Deck
    levelThree: Deck
    bonuses: Bonus[]
}

export interface GemCount {
    emerald: number
    sapphire: number
    ruby: number
    diamond: number
    obsidian: number
    wild: number
}

export interface Deck {
    shown: Card[]
    hidden: Card[]
}

export interface Card {
    id: number
    level: number
    value: number
    type: string
    cost: GemCount
}

export interface Bonus {
    id: number
    name: string
    description: string
    image: string
}

