package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time" // Importiere das "time"-Paket
)

// Definition des Aufgabe-Structs (außerhalb der main-Funktion)
type Aufgabe struct {
	ID           int
	Titel        string
	Beschreibung string
	ErstelltAm   time.Time // Geändert zu time.Time
	Erledigt     bool
}

// Globale Variablen für die ID-Verwaltung und den Datenspeicher (außerhalb der main-Funktion)
var naechsteID = 1
var aufgaben = make(map[int]Aufgabe)

const dateiName = "aufgaben.json"
const timeFormat = time.RFC3339 // Ein standardisiertes Zeitformat

func main() {

	// Aufgaben beim Start laden
	err := aufgabenLaden()
	if err != nil {
		fmt.Println("Fehler beim Laden der Aufgaben:", err)
	} else {
		// Nach dem Laden die nächste ID anpassen, falls Aufgaben vorhanden sind
		var maxID int
		for id := range aufgaben {
			if id > maxID {
				maxID = id
			}
		}
		naechsteID = maxID + 1
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/neue-aufgabe", neueAufgabeHandler)
	http.HandleFunc("/aufgabe-erledigt/", aufgabeErledigtHandler)
	http.HandleFunc("/aufgabe-loeschen/", aufgabeLoeschenHandler) // Neue Route

	// Handler für statische Dateien im 'templates'-Verzeichnis (optional, für spätere CSS/JS)
	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server startet auf http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Fehler beim Starten des Servers:", err)
		return
	}
}

// Funktion zum Hinzufügen neuer Aufgaben (außerhalb der main-Funktion)
func neueAufgabeHinzufuegen(titel string, beschreibung string) int {
	id := naechsteID
	naechsteID++
	aufgabe := Aufgabe{
		ID:           id,
		Titel:        titel,
		Beschreibung: beschreibung,
		ErstelltAm:   time.Now(), // Verwende time.Now(), um die aktuelle Zeit zu erhalten
		Erledigt:     false,
	}
	aufgaben[id] = aufgabe
	return 0 // Platzhalter, sollte die ID zurückgeben
}

func aufgabenSpeichern() error {
	aufgabenZuSpeichern := make(map[int]serialisierbareAufgabe)
	for id, aufgabe := range aufgaben {
		aufgabenZuSpeichern[id] = serialisierbareAufgabe{
			ID:           aufgabe.ID,
			Titel:        aufgabe.Titel,
			Beschreibung: aufgabe.Beschreibung,
			ErstelltAm:   aufgabe.ErstelltAm.Format(timeFormat), // Formatieren zu String
			Erledigt:     aufgabe.Erledigt,
		}
	}
	daten, err := json.MarshalIndent(aufgabenZuSpeichern, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dateiName, daten, 0644)
}

func aufgabenLaden() error {
	daten, err := os.ReadFile(dateiName)
	if err != nil {
		if os.IsNotExist(err) {
			aufgaben = make(map[int]Aufgabe)
			return nil
		}
		return err
	}
	var geladeneAufgaben map[int]serialisierbareAufgabe
	err = json.Unmarshal(daten, &geladeneAufgaben)
	if err != nil {
		return err
	}
	aufgaben = make(map[int]Aufgabe)
	for id, geladeneAufgabe := range geladeneAufgaben {
		parsedTime, err := time.Parse(timeFormat, geladeneAufgabe.ErstelltAm) // Parsen von String zu time.Time
		if err != nil {
			fmt.Println("Fehler beim Parsen des Datums:", err)
			parsedTime = time.Now() // Fallback bei Fehler
		}
		aufgaben[id] = Aufgabe{
			ID:           geladeneAufgabe.ID,
			Titel:        geladeneAufgabe.Titel,
			Beschreibung: geladeneAufgabe.Beschreibung,
			ErstelltAm:   parsedTime,
			Erledigt:     geladeneAufgabe.Erledigt,
		}
	}
	return nil
}

// Hilfs-Struct für die Serialisierung mit String-Datum
type serialisierbareAufgabe struct {
	ID           int
	Titel        string
	Beschreibung string
	ErstelltAm   string
	Erledigt     bool
}
