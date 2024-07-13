import React, { useContext, useEffect, useState } from 'react';
import 'scss/App.scss';
import {
  createBrowserRouter,
  Navigate,
  RouterProvider
} from 'react-router-dom';


import Spinner from 'components/Spinner';
// import Dashboard from 'templates/Dashboard';
import Website from 'templates/Website';
import Login from 'pages/Login';
import Games from 'pages/Games';
import NotFound from 'pages/NotFound';
import useAPI from 'api/useAPI';

const App = () => {

  const router = createBrowserRouter([{
    path: '/',
    element: <Website />,
    errorElement: <NotFound />,
    children: [{
      path: '/',
      element: <Login />
    }]
  }]);

  return <RouterProvider router={router} />;
};
export default App;
