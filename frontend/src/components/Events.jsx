import { useState, useEffect } from 'react';

const Events = () => {
    const [events, setEvents] = useState([]);
    const [userTickets, setUserTickets] = useState([]);
    const userId = localStorage.getItem('userId');

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                const response = await fetch('http://localhost:8080/events', {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setEvents(data);
            } catch (error) {
                console.error('Error fetching events:', error);
            }
        };

        const fetchUserTickets = async () => {
            try {
                const response = await fetch(`http://localhost:8080/tickets/${userId}`, {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setUserTickets(data);
            } catch (error) {
                console.error('Error fetching user tickets:', error);
            }
        };

        fetchEvents();
        fetchUserTickets();
    }, [userId]);

    const handleRegister = async (eventId) => {
        try {
            const response = await fetch('http://localhost:8080/tickets', {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    event_id: eventId,
                    user_id: userId
                })
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            console.log('Ticket created:', data);
            setUserTickets([...userTickets, data]);
        } catch (error) {
            console.error('Error creating ticket:', error);
        }
    };

    const isRegistered = (eventId) => {
        return Array.isArray(userTickets) && userTickets.some(ticket => ticket.event_id === eventId);
    };

    return (
        <div className="max-w-4xl mx-auto bg-white p-8 rounded shadow">
            <h1 className="text-2xl font-bold mb-4">Events</h1>
            <ul className="space-y-4">
                {events.map((event) => (
                    <li key={event.id} className="p-4 border rounded flex justify-between items-center">
                        <div>
                            <h2 className="text-xl font-semibold">{event.name}</h2>
                            <p>{event.time}</p>
                        </div>
                        {isRegistered(event.id) ? (
                            <span className="text-green-600 font-bold">Registered</span>
                        ) : (
                            <button onClick={() => handleRegister(event.id)} className="bg-blue-600 text-white p-2 rounded">Register</button>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Events;