package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		http.Error(w, "Fehler beim Parsen der HTML-Datei", http.StatusInternalServerError)
		fmt.Println("Fehler beim Parsen der HTML-Datei:", err)
		return
	}

	aufgabenListe := []Aufgabe{}
	for _, aufgabe := range aufgaben {
		aufgabenListe = append(aufgabenListe, aufgabe)
	}

	err = tmpl.Execute(w, aufgabenListe)
	if err != nil {
		http.Error(w, "Fehler beim Ausführen des Templates", http.StatusInternalServerError)
		fmt.Println("Fehler beim Ausführen des Templates:", err)
		return
	}
}

func neueAufgabeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Fehler beim Parsen des Formulars", http.StatusBadRequest)
			fmt.Println("Fehler beim Parsen des Formulars:", err)
			return
		}

		titel := r.FormValue("titel")
		beschreibung := r.FormValue("beschreibung")

		neueAufgabeHinzufuegen(titel, beschreibung)

		err = aufgabenSpeichern()
		if err != nil {
			fmt.Println("Fehler beim Speichern der Aufgaben:", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
}

func aufgabeErledigtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := r.URL.Path[len("/aufgabe-erledigt/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Ungültige Aufgaben-ID", http.StatusBadRequest)
			fmt.Println("Ungültige Aufgaben-ID:", err)
			return
		}

		if aufgabe, ok := aufgaben[id]; ok {
			aufgabe.Erledigt = !aufgabe.Erledigt
			aufgaben[id] = aufgabe

			err := aufgabenSpeichern()
			if err != nil {
				fmt.Println("Fehler beim Speichern der Aufgaben:", err)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
}

func aufgabeLoeschenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := r.URL.Path[len("/aufgabe-loeschen/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Ungültige Aufgaben-ID", http.StatusBadRequest)
			fmt.Println("Ungültige Aufgaben-ID:", err)
			return
		}

		// Aufgabe aus der Map entfernen
		if _, ok := aufgaben[id]; ok {
			delete(aufgaben, id)

			// Aufgaben speichern
			err := aufgabenSpeichern()
			if err != nil {
				fmt.Println("Fehler beim Speichern der Aufgaben:", err)
			}
		}

		// Zurück zur Startseite umleiten
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
}
