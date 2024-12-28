const express = require('express');
const bodyParser = require('body-parser');
const FireFly = require('@hyperledger/firefly-sdk').default; // Importa FireFly come default export
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
        const response = await firefly.invokeContractAPI(
            'logsave2', // Nome dell'API del contratto
            'AddLog', // Percorso del metodo del contratto
            {
                params: req.body // Parametri della richiesta
            }
        );
        res.json(response);
    } catch (error) {
        console.error('Errore durante l\'aggiunta del log:', error);
        res.status(500).json({ error: error.message });
    }
});

// // Route to get all logs
// app.post('/query/GetAllLogs', async (req, res) => {
//     try {
//         const response = await axios.post('http://127.0.0.1:5000/api/v1/namespaces/default/apis/logsave2/query/GetAllLogs', req.body);
//         res.json(response.data);
//     } catch (error) {
//         res.status(500).json({ error: error.message });
//     }
// });

// Route to get all logs
app.post('/query/GetAllLogs', async (req, res) => {
    try {
        const response = await firefly.queryContractAPI(
            'logsave2', // Nome dell'API del contratto
            'GetAllLogs', // Percorso del metodo del contratto
            {
                params: req.body // Parametri della richiesta
            }
        );
        res.json(response);
    } catch (error) {
        console.error('Errore durante il recupero dei log:', error);
        res.status(500).json({ error: error.message });
    }
});

// Route to get a specific log
app.post('/query/GetLog', async (req, res) => {
    try {
        const response = await firefly.queryContractAPI(
            'logsave2', // Nome dell'API del contratto
            'GetLog', // Percorso del metodo del contratto
            {
                params: req.body // Parametri della richiesta
            }
        );
        res.json(response);
    } catch (error) {
        console.error('Errore durante il recupero del log specifico:', error);
        res.status(500).json({ error: error.message });
    }
});

// Route to upload and publish a log
app.post('/invoke/UploadAndPublishLog', async (req, res) => {
    try {
        // Upload the log data to FireFly
        const uploadResponse = await firefly.uploadData({
            value: req.body.log // Send the log data under the 'value' key
        });

        // Publish the uploaded data
        const publishResponse = await firefly.publishData(uploadResponse.id, {});

        // Get the CID from the publish response
        const cid = publishResponse.public;

        // Respond with the CID
        res.json({ cid });
    } catch (error) {
        console.error('Errore durante il caricamento e la pubblicazione del log:', error);
        res.status(500).json({ error: error.message });
    }
});

// Route to send a broadcast message
app.post('/invoke/BroadcastMessage', async (req, res) => {
    try {
        // Send the broadcast message using FireFly
        const broadcastResponse = await firefly.sendBroadcast({
            header: {
                type: 'broadcast',
                topics: req.body.topics // Topics for the broadcast message
            },
            data: [
                {
                    value: req.body.message // The message to broadcast
                }
            ]
        });

        // Respond with the broadcast response
        res.json(broadcastResponse);
    } catch (error) {
        console.error('Errore durante l\'invio del messaggio in broadcast:', error);
        res.status(500).json({ error: error.message });
    }
});



app.listen(port, () => {
    console.log(`Backend server is running on http://localhost:${port}`);
});