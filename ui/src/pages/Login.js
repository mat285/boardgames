import React, { useContext, useState } from 'react';
import 'scss/Form.scss';

import { ErrorContext } from 'context/Error';
import { UserContext } from 'context/User';
import Button from 'components/Button';
import Spinner from 'components/Spinner';
import InputWithLabel from 'components/InputWithLabel';
import useAPI from 'api/useAPI';


const Login = props => {
    const { setError } = useContext(ErrorContext);
    const { setUser } = useContext(UserContext);
    const { login } = useAPI();
    const [loading, setLoading] = useState(false);
    const [username, setUsername] = useState('');
    const [errors, setErrors] = useState({});

    const onLogin = async () => {
        if (loading) return;
        setErrors({});
        setLoading(true);

        // validate
        const formErrors = {};
        if (!username) formErrors.username = 'Please enter a username';
        if (Object.keys(formErrors).length > 0) {
            setLoading(false);
            setErrors(formErrors);
            return;
        }

        try {
            const { username: uname, id: id } = await login(username);
            console.log('got' + JSON.stringify({ username: uname, id: id }))
            setUser({ username: uname, id: id });
        } catch (error) {
            setError(error);
            setLoading(false);
        }
    };

    const onKeyDown = event => {
        if (event.key !== 'Enter') return false;
        event.preventDefault();
        onLogin();
    };

    return (
        <div className='login'>
            <div className='formContainer'>
                <div className='formTitle'>Boardgames Server</div>
                <InputWithLabel
                    label='Username'
                    name='username'
                    type='username'
                    value={username}
                    error={errors.username}
                    onChange={event => setUsername(event.target.value)} />
                <Button
                    text='Login'
                    onClick={onLogin}
                    size='large' />
                <div className={'spinnerContainer' + (loading ? ' show' : '')}><Spinner /></div>
            </div>
        </div>
    );
};

export default Login;
