package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// LogRecord rappresenta la struttura di un log di sicurezza
type LogRecord struct {
	LogID           string `json:"logID"`
	ISP             string `json:"isp"`
	Timestamp       string `json:"timestamp"`
	SourceIP        string `json:"source_ip"`
	DestinationIP   string `json:"destination_ip"`
	SourcePort      int    `json:"source_port"`
	DestinationPort int    `json:"destination_port"`
	Protocol        string `json:"protocol"`
	EventType       string `json:"event_type"`
	Severity        string `json:"severity"`
	Message         string `json:"message"`
}

// SmartContract gestisce i log di sicurezza
type SmartContract struct {
	contractapi.Contract
}

// InitLedger popola il ledger con dati iniziali
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	initialLogs := []LogRecord{
		{
			LogID:           "log1",
			ISP:             "ISP1",
			Timestamp:       "2024-12-01T12:00:00Z",
			SourceIP:        "192.168.1.1",
			DestinationIP:   "10.0.0.1",
			SourcePort:      1234,
			DestinationPort: 80,
			Protocol:        "TCP",
			EventType:       "DDoS",
			Severity:        "High",
			Message:         "Suspicious traffic detected",
		},
		{
			LogID:           "log2",
			ISP:             "ISP1",
			Timestamp:       "2024-12-01T12:05:00Z",
			SourceIP:        "192.168.1.2",
			DestinationIP:   "10.0.0.2",
			SourcePort:      5678,
			DestinationPort: 443,
			Protocol:        "HTTPS",
			EventType:       "PortScan",
			Severity:        "Medium",
			Message:         "Repeated connection attempts detected",
		},
		{
			LogID:           "log3",
			ISP:             "ISP2",
			Timestamp:       "2024-12-01T12:10:00Z",
			SourceIP:        "192.168.2.1",
			DestinationIP:   "10.0.0.3",
			SourcePort:      3456,
			DestinationPort: 22,
			Protocol:        "SSH",
			EventType:       "UnauthorizedAccess",
			Severity:        "Critical",
			Message:         "Multiple failed login attempts",
		},
		{
			LogID:           "log4",
			ISP:             "ISP2",
			Timestamp:       "2024-12-01T12:15:00Z",
			SourceIP:        "192.168.2.2",
			DestinationIP:   "10.0.0.4",
			SourcePort:      7890,
			DestinationPort: 3389,
			Protocol:        "RDP",
			EventType:       "Malware",
			Severity:        "High",
			Message:         "Malicious file transfer detected",
		},
	}

	for _, log := range initialLogs {
		logKey := fmt.Sprintf("%s:%s", log.ISP, log.LogID)
		logJSON, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("errore nella serializzazione del log %s: %v", log.LogID, err)
		}

		err = ctx.GetStub().PutState(logKey, logJSON)
		if err != nil {
			return fmt.Errorf("errore nel salvataggio del log %s: %v", log.LogID, err)
		}
	}

	return nil
}

// LogExists verifica se un log con un determinato ID esiste nel ledger.
func (s *SmartContract) LogExists(ctx contractapi.TransactionContextInterface, logID string) (bool, error) {
	logJSON, err := ctx.GetStub().GetState(logID)
	if err != nil {
		return false, fmt.Errorf("errore nel recupero dello stato: %v", err)
	}

	return logJSON != nil, nil
}

// AddLog registra un log di sicurezza e genera un evento
func (s *SmartContract) AddLog(ctx contractapi.TransactionContextInterface, logID, isp, timestamp, sourceIP, destinationIP string, sourcePort, destinationPort int, protocol, eventType, severity, message string) error {

	// Controlla se il log esiste già
	exists, err := s.LogExists(ctx, logID)
	if err != nil {
		return fmt.Errorf("errore nel controllo dell'esistenza del log: %v", err)
	}
	if exists {
		return fmt.Errorf("il log con ID '%s' esiste già", logID)
	}

	// Genera la chiave del log
	logKey := fmt.Sprintf("%s:%s", isp, logID)

	// Crea il record del log
	logRecord := LogRecord{
		LogID:           logID,
		ISP:             isp,
		Timestamp:       timestamp,
		SourceIP:        sourceIP,
		DestinationIP:   destinationIP,
		SourcePort:      sourcePort,
		DestinationPort: destinationPort,
		Protocol:        protocol,
		EventType:       eventType,
		Severity:        severity,
		Message:         message,
	}
	logBytes, err := json.Marshal(logRecord)
	if err != nil {
		return fmt.Errorf("errore nella serializzazione del log: %v", err)
	}

	// Salva il log nel ledger
	if err := ctx.GetStub().PutState(logKey, logBytes); err != nil {
		return fmt.Errorf("errore nel salvataggio del log %s: %v", logID, err)
	}

	// Emissione di un evento per il log aggiunto
	eventPayload := map[string]string{
		"logID":     logID,
		"isp":       isp,
		"timestamp": timestamp,
	}
	eventBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("errore nella serializzazione dell'evento: %v", err)
	}
	if err := ctx.GetStub().SetEvent("LogAdded", eventBytes); err != nil {
		return fmt.Errorf("errore nell'emissione dell'evento LogAdded: %v", err)
	}

	return nil
}

// GetLog recupera un log di sicurezza
func (s *SmartContract) GetLog(ctx contractapi.TransactionContextInterface, logID, isp string) (*LogRecord, error) {
	// Genera la chiave del log
	logKey := fmt.Sprintf("%s:%s", isp, logID)

	// Recupera il log dal ledger
	data, err := ctx.GetStub().GetState(logKey)
	if err != nil {
		return nil, fmt.Errorf("errore nel recupero del log %s: %v", logID, err)
	}
	if data == nil {
		return nil, fmt.Errorf("log %s non trovato", logID)
	}

	var logRecord LogRecord
	if err := json.Unmarshal(data, &logRecord); err != nil {
		return nil, fmt.Errorf("errore nella deserializzazione del log: %v", err)
	}

	return &logRecord, nil
}

// UpdateLog aggiorna un log esistente nel ledger.
func (s *SmartContract) UpdateLog(ctx contractapi.TransactionContextInterface, logID string, isp string, timestamp string, sourceIP string, destinationIP string, sourcePort int, destinationPort int, protocol string, eventType string, severity string, message string) error {
	// Controlla se il log esiste
	exists, err := s.LogExists(ctx, logID)
	if err != nil {
		return fmt.Errorf("errore nel controllo dell'esistenza del log: %v", err)
	}
	if !exists {
		return fmt.Errorf("il log con ID '%s' non esiste", logID)
	}

	// Crea una nuova istanza di Log con i dati aggiornati
	updatedLog := LogRecord{
		LogID:           logID,
		ISP:             isp,
		Timestamp:       timestamp,
		SourceIP:        sourceIP,
		DestinationIP:   destinationIP,
		SourcePort:      sourcePort,
		DestinationPort: destinationPort,
		Protocol:        protocol,
		EventType:       eventType,
		Severity:        severity,
		Message:         message,
	}

	// Serializza il log aggiornato in JSON
	logJSON, err := json.Marshal(updatedLog)
	if err != nil {
		return fmt.Errorf("errore nella serializzazione del log aggiornato: %v", err)
	}

	// Salva il log aggiornato nel ledger
	return ctx.GetStub().PutState(logID, logJSON)
}

// GetAllLogs returns all logs found in the world state
func (s *SmartContract) GetAllLogs(ctx contractapi.TransactionContextInterface) ([]*LogRecord, error) {
	// Esegue una query aperta per tutti i log nel namespace del chaincode
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("errore nell'esecuzione della query: %v", err)
	}
	defer resultsIterator.Close()

	var logs []*LogRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("errore nel recupero del prossimo log: %v", err)
		}

		var log LogRecord
		// Deserializza il valore JSON nel tipo Log
		err = json.Unmarshal(queryResponse.Value, &log)
		if err != nil {
			return nil, fmt.Errorf("errore nella deserializzazione del log con ID '%s': %v", queryResponse.Key, err)
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

// GetLogHistory recupera la cronologia delle modifiche di un log
func (s *SmartContract) GetLogHistory(ctx contractapi.TransactionContextInterface, logID, isp string) ([]LogRecord, error) {
	// Genera la chiave del log
	logKey := fmt.Sprintf("%s:%s", isp, logID)

	// Recupera la cronologia delle modifiche
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(logKey)
	if err != nil {
		return nil, fmt.Errorf("errore nel recupero della cronologia per il log %s: %v", logID, err)
	}
	defer resultsIterator.Close()

	var history []LogRecord
	for resultsIterator.HasNext() {
		modification, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("errore durante l'iterazione della cronologia: %v", err)
		}

		var logRecord LogRecord
		if err := json.Unmarshal(modification.Value, &logRecord); err != nil {
			return nil, fmt.Errorf("errore nella deserializzazione del log: %v", err)
		}
		history = append(history, logRecord)
	}

	return history, nil
}

// DeleteLog elimina un log di sicurezza specifico
func (s *SmartContract) DeleteLog(ctx contractapi.TransactionContextInterface, logID, isp string) error {

	// Controlla se il log esiste già
	exists, err := s.LogExists(ctx, logID)
	if err != nil {
		return fmt.Errorf("errore nel controllo dell'esistenza del log: %v", err)
	}
	if exists {
		return fmt.Errorf("il log con ID '%s' esiste già", logID)
	}

	// Genera la chiave del log
	logKey := fmt.Sprintf("%s:%s", isp, logID)

	// Recupera il log per assicurarsi che esista
	existing, err := ctx.GetStub().GetState(logKey)
	if err != nil {
		return fmt.Errorf("errore nel controllo del log %s: %v", logID, err)
	}
	if existing == nil {
		return fmt.Errorf("log %s non trovato", logID)
	}

	// Elimina il log
	if err := ctx.GetStub().DelState(logKey); err != nil {
		return fmt.Errorf("errore nell'eliminazione del log %s: %v", logID, err)
	}

	// Emissione di un evento per l'eliminazione del log
	eventPayload := map[string]string{
		"logID": logID,
		"isp":   isp,
	}
	eventBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("errore nella serializzazione dell'evento: %v", err)
	}
	if err := ctx.GetStub().SetEvent("LogDeleted", eventBytes); err != nil {
		return fmt.Errorf("errore nell'emissione dell'evento LogDeleted: %v", err)
	}

	return nil
}

// QueryLogsByISP restituisce tutti i log di un determinato ISP
func (s *SmartContract) QueryLogsByISP(ctx contractapi.TransactionContextInterface, isp string) ([]LogRecord, error) {
	// Creazione di un selettore per la query
	query := fmt.Sprintf(`{"selector":{"isp":"%s"}}`, isp)

	// Esecuzione della query
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("errore nell'esecuzione della query: %v", err)
	}
	defer resultsIterator.Close()

	var logs []LogRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("errore durante l'iterazione dei risultati: %v", err)
		}

		var log LogRecord
		if err := json.Unmarshal(queryResponse.Value, &log); err != nil {
			return nil, fmt.Errorf("errore nella deserializzazione del log: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// QueryLogsBySeverity restituisce i log filtrati per livello di gravità
func (s *SmartContract) QueryLogsBySeverity(ctx contractapi.TransactionContextInterface, severity string) ([]LogRecord, error) {
	// Creazione di un selettore per la query
	query := fmt.Sprintf(`{"selector":{"severity":"%s"}}`, severity)

	// Esecuzione della query
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("errore nell'esecuzione della query: %v", err)
	}
	defer resultsIterator.Close()

	var logs []LogRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("errore durante l'iterazione dei risultati: %v", err)
		}

		var log LogRecord
		if err := json.Unmarshal(queryResponse.Value, &log); err != nil {
			return nil, fmt.Errorf("errore nella deserializzazione del log: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// QueryLogsByEventType restituisce i log filtrati per tipo di evento
func (s *SmartContract) QueryLogsByEventType(ctx contractapi.TransactionContextInterface, eventType string) ([]LogRecord, error) {
	// Creazione di un selettore per la query
	query := fmt.Sprintf(`{"selector":{"event_type":"%s"}}`, eventType)

	// Esecuzione della query
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("errore nell'esecuzione della query: %v", err)
	}
	defer resultsIterator.Close()

	var logs []LogRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("errore durante l'iterazione dei risultati: %v", err)
		}

		var log LogRecord
		if err := json.Unmarshal(queryResponse.Value, &log); err != nil {
			return nil, fmt.Errorf("errore nella deserializzazione del log: %v", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Errore nella creazione del chaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Errore nell'avvio del chaincode: %v\n", err)
	}
}
