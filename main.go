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

var destinations = []string{"Paris", "Miami", "Isle of Sgàil", "New York"}

var targets = map[string][]string{
	"Paris":         {"Viktor Novikov", "Dalia Margolis"},
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
		"Stylist",
		"Tech Crew",
		"Vampire Magician",
	},

	"Miami": {
		"Mascot",
		"Florida Man",
		"Medic",
		"Food Vendor",
		"Kronstadt Engineer",
		"Ted Mendez",
		"Race Marshal",
		"Street Musician",
		"Pale Rider",
	},

	"Isle of Sgàil": {
		"Architect",
		"Ark Member",
		"Chef",
		"Entertainer",
		"Initiate",
		"Jebediah Block",
		"Burial Robes",
		"Master of Ceremonies",
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
		"Fire Axe",
		"Kitchen Knife",
		"Saber",
		"Hatchet",
		"Lethal Poison",
		"Scissors",
		"Bare hands (snap neck)",
		"Letter Opener",
	},

	"Miami": {
		"Starfish",
		"Battle Axe",
		"Amputation Knife",
		"Fire Axe",
		"Kitchen Knife",
		"Lethal Pills",
		"Modern Lethal Syringe",
		"Old Axe",
		"Scissors",
	},

	"Isle of Sgàil": {
		"Battle Axe",
		"Saber",
		"Scalpel",
		"Viking Axe",
		"Cleaver",
		"Hatchet",
	},

	"New York": {
		"Antique Curved Knife",
		"Burial Dagger",
		"Fire Axe",
		"Folding Knife",
		"Letter Opener",
		"Scissors",
		"Hobby Knife",
	},
}

var wildCards = map[string][]string{
	"Paris": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"Drop a speaker.",
		"Drop a chandelier.",
		"Throw an explosive in the catwalk room.",
		"Set off the fireworks.",
		"Create a Bare Knuckle Boxer (poisonous cocktail).",
		"Escape through the basement.",
		"Knock someone out with a piano.",
	},

	"Miami": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
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
		"Find a fish trophy.",
	},

	"Isle of Sgàil": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"Take out an NPC with a filigree egg.",
		"Collect 10 commemorative tokens.",
		"Fire the cannon.",
		"Knock out three people with a fish.",
		"Poison the Constant.",
	},

	"New York": {
		"Choose the starting location.",
		"Finish with 'No Recordings'.",
		"Finish with 'No Bodies Found'.",
		"Pull a fire alarm.",
		"Replace one outfit with one of your choosing.",
		"Replace one weapon with one of your choosing.",

		"Obtain a janitor's key.",
		"Escape in the armored car found in the basement.",
		"Knock out three NPCs with a goldbar.",
		"Find and explode the letterbomb parcel.",
		"Throw a cheeseburger at someone.",
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

	// Randomly select an outfit and weapon for each target and a wild card.
	// Create a file for each player, and write the selections to the file.
	// NOTE: All picked options are removed for other players.
	rand.Seed(time.Now().UnixNano())

	currOutfits := outfits[d]
	currWeapons := weapons[d]
	currWildCards := wildCards[d]
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
