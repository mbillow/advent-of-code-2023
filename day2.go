package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BlockCount struct {
	red   int
	green int
	blue  int
}

type Game struct {
	id       int
	results  []BlockCount
	minPower int
}

func (g *Game) isPossibleWith(max BlockCount) bool {
	possible := true
	for _, res := range g.results {
		if res.red > max.red || res.green > max.green || res.blue > max.blue {
			possible = false
		}
	}
	return possible
}

func parseGame(line string) Game {
	// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
	// to
	// {1 [{4 0 3} {1 2 6} {0 2 0}]}
	s := strings.Split(line, ": ")
	id, _ := strconv.Atoi(strings.Replace(s[0], "Game ", "", 1))

	var results []BlockCount
	var maxRed, maxGreen, maxBlue int

	rawRes := strings.Split(s[1], "; ")
	for _, res := range rawRes {
		var red, green, blue int
		colorCounts := strings.Split(res, ", ")
		for _, count := range colorCounts {
			switch sc := strings.Split(count, " "); sc[1] {
			case "red":
				red, _ = strconv.Atoi(sc[0])
				if red > maxRed {
					maxRed = red
				}
			case "green":
				green, _ = strconv.Atoi(sc[0])
				if green > maxGreen {
					maxGreen = green
				}
			case "blue":
				blue, _ = strconv.Atoi(sc[0])
				if blue > maxBlue {
					maxBlue = blue
				}
			}
		}
		results = append(results, BlockCount{
			red:   red,
			green: green,
			blue:  blue,
		})
	}
	minPower := maxRed * maxGreen * maxBlue
	return Game{id: id, results: results, minPower: minPower}
}

func main() {
	var sum, powerSum int

	maxBlocks := BlockCount{
		red:   12,
		green: 13,
		blue:  14,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		game := parseGame(scanner.Text())
		powerSum += game.minPower
		fmt.Println(game)
		isPossible := game.isPossibleWith(maxBlocks)
		fmt.Println(isPossible)
		fmt.Println("")

		if isPossible {
			sum += game.id
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Sum of Posible Game IDs: %d\n", sum)
	fmt.Printf("Minimum Power Sum: %d\n", powerSum)
}
