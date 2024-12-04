const express = require('express');
const bodyParser = require('body-parser');
const FireFly = require('@hyperledger/firefly-sdk').default; // Importa FireFly come default export
const axios = require('axios');
const cors = require('cors');

const app = express();
const port = 3001;

// Middleware
app.use(bodyParser.json());
app.use(cors());

// Initialize FireFly SDK
const firefly = new FireFly({
    host: 'http://localhost:5000', // FireFly API endpoint
    namespace: 'default', // FireFly namespace
});

// Route to add a log
app.post('/invoke/AddLog', async (req, res) => {
    try {
        const response = await axios.post('http://127.0.0.1:5000/api/v1/namespaces/default/apis/logsave2/invoke/AddLog', req.body);
        res.json(response.data);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Route to get all logs
app.post('/query/GetAllLogs', async (req, res) => {
    try {
        const response = await axios.post('http://127.0.0.1:5000/api/v1/namespaces/default/apis/logsave2/query/GetAllLogs', req.body);
        res.json(response.data);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Route to get a specific log
app.post('/query/GetLog', async (req, res) => {
    try {
        const response = await axios.post('http://127.0.0.1:5000/api/v1/namespaces/default/apis/logsave2/query/GetLog', req.body);
        res.json(response.data);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.listen(port, () => {
    console.log(`Backend server is running on http://localhost:${port}`);
});