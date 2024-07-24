const http = require("http");
const fs = require('fs');

const indexHTML = fs.readFileSync('./index.html');

http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/html' });
    res.end(indexHTML);
}).listen(3000, () => {
    console.log('Server is running on http://localhost:3000');
})