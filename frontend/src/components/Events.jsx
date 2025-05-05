import { useState, useEffect, useMemo } from 'react';

const EventModal = ({ event, isOpen, onClose, isRegistered, onRegister }) => {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                {/* Modal Header */}
                <div className="p-6 border-b border-gray-200">
                    <div className="flex justify-between items-start">
                        <h2 className="text-2xl font-bold text-gray-900">{event.name}</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-gray-500 transition-colors"
                        >
                            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                {/* Modal Content */}
                <div className="p-6 space-y-6">
                    {/* Date and Time */}
                    <div className="flex items-center text-gray-600">
                        <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>
                            {new Date(event.time).toLocaleDateString('en-US', {
                                weekday: 'long',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                            })}
                        </span>
                    </div>

                    {/* Location */}
                    <div className="flex items-center text-gray-600">
                        <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        <span>{event.location}</span>
                    </div>

                    {/* Price */}
                    <div className="flex items-center text-gray-600">
                        <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>Rs.{event.price}</span>
                    </div>

                    {/* Available Seats */}
                    <div className="flex items-center text-gray-600">
                        <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                        </svg>
                        <span>{event.seatsAvailable} seats available</span>
                    </div>

                    {/* Description */}
                    {event.description && (
                        <div className="space-y-2">
                            <h3 className="text-lg font-medium text-gray-900">About this event</h3>
                            <p className="text-gray-600">{event.description}</p>
                        </div>
                    )}

                    {/* Tags */}
                    {event.tags && event.tags.length > 0 && (
                        <div className="space-y-2">
                            <h3 className="text-lg font-medium text-gray-900">Tags</h3>
                            <div className="flex flex-wrap gap-2">
                                {event.tags.map((tag, index) => (
                                    <span
                                        key={index}
                                        className="px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
                                    >
                                        {tag}
                                    </span>
                                ))}
                            </div>
                        </div>
                    )}
                </div>

                {/* Modal Footer */}
                <div className="p-6 border-t border-gray-200">
                    {isRegistered(event.id) ? (
                        <div className="flex items-center justify-center px-4 py-3 bg-green-100 text-green-700 rounded-lg">
                            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                            </svg>
                            <span className="font-medium">You're registered for this event</span>
                        </div>
                    ) : (
                        <button
                            onClick={() => {
                                onRegister(event.id);
                                onClose();
                            }}
                            className="w-full flex items-center justify-center px-4 py-3 border border-transparent 
                                     text-base font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 
                                     focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 
                                     transition-colors duration-200"
                        >
                            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                            </svg>
                            Register for Event
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

const CreateEventModal = ({ isOpen, onClose }) => {
    const [tags, setTags] = useState([]);
    const [formData, setFormData] = useState({
        name: '',
        time: '',
        description: '',
        location: '',
        seatsAvailable: '',
        price: ''
    });

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch('http://localhost:8080/events', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    ...formData,
                    price: parseFloat(formData.price),
                    seatsAvailable: parseInt(formData.seatsAvailable),
                    tags: tags,
                    createdBy: localStorage.getItem('userId')
                })
            });
            if (response.ok) {
                onClose(true); // Pass true to indicate successful creation
                setFormData({
                    name: '',
                    time: '',
                    description: '',
                    location: '',
                    seatsAvailable: '',
                    price: ''
                });
                setTags([]);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                <div className="p-6 border-b border-gray-200">
                    <div className="flex justify-between items-center">
                        <h2 className="text-2xl font-bold text-gray-900">Create New Event</h2>
                        <button
                            onClick={() => onClose(false)}
                            className="text-gray-400 hover:text-gray-500 transition-colors"
                        >
                            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">Event Name</label>
                        <input
                            required
                            type="text"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            placeholder="Enter event name"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">Date & Time</label>
                        <input
                            required
                            type="datetime-local"
                            name="time"
                            value={formData.time}
                            onChange={handleChange}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">Description</label>
                        <textarea
                            required
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            rows="3"
                            placeholder="Enter event description"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">Location</label>
                        <input
                            required
                            type="text"
                            name="location"
                            value={formData.location}
                            onChange={handleChange}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            placeholder="Enter event location"
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="block text-sm font-medium text-gray-700">Available Seats</label>
                            <input
                                required
                                type="number"
                                name="seatsAvailable"
                                value={formData.seatsAvailable}
                                onChange={handleChange}
                                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="Number of seats"
                                min="1"
                            />
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium text-gray-700">Price (Rs.)</label>
                            <input
                                required
                                type="number"
                                name="price"
                                value={formData.price}
                                onChange={handleChange}
                                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="Event price"
                                min="0"
                                step="0.01"
                            />
                        </div>
                    </div>

                    <div className="space-y-2">
                        <label className="block text-sm font-medium text-gray-700">Tags</label>
                        <div className="flex flex-wrap gap-2">
                            {tags.map((tag, index) => (
                                <span
                                    key={index}
                                    className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm flex items-center"
                                >
                                    {tag}
                                    <button
                                        type="button"
                                        onClick={() => setTags(tags.filter((_, i) => i !== index))}
                                        className="ml-2 text-blue-600 hover:text-blue-800"
                                    >
                                        ×
                                    </button>
                                </span>
                            ))}
                            <input
                                type="text"
                                placeholder="Add tag and press Enter"
                                className="flex-grow p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        e.preventDefault();
                                        const value = e.target.value.trim();
                                        if (value && !tags.includes(value)) {
                                            setTags([...tags, value]);
                                            e.target.value = '';
                                        }
                                    }
                                }}
                            />
                        </div>
                    </div>

                    <div className="pt-4">
                        <button
                            type="submit"
                            className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 
                                     focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 
                                     transition-colors duration-200"
                        >
                            Create Event
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

const EventFilters = ({ onSearch, onFilterChange }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [filters, setFilters] = useState({
        priceRange: 'all', // all, free, paid
        dateRange: 'all', // all, today, week, month
        seatsAvailable: false,
    });

    const handleSearchChange = (e) => {
        const value = e.target.value;
        setSearchTerm(value);
        onSearch(value);
    };

    const handleFilterChange = (key, value) => {
        const newFilters = { ...filters, [key]: value };
        setFilters(newFilters);
        onFilterChange(newFilters);
    };

    return (
        <div className="bg-white rounded-xl shadow-md p-6 mb-6">
            <div className="space-y-4">
                {/* Search Bar */}
                <div className="relative">
                    <input
                        type="text"
                        placeholder="Search events by name, location, or tags..."
                        value={searchTerm}
                        onChange={handleSearchChange}
                        className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                    <div className="absolute left-3 top-2.5 text-gray-400">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </div>
                </div>

                {/* Filters */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Price Filter */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Price</label>
                        <select
                            value={filters.priceRange}
                            onChange={(e) => handleFilterChange('priceRange', e.target.value)}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="all">All Prices</option>
                            <option value="free">Free Events</option>
                            <option value="paid">Paid Events</option>
                        </select>
                    </div>

                    {/* Date Filter */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
                        <select
                            value={filters.dateRange}
                            onChange={(e) => handleFilterChange('dateRange', e.target.value)}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="all">All Dates</option>
                            <option value="today">Today</option>
                            <option value="week">This Week</option>
                            <option value="month">This Month</option>
                        </select>
                    </div>

                    {/* Available Seats Filter */}
                    <div className="flex items-center space-x-2">
                        <input
                            type="checkbox"
                            id="seatsAvailable"
                            checked={filters.seatsAvailable}
                            onChange={(e) => handleFilterChange('seatsAvailable', e.target.checked)}
                            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                        />
                        <label htmlFor="seatsAvailable" className="text-sm font-medium text-gray-700">
                            Show only events with available seats
                        </label>
                    </div>
                </div>
            </div>
        </div>
    );
};

const Events = () => {
    const [events, setEvents] = useState([]);
    const [userTickets, setUserTickets] = useState([]);
    const [selectedEvent, setSelectedEvent] = useState(null);
    const userId = localStorage.getItem('userId');
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');
    const [activeFilters, setActiveFilters] = useState({
        priceRange: 'all',
        dateRange: 'all',
        seatsAvailable: false,
    });

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
            setUserTickets([...userTickets, data]);
        } catch (error) {
            console.error('Error creating ticket:', error);
        }
    };

    const isRegistered = (eventId) => {
        return Array.isArray(userTickets) && userTickets.some(ticket => ticket.event_id === eventId);
    };

    // Filter events based on search and filters
    const filteredEvents = useMemo(() => {
        return events.filter(event => {
            // Search term filter
            const searchLower = searchTerm.toLowerCase();
            const matchesSearch = searchTerm === '' ||
                event.name.toLowerCase().includes(searchLower) ||
                event.location.toLowerCase().includes(searchLower) ||
                (event.tags && event.tags.some(tag => tag.toLowerCase().includes(searchLower)));

            // Price filter
            const matchesPrice = activeFilters.priceRange === 'all' ||
                (activeFilters.priceRange === 'free' && event.price === 0) ||
                (activeFilters.priceRange === 'paid' && event.price > 0);

            // Date filter
            const eventDate = new Date(event.time);
            const today = new Date();
            const matchesDate = activeFilters.dateRange === 'all' ||
                (activeFilters.dateRange === 'today' &&
                    eventDate.toDateString() === today.toDateString()) ||
                (activeFilters.dateRange === 'week' &&
                    eventDate <= new Date(today.getTime() + 7 * 24 * 60 * 60 * 1000)) ||
                (activeFilters.dateRange === 'month' &&
                    eventDate.getMonth() === today.getMonth() &&
                    eventDate.getFullYear() === today.getFullYear());

            // Available seats filter
            const matchesSeats = !activeFilters.seatsAvailable || event.seatsAvailable > 0;

            return matchesSearch && matchesPrice && matchesDate && matchesSeats;
        });
    }, [events, searchTerm, activeFilters]);

    return (
        <div className="page-container animate-fade-in bg-gray-50 min-h-screen">
            <div className="flex justify-between items-center mb-6">
                <div className="text-center">
                    <h1 className="text-4xl font-bold text-gray-800 mb-2">
                        Upcoming Events
                    </h1>
                    <div className="h-1 w-24 bg-blue-600 rounded-full"></div>
                </div>
                <button
                    onClick={() => setIsCreateModalOpen(true)}
                    className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 
                             transition-colors duration-200 flex items-center space-x-2"
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4" />
                    </svg>
                    <span>Create Event</span>
                </button>
            </div>

            {/* Add EventFilters component */}
            <EventFilters
                onSearch={setSearchTerm}
                onFilterChange={setActiveFilters}
            />

            {/* Show results count */}
            <div className="mb-6 text-gray-600">
                Found {filteredEvents.length} event{filteredEvents.length !== 1 ? 's' : ''}
            </div>

            {/* Display filtered events */}
            <div className="grid gap-6">
                {filteredEvents.length > 0 ? (
                    filteredEvents.map((event) => (
                        <div
                            key={event.id}
                            className="bg-white rounded-xl shadow-md overflow-hidden cursor-pointer transform transition-all duration-300 hover:shadow-xl hover:-translate-y-1"
                            onClick={() => setSelectedEvent(event)}
                        >
                            <div className="p-6">
                                <div className="flex justify-between items-start">
                                    <div className="space-y-3">
                                        <h2 className="text-2xl font-semibold text-gray-800">
                                            {event.name}
                                        </h2>
                                        <div className="flex items-center text-gray-600">
                                            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                            </svg>
                                            <span className="text-sm">
                                                {new Date(event.time).toLocaleDateString('en-US', {
                                                    weekday: 'long',
                                                    year: 'numeric',
                                                    month: 'long',
                                                    day: 'numeric',
                                                    hour: '2-digit',
                                                    minute: '2-digit'
                                                })}
                                            </span>
                                        </div>
                                    </div>

                                    <div className="ml-4">
                                        {isRegistered(event.id) ? (
                                            <div className="flex items-center px-4 py-2 bg-green-100 text-green-700 rounded-lg">
                                                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                                                </svg>
                                                <span className="font-medium">Registered</span>
                                            </div>
                                        ) : (
                                            <button
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    handleRegister(event.id);
                                                }}
                                                className="btn-primary flex items-center"
                                            >
                                                Register Now
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    <div className="text-center py-12 bg-white rounded-xl shadow-md">
                        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <h3 className="mt-2 text-sm font-medium text-gray-900">No events found</h3>
                        <p className="mt-1 text-sm text-gray-500">
                            Try adjusting your search or filters to find what you're looking for.
                        </p>
                    </div>
                )}
            </div>

            <EventModal
                event={selectedEvent}
                isOpen={!!selectedEvent}
                onClose={() => setSelectedEvent(null)}
                isRegistered={isRegistered}
                onRegister={handleRegister}
            />

            <CreateEventModal
                isOpen={isCreateModalOpen}
                onClose={(wasSuccessful) => {
                    setIsCreateModalOpen(false);
                    if (wasSuccessful) {
                        // Refresh events after successful creation
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
                        fetchEvents();
                    }
                }}
            />
        </div>
    );
};

export default Events;