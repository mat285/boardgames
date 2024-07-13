import React from 'react';
import { Link } from 'react-router-dom';
import 'scss/NotFound.scss';

const NotFound = () => {
  return (
    <div className='notFound'>
      <h1>Not Found</h1>
      <p>Something went wrong! Try refreshing the page or <Link to='/'>go home</Link>.</p>
    </div>
  );
};

export default NotFound;
