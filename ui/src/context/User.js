import React, { createContext, useState } from 'react';

const UserContext = createContext({
    user: null,
    setUser: (user) => { }
});

const UserProvider = props => {
    const [user, setUser] = useState(null);

    return (
        <UserContext.Provider
            value={{ user, setUser }}
            {...props} />
    );
};

export { UserContext, UserProvider };
