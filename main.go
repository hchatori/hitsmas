package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var outputDir string

var destinations = []string{"Paris", "Bangkok", "Miami", "Isle of Sgàil", "New York"}

var targets = map[string][]string{
	"Paris":         {"Viktor Novikov", "Dalia Margolis"},
	"Bangkok":       {"Jordan Cross", "Ken Morgan"},
	"Miami":         {"Robert Knox", "Sierra Knox"},
	"Isle of Sgàil": {"Zoe Washington", "Sophia Washington"},
	"New York":      {"Athena Savalas"},
}

var outfits = map[string][]string{
	"Paris": {
		"Chef",
		"Helmut Kruger",
		"Sheikh",
		"Auction Staff",
		"Palace Staff",
		"Security Guard",
		"Stylist",
		"Tech Crew",
		"Vampire Magician",
	},

	"Bangkok": {
		"Exterminator",
		"Hotel Security",
		"Hotel Staff",
		"Kitchen Staff",
		"Jordan Cross' Bodyguard",
		"Morgan's Bodyguard",
		"Waiter",
	},

	"Miami": {
		"Mascot",
		"Florida Man",
		"Medic",
		"Event Security",
		"Food Vendor",
		"Kitchen Staff",
		"Kronstadt Engineer",
		"Kronstadt Mechanic",
		"Ted Mendez",
		"Race Marshal",
	},

	"Isle of Sgàil": {
		"Architect",
		"Ark Member",
		"Castle Staff",
		"Chef",
		"Guard",
		"Entertainer",
		"Initiate",
		"Jebediah Block",
		"Knight's Armour",
		"Raider",
	},

	"New York": {
		"Bank Robber",
		"Bank Teller",
		"High Security Guard",
		"Janitor",
		"Investment Banker",
		"IT Worker",
		"Fired Banker",
		"Job Applicant",
	},
}

var weapons = map[string][]string{
	"Paris": {
		"Gun",
		"Fire Axe",
		"Kitchen Knife",
		"Saber",
		"Hatchet",
		"Fire Extinguisher",
	},

	"Bangkok": {
		"Gun",
		"Hatchet",
		"Katana",
		"Letter Opener",
		"Sapper's Axe",
		"Cleaver",
		"Fire Extinguisher",
	},

	"Miami": {
		"Gun",
		"Starfish",
		"Battle Axe",
		"Amputation Knife",
		"Fire Axe",
		"Kitchen Knife",
		"Fire Extinguisher",
	},

	"Isle of Sgàil": {
		"Gun",
		"Battle Axe",
		"Burial Dagger",
		"Circumcision Knife",
		"Katana",
		"Saber",
		"Scalpel",
		"Viking Axe",
		"Fire Extinguisher",
	},

	"New York": {
		"Gun",
		"Antique Curved Knife",
		"Fire Axe",
		"Folding Knife",
		"Letter Opener",
		"Scissors",
		"Hobby Knife",
		"Fire Extinguisher",
	},
}

var wildCards = map[string][]string{
	"Paris": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Start with no coins.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"Drop a speaker.",
		"Drop a chandelier.",
		"Throw an explosive in the catwalk room.",
		"TODO",
	},

	"Bangkok": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Start with no coins.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"TODO",
	},

	"Miami": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Start with no coins.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"Kill the mascot.",
		"Throw android arm at two non-targets.",
		"Drop the shark.",
		"Destroy evidence of a non-target kill with the wood chipper.",
		"Find a didgeridoo.",
		"Thwack a Thwack Mechanic.",
		"Steal coins from the street musician.",
	},

	"Isle of Sgàil": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Start with no coins.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"TODO",
	},

	"New York": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Start with no coins.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"TODO",
	},
}

func main() {
	// Set output directory for each player's Hitsman run.
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	outputDir = wd + string(os.PathSeparator) + "out"

	fmt.Println("Welcome to Hitsmas!")
	fmt.Println()

	// Ask user for names of players.
	fmt.Println("Hit [Enter] after the name of each player. Hit [Enter] again when done.")
	players := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			break
		}
		players = append(players, input)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Provide a list of destinations and allow user to select one.
	var d string
	for {
		fmt.Println("Select a destination from the following:")
		for i, d := range destinations {
			fmt.Printf("%d. %s\n", i+1, d)
		}
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		i, err := strconv.Atoi(input)
		if err != nil || i < 1 || i > len(destinations) {
			continue
		}

		d = destinations[i-1]
		break
	}

	// Create a file for each player. Then, randomly select an
	// outfit and weapon for each target and a wild card. Write the selections
	// to the file. NOTE: All picked options are removed for other players.

	rand.Seed(time.Now().UnixNano())

	currOutfits := outfits[d]
	currWeapons := weapons[d]
	currWildCards := wildCards[d]

	// Clean any existing outputs.
	err = cleanOutputDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Mkdir(outputDir, os.ModeDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pp := range players {
		data := fmt.Sprintf("Destination: %s\n\n", d)
		for _, tt := range targets[d] {
			o := rand.Intn(len(currOutfits))
			w := rand.Intn(len(currWeapons))
			data += fmt.Sprintf("Eliminate %s as [%s] with [%s].\n", tt, currOutfits[o], currWeapons[w])
			currOutfits = removeEntry(o, currOutfits)
			currWeapons = removeEntry(w, currWeapons)
		}
		c := rand.Intn(len(currWildCards))
		data += fmt.Sprintf("\nWild Card: %s", currWildCards[c])
		currWildCards = removeEntry(c, currWildCards)

		fileName := strings.ToLower(pp) + ".txt"
		f, err := os.Create(outputDir + string(os.PathSeparator) + fileName)
		if err != nil {
			fmt.Println(err)
			_ = cleanOutputDir()
			return
		}

		_, err = f.WriteString(data)
		if err != nil {
			fmt.Println(err)
			f.Close()
			_ = cleanOutputDir()
			return
		}
		f.Close()

		fmt.Printf("%s: Hitsmas run saved to %s\n", pp, outputDir+string(os.PathSeparator)+fileName)
	}
}

func removeEntry(i int, category []string) []string {
	newCatLength := len(category) - 1
	category[i] = category[newCatLength]
	category[newCatLength] = ""
	return category[:newCatLength]
}

func cleanOutputDir() error {
	err := os.RemoveAll(outputDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
