// frontend/src/components/Login.tsx
import React, { useState } from 'react';
import axios from 'axios';

const AUTH_SERVICE_URL = 'http://localhost:8081';

export const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setError('');
      setMessage('');
      await axios.post(`${AUTH_SERVICE_URL}/register`, { username, password });
      setMessage('Registration successful! Please log in.');
    } catch (err) {
      setError('Registration failed. Username may already be in use.');
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setError('');
      setMessage('');
      const response = await axios.post(`${AUTH_SERVICE_URL}/login`, { username, password });
      localStorage.setItem('authToken', response.data.token);
      setMessage('Login successful!');
      window.location.reload();
    } catch (err) {
      setError('Login failed. Please check your credentials.');
    }
  };

  return (
    <div className="login-container">
      <h2>NexusBoard Login</h2>
      <form>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <div className="button-group">
            <button type="submit" onClick={handleLogin}>Login</button>
            <button type="button" onClick={handleRegister}>Register</button>
        </div>
      </form>
      {error && <p className="error-message">{error}</p>}
      {message && <p className="success-message">{message}</p>}
    </div>
  );
};