import React, { createContext, useState } from 'react';

const ErrorContext = createContext({
  error: null,
  setError: error => { }
});

const ErrorProvider = props => {
  const [error, setError] = useState(null);

  return (
    <ErrorContext.Provider
      value={{ error, setError }}
      {...props} />
  );
};

export { ErrorContext, ErrorProvider };
