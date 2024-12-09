// server.js
const WebSocket = require('ws');

// Create a WebSocket server on port 5000
const server = new WebSocket.Server({ port: 5000 });

server.on('connection', (socket) => {
    console.log('Client connected');

    // Listen for messages from the client
    socket.on('message', (message) => {
        // Log the command received from the client
        console.log(`Your Command: ${message}`);
        
        // Prepare a response
        const response = `Server Response: Server received your command: "${message}"`;
        
        // Send the response back to the client
        socket.send(response);
    });

    // Log when a client disconnects
    socket.on('close', () => {
        console.log('Client disconnected');
    });
});

// Log that the server is running
console.log('WebSocket server is running on ws://localhost:5000');