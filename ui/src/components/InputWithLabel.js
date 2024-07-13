import React, { useState } from 'react';
import PropTypes from 'prop-types';
import 'scss/InputWithLabel.scss';

const propTypes = {
  label: PropTypes.string.isRequired,
  error: PropTypes.string
  // any other <input /> props accepted
};

const InputWithLabel = props => {
  const [isFocused, setIsFocused] = useState(false);

  return (
    <div className={'inputWithLabel' + (isFocused ? ' focused' : '') + (String(props.value).length > 0 ? ' filled' : '') + (props.error ? ' error' : '')}>
      <input {...props} onFocus={() => setIsFocused(true)} onBlur={() => setIsFocused(false)} />
      <div className='label'>{props.error || props.label}</div>
    </div>
  );
};

InputWithLabel.propTypes = propTypes;
export default InputWithLabel;
