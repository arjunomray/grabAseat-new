import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

export const Profile = ({ setUserId }) => {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [userEvents, setUserEvents] = useState([]);
  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem('token');
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
    };
    fetchUser();
  }, []);

  useEffect(() => {
    const fetchUserEvents = async () => {
      const token = localStorage.getItem('token');
      const response = await fetch(
        `http://localhost:8080/users/${localStorage.getItem('userId')}/events`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );
      if (!response.ok) {
        setUserEvents([]);
        return;
      }
      const data = await response.json();
      console.log(data);
      setUserEvents(data);
    };
    console.log(userEvents);
    fetchUserEvents();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    setUserId(null);
    navigate('/login');
  };

  if (!user) {
    return <div>Loading or not authorized.</div>;
  }

  return (
    <div className="p-10">
      <h1 className="text-4xl font-bold">Profile</h1>
      <div className="px-10 py-5">
        <h2 className="text-2xl font-bold">Name</h2>
        <p className="text-lg">{user.name}</p>
      </div>
      <div className="px-10 py-5">
        <h2 className="text-2xl font-bold">Email</h2>
        <p className="text-lg">{user.email}</p>
      </div>
      <div className="px-10 py-5">
        <h2 className="text-2xl font-bold">My Events</h2>
        <div className="flex flex-col items-center justify-center">
          {userEvents.length > 0 ? (
            userEvents.map((event) => (
              <div key={event.id} className="w-full h-10 bg-gray-200 rounded-md">
                <div className="flex flex-row items-center justify-between">
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.name}</p>
                  </div>
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.time}</p>
                  </div>
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.location}</p>

                  </div>
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.price}</p>
                  </div>
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.seatsAvailable}</p>
                  </div>
                  <div className="flex flex-row items-center justify-center">
                    <p className="text-lg">{event.tags}</p>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="w-full h-10 bg-gray-200 rounded-md">
              No events found
            </div>
          )}
        </div>
      </div>
      <button onClick={handleLogout} className="bg-red-600 text-white p-2 rounded">Logout</button>
    </div>
  );
};

export default Profile;