import React, { useCallback, useContext } from 'react';
import PropTypes from 'prop-types';
import 'scss/Contractor.scss';
import { useNavigate } from 'react-router-dom';

import { TokenContext } from 'context/Token';
import { ErrorContext } from 'context/Error';
import Image from 'components/Image';
import useAPI from 'hooks/useAPI';

const propTypes = {
  id: PropTypes.string.isRequired,
  displayName: PropTypes.string.isRequired,
  picture: PropTypes.shape({
    url: PropTypes.string,
    width: PropTypes.number,
    height: PropTypes.number
  }),
  location: PropTypes.string.isRequired,
  age: PropTypes.number.isRequired
};

const Contractor = props => {
  const navigate = useNavigate();
  const { token } = useContext(TokenContext);
  const { setError } = useContext(ErrorContext);
  const { createConversation } = useAPI();

  const startConversation = useCallback(async () => {
    try {
      const { id } = await createConversation({ user: token.userID, contractor: props.id });
      navigate(`/conversation/${id}`);
    } catch (error) {
      setError(error);
    }
  }, [navigate, setError, token, props.id]);

  return (
    <div className='contractor' onClick={startConversation}>
      <div className='imageContainer'><Image src={props.picture ? props.picture.url : null} alt={props.displayName} /></div>
      <div className='contractorInfo'>
        <div className='contractorName'>{props.displayName}</div>
        <div className='contractorAge'>{`${props.age} years old`}</div>
        <div className='contractorLocation'>{props.location || 'Unknown Location'}</div>
      </div>
    </div>
  );
};

Contractor.propTypes = propTypes;
export default Contractor;
