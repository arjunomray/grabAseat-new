import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Login = ({ setUserId }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [error, setError] = useState('');

  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(''); // Clear any previous error
    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      const data = await response.json();
      if (response.ok) {
        localStorage.setItem('token', data.token);
        localStorage.setItem('userId', data.userId);
        setUserId(data.userId);
        navigate('/events');
      } else {
        setError(data.error || 'Login failed. Please try again.');
      }
    } catch (error) {
      setError('An error occurred. Please try again.');
      console.error('Error:', error);
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white p-8 rounded shadow">
      <h1 className="text-2xl font-bold mb-4">Login</h1>
      {error && <p className="text-red-600 mb-4">{error}</p>}
      <form onSubmit={handleSubmit}>
        <input type="email" name="email" placeholder="Email" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
        <input type="password" name="password" placeholder="Password" onChange={handleChange} className="w-full p-2 mb-4 border rounded" />
        <button type="submit" className="w-full bg-blue-600 text-white p-2 rounded">Login</button>
      </form>
    </div>
  );
};

export default Login;