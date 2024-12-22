import { useState } from 'react';

const CreateEvent = () => {
    const [formData, setFormData] = useState({
        name: '',
        time: '',
    });
    const [successMessage, setSuccessMessage] = useState('');

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSuccessMessage(''); // Clear any previous success message
        try {
            const response = await fetch('http://localhost:8080/events', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            });
            if (response.ok) {
                setSuccessMessage('Event created successfully!');
                setFormData({ name: '', time: '' }); // Clear the form
            } else {
                console.error('Event creation failed');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div className="max-w-md mx-auto bg-white p-8 rounded shadow">
            <h1 className="text-2xl font-bold mb-4">Create Event</h1>
            {successMessage && <p className="text-green-600 mb-4">{successMessage}</p>}
            <form onSubmit={handleSubmit}>
                <input type="text" name="name" placeholder="Event Name" value={formData.name} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input type="datetime-local" name="time" placeholder="Event Time" value={formData.time} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <button type="submit" className="w-full bg-blue-600 text-white p-2 rounded">Create Event</button>
            </form>
        </div>
    );
};

export default CreateEvent;