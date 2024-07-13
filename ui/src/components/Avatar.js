import React from 'react';
import PropTypes from 'prop-types';
import 'scss/Avatar.scss';

const propTypes = {
  displayName: PropTypes.string.isRequired,
  picture: PropTypes.shape({
    id: PropTypes.string,
    url: PropTypes.string,
    width: PropTypes.number
  }),
  size: PropTypes.oneOf(['small', 'large'])
};

const Avatar = ({ displayName, picture, size = 'small' }) => {
  return (
    <div className={`avatar${size === 'large' ? ' large' : ''}`} style={picture ? { backgroundImage: `url('${picture.url}')` } : {}}>
      {!picture && <span>{displayName ? displayName.charAt(0) : ''}</span>}
    </div>
  );
};

Avatar.propTypes = propTypes;
export default Avatar;
