import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { QRCodeSVG } from 'qrcode.react';

const UserTickets = () => {
    const { id } = useParams();
    const [tickets, setTickets] = useState([]);
    const [events, setEvents] = useState({});

    useEffect(() => {
        const fetchTickets = async () => {
            const token = localStorage.getItem('token');
            if (!token) return;

            try {
                const response = await fetch(`http://localhost:8080/tickets/${id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                const data = await response.json();
                setTickets(data);

                // Fetch event details for each ticket
                const eventDetails = {};
                for (const ticket of data) {
                    const eventResponse = await fetch(`http://localhost:8080/events/${ticket.event_id}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                    const eventData = await eventResponse.json();
                    eventDetails[ticket.event_id] = eventData.name;
                }
                setEvents(eventDetails);
            } catch (error) {
                console.error('Error fetching tickets or events:', error);
            }
        };

        fetchTickets();
    }, [id]);

    return (
        <div className="max-w-4xl mx-auto bg-white p-8 rounded shadow">
            <h1 className="text-2xl font-bold mb-4">Your Tickets</h1>
            <ul className="space-y-4">
                {tickets.map((ticket) => (
                    <li key={ticket.id} className="p-4 border rounded flex justify-between items-center">
                        <div>
                            <h2 className="text-xl font-semibold">{events[ticket.event_id]}</h2>
                            <p>Ticket ID: {ticket.id}</p>
                        </div>
                        <QRCodeSVG value={ticket.id} />
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserTickets;