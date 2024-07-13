import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import 'scss/ConversationPreview.scss';

import Avatar from 'components/Avatar';

const propTypes = {
  isActive: PropTypes.bool.isRequired,
  url: PropTypes.string.isRequired,
  partner: PropTypes.shape({
    displayName: PropTypes.string.isRequired,
    picture: PropTypes.shape({
      url: PropTypes.string
    })
  }).isRequired,
  lastMessage: PropTypes.string
};

const ConversationPreview = props => {
  return (
    <Link to={props.url} className={'conversationPreview' + (props.isActive ? ' active' : '')}>
      <Avatar displayName={props.partner.displayName} picture={props.partner.picture} />
      <div className={'messageContainer' + (props.isUnread && !props.isActive ? ' unread' : '')}>
        <div className='partnerName'>{props.partner.displayName}</div>
        <div className='messagePreview'>{props.lastMessage || 'no messages yet'}</div>
      </div>
    </Link>
  );
};

ConversationPreview.propTypes = propTypes;
export default ConversationPreview;
