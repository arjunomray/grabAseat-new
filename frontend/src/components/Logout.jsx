import { useNavigate } from "react-router";

export const Logout = ({ setUserId }) => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    setUserId(null);
    navigate('/login');
  };

  return (
    <button onClick={handleLogout} className="bg-red-600 text-white p-2 rounded">Logout</button>
  );
};