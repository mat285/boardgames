import React, { useCallback, useContext, useEffect, useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import 'scss/Dashboard.scss';

import { TokenContext } from 'context/Token';
import { ErrorContext } from 'context/Error';
import ConversationPreview from 'components/ConversationPreview';
import Error from 'components/Error';
import Image from 'components/Image';
import Spinner from 'components/Spinner';

import useAPI from 'hooks/useAPI';
import useSocket from 'hooks/useSocket';

import logo from 'images/pickl.svg';
import dashboardIcon from 'images/icons/dashboard.svg';
import profileIcon from 'images/icons/profile.svg';
import logoutIcon from 'images/icons/logout.svg';
import menuIcon from 'images/icons/menu.svg';

const Dashboard = () => {
  const { token, setToken } = useContext(TokenContext);
  const { error, setError } = useContext(ErrorContext);
  const { getConversations, logout } = useAPI();
  const { socket, PACKET_TYPES } = useSocket();
  // grab conversation id to match active conversation
  const location = useLocation();
  const locationParts = location.pathname.split('/');
  const conversationID = locationParts[locationParts.length - 1];

  const [loading, setLoading] = useState(true);
  const [conversations, setConversations] = useState([]);
  const [search, setSearch] = useState('');
  const [isConversationMenuOpen, setIsConversationMenuOpen] = useState(false);

  const loadConversations = useCallback(async () => {
    try {
      const response = await getConversations(token.userID);
      setConversations(response);
      setLoading(false);
    } catch (error) {
      setError(error);
    }
  }, []);

  useEffect(() => {
    loadConversations();
  }, [loadConversations]);

  // when a conversation is opened, update the preview to mark as read
  useEffect(() => {
    if (!conversationID || conversations.length === 0) return;
    conversations.forEach((conversation, index) => {
      if (conversation.id !== conversationID || !conversation.preview.latestMessage || conversation.preview.latestMessage.read) return;
      setConversations([
        ...conversations.slice(0, index),
        { ...conversation, preview: { ...conversation.preview, latestMessage: { ...conversation.preview.latestMessage, read: true } } },
        ...conversations.slice(index + 1)
      ]);
    });
  }, [conversationID]);

  // when a new message is received, update the conversation list to put the conversation w/ the new message at the top and display the new preview
  const updateConversationPreviews = useCallback(message => {
    const conversationIndex = conversations.findIndex(conversation => conversation.id === message.conversation);
    if (conversationIndex !== -1) setConversations([
      { ...conversations[conversationIndex], preview: { ...conversations[conversationIndex].preview, latestMessage: message } },
      ...conversations.slice(0, conversationIndex),
      ...conversations.slice(conversationIndex + 1)
    ]);
  }, [conversations]);

  const onMessage = useCallback(message => {
    const data = JSON.parse(message.data);
    socket.send(JSON.stringify({ type: PACKET_TYPES.ACKNOWLEDGE_DELIVERED, version: 'v1alpha1', body: { id: data.body.id, timestamp: new Date().toUTCString() } })); // received
    if (data.type === PACKET_TYPES.MESSAGE_RECEIVED) {
      updateConversationPreviews(data.body);
    } else if (data.type === PACKET_TYPES.CONVERSATION_CREATED) {
      // when a new conversation is created, add it to the list
      setConversations([data.body, ...conversations]);
    }
  }, [socket, conversations, updateConversationPreviews]);

  useEffect(() => {
    if (!socket) return;
    socket.addEventListener('message', onMessage);
    return () => {
      socket.removeEventListener('message', onMessage);
    };
  }, [socket, onMessage]);

  const onLogout = async () => {
    await logout();
    setToken(null);
  };

  return (
    <div className='dashboard'>
      <Error error={error} />
      <div className='chatApp'>
        <div className='actions'>
          <div className='logo'><Image src={logo} alt='Pickl Logo' /></div>
          <div className='menu' onClick={() => setIsConversationMenuOpen(!isConversationMenuOpen)}><Image src={menuIcon} alt='Open menu' /></div>
          <div className='actionButtons'>
            <Link to='/' className='action'><Image src={dashboardIcon} alt='dashboard' /></Link>
            <Link to='/profile' className='action'><Image src={profileIcon} alt='profile' /></Link>
            <div className='action' onClick={onLogout}><Image src={logoutIcon} alt='logout' /></div>
          </div>
          <div className='actionFooter'></div>
        </div>
        <div className={'conversationList' + (isConversationMenuOpen ? ' open' : '')}>
          <div className='listHeader'>
            <div className='searchContainer'><input type='text' placeholder='Search' value={search} onChange={event => setSearch(event.target.value)} /></div>
            {!token.isContractor && <Link to='/new' className='listHeaderAction'><svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 448 512'><path d='M256 80c0-17.7-14.3-32-32-32s-32 14.3-32 32V224H48c-17.7 0-32 14.3-32 32s14.3 32 32 32H192V432c0 17.7 14.3 32 32 32s32-14.3 32-32V288H400c17.7 0 32-14.3 32-32s-14.3-32-32-32H256V80z' /></svg></Link>}
          </div>
          {loading && <div className='spinnerContainer'><Spinner /></div>}
          {!loading && conversations.filter(conversation => conversation.preview.receiver.displayName.includes(search)).map(conversation => (
            <ConversationPreview
              key={conversation.id}
              isUnread={conversation.preview.latestMessage && conversation.preview.latestMessage.sender !== token.userID && !conversation.preview.latestMessage.read}
              isActive={conversationID === conversation.id}
              url={`/conversation/${conversation.id}`}
              partner={conversation.preview.receiver}
              lastMessage={conversation.preview.latestMessage ? conversation.preview.latestMessage.text : null} />
          ))}
          {(!loading && conversations.length === 0 && !token.isContractor) && <Link to='/new' className='noConversations'><svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 448 512'><path d='M256 80c0-17.7-14.3-32-32-32s-32 14.3-32 32V224H48c-17.7 0-32 14.3-32 32s14.3 32 32 32H192V432c0 17.7 14.3 32 32 32s32-14.3 32-32V288H400c17.7 0 32-14.3 32-32s-14.3-32-32-32H256V80z' /></svg></Link>}
          {(!loading && conversations.length === 0 && token.isContractor) && <div className='noMessagesContractor'>No one has sent you their pickl yet</div>}
        </div>
        <div className='contentContainer'>
          <Outlet context={{ onMessageSent: updateConversationPreviews }} />
        </div>
      </div>
    </div>
  );
};


export default Dashboard;
