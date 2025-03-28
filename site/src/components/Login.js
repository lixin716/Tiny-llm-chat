import React, { useState } from 'react';
import './Login.css';

function Login({ onLogin }) {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const endpoint = isLogin ? '/api/auth/login' : '/api/auth/register';
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (data.code === 200) {
        if (isLogin) {
          // 登录成功
          onLogin(data.data.user, data.data.token);
        } else {
          // 注册成功，切换到登录
          setIsLogin(true);
          setUsername('');
          setPassword('');
          setError('注册成功，请登录');
        }
      } else {
        setError(data.message || '操作失败');
      }
    } catch (err) {
      setError('网络错误，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-form">
        <h1 className="login-title">Chat Llama</h1>
        <h2>{isLogin ? '登录' : '注册'}</h2>
        
        {error && <div className="login-error">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">用户名</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">密码</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          
          <button 
            type="submit" 
            className="login-button"
            disabled={loading}
          >
            {loading ? '处理中...' : isLogin ? '登录' : '注册'}
          </button>
        </form>
        
        <div className="login-switch">
          {isLogin ? (
            <p>还没有账号？ <button onClick={() => setIsLogin(false)}>注册</button></p>
          ) : (
            <p>已有账号？ <button onClick={() => setIsLogin(true)}>登录</button></p>
          )}
        </div>
      </div>
    </div>
  );
}

export default Login; 