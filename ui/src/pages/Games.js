import React, { useCallback, useContext, useEffect, useState, useRef } from 'react';
import { useOutletContext, useParams } from 'react-router-dom';
import 'scss/Conversation.scss';

import { UserContext } from 'context/User';
import { ErrorContext } from 'context/Error';
import Avatar from 'components/Avatar';
import Spinner from 'components/Spinner';
import useAPI from 'api/useAPI';


const Games = () => {
    const { setUser } = useContext(UserContext);
    const [loading, setLoading] = useState(false);
    const { setError } = useContext(ErrorContext);
    const { listGames } = useAPI();

    const loadGames = useCallback(async () => {
        try {
            setLoading(true);
            const response = await listGames({});
            console.log(JSON.stringify(response))
            setLoading(false);
        } catch (error) {
            setError(error);
            setLoading(false);
        }
    }, []);

    const onKeyDown = useCallback(event => {
        if (event.key !== 'Enter') return false;
        event.preventDefault();
    }, []);

    return (
        <div className='games'>
            {loading && <div className='spinnerContainer'><Spinner /></div>}
            {!loading && <div className='gamesContentContainer'>
                <div className='gameInformation'>
                    {/* <Avatar gameName={game.name} /> */}
                    <div className='userDetails'>
                    </div>
                </div>
            </div>}
        </div>
    );
};

export default Games;
