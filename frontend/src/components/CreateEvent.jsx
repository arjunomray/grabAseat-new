import { useState } from 'react';
import { useRef } from 'react';

const EventTags = ({ tags, setTags }) => {
    const [inputValue, setInputValue] = useState('');
    const tagInputRef = useRef(null);

    const handleAddTag = () => {
        if (inputValue.trim() && !tags.includes(inputValue.trim())) {
            setTags([...tags, inputValue.trim()]);
            setInputValue('');
        }
    };

    const handleRemoveTag = (tagToRemove) => {
        setTags(tags.filter(tag => tag !== tagToRemove));
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleAddTag();
        }
    };

    return (
        <div className="flex flex-col gap-2 mb-4">
            <div className="flex flex-row items-center gap-2 justify-between w-full">
                <input
                    type="text"
                    name="tags"
                    placeholder="Event Tags"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyPress={handleKeyPress}
                    className="w-full p-2 min-h-8 border rounded"
                />
                <button
                    type="button"
                    onClick={handleAddTag}
                    className="bg-blue-600 text-white min-h-8 text-xs py-2 px-4 rounded hover:bg-blue-700"
                >
                    Add
                </button>
            </div>
            <div ref={tagInputRef} className="flex flex-wrap gap-2">
                {tags.map((tag) => (
                    <div
                        key={tag}
                        className="flex items-center gap-1 bg-gray-100 px-3 py-1 rounded-full text-sm"
                    >
                        <span>{tag}</span>
                        <button
                            onClick={() => handleRemoveTag(tag)}
                            className="text-gray-500 hover:text-gray-700"
                        >
                            ×
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
};

const CreateEvent = () => {

    const [tags, setTags] = useState([]);
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
        setSuccessMessage('');
        try {
            const response = await fetch('http://localhost:8080/events', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: formData.name,
                    time: formData.time,
                    location: formData.location,
                    price: parseFloat(formData.price),
                    seatsAvailable: parseInt(formData.seatsAvailable),
                    tags: tags
                })
            });
            if (response.ok) {
                setSuccessMessage('Event created successfully!');
                setFormData({ name: '', time: '' });
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
                <input required type="text" name="name" placeholder="Event Name" value={formData.name} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input required type="datetime-local" name="time" placeholder="Event Time" value={formData.time} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input required type="text" name="description" placeholder="Event Description" value={formData.description} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input required type="text" name="location" placeholder="Event Location" value={formData.location} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input required type="number" name="seatsAvailable" placeholder="Seats Available" value={formData.seatsAvailable} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <input required type="number" name="price" placeholder="Price" value={formData.price} onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
                <EventTags tags={tags} setTags={setTags} />
                <button type="submit" className="w-full bg-blue-600 text-white p-2 rounded">Create Event</button>
            </form>
        </div>
    );
};

export default CreateEvent;