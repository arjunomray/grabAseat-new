import { BrowserRouter as Router, Routes, Route, Link, useLocation, Navigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Register from './components/Register';
import Login from './components/Login';
import CreateEvent from './components/CreateEvent';
import Events from './components/Events';
import UserTickets from './components/UserTickets';
import Profile from './components/Profile';
import LandingPage from './components/LandingPage';

const NavLink = ({ to, children }) => {
  const location = useLocation();
  const isActive = location.pathname === to;

  return (
    <Link
      to={to}
      className={`px-4 py-2 rounded-lg transition-all duration-200 ${isActive
        ? 'bg-white text-blue-600 shadow-sm'
        : 'text-white hover:bg-blue-700'
        }`}
    >
      {children}
    </Link>
  );
};

const App = () => {
  const [userId, setUserId] = useState(localStorage.getItem('userId'));
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleStorageChange = () => {
      setUserId(localStorage.getItem('userId'));
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    setUserId(null);
  };

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        {userId && (
          <nav className="bg-blue-600 shadow-lg">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex justify-between h-16">
                <div className="flex items-center">
                  <Link to="/" className="flex items-center space-x-3">
                    <svg className="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                        d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                    </svg>
                    <span className="text-white text-xl font-bold tracking-wide">GrabMySeat</span>
                  </Link>
                </div>

                <div className="hidden md:flex items-center space-x-2">
                  {!userId ? (
                    <>
                      <NavLink to="/register">Register</NavLink>
                      <NavLink to="/login">Login</NavLink>
                    </>
                  ) : (
                    <>

                      <NavLink to="/events/create">Create Events</NavLink>
                      <NavLink to="/events">Events</NavLink>
                      <NavLink to={`/user-tickets/${userId}`}>Your Tickets</NavLink>
                      <NavLink to="/profile">
                        <div className="flex items-center space-x-2">
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                          </svg>
                          <span>Profile</span>
                        </div>
                      </NavLink>
                      <button
                        onClick={handleLogout}
                        className="px-4 py-2 text-white hover:bg-blue-700 rounded-lg transition-all duration-200"
                      >
                        Logout
                      </button>
                    </>
                  )}
                </div>

                <div className="md:hidden flex items-center">
                  <button
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                    className="text-white hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-white"
                  >
                    <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      {isMobileMenuOpen ? (
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                      ) : (
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
                      )}
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            {isMobileMenuOpen && (
              <div className="md:hidden bg-blue-700">
                <div className="px-2 pt-2 pb-3 space-y-1">
                  {!userId ? (
                    <>
                      <Link to="/register" className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Register
                      </Link>
                      <Link to="/login" className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Login
                      </Link>
                    </>
                  ) : (
                    <>

                      <Link to="/events/create" className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Create Event
                      </Link>
                      <Link to="/events" className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Events
                      </Link>
                      <Link to={`/user-tickets/${userId}`} className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Your Tickets
                      </Link>
                      <Link to="/profile" className="block px-3 py-2 rounded-md text-white hover:bg-blue-800">
                        Profile
                      </Link>
                      <button
                        onClick={handleLogout}
                        className="block w-full text-left px-3 py-2 rounded-md text-white hover:bg-blue-800"
                      >
                        Logout
                      </button>
                    </>
                  )}
                </div>
              </div>
            )}
          </nav>
        )}

        <main className={userId ? "max-w-7xl mx-auto py-6 sm:px-6 lg:px-8" : ""}>
          <div className={userId ? "px-4 py-6 sm:px-0" : ""}>
            <Routes>
              <Route path="/" element={userId ? <Events /> : <LandingPage />} />
              <Route path="/register" element={<Register />} />
              <Route path="/login" element={<Login setUserId={setUserId} />} />

              <Route
                path="/events/create"
                element={userId ? <CreateEvent /> : <Navigate to="/" />}
              />
              <Route
                path="/events"
                element={userId ? <Events /> : <Navigate to="/" />}
              />
              <Route
                path="/user-tickets/:id"
                element={userId ? <UserTickets /> : <Navigate to="/" />}
              />
              <Route
                path="/profile"
                element={userId ? <Profile setUserId={setUserId} /> : <Navigate to="/" />}
              />
              <Route
                path="*"
                element={
                  <div className="text-center py-12">
                    <h1 className="text-4xl font-bold text-gray-900 mb-4">
                      404 - Page Not Found
                    </h1>
                    <p className="text-gray-600 mb-8">
                      The page you're looking for doesn't exist.
                    </p>
                    <Link
                      to="/"
                      className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                    >
                      Go Home
                    </Link>
                  </div>
                }
              />
            </Routes>
          </div>
        </main>

        {userId && (
          <footer className="bg-white border-t border-gray-200">
            <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
              <div className="flex justify-between items-center">
                <div className="text-gray-500 text-sm">
                  © 2024 GrabMySeat. All rights reserved.
                </div>
                <div className="flex space-x-6">
                  <a href="#" className="text-gray-400 hover:text-gray-500">
                    Terms
                  </a>
                  <a href="#" className="text-gray-400 hover:text-gray-500">
                    Privacy
                  </a>
                  <a href="#" className="text-gray-400 hover:text-gray-500">
                    Contact
                  </a>
                </div>
              </div>
            </div>
          </footer>
        )}
      </div>
    </Router>
  );
};

export default App;