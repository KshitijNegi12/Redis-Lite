import React, { useEffect, useRef, useState } from "react";
import "./App.css";

const App = () => {
    const [input, setInput] = useState("");
    const [output, setOutput] = useState("");
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState(null);
    const socketRef = useRef(null);

    useEffect(() => {
        if (!socketRef.current) { 
            const ws = new WebSocket("ws://localhost:8080");

            ws.onopen = () => {
                console.log("WebSocket connection established");
                setIsConnected(true);
                setError(null);
            };

            ws.onmessage = (event) => {
                setOutput(event.data);
            };

            ws.onerror = (event) => {
                console.error("WebSocket error observed:", event);
                setError("WebSocket error occurred");
            };

            ws.onclose = () => {
                console.log("WebSocket connection closed");
                setIsConnected(false);
            };

            socketRef.current = ws;
        }

        return () => {};
    }, []);

    const handleSend = () => {
        if (socketRef.current && input.trim() !== "") {
            socketRef.current.send(input);
            setInput("");
            setError(null);
        } else {
            setError("Input cannot be empty");
        }
    };

    return (
        <div className="container">
            <div className="input-area">
                <h2>Input Area</h2>
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    placeholder="Type your message here..."
                    className="input"
                />
                <button onClick={handleSend} disabled={!isConnected} className="button">
                    Send
                </button>
                {error && <p className="error">{error}</p>}
            </div>
            <div className="output-area">
                <h2>Output Area</h2>
                <p>{output && `Server Response: ${output}`}</p>
                {!isConnected && <p className="status">Disconnected</p>}
            </div>
        </div>
    );
};

export default App;