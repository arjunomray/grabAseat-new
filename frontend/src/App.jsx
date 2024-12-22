import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Register from './components/Register';
import Login from './components/Login';
import CreateEvent from './components/CreateEvent';
import Events from './components/Events';
import UserTickets from './components/UserTickets';
import { Logout } from './components/Logout';

const App = () => {
  const [userId, setUserId] = useState(localStorage.getItem('userId'));

  useEffect(() => {
    const handleStorageChange = () => {
      setUserId(localStorage.getItem('userId'));
    };

    window.addEventListener('storage', handleStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
    };
  }, []);

  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <nav className="bg-blue-600 p-4 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <span className="text-white text-xl font-bold">GrabMySeat</span>
          </div>
          <ul className="flex space-x-4 text-white">
            {!userId && <li><Link to="/register" className="hover:underline">Register</Link></li>}
            {!userId && <li><Link to="/login" className="hover:underline">Login</Link></li>}
            {userId && <li><Link to="/create-event" className="hover:underline">Create Event</Link></li>}
            {userId && <li><Link to="/events" className="hover:underline">Events</Link></li>}
            {userId && <li><Link to={`/user-tickets/${userId}`} className="hover:underline">Your Tickets</Link></li>}
            {userId && <li><Link to="/logout" className="hover:underline">Logout</Link></li>}
          </ul>
        </nav>
        <div className="container mx-auto p-4">
          <Routes>
            <Route path="/register" element={<Register />} />
            <Route path="/login" element={<Login setUserId={setUserId} />} />
            <Route path="/create-event" element={<CreateEvent />} />
            <Route path="/events" element={<Events />} />
            <Route path="/user-tickets/:id" element={<UserTickets />} />
            <Route path="/logout" element={<Logout setUserId={setUserId} />} />
            <Route path="*" element={<h1 className="text-center text-2xl">Not Found</h1>} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

export default App;