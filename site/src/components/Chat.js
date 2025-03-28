import React, { useState, useEffect, useRef } from 'react';
import Sidebar from './Sidebar';
import Message from './Message';
import './Chat.css';

function Chat({ user, onLogout }) {
  const [conversations, setConversations] = useState([]);
  const [currentConversation, setCurrentConversation] = useState(null);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  
  const messageEndRef = useRef(null);
  const token = localStorage.getItem('token');

  // 获取所有会话
  const fetchConversations = async () => {
    try {
      const response = await fetch('/api/conversations', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      const data = await response.json();
      if (data.code === 200) {
        setConversations(data.data);
      }
    } catch (err) {
      console.error('获取会话列表失败', err);
    }
  };

  // 获取会话历史
  const fetchMessages = async (conversationId) => {
    try {
      const response = await fetch(`/api/conversations/${conversationId}/messages`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      const data = await response.json();
      if (data.code === 200) {
        setMessages(data.data);
      }
    } catch (err) {
      console.error('获取消息历史失败', err);
    }
  };

  // 创建新会话
  const createNewChat = () => {
    setCurrentConversation(null);
    setMessages([]);
    setMobileSidebarOpen(false);
  };

  // 选择会话
  const selectConversation = (conversation) => {
    setCurrentConversation(conversation);
    fetchMessages(conversation.id);
    setMobileSidebarOpen(false);
  };

  // 删除会话
  const deleteConversation = async (conversationId) => {
    try {
      const response = await fetch(`/api/conversations/${conversationId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        await fetchConversations();
        if (currentConversation && currentConversation.id === conversationId) {
          createNewChat();
        }
      }
    } catch (err) {
      console.error('删除会话失败', err);
    }
  };

  // 发送消息
  const sendMessage = async (e) => {
    e.preventDefault();
    
    if (!inputMessage.trim()) return;
    
    const userMessage = {
      role: 'user',
      content: inputMessage,
    };
    
    // 添加用户消息到界面
    setMessages(prevMessages => [...prevMessages, userMessage]);
    
    // 清空输入框
    setInputMessage('');
    
    // 显示加载状态
    setLoading(true);
    
    try {
      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          conversation_id: currentConversation ? currentConversation.id : '',
          message: userMessage.content
        })
      });
      
      const data = await response.json();
      
      if (data.code === 200) {
        // 如果是新会话，更新当前会话ID和获取会话列表
        if (!currentConversation) {
          setCurrentConversation({ id: data.data.conversation_id });
          fetchConversations();
        }
        
        // 添加AI回复到消息列表
        const assistantMessage = {
          role: 'assistant',
          content: data.data.message
        };
        
        setMessages(prevMessages => [...prevMessages, assistantMessage]);
      } else {
        // 显示错误消息
        const errorMessage = {
          role: 'assistant',
          content: '抱歉，发生了错误: ' + data.message
        };
        setMessages(prevMessages => [...prevMessages, errorMessage]);
      }
    } catch (err) {
      console.error('发送消息失败', err);
      const errorMessage = {
        role: 'assistant',
        content: '抱歉，发送消息失败，请稍后再试。'
      };
      setMessages(prevMessages => [...prevMessages, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    fetchConversations();
  }, []);

  // 消息滚动到底部
  useEffect(() => {
    messageEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <div className="chat-container">
      <div className={`sidebar ${mobileSidebarOpen ? 'sidebar-mobile-open' : ''}`}>
        <Sidebar
          conversations={conversations}
          onSelectConversation={selectConversation}
          onNewChat={createNewChat}
          onDeleteConversation={deleteConversation}
          currentConversationId={currentConversation?.id}
          onLogout={onLogout}
        />
      </div>
      
      <div className="chat-main">
        <button 
          className="mobile-sidebar-toggle"
          onClick={() => setMobileSidebarOpen(!mobileSidebarOpen)}
        >
          ☰
        </button>
        
        <div className="chat-messages">
          {messages.length === 0 ? (
            <div className="empty-chat">
              <h1>Chat Llama</h1>
              <p>与AI助手对话</p>
            </div>
          ) : (
            messages.map((message, index) => (
              <Message 
                key={index} 
                role={message.role} 
                content={message.content} 
              />
            ))
          )}
          {loading && (
            <div className="message assistant">
              <div className="message-content">
                <div className="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          )}
          <div ref={messageEndRef} />
        </div>
        
        <form className="chat-input-form" onSubmit={sendMessage}>
          <div className="chat-input-container">
            <input
              type="text"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              placeholder="发送消息..."
              disabled={loading}
            />
            <button type="submit" disabled={loading || !inputMessage.trim()}>
              发送
            </button>
          </div>
          <p className="disclaimer">Chat Llama 可能会产生不准确的信息。</p>
        </form>
      </div>
    </div>
  );
}

export default Chat; 