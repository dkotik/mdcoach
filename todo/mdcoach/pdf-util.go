package main

import (
	"fmt"

	"github.com/dcu/pdf"
)

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	// var buf bytes.Buffer
	// b, err := r.GetPlainText()
	// if err != nil {
	// 	fmt.Println(path)
	// 	panic(err)
	// 	return "", err
	// }
	// buf.ReadFrom(b)
	// fmt.Print(buf.String())
	//
	// return "", nil

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		s, _ := p.GetPlainText(nil)
		fmt.Printf("[page %d] %s\n", pageIndex, s)
		// rows, _ := p.GetTextByRow()
		// for _, row := range rows {
		// 	println("\n")
		// 	// row.Content
		// 	spaceStreak := 0
		// 	for _, word := range row.Content {
		// 		if strings.HasSuffix(word.S, "\n") {
		// 			fmt.Println("\n--------", word.S)
		// 		} else if strings.TrimSpace(word.S) == "" {
		// 			spaceStreak += len(word.S)
		// 		} else {
		// 			if spaceStreak == 1 {
		// 				fmt.Print(" ")
		// 				spaceStreak = 0
		// 			} else if spaceStreak > 1 {
		// 				fmt.Printf("\n\n!!!%d!!! ", spaceStreak)
		// 				spaceStreak = 0
		// 			}
		// 			fmt.Print(word.S)
		// 		}
		// 	}
		// }
		fmt.Print("\n")
	}
	return "", nil
}
