package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	regexp "github.com/dlclark/regexp2"
)

var strNumbers = "(one|two|three|four|five|six|seven|eight|nine)"
var strToDigit = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func findAllStringWithLookAhead(re *regexp.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		// Positive look ahead returns an empty match with the index at the start
		// of the character.
		if m.String() == "" {
			nlr := regexp.MustCompile(strNumbers, regexp.None)
			n, _ := nlr.FindStringMatchStartingAt(s, m.Index)
			matches = append(matches, n.String())
		} else {
			matches = append(matches, m.String())
		}
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func getValueForLine(line string) int {
	fmt.Println(line)
	// One of the lines in the sample imput is "zoneight234"
	// This regex gives us [one, 2, 3, 4] but the string also includes "eight"
	// http://golang.org/pkg/regexp/ explicitly says "non-overlapping" which means
	// we can't use positive look-ahead in our pattern... booo
	// Swapping to https://github.com/dlclark/regexp2
	r := regexp.MustCompile(`(\d|(?=`+strNumbers+"))", regexp.None)
	digits := findAllStringWithLookAhead(r, line)
	fmt.Println(digits)
	for i, digit := range digits {
		for str, num := range strToDigit {
			if digit == str {
				digits[i] = num
			}
		}
	}
	fmt.Println(digits)
	if len(digits) == 1 {
		number, _ := strconv.Atoi(digits[0])
		return number * 11
	}
	first, _ := strconv.Atoi(digits[0])
	last, _ := strconv.Atoi(digits[len(digits)-1])
	return (first * 10) + last
}

func main() {
	sum := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		number := getValueForLine(scanner.Text())
		fmt.Println(number)
		fmt.Println("")
		sum += number
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(sum)
}
