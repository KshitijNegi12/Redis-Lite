import React, { useEffect, useRef, useState } from "react";
import "./App.css";

const App = () => {
    const [input, setInput] = useState("");
    const [output, setOutput] = useState([]); // [{ timestamp, command, response }]
    const [isConnected, setIsConnected] = useState(false);
    const [isConnecting, setIsConnecting] = useState(true);
    const [error, setError] = useState(null);
    const socketRef = useRef(null);
    const outputEndRef = useRef(null);
    const [redisPort, setRedisPort] = useState("6379");
    const [hasSentInitial, setHasSentInitial] = useState(false);

    useEffect(() => {
        const connectWebSocket = () => {
            const ws = new WebSocket("ws://localhost:8080");

            ws.onopen = () => {
                console.log("WebSocket connected");
                setIsConnected(true);
                setIsConnecting(false);
                setError(null);
                setHasSentInitial(false);
            };

            ws.onmessage = (event) => {
                const timestamp = new Date().toLocaleTimeString();
                try {
                    const data = JSON.parse(event.data);
                    if (data.type === "error") {
                        setError(data.message);
                    } else if (data.type === "response") {
                        setOutput((prev) => {
                            const updated = [...prev];
                            const last = updated.pop();
                            if (last && last.response === null) {
                                updated.push({ ...last, response: data.message });
                            } else {
                                updated.push(last);
                                updated.push({ timestamp, command: null, response: data.message });
                            }
                            return updated;
                        });
                    } else {
                        setOutput((prev) => [...prev, { timestamp, command: null, response: JSON.stringify(data) }]);
                    }
                } catch {
                    // Handle raw message
                    setOutput((prev) => {
                        const updated = [...prev];
                        const last = updated.pop();
                        if (last && last.response === null) {
                            updated.push({ ...last, response: event.data });
                        } else {
                            updated.push(last);
                            updated.push({ timestamp, command: null, response: event.data });
                        }
                        return updated;
                    });
                }
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
            if (!hasSentInitial) {
                socketRef.current.send(JSON.stringify({ type: "init", port: redisPort }));
                setHasSentInitial(true);
            }

            socketRef.current.send(JSON.stringify({ type: "data", data: input }));
            const timestamp = new Date().toLocaleTimeString();
            setOutput((prev) => [...prev, { timestamp, command: input, response: null }]);
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
        const textToCopy = output.map(
            (item) => `${item.timestamp}${item.command ? ` > ${item.command}` : ""}\n${item.response ?? ""}`
        ).join("\n\n");
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
                <div className="port-config">
                    <label htmlFor="port">Redis Port:</label>
                    <input
                        id="port"
                        type="number"
                        value={redisPort}
                        onChange={(e) => setRedisPort(e.target.value)}
                        className="port-input"
                        placeholder="6379"
                    />
                </div>

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
                                    {item.command && (
                                        <span className="timestamp">
                                            {item.timestamp} &gt; 
                                        </span>
                                    )}
                                    <span>{item.command}</span>
                                    {item.response && <pre>{formatMessage(item.response)}</pre>}
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