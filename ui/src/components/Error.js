import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import 'scss/Error.scss';

import { ErrorContext } from 'context/Error';

const propTypes = {
  error: PropTypes.object
};

const Error = props => {
  const { setError } = useContext(ErrorContext);

  if (!props.error) return <div className='errorEyebrow' />;

  return (
    <div className='errorEyebrow hasError'>
      <div className='errorMessage'>{props.error.message}</div>
      <div className='closeError' onClick={() => setError(null)}>&times;</div>
    </div>
  );
};

Error.propTypes = propTypes;
export default Error;
