package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var emptySpace = "."
var numbers = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var gear = "*"

type Coordinate struct {
	x, y int
}

type PartNumber struct {
	number   int
	location Coordinate
	length   int
}

func checkAdjacentRowSlice(rs []string) bool {
	for _, char := range rs {
		// If the character isn't a number or empty space, we have a hit!
		if _, err := strconv.Atoi(char); err != nil && char != emptySpace {
			//fmt.Printf("found an adjacent '%s'\n", char)
			return true
		}
	}
	return false
}

func findPartNumbers(grid [][]string, index int) ([]PartNumber, []Coordinate) {
	var parts []PartNumber
	var gears []Coordinate
	line := grid[index]

	var currNumber int
	currNumberStartIndex := -1
	for charIndex, char := range line {
		if char == gear {
			gears = append(gears, Coordinate{x: charIndex, y: index})
		}
		// fmt.Printf("currChar: %s currNumber: %d startIndex: %d\n", char, currNumber, currNumberStartIndex)
		if (char == emptySpace || !slices.Contains(numbers, char)) && currNumberStartIndex == -1 {
			fmt.Printf("%s", char)
			continue
		}
		if number, err := strconv.Atoi(char); err == nil {
			if currNumberStartIndex == -1 {
				currNumber = number
				currNumberStartIndex = charIndex
			} else {
				currNumber = (currNumber * 10) + number
				//fmt.Printf("found a consecutive number, now: %d\n", currNumber)
			}
			// If we it the end of the line, don't continue, we've hit the end of a number.
			if charIndex != len(line)-1 {
				continue
			}
		}

		// We've hit the end of a number.
		hasAdjacentSymbol := false

		// If the next character is a symbol, and we aren't at the
		// end of the line: we know we need to add.
		if char != emptySpace && charIndex < len(line)-1 {
			//fmt.Println("symbol found after number, adding")
			hasAdjacentSymbol = true
		}

		// If the start index wasn't zero, check the character before.
		if !hasAdjacentSymbol && currNumberStartIndex != 0 && line[currNumberStartIndex-1] != emptySpace {
			//fmt.Println("symbol found before number, adding")
			hasAdjacentSymbol = true
		}

		adjacentStartIndex := currNumberStartIndex - 1
		if adjacentStartIndex < 0 {
			adjacentStartIndex = 0
		}
		adjacentEndIndex := charIndex + 1

		if !hasAdjacentSymbol && index > 0 {
			// If we aren't at the top, check above.
			above := checkAdjacentRowSlice(grid[index-1][adjacentStartIndex:adjacentEndIndex])
			if above {
				hasAdjacentSymbol = true
			}
		}

		if !hasAdjacentSymbol && index != len(grid)-1 {
			// If we aren't at the bottom, check below.
			below := checkAdjacentRowSlice(grid[index+1][adjacentStartIndex:adjacentEndIndex])
			if below {
				hasAdjacentSymbol = true
			}
		}

		if hasAdjacentSymbol {
			fmt.Printf("%d", currNumber)
			parts = append(parts, PartNumber{
				number:   currNumber,
				location: Coordinate{currNumberStartIndex, charIndex - 1},
				length:   len(strconv.Itoa(currNumber)),
			})
		} else {
			fmt.Printf("%s", strings.Repeat(" ", len(strconv.Itoa(currNumber))))
		}
		if charIndex != len(line)-1 || (charIndex == len(line)-1 && char == emptySpace) {
			fmt.Printf("%s", char)
		}

		// Reset vars
		currNumber = 0
		currNumberStartIndex = -1
	}
	return parts, gears
}

func main() {
	var grid [][]string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		chars := strings.Split(scanner.Text(), "")
		grid = append(grid, chars)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	var sum, gearSum int
	var partGrid [][]*PartNumber
	var gears []Coordinate
	lineLength := len(grid[0])
	for index := range grid {
		var ls int
		partsInRow, gearsInRow := findPartNumbers(grid, index)
		partGridLine := make([]*PartNumber, lineLength)
		for _, part := range partsInRow {
			sum += part.number
			ls += part.number
			// I hate Go
			p := part

			for i := part.location.x; i < part.location.x+part.length; i++ {
				partGridLine[i] = &p
			}
		}
		partGrid = append(partGrid, partGridLine)
		gears = append(gears, gearsInRow...)
		fmt.Printf("   %d\n", ls)
	}

	fmt.Printf("\n\nSUM: %d\n", sum)

	for _, gear := range gears {
		fmt.Println(gear)
		var adjacentParts []*PartNumber
		// Check to the left
		if gear.x != 0 {
			adjacentPart := partGrid[gear.y][gear.x-1]
			if adjacentPart != nil && !slices.Contains(adjacentParts, adjacentPart) {
				fmt.Printf("left: %p\n", adjacentPart)
				adjacentParts = append(adjacentParts, adjacentPart)
			}

		}
		// Check to the right
		if gear.x < lineLength-1 {
			adjacentPart := partGrid[gear.y][gear.x+1]
			if adjacentPart != nil && !slices.Contains(adjacentParts, adjacentPart) {
				fmt.Printf("right: %p\n", adjacentPart)
				adjacentParts = append(adjacentParts, adjacentPart)
			}
		}

		adjacentStart := gear.x - 1
		if adjacentStart < 0 {
			adjacentStart = 0
		}
		adjacentEnd := gear.x + 2
		if adjacentEnd > lineLength {
			adjacentEnd = lineLength - 1
		}

		// Look up and down
		if gear.y != 0 {
			above := partGrid[gear.y-1][adjacentStart:adjacentEnd]
			fmt.Printf("above: [%d:%d] %v\n", adjacentStart, adjacentEnd, above)
			for _, part := range above {
				if part != nil && !slices.Contains(adjacentParts, part) {
					adjacentParts = append(adjacentParts, part)
				}

			}
		}
		if gear.y != len(partGrid)-1 {
			below := partGrid[gear.y+1][adjacentStart:adjacentEnd]
			fmt.Printf("below: [%d:%d] %v\n", adjacentStart, adjacentEnd, below)
			for _, part := range below {
				if part != nil && !slices.Contains(adjacentParts, part) {
					adjacentParts = append(adjacentParts, part)
				}
			}
		}
		fmt.Println(adjacentParts)
		if len(adjacentParts) == 2 {
			power := adjacentParts[0].number * adjacentParts[1].number
			fmt.Printf("%d * %d = %d\n", adjacentParts[0].number, adjacentParts[1].number, power)
			gearSum += power
		}
	}
	fmt.Printf("\n\nGEARS: %d\n", gearSum)
}
