import { useCallback, useEffect, useRef, useState } from 'react';
import { BrowserQRCodeReader } from '@zxing/browser';
import { apiFetch } from '../api';

const CheckIn = () => {
    const [events, setEvents] = useState([]);
    const [selectedEventId, setSelectedEventId] = useState('');
    const [manualCode, setManualCode] = useState('');
    const [result, setResult] = useState(null);
    const [isScanning, setIsScanning] = useState(false);
    const [cameraError, setCameraError] = useState('');

    const videoRef = useRef(null);
    const controlsRef = useRef(null);
    const lastCodeRef = useRef('');

    useEffect(() => {
        const loadEvents = async () => {
            const userId = localStorage.getItem('userId');
            if (!userId) return;
            try {
                const response = await apiFetch(`/users/${userId}/events`);
                if (!response.ok) return;
                const data = await response.json();
                const merged = [...(data.created_events || []), ...(data.admin_events || [])];
                const byId = new Map(merged.map((event) => [event.id, event]));
                setEvents([...byId.values()]);
            } catch {
                // Leave the list empty; scanning still works.
            }
        };
        loadEvents();
    }, []);

    const verify = useCallback(async (code) => {
        if (!code) return;
        setResult(null);
        try {
            const response = await apiFetch('/tickets/verify', {
                method: 'POST',
                body: JSON.stringify({ code }),
            });
            const body = await response.json().catch(() => ({}));

            if (response.ok) {
                setResult({ status: 'valid', message: 'Valid ticket', ticket: body });
            } else if (response.status === 409) {
                setResult({ status: 'used', message: body.error || 'Ticket already used' });
            } else {
                setResult({ status: 'error', message: body.error || 'Verification failed' });
            }
        } catch {
            setResult({ status: 'error', message: 'Verification failed' });
        }
    }, []);

    const stopScanning = useCallback(() => {
        controlsRef.current?.stop();
        controlsRef.current = null;
        setIsScanning(false);
    }, []);

    useEffect(() => () => controlsRef.current?.stop(), []);

    const startScanning = async () => {
        setCameraError('');
        setResult(null);
        lastCodeRef.current = '';
        try {
            const reader = new BrowserQRCodeReader();
            controlsRef.current = await reader.decodeFromVideoDevice(undefined, videoRef.current, (scanResult) => {
                if (!scanResult) return;
                const text = scanResult.getText();
                if (text === lastCodeRef.current) return;
                lastCodeRef.current = text;
                verify(text);
            });
            setIsScanning(true);
        } catch {
            setCameraError('Could not access the camera. Check permissions and try again.');
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
            <div className="max-w-2xl mx-auto space-y-6">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Ticket check-in</h1>
                    <p className="mt-2 text-sm text-gray-600">
                        Scan a ticket's QR code to check it in for one of your events.
                    </p>
                </div>

                {events.length > 0 && (
                    <div className="bg-white rounded-xl shadow-md p-6 space-y-2">
                        <label htmlFor="event" className="block text-sm font-medium text-gray-700">
                            Event (for reference)
                        </label>
                        <select
                            id="event"
                            value={selectedEventId}
                            onChange={(e) => setSelectedEventId(e.target.value)}
                            className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                            <option value="">Any event I administer</option>
                            {events.map((event) => (
                                <option key={event.id} value={event.id}>{event.name}</option>
                            ))}
                        </select>
                    </div>
                )}

                <div className="bg-white rounded-xl shadow-md p-6 space-y-4">
                    <div className="relative overflow-hidden rounded-lg bg-gray-900 aspect-video">
                        <video ref={videoRef} className="w-full h-full object-cover" muted playsInline />
                        {!isScanning && (
                            <div className="absolute inset-0 flex items-center justify-center text-gray-400 text-sm">
                                Camera preview
                            </div>
                        )}
                    </div>

                    <div className="flex gap-2">
                        {isScanning ? (
                            <button
                                type="button"
                                onClick={stopScanning}
                                className="flex-1 bg-gray-700 text-white py-2 rounded-lg hover:bg-gray-800"
                            >
                                Stop scanning
                            </button>
                        ) : (
                            <button
                                type="button"
                                onClick={startScanning}
                                className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700"
                            >
                                Start camera scan
                            </button>
                        )}
                    </div>

                    {cameraError && <p className="text-sm text-red-600">{cameraError}</p>}

                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            verify(manualCode.trim());
                        }}
                        className="flex gap-2"
                    >
                        <input
                            type="text"
                            value={manualCode}
                            onChange={(e) => setManualCode(e.target.value)}
                            placeholder="Or paste a ticket code"
                            className="flex-grow p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                        <button
                            type="submit"
                            className="bg-blue-600 text-white px-4 rounded-lg hover:bg-blue-700"
                        >
                            Verify
                        </button>
                    </form>
                </div>

                {result && (
                    <div
                        className={`rounded-xl shadow-md p-6 ${
                            result.status === 'valid'
                                ? 'bg-green-50 border border-green-200'
                                : result.status === 'used'
                                    ? 'bg-amber-50 border border-amber-200'
                                    : 'bg-red-50 border border-red-200'
                        }`}
                    >
                        <p
                            className={`text-lg font-semibold ${
                                result.status === 'valid'
                                    ? 'text-green-700'
                                    : result.status === 'used'
                                        ? 'text-amber-700'
                                        : 'text-red-700'
                            }`}
                        >
                            {result.message}
                        </p>
                        {result.ticket?.event && (
                            <p className="mt-1 text-sm text-gray-700">
                                {result.ticket.event.name} · {new Date(result.ticket.event.time).toLocaleString()}
                            </p>
                        )}
                        {result.ticket?.id && (
                            <p className="mt-1 text-xs text-gray-500">Ticket #{result.ticket.id}</p>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default CheckIn;
