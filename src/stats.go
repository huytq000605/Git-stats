package main

import (
	"fmt"
	"sort"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type col []int

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

func stats(email string) {
	commits := processRepositories(email)
	// commits := make(map[int]int)
	printCommits(commits)
	// fmt.Println(commits)
}

func processRepositories(email string) map[int]int {
	statFilePath := getStatFilePath()
	repos := readFileToSlices(statFilePath)

	commits := make(map[int]int, daysInLastSixMonths)
	for i := daysInLastSixMonths; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits

}

func printCommits(commits map[int]int) {
	keys := []int{}
	for key := range commits {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	cols := buildColumns(keys, commits)

	printCells(cols)
}

func printCells(columns map[int]col) {
	printMonths()
	for i := 6; i >= 0; i-- {
		for j := weeksInLastSixMonths + 1; j >= 0; j-- {
			if j == weeksInLastSixMonths+1 {
				printDayCol(i)
			}
			if col, ok := columns[j]; ok {
				//special case today
				if j == 0 && i == calcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > i {
						printCell(col[i], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

func printMonths() {
	yearNow, monthNow, dayNow := time.Now().Date()
	week := time.Date(yearNow, monthNow, dayNow, 0, 0, 0, 0, time.Now().Location()).Add((-time.Hour * 24 * daysInLastSixMonths))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(time.Hour * 24 * 7)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")

}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 0:
		out = " Sun "
	case 1:
		out = " Mon "
	case 2:
		out = " Tue "
	case 3:
		out = " Wed "
	case 4:
		out = " Thu "
	case 5:
		out = " Fri "
	case 6:
		out = " Sat "
	}

	fmt.Print(out)
}

func buildColumns(keys []int, commits map[int]int) map[int]col {
	columns := make(map[int]col)
	column := col{}
	for _, key := range keys {
		week := key / 7      //26,25...1
		dayinweek := key % 7 // 0,1,2,3,4,5,6

		if dayinweek == 0 { //reset
			column = col{}
		}

		column = append(column, commits[key])

		if dayinweek == 6 {
			columns[week] = column
		}
	}
	return columns
}

func fillCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		// fmt.Println(c.Author.Name)
		daysAgo := countDaysSinceDate(c.Author.When) + offset

		if c.Author.Name != email && email != "" {
			return nil
		}

		if daysAgo <= daysInLastSixMonths {
			commits[daysAgo]++
		}
		return nil

	})

	if err != nil {
		panic(err)
	}
	return commits

}

func countDaysSinceDate(date time.Time) int {
	days := 0
	year, month, day := time.Now().Date()
	startOfToday := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	for date.Before(startOfToday) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}
