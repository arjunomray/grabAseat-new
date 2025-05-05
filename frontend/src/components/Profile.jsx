import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const EventCard = ({ event }) => (
  <div className="bg-gray-50 rounded-xl p-4 hover:shadow-md transition-shadow duration-200">
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
      <div>
        <h3 className="text-sm font-medium text-gray-500">Event Name</h3>
        <p className="mt-1 text-sm text-gray-900">{event.name}</p>
      </div>
      <div>
        <h3 className="text-sm font-medium text-gray-500">Date & Time</h3>
        <p className="mt-1 text-sm text-gray-900">
          {new Date(event.time).toLocaleDateString('en-US', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          })}
        </p>
      </div>
      <div>
        <h3 className="text-sm font-medium text-gray-500">Location</h3>
        <p className="mt-1 text-sm text-gray-900">{event.location}</p>
      </div>
      <div>
        <h3 className="text-sm font-medium text-gray-500">Price</h3>
        <p className="mt-1 text-sm text-gray-900">${event.price}</p>
      </div>
      <div>
        <h3 className="text-sm font-medium text-gray-500">Available Seats</h3>
        <p className="mt-1 text-sm text-gray-900">{event.seatsAvailable}</p>
      </div>
      {event.tags && event.tags.length > 0 && (
        <div>
          <h3 className="text-sm font-medium text-gray-500">Tags</h3>
          <div className="mt-1 flex flex-wrap gap-1">
            {event.tags.map((tag, index) => (
              <span
                key={index}
                className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  </div>
);

const EmptyState = ({ message, actionText, onAction }) => (
  <div className="text-center py-12 bg-gray-50 rounded-xl">
    <svg
      className="mx-auto h-12 w-12 text-gray-400"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
      />
    </svg>
    <h3 className="mt-2 text-sm font-medium text-gray-900">{message}</h3>
    {actionText && (
      <div className="mt-6">
        <button
          onClick={onAction}
          className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          {actionText}
        </button>
      </div>
    )}
  </div>
);

export const Profile = ({ setUserId }) => {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [createdEvents, setCreatedEvents] = useState([]);
  const [registeredEvents, setRegisteredEvents] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem('token');
      try {
        const response = await fetch(
          `http://localhost:8080/users/${localStorage.getItem('userId')}`,
          {
            headers: {
              'Authorization': `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) {
          setUser(null);
          return;
        }
        const data = await response.json();
        setUser(data);
      } finally {
        setIsLoading(false);
      }
    };
    fetchUser();
  }, []);

  useEffect(() => {
    const fetchUserEvents = async () => {
      const token = localStorage.getItem('token');
      try {
        const response = await fetch(
          `http://localhost:8080/users/${localStorage.getItem('userId')}/events`,
          {
            headers: {
              'Authorization': `Bearer ${token}`,
            },
          }
        );
        if (!response.ok) {
          setCreatedEvents([]);
          setRegisteredEvents([]);
          return;
        }
        const data = await response.json();
        setCreatedEvents(data.created_events || []);
        setRegisteredEvents(data.registered_events || []);
      } catch (error) {
        console.error('Error fetching events:', error);
        setCreatedEvents([]);
        setRegisteredEvents([]);
      }
    };
    fetchUserEvents();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    setUserId(null);
    navigate('/login');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center p-8 bg-white rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-gray-800">Not Authorized</h2>
          <p className="mt-2 text-gray-600">Please log in to view your profile.</p>
          <button
            onClick={() => navigate('/login')}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Go to Login
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Profile Header */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">My Profile</h1>
              <p className="mt-1 text-sm text-gray-500">Manage your account and view your events</p>
            </div>
            <button
              onClick={handleLogout}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Logout
            </button>
          </div>

          <div className="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div className="bg-gray-50 p-6 rounded-xl">
              <h2 className="text-sm font-medium text-gray-500">Name</h2>
              <p className="mt-2 text-lg font-medium text-gray-900">{user.name}</p>
            </div>
            <div className="bg-gray-50 p-6 rounded-xl">
              <h2 className="text-sm font-medium text-gray-500">Email</h2>
              <p className="mt-2 text-lg font-medium text-gray-900">{user.email}</p>
            </div>
          </div>
        </div>

        {/* Created Events Section */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Events I've Created</h2>
            <button
              onClick={() => navigate('/events/create')}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Create New Event
            </button>
          </div>
          <div className="space-y-4">
            {createdEvents.length > 0 ? (
              createdEvents.map((event) => (
                <EventCard key={event.id} event={event} />
              ))
            ) : (
              <EmptyState
                message="You haven't created any events yet"
                actionText="Create Your First Event"
                onAction={() => navigate('/events/create')}
              />
            )}
          </div>
        </div>

        {/* Registered Events Section */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Events I'm Registered For</h2>
            <button
              onClick={() => navigate('/events')}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Browse Events
            </button>
          </div>
          <div className="space-y-4">
            {registeredEvents.length > 0 ? (
              registeredEvents.map((event) => (
                <EventCard key={event.id} event={event} />
              ))
            ) : (
              <EmptyState
                message="You haven't registered for any events yet"
                actionText="Browse Available Events"
                onAction={() => navigate('/events')}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;