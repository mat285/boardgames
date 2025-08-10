'use client'

import type { Dispatch, ReactNode, SetStateAction } from 'react'
import CircularProgress from '@mui/material/CircularProgress'
import { createContext, useEffect, useState } from 'react'

import useApi from '@hooks/splendor/api/useApi'

import type { User } from '@types/splendor/games'
import { redirect } from 'next/navigation'

export interface UserContextInterface {
    user: User | null,
    setUser: Dispatch<SetStateAction<User | null>>
}


const UserContext = createContext(({
    user: null,
    setUser: () => { },
}) as UserContextInterface)

const UserProvider = ({ children }: { children: ReactNode }) => {
    const [loading, setLoading] = useState(true)
    const [user, setUser] = useState<User | null>(null)

    const api = useApi('', '')

    useEffect(() => {
        setLoading(true)
        api.user.getCurrentUser().then((user) => {
            setUser(user)

        }).catch((error) => {
            console.error(error)
        }).finally(() => {
            setLoading(false)
        })
    }, [])

    if (!loading && !user) {
        redirect('/login');

        return;
    }

    return (
        <>
            {!loading && user && <UserContext.Provider value={{ user, setUser }}>
                {children}
            </UserContext.Provider>}
            {loading && <CircularProgress />}
        </>
    );
};

export { UserContext, UserProvider };
