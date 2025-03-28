import React from 'react';
import './Message.css';

function Message({ role, content }) {
  return (
    <div className={`message ${role}`}>
      <div className="message-avatar">
        {role === 'user' ? '👤' : '🤖'}
      </div>
      <div className="message-content">
        {content}
      </div>
    </div>
  );
}

export default Message; 