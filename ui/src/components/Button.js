import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import 'scss/Button.scss';

const propTypes = {
  text: PropTypes.string.isRequired,
  url: PropTypes.string,
  size: PropTypes.oneOf(['small', 'large']),
  disabled: PropTypes.bool,
  onClick: PropTypes.func
};

const Button = ({ text, url, size = 'small', disabled, onClick }) => {
  const classes = 'button' + (size === 'small' ? ' small' : '') + (disabled ? ' disabled' : '');

  if (url) {
    return (
      <Link to={url} className={classes}>
        <div className='buttonText'>{text}</div>
      </Link>
    );
  }

  return (
    <div onClick={onClick} className={classes}>
      <div className='buttonText'>{text}</div>
    </div>
  );
};

Button.propTypes = propTypes;
export default Button;
