import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [log, setLog] = useState({
    logID: '12345',
    isp: 'Example ISP',
    timestamp: '2023-10-01T12:00:00Z',
    source_ip: '192.168.1.1',
    destination_ip: '192.168.1.2',
    source_port: 8080,
    destination_port: 80,
    protocol: 'TCP',
    event_type: 'Access',
    severity: 'High',
    message: 'Unauthorized access attempt detected'
  });

  const [logs, setLogs] = useState([]);
  const [broadcastMessage, setBroadcastMessage] = useState('');
  const [broadcastTopics, setBroadcastTopics] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setLog((prevLog) => ({
      ...prevLog,
      [name]: value
    }));
  };

  const handleBroadcastChange = (e) => {
    const { name, value } = e.target;
    if (name === 'message') {
      setBroadcastMessage(value);
    } else if (name === 'topics') {
      setBroadcastTopics(value);
    }
  };

  const uploadAndPublishLog = async () => {
    try {
      console.log('Sending log data:', log); // Debug: log the data being sent
      const response = await axios.post('http://localhost:3001/invoke/UploadAndPublishLog', { log });
      console.log('Log uploaded and published:', response.data);
      alert(`Log uploaded and published with CID: ${response.data.cid}`);
    } catch (error) {
      console.error('Error uploading and publishing log:', error);
      alert('Error uploading and publishing log');
    }
  };

  const getAllLogs = async () => {
    try {
      const response = await axios.post('http://localhost:3001/query/GetAllLogs', {});
      setLogs(response.data);
    } catch (error) {
      console.error('Error fetching logs:', error);
    }
  };

  const sendBroadcastMessage = async () => {
    try {
      const response = await axios.post('http://localhost:3001/invoke/BroadcastMessage', {
        message: broadcastMessage,
        topics: broadcastTopics.split(',').map(topic => topic.trim())
      });
      console.log('Broadcast message sent:', response.data);
      alert('Broadcast message sent successfully');
    } catch (error) {
      console.error('Error sending broadcast message:', error);
      alert('Error sending broadcast message');
    }
  };

  return (
    <div className="App">
      <h1>Upload and Publish Security Log</h1>
      <div>
        <input type="text" name="logID" value={log.logID} onChange={handleChange} placeholder="Log ID" />
        <input type="text" name="isp" value={log.isp} onChange={handleChange} placeholder="ISP" />
        <input type="text" name="timestamp" value={log.timestamp} onChange={handleChange} placeholder="Timestamp" />
        <input type="text" name="source_ip" value={log.source_ip} onChange={handleChange} placeholder="Source IP" />
        <input type="text" name="destination_ip" value={log.destination_ip} onChange={handleChange} placeholder="Destination IP" />
        <input type="number" name="source_port" value={log.source_port} onChange={handleChange} placeholder="Source Port" />
        <input type="number" name="destination_port" value={log.destination_port} onChange={handleChange} placeholder="Destination Port" />
        <input type="text" name="protocol" value={log.protocol} onChange={handleChange} placeholder="Protocol" />
        <input type="text" name="event_type" value={log.event_type} onChange={handleChange} placeholder="Event Type" />
        <input type="text" name="severity" value={log.severity} onChange={handleChange} placeholder="Severity" />
        <input type="text" name="message" value={log.message} onChange={handleChange} placeholder="Message" />
        <button onClick={uploadAndPublishLog}>Upload and Publish Log</button>
      </div>
      <div>
        <h2>Get All Logs</h2>
        <button onClick={getAllLogs}>Get All Logs</button>
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
        <h2>Send Broadcast Message</h2>
        <input type="text" name="message" value={broadcastMessage} onChange={handleBroadcastChange} placeholder="Message" />
        <input type="text" name="topics" value={broadcastTopics} onChange={handleBroadcastChange} placeholder="Topics (comma separated)" />
        <button onClick={sendBroadcastMessage}>Send Broadcast Message</button>
      </div>
    </div>
  );
}

export default App;