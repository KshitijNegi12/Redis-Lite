import React, { useEffect, useRef, useState } from "react";
import "./App.css";

const App = () => {
    const [input, setInput] = useState("");
    const [output, setOutput] = useState([]);
    const [isConnected, setIsConnected] = useState(false);
    const [isConnecting, setIsConnecting] = useState(true);
    const [error, setError] = useState(null);
    const socketRef = useRef(null);
    const outputEndRef = useRef(null);

    useEffect(() => {
        const connectWebSocket = () => {
            const ws = new WebSocket("ws://localhost:8080");

            ws.onopen = () => {
                console.log("WebSocket connected");
                setIsConnected(true);
                setIsConnecting(false);
                setError(null);
            };

            ws.onmessage = (event) => {
                const timestamp = new Date().toLocaleTimeString();
                setOutput((prev) => [...prev, { message: event.data, timestamp }]);
            };

            ws.onerror = (event) => {
                console.error("WebSocket error:", event);
                setError("WebSocket error occurred");
                setIsConnecting(false);
            };

            ws.onclose = () => {
                console.log("WebSocket disconnected");
                setIsConnected(false);
                setIsConnecting(false);
            };

            socketRef.current = ws;
        };

        if (!socketRef.current) connectWebSocket();

        return () => {
            if (socketRef.current) {
                socketRef.current.close();
                socketRef.current = null;
            }
        };
    }, []);

    useEffect(() => {
        if (outputEndRef.current) {
            outputEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [output]);

    const handleSend = () => {
        if (!input.trim()) {
            setError("Input cannot be empty");
            return;
        }
        if (socketRef.current) {
            socketRef.current.send(input);
            const timestamp = new Date().toLocaleTimeString();
            setOutput((prev) => [...prev, { message: input, timestamp }]);
            setInput("");
            setError(null);
        }
    };

    const handleKeyPress = (event) => {
        if (event.key === "Enter") {
            handleSend();
        }
    };

    const handleClearOutput = () => {
        setOutput([]);
    };

    const handleCopyOutput = () => {
        const textToCopy = output.map((item) => `[${item.timestamp}] ${item.message}`).join("\n");
        navigator.clipboard.writeText(textToCopy);
    };

    const handleRetryConnection = () => {
        setIsConnecting(true);
        if (socketRef.current) {
            socketRef.current.close();
            socketRef.current = null;
        }
        const ws = new WebSocket("ws://localhost:8080");
        socketRef.current = ws;
    };

    const formatMessage = (message) => {
        try {
            const parsed = JSON.parse(message);
            return JSON.stringify(parsed, null, 2);
        } catch (e) {
            return message;
        }
    };

    return (
        <div className="app-container">
            <div className="header">
                <h1>REDIS LITE ADMIN PANEL</h1>
                <div className="connection-info">
                    {isConnecting ? (
                        <div className="connection-status connecting">
                            <span className="spinner"></span> Connecting...
                        </div>
                    ) : (
                        <div className={`connection-status ${isConnected ? "connected" : "disconnected"}`}>
                            {isConnected ? "Connected " : "Disconnected "}
                        </div>
                    )}
                    {!isConnected && !isConnecting && (
                        <button onClick={handleRetryConnection} className="retry-button">
                            Retry
                        </button>
                    )}
                </div>
            </div>
            <div className="main-content">
                <div className="input-section">
                    <h2>Send a Message</h2>
                    <div className="input-group">
                        <input
                            type="text"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            onKeyPress={handleKeyPress}
                            placeholder="Type your Command HERE.."
                            className="input-field"
                        />
                        <button onClick={handleSend} disabled={!isConnected} className="send-button">
                            Send
                        </button>
                    </div>
                    {error && <p className="error-message">{error}</p>}
                </div>
                <div className="output-section">
                    <div className="output-header">
                        <h2>Server Response</h2>
                        <div className="output-actions">
                            <button onClick={handleCopyOutput} className="action-button">
                                Copy
                            </button>
                            <button onClick={handleClearOutput} className="action-button">
                                Clear
                            </button>
                        </div>
                    </div>
                    <div className="output-box">
                        {output.length > 0 ? (
                            output.map((item, index) => (
                                <div key={index} className="output-message">
                                    <span className="timestamp">[{item.timestamp}]</span>
                                    <pre>{formatMessage(item.message)}</pre>
                                </div>
                            ))
                        ) : (
                            <p className="placeholder">Waiting for server response...</p>
                        )}
                        <div ref={outputEndRef}></div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default App;