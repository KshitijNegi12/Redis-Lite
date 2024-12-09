// src/App.js
import React, { useEffect, useState } from 'react';

const App = () => {
    const [input, setInput] = useState('');
    const [output, setOutput] = useState('');
    const [socket, setSocket] = useState(null);
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        // Create a WebSocket connection
        const ws = new WebSocket('ws://localhost:8080'); // Change to your WebSocket server URL

        ws.onopen = () => {
            console.log('WebSocket connection established');
            setIsConnected(true);
            setError(null);
        };

        ws.onmessage = (event) => {
            setOutput(event.data); // Update output with the received message
        };

        ws.onerror = (event) => {
            console.error('WebSocket error observed:', event);
            setError('WebSocket error occurred');
        };

        ws.onclose = () => {
            console.log('WebSocket connection closed');
            setIsConnected(false);
        };

        setSocket(ws);

        // Cleanup on component unmount
        return () => {
            ws.close();
        };
    }, []);

    const handleSend = () => {
        if (socket && input) {
            socket.send(input); // Send input to the server
            setInput(''); // Clear input field
        } else {
            setError('Input cannot be empty');
        }
    };

    // Inline styles
    const styles = {
        container: {
            display: 'flex',
            height: '100vh',
            backgroundColor: '#f0f0f0',
        },
        inputArea: {
            flex: 1,
            padding: '20px',
            backgroundColor: '#ffffff',
            borderRight: '2px solid #ccc',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
        },
        outputArea: {
            flex: 1,
            padding: '20px',
            backgroundColor: '#e0e0e0',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
        },
        input: {
            padding: '10px',
            fontSize: '16px',
            marginBottom: '10px',
        },
        button: {
            padding: '10px',
            fontSize: '16px',
            cursor: 'pointer',
        },
        error: {
            color: 'red',
            fontWeight: 'bold',
        },
        status: {
            color: 'orange',
            fontWeight: 'bold',
        },
    };

    return (
        <div style={styles.container}>
            <div style={styles.inputArea}>
                <h2>Input Area</h2>
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Type your message here..."
                    style={styles.input}
                />
                <button onClick={handleSend} disabled={!isConnected} style={styles.button}>
                    Send
                </button>
                {error && <p style={styles.error}>{error}</p>}
            </div>
            <div style={styles.outputArea}>
                <h2>Output Area</h2>
                <p>{output}</p>
                {!isConnected && <p style={styles.status}>Disconnected</p>}
            </div>
        </div>
    );
};

export default App;