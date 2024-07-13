import React, { useContext, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import 'scss/Message.scss';

import { TokenContext } from 'context/Token';
import Image from 'components/Image';

import useAPI from 'hooks/useAPI';

const propTypes = {
  message: PropTypes.shape({
    id: PropTypes.string,
    sender: PropTypes.string,
    text: PropTypes.string,
    image: PropTypes.string
  }),
  isSending: PropTypes.bool
};

const Message = ({ message, isSending }) => {
  const { token } = useContext(TokenContext);
  const [imageURL, setImageURL] = useState(null);
  const [showBigImage, setShowBigImage] = useState(false);

  const { getMessageImage } = useAPI();

  const loadImage = async () => {
    if (!message.image) return;
    try {
      const response = await getMessageImage(message.id);
      setImageURL(response);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    loadImage();
  }, []);

  return (
    <>
      <div className={'message' + (message.sender === token.userID ? ' mine' : ' theirs') + (isSending ? ' sending' : '') + (message.image ? ' hasImage' : '')}>
        {message.image && <div className='messageImage' onClick={() => setShowBigImage(true)}>
          <div className='imagePlaceholder' style={{ paddingBottom: '100%' }} />
          <Image className={!message.text ? 'noText' : ''} src={imageURL} alt='user image' />
        </div>}
        {message.text && <div className='messageText'>{message.text}</div>}
      </div>
      {showBigImage && <div className='bigImage' onClick={() => setShowBigImage(false)}><Image src={imageURL} alt='user image' /></div>}
    </>
  );
};

Message.propTypes = propTypes;
export default Message;
