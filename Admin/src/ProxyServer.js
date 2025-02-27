const net = require('net');
const WebSocket = require('ws');
const { toRESP } = require('./resp/encode');
const { parseRESP } = require('./resp/decode');

const port = 8080
const wss = new WebSocket.Server({ port: port });

const formatData = (data) =>{
    data = data.split(' ');
    if(data.length == 1){
        return data[0];
    }
    data = data.map(ele => {
        const number = Number(ele);
        return !isNaN(number) ? number : ele;
    });

    return data;
}

wss.on('connection', (ws) => {
    console.log("connected");
    const client = net.createConnection({ port: 6379, host: '127.0.0.1' }, () => {});

    ws.on('message', (data) => {
        // console.log(data);
        data = formatData(data)
        console.log(data);
        const encodeData = toRESP(data)
        // console.log(encodeData);
        client.write(encodeData);
    });

    client.on('data', (data) => {
        const parsedData = parseRESP(data.toString().split('\r\n'))[0];
        ws.send(parsedData.toString());
    });

    client.on('error', (err) => {
        console.log("Error during Connection: ", err);
    });

    ws.on('close', () => {
        client.end();
    });

    client.on('close', ()=>{
        console.log("Connection Terminated.");
    })
});

console.log(`Proxy Server running on Port ${port}`);
