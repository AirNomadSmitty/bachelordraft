package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

const teamSize = 6

type pair struct {
	Name string
	Seed int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Seed < p[j].Seed }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func prettyPrint(m map[string][]string) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func reverseDraftOrder(order pairList) {
	for front, back := 0, len(order)-1; front < back; front, back = front+1, back-1 {
		order[front], order[back] = order[back], order[front]
	}
}

func generateRankings(csvName string) map[string][]string {
	csvFile, _ := os.Open(csvName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'

	rankings := make(map[string][]string)
	currentPlayer := ""
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		player := line[0]
		contestant := strings.Replace(line[1], ".", "", -1)
		if player != currentPlayer {
			currentPlayer = player
			rankings[currentPlayer] = []string{}
		}

		rankings[currentPlayer] = append(rankings[currentPlayer], contestant)
	}
	return rankings
}

func doDraft(rankings map[string][]string, contestants map[string]int) map[string][]string {
	teams := make(map[string][]string)
	playerRngs := make(pairList, len(rankings))

	// Generating draft order
	rand.Seed(time.Now().UnixNano())
	i := 0
	for player := range rankings {
		playerRngs[i] = pair{player, rand.Intn(100)}
		i++
		teams[player] = []string{}
	}
	sort.Sort(playerRngs)
	fmt.Println("Draft order:")
	fmt.Println(playerRngs)
	var curRankings []string
	var player pair
	// Drafting round by round
	for round := 0; round < teamSize; round++ {
		// Going down the draft order
		for pos := 0; pos < len(playerRngs); pos++ {
			player = playerRngs[pos]
			curRankings = rankings[player.Name]
			// Loop down current rankings till we find someone draftable
			for rank, contestant := range curRankings {
				if contestants[contestant] > 0 {
					teams[player.Name] = append(teams[player.Name], contestant)
					contestants[contestant]--
					curRankings = append(curRankings[:rank], curRankings[rank+1:]...)
					break
				}
			}
		}
		reverseDraftOrder(playerRngs)
	}

	return teams
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("This program expects exactly 1 argument: an ordered, tab delineated rankings file with one Player-Contestant pair per line")
		os.Exit(1)
	}
	csv := os.Args[1]
	rankings := generateRankings(csv)
	contestantList := make(map[string]int)
	var rankingsTemplate []string
	for _, rankingsTemplate = range rankings {
		break
	}
	occurrances := int(math.Ceil(float64(teamSize*len(rankingsTemplate)) / float64(len(rankingsTemplate))))

	for _, contestant := range rankingsTemplate {
		contestantList[contestant] = occurrances
	}

	teams := doDraft(rankings, contestantList)
	fmt.Println("Teams:")
	prettyPrint(teams)
}
