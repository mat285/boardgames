import React, { useContext } from 'react';
import 'scss/Website.scss';
import { Outlet } from 'react-router-dom';

import { ErrorContext } from 'context/Error';

import Error from 'components/Error';
import Header from 'components/Header';
import Footer from 'components/Footer';

const LoggedOutRouter = () => {
  const { error } = useContext(ErrorContext);

  return (
    <div className='website'>
      <Header />
      <Error error={error} />
      <div className='page'>
        <Outlet />
      </div>
      <Footer />
    </div>
  );
};
export default LoggedOutRouter;
