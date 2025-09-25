package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

/* ------------ Modèle JSON ------------- */

type Event struct {
	ID        int       `json:"id"`
	Source    string    `json:"source"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

const historyFile = "events.json"

/* ------------ main (subcommands) ------------- */

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go run   # lancer pipeline concurrent + sauvegarde JSON")
		fmt.Println("  go run main.go show  # afficher l'historique JSON")
		fmt.Println("  go run main.go race  # démo race condition + mutex")
		return
	}
	switch os.Args[1] {
	case "run":
		runPipeline()
	case "show":
		showHistory()
	case "race":
		demoRaceAndMutex() // On entoure la modification par un verrou (Lock/Unlock).
		// Ainsi, une seule goroutine à la fois peut exécuter counter++.
		demoRaceSansMutex() // pour tester la version non protégée
	default:
		fmt.Println("Sous-commande inconnue:", os.Args[1])
	}
}

/* ------------ Pipeline concurrent ------------- */

func runPipeline() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	out := make(chan Event, 8) // channel tamponné pour fluidifier
	var wg sync.WaitGroup
	totalPerProducer := 5

	wg.Add(2)
	go producer("alpha", out, totalPerProducer, &wg)
	go producer("beta", out, totalPerProducer, &wg)

	// Fermera le channel une fois les producteurs terminés
	go func() {
		wg.Wait()
		close(out)
	}()

	var events []Event
	timeout := time.After(3 * time.Second)

	for {
		select {
		case e, ok := <-out:
			if !ok {
				fmt.Println("✅ Flux terminé (channel fermé)")
				saveEvents(events)
				fmt.Printf("Enregistrés: %d évènement(s)\n", len(events))
				return
			}
			fmt.Printf("⬇️  Reçu: %-5s | %s\n", e.Source, e.Payload)
			events = append(events, e)
		case <-timeout:
			fmt.Println("⏰ Timeout atteint, on enregistre ce qu'on a")
			saveEvents(events)
			fmt.Printf("Enregistrés: %d évènement(s)\n", len(events))
			return
		}
	}
}

func producer(name string, out chan<- Event, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= n; i++ {
		time.Sleep(time.Duration(200+rand.Intn(400)) * time.Millisecond)
		out <- Event{
			ID:        i,
			Source:    name,
			Payload:   fmt.Sprintf("message #%d de %s", i, name),
			CreatedAt: time.Now(),
		}
	}
}

/* ------------ JSON I/O ------------- */

func saveEvents(events []Event) {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		fmt.Println("Erreur JSON (marshal):", err)
		return
	}
	if err := os.WriteFile(historyFile, data, 0644); err != nil {
		fmt.Println("Erreur d'écriture fichier:", err)
		return
	}
	fmt.Println("💾 Historique écrit dans", historyFile)
}

func showHistory() {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Aucun historique. Lance d'abord: go run main.go run")
			return
		}
		fmt.Println("Erreur de lecture:", err)
		return
	}
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		fmt.Println("Erreur JSON (unmarshal):", err)
		return
	}
	if len(events) == 0 {
		fmt.Println("Historique vide.")
		return
	}
	fmt.Println("📜 Historique :")
	for i, e := range events {
		fmt.Printf("%2d) %-5s | %s | %s\n", i+1, e.Source, e.CreatedAt.Format(time.RFC3339), e.Payload)
	}
}

/* ------------ Race + Mutex ------------- */

func demoRaceAndMutex() {
	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	// Version protégée par mutex (active par défaut)
	counter = 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Résultat (avec mutex)  =", counter)

	// 👉 Pour tester la version non protégée, commente le bloc ci-dessus et
	//    essaye d'incrémenter sans mu.Lock/mu.Unlock. Lance aussi:
	//    go run -race main.go race

}

func demoRaceSansMutex() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// ❌ accès non protégé
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Résultat (sans mutex) =", counter)
}

// go run main.go run   # génère et sauvegarde des évènements (concurrence + select + JSON)
// go run main.go show  # relit et affiche l'historique
// go run main.go race  # démo mutex (vs. race)
