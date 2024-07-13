import React from 'react';
import 'scss/Footer.scss';

const Footer = () => {
  return (
    <div className='footer'>
      <div className='copyright'>&copy;{new Date().getFullYear()} pickl.io</div>
    </div>
  );
};

export default Footer;
