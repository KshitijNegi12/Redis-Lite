const net = require('net');
const WebSocket = require('ws');
const { toRESP } = require('../resp/encode');
const { parseRESP } = require('../resp/decode');

const port = 8080;
const wss = new WebSocket.Server({ port });

const formatData = (data) => {
    data = data.split(' ');
    if (data.length === 1) {
        return data[0];
    }
    return data.map((ele) => {
        const number = Number(ele);
        return !isNaN(number) ? number : ele;
    });
};

wss.on('connection', (ws) => {
    let client = null;

    ws.on('message', (data) => {
        const parsed = JSON.parse(data); 

        if (parsed.type === "init") {
            const redisPort = parseInt(parsed.port) || 6379;

            client = net.createConnection({ port: redisPort, host: '127.0.0.1' }, () => {
                console.log(`Connected to Redis on port ${redisPort}`);
            });

            client.on('data', (data) => {
                const parsedData = parseRESP(data.toString().split('\r\n'))[0];
                ws.send(JSON.stringify({
                    type: 'response',
                    source: 'redis',
                    message: parsedData.toString()
                }));
            });

            client.on('error', (err) => {
                console.log("Error during Redis connection:", err);
                ws.send(JSON.stringify({
                    type: 'error',
                    source: 'proxy',
                    message: "Redis connection is not available. Please reconnect."
                }));
                client.destroy();
            });

            client.on('close', () => {
                console.log("Redis connection closed");
            });

            ws.on('close', () => {
                client.end();
            });

        } else if (parsed.type === "data"){
            const command = formatData(parsed.data);
            const encodeData = toRESP(command);
            client.write(encodeData);
        }
        else{
            ws.send(JSON.stringify({
                type: 'error',
                source: 'proxy',
                message: "Invalid message type received."
            }));
        }
    });
});

console.log(`Proxy Server running on Port ${port}`);
