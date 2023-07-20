const express = require('express');
const bodyParser = require('body-parser');
const app = express();

app.use(bodyParser.json());

let counter = 0;
let restore = false;

origLog = console.log;
console.log = (...args) => {
    origLog(`${new Date().toISOString()}\t`, ...args)
}

app.get('/state', (req, res) => {
    console.log('State was requested!')
    res.json({ counter });
});

app.post('/state', (req, res) => {
    counter = req.body.counter;
    restore = true;
    console.log(`State was restored to: ${counter}`);
    res.end();
});

app.get('/kill', (req, res) => {
    console.log('Kill was requested!')
    res.json({ counter });
    process.exit(1);
});

app.get('/restore', (req, res) => {
    res.json({ restore, counter });
});

app.get('/health', (req, res) => {
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end('200 OK');
    
});

setInterval(() => {
    counter++;
    // don't reset the counter for now. It is easier to measure lost state that way
    /*if (Math.random() > .9) {
        console.log(`Oh nooooo I lost my count! Let's begin again...`);
        counter = 0;
    }*/
    console.log(`${counter} Numbers and I keep counting! For my name is the Count!`);
}, 100);

const server = app.listen(8080, () => {
    const host = server.address().address
    const port = server.address().port
    console.log("I am the Count and will count anything. You can reach me at http://%s:%s", host, port)
});