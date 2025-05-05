import { Link } from 'react-router-dom';

const LandingPage = () => {
    return (
        <div className="min-h-screen bg-gradient-to-b from-blue-500 to-blue-700">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-16">
                {/* Hero Section */}
                <div className="text-center">
                    <h1 className="text-5xl font-extrabold text-white mb-6 animate-fade-in">
                        Welcome to GrabMySeat
                    </h1>
                    <p className="text-xl text-blue-100 mb-12 max-w-2xl mx-auto">
                        Your one-stop destination for discovering and booking amazing events.
                        Join us to explore concerts, workshops, conferences, and more!
                    </p>
                    <div className="space-x-4">
                        <Link
                            to="/register"
                            className="inline-block px-8 py-3 bg-white text-blue-600 font-semibold rounded-lg 
                                     hover:bg-blue-50 transition-colors duration-200 transform hover:scale-105"
                        >
                            Get Started
                        </Link>
                        <Link
                            to="/login"
                            className="inline-block px-8 py-3 border-2 border-white text-white font-semibold 
                                     rounded-lg hover:bg-white hover:text-blue-600 transition-colors duration-200"
                        >
                            Sign In
                        </Link>
                    </div>
                </div>

                {/* Features Section */}
                <div className="mt-24 grid grid-cols-1 md:grid-cols-3 gap-12">
                    <div className="bg-white bg-opacity-10 rounded-xl p-8 text-center transform hover:-translate-y-1 transition-transform duration-200">
                        <div className="text-white mb-4">
                            <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                            </svg>
                        </div>
                        <h3 className="text-xl font-semibold text-white mb-2">Discover Events</h3>
                        <p className="text-blue-100">
                            Find events that match your interests with our powerful search and filter system.
                        </p>
                    </div>

                    <div className="bg-white bg-opacity-10 rounded-xl p-8 text-center transform hover:-translate-y-1 transition-transform duration-200">
                        <div className="text-white mb-4">
                            <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                    d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                            </svg>
                        </div>
                        <h3 className="text-xl font-semibold text-white mb-2">Easy Booking</h3>
                        <p className="text-blue-100">
                            Book your tickets with just a few clicks and secure your spot instantly.
                        </p>
                    </div>

                    <div className="bg-white bg-opacity-10 rounded-xl p-8 text-center transform hover:-translate-y-1 transition-transform duration-200">
                        <div className="text-white mb-4">
                            <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </div>
                        <h3 className="text-xl font-semibold text-white mb-2">Real-time Updates</h3>
                        <p className="text-blue-100">
                            Get instant notifications about event updates and ticket confirmations.
                        </p>
                    </div>
                </div>

                {/* Call to Action */}
                <div className="mt-24 text-center">
                    <h2 className="text-3xl font-bold text-white mb-6">
                        Ready to Experience Amazing Events?
                    </h2>
                    <p className="text-blue-100 mb-8 max-w-2xl mx-auto">
                        Join thousands of users who are already discovering and attending incredible events through GrabMySeat.
                    </p>
                    <Link
                        to="/register"
                        className="inline-block px-8 py-3 bg-white text-blue-600 font-semibold rounded-lg 
                                 hover:bg-blue-50 transition-colors duration-200 transform hover:scale-105"
                    >
                        Create Your Account
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default LandingPage;