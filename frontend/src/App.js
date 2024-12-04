// src/App.js
import React, { useState } from 'react';
import axios from 'axios';

function App() {
  const [logs, setLogs] = useState([]);
  const [log, setLog] = useState('');
  const [logId, setLogId] = useState('');

  const getAllLogs = async () => {
    try {
      const response = await axios.post('http://localhost:3001/query/GetAllLogs', {});
      setLogs(response.data);
    } catch (error) {
      console.error('Error fetching logs:', error);
    }
  };

  const addLog = async () => {
    try {
      const response = await axios.post('http://localhost:3001/invoke/AddLog', { log });
      console.log('Log added:', response.data);
      getAllLogs(); // Refresh logs
    } catch (error) {
      console.error('Error adding log:', error);
    }
  };

  const getLog = async () => {
    try {
      const response = await axios.post('http://localhost:3001/query/GetLog', { id: logId });
      console.log('Log fetched:', response.data);
    } catch (error) {
      console.error('Error fetching log:', error);
    }
  };

  return (
    <div className="App">
      <h1>FireFly Chaincode Interface</h1>
      <div>
        <h2>Add Log</h2>
        <input
          type="text"
          value={log}
          onChange={(e) => setLog(e.target.value)}
          placeholder="Enter log"
        />
        <button onClick={addLog}>Add Log</button>
      </div>
      <div>
        <h2>Get All Logs</h2>
        <button onClick={getAllLogs}>Get All Logs</button>
        {/* <ul>
          {logs.map((log, index) => (
            <li key={index}>{JSON.stringify(log)}</li>
          ))}
        </ul> */}
        <ul>
          {logs.map((log, index) => (
            <li key={index}>
              <p>Destination IP: {log.destination_ip}</p>
              <p>Destination Port: {log.destination_port}</p>
              <p>Event Type: {log.event_type}</p>
              <p>ISP: {log.isp}</p>
              <p>Log ID: {log.logID}</p>
              <p>Message: {log.message}</p>
              <p>Protocol: {log.protocol}</p>
              <p>Severity: {log.severity}</p>
              <p>Source IP: {log.source_ip}</p>
              <p>Source Port: {log.source_port}</p>
              <p>Timestamp: {log.timestamp}</p>
            </li>
          ))}
        </ul>
      </div>
      <div>
        <h2>Get Log by ID</h2>
        <input
          type="text"
          value={logId}
          onChange={(e) => setLogId(e.target.value)}
          placeholder="Enter log ID"
        />
        <button onClick={getLog}>Get Log</button>
      </div>
    </div>
  );
}

export default App;