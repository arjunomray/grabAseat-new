import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { QRCodeSVG } from 'qrcode.react';

const UserTickets = () => {
    const { id } = useParams();
    const [tickets, setTickets] = useState([]);
    const [events, setEvents] = useState({});
    const [isLoading, setIsLoading] = useState(true);

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
                    eventDetails[ticket.event_id] = eventData;
                }
                setEvents(eventDetails);
            } catch (error) {
                console.error('Error fetching tickets or events:', error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchTickets();
    }, [id]);

    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto">
                <div className="text-center mb-8">
                    <h1 className="text-3xl font-bold text-gray-900">Your Tickets</h1>
                    <p className="mt-2 text-sm text-gray-600">
                        Manage and view your event tickets
                    </p>
                </div>

                {tickets.length === 0 ? (
                    <div className="bg-white rounded-xl shadow-md p-8 text-center">
                        <div className="mx-auto h-12 w-12 text-gray-400">
                            <svg className="h-full w-full" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                            </svg>
                        </div>
                        <h3 className="mt-2 text-sm font-medium text-gray-900">No tickets found</h3>
                        <p className="mt-1 text-sm text-gray-500">Get started by registering for an event.</p>
                    </div>
                ) : (
                    <div className="space-y-6">
                        {tickets.map((ticket) => {
                            const event = events[ticket.event_id];
                            return (
                                <div
                                    key={ticket.id}
                                    className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300"
                                >
                                    <div className="p-6">
                                        <div className="flex justify-between items-start">
                                            <div className="flex-grow">
                                                <div className="flex items-center space-x-3">
                                                    <h2 className="text-xl font-semibold text-gray-900">
                                                        {event?.name}
                                                    </h2>
                                                    <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                                        Active
                                                    </span>
                                                </div>

                                                <div className="mt-4 grid grid-cols-2 gap-4">
                                                    <div>
                                                        <p className="text-sm text-gray-500">Date & Time</p>
                                                        <p className="mt-1 text-sm font-medium text-gray-900">
                                                            {new Date(event?.time).toLocaleDateString('en-US', {
                                                                weekday: 'long',
                                                                year: 'numeric',
                                                                month: 'long',
                                                                day: 'numeric',
                                                                hour: '2-digit',
                                                                minute: '2-digit'
                                                            })}
                                                        </p>
                                                    </div>
                                                    <div>
                                                        <p className="text-sm text-gray-500">Ticket ID</p>
                                                        <p className="mt-1 text-sm font-medium text-gray-900">
                                                            #{ticket.id}
                                                        </p>
                                                    </div>
                                                    {event?.location && (
                                                        <div>
                                                            <p className="text-sm text-gray-500">Location</p>
                                                            <p className="mt-1 text-sm font-medium text-gray-900">
                                                                {event.location}
                                                            </p>
                                                        </div>
                                                    )}
                                                </div>
                                            </div>

                                            <div className="ml-6 flex-shrink-0">
                                                <div className="p-2 bg-white rounded-lg shadow-sm border border-gray-200">
                                                    <QRCodeSVG
                                                        value={ticket.id.toString()}
                                                        size={120}
                                                        level="H"
                                                        includeMargin={true}
                                                    />
                                                </div>
                                                <p className="mt-2 text-xs text-center text-gray-500">
                                                    Scan for entry
                                                </p>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="bg-gray-50 px-6 py-3 flex items-center justify-between">
                                        <div className="flex items-center space-x-2 text-sm text-gray-500">
                                            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                            </svg>
                                            <span>Show this ticket at the event entrance</span>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
};

export default UserTickets;