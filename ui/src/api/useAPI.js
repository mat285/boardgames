import { useCallback, useContext } from 'react';
import { UserContext } from 'context/User';
// import decode from 'jwt-decode';

const API_URL = '/api/v1alpha1';

const useAPI = () => {
    const { user, setUser } = useContext(UserContext);

    const login = useCallback(async (uname) => {
        console.log('trying to call login')
        console.log(JSON.stringify({ username: uname }))
        const response = await fetch(`${API_URL}/user/login`, { method: 'POST', body: JSON.stringify({ username: uname }), });
        if (response.status !== 200) {
            const error = await response.json();
            console.error(error);
            throw new Error(error);
        }
        const user = await response.json();
        console.log(JSON.stringify(user))
        return user
        return await response.json();
    });

    const post = useCallback(async (path, request) => {
        // const { username, userID } = userID ? await login() : null;
        const response = await fetch(`${API_URL}${path}`, { method: 'POST', body: JSON.stringify(request), credentials: 'include', headers: { 'X-Username': user.username, 'X-UserID': user.id } });
        if (response.status !== 200) {
            const error = await response.json();
            console.error(error);
            throw new Error(error);
        }
        const user = await response.json();
        return user
    }, [user, setUser]);

    const get = useCallback(async path => {
        const response = await fetch(`${API_URL}${path}`, { method: 'GET', credentials: 'include', headers: { 'X-Username': user.username, 'X-UserID': user.id } });
        if (response.status !== 200) {
            const error = await response.json();
            console.error(error);
            throw new Error(error);
        }
        return await response.json();
    }, [user, setUser]);

    return {
        login: async (username) => await login(username),
        // post
        newGame: async ({ name }) => await post(`/games/${name}/new`, { name }),
        joinGame: async ({ id }) => await post(`/game/${id}/join`, { id }),
        startGame: async ({ id }) => await post(`/game/${id}/start`, { id }),
        getGameState: async ({ id }) => await get(`/game/${id}/state`, { id }),
        sendGamePacket: async ({ id, packet }) => await post(`/game/${id}/packet`, packet),

        // login: async ({ username }) => await post('/user/login', { username }),
        listGames: async update => await get(`/user/games`),
    };
};

export default useAPI;
