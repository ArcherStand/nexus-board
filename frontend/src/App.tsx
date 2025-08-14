// frontend/src/App.tsx
import { Login } from './components/Login';
import './App.css';

function App() {
  const token = localStorage.getItem('authToken');

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    window.location.reload();
  };

  if (!token) {
    return <Login />;
  }

  return (
    <div className="app-container">
      <header>
        <h1>Welcome to NexusBoard</h1>
        <button onClick={handleLogout}>Logout</button>
      </header>
      <main>
        <p>Your collaborative board will be displayed here.</p>
      </main>
    </div>
  );
}

export default App;