import React from 'react';
import { Link } from 'react-router-dom';
import 'scss/Header.scss';

import Button from 'components/Button';
import Image from 'components/Image';


const Header = props => {
  return (
    <div className='header'>
      <div className='headerLeft'>
      </div>
      <div className='headerRight'>
        <Link to='/login'>Log In</Link>
        <Button text='Sign Up' url='/register' size='small' />
      </div>
    </div>
  );
};

export default Header;
