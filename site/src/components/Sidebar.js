import React from 'react';
import './Sidebar.css';

function Sidebar({ 
  conversations, 
  onSelectConversation, 
  onNewChat, 
  onDeleteConversation,
  currentConversationId,
  onLogout
}) {
  return (
    <>
      <div className="sidebar-header">
        <button className="new-chat-button" onClick={onNewChat}>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
          </svg>
          新对话
        </button>
      </div>
      
      <div className="conversations-list">
        {conversations.map(conversation => (
          <div 
            key={conversation.id} 
            className={`conversation-item ${currentConversationId === conversation.id ? 'active' : ''}`}
            onClick={() => onSelectConversation(conversation)}
          >
            <div className="conversation-title">
              {conversation.title}
              <div className="conversation-date">
                {formatDate(conversation.created_at || new Date().toISOString())}
              </div>
            </div>
            <button 
              className="delete-button"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteConversation(conversation.id);
              }}
              aria-label="删除对话"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5Zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5Zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6Z"/>
                <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1ZM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118ZM2.5 3h11V2h-11v1Z"/>
              </svg>
            </button>
          </div>
        ))}
      </div>
      
      <div className="sidebar-footer">
        <button className="logout-button" onClick={onLogout}>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path fillRule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0v2z"/>
            <path fillRule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3z"/>
          </svg>
          退出登录
        </button>
      </div>
    </>
  );
}

// 格式化日期函数
function formatDate(dateString) {
  try {
    const date = new Date(dateString);
    const now = new Date();
    const yesterday = new Date(now);
    yesterday.setDate(yesterday.getDate() - 1);
    
    // 如果是今天
    if (date.toDateString() === now.toDateString()) {
      return '今天';
    }
    // 如果是昨天
    else if (date.toDateString() === yesterday.toDateString()) {
      return '昨天';
    }
    // 其他日期
    else {
      return `${date.getMonth() + 1}月${date.getDate()}日`;
    }
  } catch (e) {
    return '';
  }
}

export default Sidebar; 