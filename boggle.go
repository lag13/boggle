package main

import (
	"fmt"
)

// mooreNeighborhood returns the moore neighborhood of an array.
func mooreNeighborhood(numRows int, numCols int, pos int) []int {
	row := pos / numRows
	col := pos % numCols
	neighborCoords := []int{}
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r < 0 || r >= numRows || c < 0 || c >= numCols || (r == row && c == col) {
				continue
			}
			neighborCoords = append(neighborCoords, r*numCols+c)
		}
	}
	return neighborCoords
}

func makeCopy(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func generateBogglePuzzlesWithOneWordStartingFrom(acc [][]string, puzzle []string, lettersLeft string, numRows int, numCols int, coord int) [][]string {
	if lettersLeft == "" {
		return append(acc, makeCopy(puzzle))
	}
	for _, neighboringCoord := range mooreNeighborhood(numRows, numCols, coord) {
		// Every spot not filled will be the empty string and so that
		// means this spot has already been filled and shouldn't be
		// filled again. TODO: I don't like how I'm being a bit lazy
		// here with this (since I feel like there's an unspoken
		// contract that the non-filled spots in the puzzle must be
		// empty strings) but eff it, I don't know if I really like
		// this code at all anyway.
		if puzzle[neighboringCoord] != "" {
			continue
		}
		copiedPuzzle := makeCopy(puzzle)
		copiedPuzzle[neighboringCoord] = lettersLeft[0:1]
		acc = generateBogglePuzzlesWithOneWordStartingFrom(acc, copiedPuzzle, lettersLeft[1:], numRows, numCols, neighboringCoord)
	}
	return acc
}

func generateBogglePuzzlesWithOneWord(numRows int, numCols int, word string) [][]string {
	res := [][]string{}
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			emptyPuzzle := make([]string, numRows*numCols)
			for i := 0; i < numRows*numCols; i++ {
				emptyPuzzle[i] = ""
			}
			coord := r*numRows + c
			emptyPuzzle[coord] = word[0:1]
			res = generateBogglePuzzlesWithOneWordStartingFrom(res, emptyPuzzle, word[1:], numRows, numCols, coord)
		}
	}
	return res
}

func emptyBogglePuzzle(p []string) {
	for i := range p {
		p[i] = ""
	}
}

func puzzlesEqual(p1 []string, p2 []string) bool {
	for i := range p1 {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}

// Rather confusing bit of logic where we take the single number which
// represents a position within our 2d puzzle represented as a 1d
// slice, convert that into a row and column value, convert that into
// a x and y value, then translate those x and y coordinates in such a
// way that the center of the puzzle will be on the origin so we can
// do rotations and all that good stuff.
func pointInPuzzleToCoord(coord int, dim int) (int, int) {
	row := coord / dim
	col := coord % dim
	x := col
	y := -row
	scaleAmt := 1
	translateX := -dim / 2
	translateY := dim / 2
	if dim%2 == 0 {
		scaleAmt = 2
		translateX = -(dim - 1)
		translateY = dim - 1
	}
	x = x*scaleAmt + translateX
	y = y*scaleAmt + translateY
	return x, y
}

// The inverse of the above function
func pointInCoordToPuzzle(x int, y int, dim int) int {
	scaleAmt := 1
	translateX := dim / 2
	translateY := -dim / 2
	if dim%2 == 0 {
		scaleAmt = 2
		translateX = dim - 1
		translateY = -(dim - 1)
	}
	x = (x + translateX) / scaleAmt
	y = (y + translateY) / scaleAmt
	row := -y
	col := x
	return row*dim + col%(dim*dim)
}

// Returns true if the two puzzles are equivalent after rotations
func classicBoggleTwoPuzzlesAreEquivalent(p1 []string, p2 []string) bool {
	dim := 4
	tmp := make([]string, len(p1))
	// Rotate 90 degrees clockwise (x,y) -> (y,-x)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(y, -x, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// Rotate 180 degrees (x,y) -> (-x,-y)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(-x, -y, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// Rotate 270 degrees clockwise (x,y) -> (-y,x)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(-y, x, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// reflection across the x axis (x,y) -> (x,-y)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(x, -y, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// reflection across the y axis (x,y) -> (-x,y)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(-x, y, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// reflection across the line y=x (x,y) -> (y,x)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(y, x, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	// reflection across the line y=-x (x,y) -> (-y,-x)
	for i := range p2 {
		x, y := pointInPuzzleToCoord(i, dim)
		tmp[pointInCoordToPuzzle(-y, -x, dim)] = p2[i]
	}
	if puzzlesEqual(p1, tmp) {
		return true
	}
	return false
}

// Returns true of the puzzle P is contained within PS.
func containsEquivalentPuzzle(xs [][]string, y []string) bool {
	for _, x := range xs {
		if classicBoggleTwoPuzzlesAreEquivalent(x, y) {
			return true
		}
	}
	return false
}

type Coord struct {
	X int
	Y int
}

func coordsAreSymmetricallyEquivalent(pos1 int, pos2 int, dim int) bool {
	x1, y1 := pointInPuzzleToCoord(pos1, dim)
	x2, y2 := pointInPuzzleToCoord(pos2, dim)
	coords := []Coord{
		{X: x1, Y: y1},
		// Rotate 90 degrees clockwise (x,y) -> (y,-x)
		{X: y1, Y: -x1},
		// Rotate 180 degrees (x,y) -> (-x,-y)
		{X: -x1, Y: -y1},
		// Rotate 270 degrees clockwise (x,y) -> (-y,x)
		{X: -y1, Y: x1},
		// reflection across the x axis (x,y) -> (x,-y)
		{X: x1, Y: -y1},
		// reflection across the y axis (x,y) -> (-x,y)
		{X: -x1, Y: y1},
		// reflection across the line y=x (x,y) -> (y,x)
		{X: y1, Y: x1},
		// reflection across the line y=-x (x,y) -> (-y,-x)
		{X: -y1, Y: -x1},
	}
	for _, coord := range coords {
		if coord.X == x2 && coord.Y == y2 {
			return true
		}
	}
	return false
}

func mooreNeighborhoodMinusOccupiedSpots(numRows int, numCols int, pos int, puzzle []string) []int {
	neighborhoodCoords := []int{}
	for _, neighboringCoord := range mooreNeighborhood(numRows, numCols, pos) {
		if puzzle[neighboringCoord] != "" {
			continue
		}
		neighborhoodCoords = append(neighborhoodCoords, neighboringCoord)
	}
	return neighborhoodCoords
}

func mooreNeighborhoodMinusOccupiedSpotsAndSymmetricallyEquivalentPts(numRows int, numCols int, pos int, puzzle []string) []int {
	resCoords := []int{}
	for _, neighboringCoord := range mooreNeighborhoodMinusOccupiedSpots(numRows, numCols, pos, puzzle) {
		addCoordToRes := true
		for _, resCoord := range resCoords {
			if coordsAreSymmetricallyEquivalent(resCoord, neighboringCoord, numRows) {
				addCoordToRes = false
				break
			}
		}
		if addCoordToRes {
			resCoords = append(resCoords, neighboringCoord)
		}
	}
	return resCoords
}

func generateBogglePuzzlesWithOneWordStartingFromEfficient(acc [][]string, puzzle []string, lettersLeft string, numRows int, numCols int, lastInsertedLetterCoord int) [][]string {
	if lettersLeft == "" {
		return append(acc, makeCopy(puzzle))
	}
	copiedPuzzles := [][]string{}
	lastInsertedLetterCoords := []int{}
	for _, neighboringCoord := range mooreNeighborhoodMinusOccupiedSpots(numRows, numCols, lastInsertedLetterCoord, puzzle) {
		// TODO: All this copying shit is exactly why I think I'd
		// prefer something like clojure where it does the copying for
		// you but also does it efficiently by sharing things and
		// whatnot. There's just too much room for error here. I feel
		// like it also gets me worried about performance and
		// efficient memory usage before I want to worry about such
		// things. I believe it was the little schemer that said to
		// first make something correct and then worry about making it
		// the best version of that correctness (it was worded much
		// better of course in the book).
		copiedPuzzle := makeCopy(puzzle)
		copiedPuzzle[neighboringCoord] = lettersLeft[0:1]
		alreadySeen := false
		for _, cp := range copiedPuzzles {
			if classicBoggleTwoPuzzlesAreEquivalent(copiedPuzzle, cp) {
				alreadySeen = true
			}
		}
		if !alreadySeen {
			copiedPuzzles = append(copiedPuzzles, copiedPuzzle)
			lastInsertedLetterCoords = append(lastInsertedLetterCoords, neighboringCoord)
		}
	}
	for i := range copiedPuzzles {
		acc = generateBogglePuzzlesWithOneWordStartingFromEfficient(acc, copiedPuzzles[i], lettersLeft[1:], numRows, numCols, lastInsertedLetterCoords[i])
	}
	return acc
}

// TODO: If this is given a word with duplicate letters then it will
// generate some symmetrically duplicate boards so I think we'll have
// to prune those after the fact which is unfortunate since that is a
// slow process for for words with lengths >= 6 is slow (6 letters = 5
// seconds, 7 letters = 1 minute, 8 letters = minutes). Luckily I
// don't have to worry about it yet since I currently only want to use
// a seed word "veronica" which has unique letters so I'll worry about
// it later.
func generateClassicBogglePuzzlesWithOneWordEfficient(word string, dim int) [][]string {
	tmp := [][]string{}
	startingCoords := []int{}
	for i := 0; i < dim*dim; i++ {
		haveNotSeenBefore := true
		for _, coord := range startingCoords {
			if coordsAreSymmetricallyEquivalent(i, coord, dim) {
				haveNotSeenBefore = true
				break
			}
		}
		if haveNotSeenBefore {
			startingCoords = append(startingCoords, i)
		}
	}
	for _, coord := range startingCoords {
		emptyPuzzle := make([]string, dim*dim)
		for i := 0; i < dim*dim; i++ {
			emptyPuzzle[i] = ""
		}
		emptyPuzzle[coord] = word[0:1]
		tmp = generateBogglePuzzlesWithOneWordStartingFromEfficient(tmp, emptyPuzzle, word[1:], dim, dim, coord)
	}
	return tmp
}

func canOverlayPuzzles(p1 []string, p2 []string) bool {
	for i := range p1 {
		if p1[i] == "" || p2[i] == "" || p1[i] == p2[i] {
			continue
		}
		return false
	}
	return true
}

func overlayPuzzles(p1 []string, p2 []string) []string {
	res := makeCopy(p1)
	for i := range p2 {
		if p2[i] == "" {
			continue
		}
		res[i] = p2[i]
	}
	return res
}

func reverseBoggleHelper(seedBoard []string, listOfListsOfPuzzlesToOverlay [][][]string) []string {
	if len(listOfListsOfPuzzlesToOverlay) == 0 {
		return seedBoard
	}
	for _, toOverlay := range listOfListsOfPuzzlesToOverlay[0] {
		if !canOverlayPuzzles(seedBoard, toOverlay) {
			continue
		}
		return reverseBoggleHelper(overlayPuzzles(seedBoard, toOverlay), listOfListsOfPuzzlesToOverlay[1:])
	}
	return nil
}

// reverseBoggle constructs a boggle puzzle from a set of words
func reverseBoggle(words []string) []string {
	seedBoards := generateClassicBogglePuzzlesWithOneWordEfficient(words[0], 4)
	listOfListsOfPuzzlesToOverlay := make([][][]string, len(words)-1)
	for i := 1; i < len(words); i++ {
		listOfListsOfPuzzlesToOverlay[i-1] = generateBogglePuzzlesWithOneWord(4, 4, words[i])
	}
	for _, seedBoard := range seedBoards {
		res := reverseBoggleHelper(seedBoard, listOfListsOfPuzzlesToOverlay)
		if res != nil {
			return res
		}
	}
	return nil
}

func printPuzzle(puzzle []string, numRows int, numCols int) {
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if puzzle[row*numRows+col] == "" {
				fmt.Printf("- ")
			} else {
				fmt.Printf("%s ", puzzle[row*numRows+col])
			}
		}
		fmt.Println()
	}
}

func main() {
	// for i := 0; i < 16; i++  {
	// 	// pointInCoordToPuzzle
	// 	x, y :=	pointInPuzzleToCoord(i, 4)
	// 	fmt.Println(pointInCoordToPuzzle(x, y, 4))
	// }

	res := reverseBoggle([]string{"veronica", "lucas", "danny", "anais", "janey", "joey"})
	if res == nil {
		fmt.Println("no solution")
	} else {
		printPuzzle(res, 4, 4)
	}
	// for i := range res {
	// 	if i % 4 == 0 {
	// 		fmt.Println()
	// 	}
	// 	if res[i] == "" {
	// 		fmt.Printf("  ")
	// 		continue
	// 	}
	// 	fmt.Printf("%s ", res[i])
	// }
	// fmt.Println(reverseBoggle([]string{"lucas", "danny", "anais"}))

	// for i := 1; i <= 8; i++ {
	// 	start := time.Now()
	// 	fmt.Println(i, len(generateClassicBogglePuzzlesWithOneWordEfficient("veronica"[0:i])))
	// 	fmt.Println(time.Since(start))
	// }

	// generateClassicBogglePuzzlesWithOneWordEfficient("ver")

	// fmt.Println(overlayPuzzles([]string{"a", "c", "", ""}, []string{"a", "", "d", ""}))
	// // for r := 0; r < 16; r++ {
	// // 	fmt.Println(mooreNeighborhood(4, 4, r))
	// // }
	// fmt.Printf("%q\n", "hello"[0:1])
}

// Before I optimized the generation of a set of boards with a single
// word that ONLY contained symmetrically unique words I would
// basically generate all boards and then remove any symmetric
// duplicates. The timings for words of various lengths escalated
// pretty quickly!
//
// 1 3
// 350.9nanoseconds
// 2 12
// 357.3nanoseconds
// 3 52
// 2.9824ms
// 4 221
// 39.1982ms
// 5 839
// 548.2069ms
// 6 2834
// 5.7015548s
// 7 8534
// 59.1626271s
// 8 22934
// 6m18.5241482s
//
// After my optimization which prunes boards that are symmetrically
// equivalent to other boards on the same depth, we get these MUCH
// improved timings:
//
// 1 3
// 0s
// 2 12
// 534.3nanoseconds
// 3 52
// 530.5nanoseconds
// 4 221
// 527.3nanoseconds
// 5 839
// 2.3188ms
// 6 2834
// 9.9191ms
// 7 8534
// 21.5829ms
// 8 22934
// 54.8924ms

// TODO: It would be cool to generate a boggle puzzle that has a
// contains a set of words that could be found toroidally. That would
// honestly be an interesting modification addition to boggle
// actually.
