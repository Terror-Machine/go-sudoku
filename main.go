// Copyright (c) 2025 HyHy. All rights reserved.
//
// SudokuCLI adalah permainan Sudoku interaktif yang dimainkan
// melalui terminal, dengan output visual berupa gambar PNG.

package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/fogleman/gg"
)

type GameState struct {
	Puzzle    [81]int
	Solution  [81]int
	Board     [81]int
	HintsUsed int
	TimeoutID *time.Timer
}

func findEmpty(board *[9][9]int) (int, int, bool) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] == 0 {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func isValid(board *[9][9]int, row, col, num int) bool {
	for c := 0; c < 9; c++ {
		if board[row][c] == num {
			return false
		}
	}
	for r := 0; r < 9; r++ {
		if board[r][col] == num {
			return false
		}
	}
	startRow, startCol := row-row%3, col-col%3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if board[r+startRow][c+startCol] == num {
				return false
			}
		}
	}
	return true
}

func solveSudoku(board *[9][9]int) bool {
	row, col, found := findEmpty(board)
	if !found {
		return true
	}
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	for _, num := range nums {
		if isValid(board, row, col, num) {
			board[row][col] = num
			if solveSudoku(board) {
				return true
			}
			board[row][col] = 0
		}
	}
	return false
}

func generateSudoku(difficulty float64) ([81]int, [81]int) {
	solutionGrid := [9][9]int{}
	solveSudoku(&solutionGrid)
	puzzleGrid := solutionGrid
	cellsToRemove := int(81 * (1.0 - difficulty))
	for i := 0; i < cellsToRemove; {
		row := rand.Intn(9)
		col := rand.Intn(9)
		if puzzleGrid[row][col] != 0 {
			puzzleGrid[row][col] = 0
			i++
		}
	}
	puzzleSlice := make([]int, 81)
	solutionSlice := make([]int, 81)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			index := r*9 + c
			puzzleSlice[index] = puzzleGrid[r][c]
			solutionSlice[index] = solutionGrid[r][c]
		}
	}
	var p, s [81]int
	copy(p[:], puzzleSlice)
	copy(s[:], solutionSlice)
	return p, s
}

func generateSudokuImage(puzzle, board [81]int, errorIndices map[int]bool) {
	const squareSize = 60.0
	const boardSize = squareSize * 9
	const fontPath = "font/arial.ttf"
	dc := gg.NewContext(int(boardSize), int(boardSize))
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	colLabels := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	rowLabels := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			index := r*9 + c
			x := float64(c * squareSize)
			y := float64(r * squareSize)
			dc.DrawRectangle(x, y, squareSize, squareSize)
			dc.SetColor(color.RGBA{221, 221, 221, 255})
			dc.SetLineWidth(1)
			dc.Stroke()
			if board[index] != 0 {
				isPuzzleNum := puzzle[index] != 0
				isError := errorIndices[index]
				fontSize := 34.0
				if isPuzzleNum {
					fontSize = 36.0
					dc.SetRGB(0, 0, 0)
				} else if isError {
					dc.SetColor(color.RGBA{211, 47, 47, 255})
				} else {
					dc.SetColor(color.RGBA{0, 85, 204, 255})
				}
				if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
					dc.DrawStringAnchored(strconv.Itoa(board[index]), x+squareSize/2, y+squareSize/2, 0.5, 0.5)
				}
			}
		}
	}
	if err := dc.LoadFontFace(fontPath, 14); err == nil {
		dc.SetColor(color.RGBA{0, 0, 0, 150})
		for i := 0; i < 9; i++ {
			row_x := float64(0*squareSize) + 5
			row_y := float64(i*squareSize) + 15
			dc.DrawStringAnchored(rowLabels[i], row_x, row_y, 0, 0)
			col_x := float64(i*squareSize) + squareSize - 5
			col_y := float64(0*squareSize) + 15
			dc.DrawStringAnchored(colLabels[i], col_x, col_y, 1, 0)
		}
	} else {
		log.Println("Gagal memuat font untuk label.")
	}
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(3.0)
	for i := 0; i <= 9; i++ {
		if i%3 == 0 {
			v := float64(i * squareSize)
			dc.DrawLine(v, 0, v, boardSize)
			dc.DrawLine(0, v, boardSize, v)
			dc.Stroke()
		}
	}
	if err := dc.SavePNG("sudoku.png"); err != nil {
		log.Fatalf("Gagal menyimpan gambar: %v", err)
	}
}

func parseCoord(s string) (int, bool) {
	if len(s) != 2 { return 0, false }
	col := int(s[0] - 'a')
	row, err := strconv.Atoi(string(s[1]))
	if err != nil || col < 0 || col > 8 || row < 1 || row > 9 { return 0, false }
	return (row-1)*9 + col, true
}

func newGame(difficulty float64) GameState {
	puzzle, solution := generateSudoku(difficulty)
	var board [81]int
	copy(board[:], puzzle[:])
	return GameState{ Puzzle: puzzle, Solution: solution, Board: board, HintsUsed: 0 }
}

func main() {
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Selamat datang di Game Sudoku CLI!")
	difficultyMap := map[string]float64{"easy": 0.85, "normal": 0.65, "hard": 0.55, "extreme": 0.45}
	game := newGame(difficultyMap["extreme"])
	generateSudokuImage(game.Puzzle, game.Board, make(map[int]bool))
	fmt.Println("Papan Sudoku (extreme) telah dibuat dan disimpan sebagai 'sudoku.png'.")
	fmt.Println("\nPerintah:")
	fmt.Println("  a1 5		-> Mengisi kotak A1 dengan angka 5.")
	fmt.Println("  a1 0		-> Menghapus isian di kotak A1.")
	fmt.Println("  cek		-> Memeriksa jawabanmu.")
	fmt.Println("  hint		-> Mendapatkan bantuan (maks. 3).")
	fmt.Println("  new <level>	-> Memulai game baru (level: easy, normal, hard, extreme).")
	fmt.Println("  exit		-> Keluar dari permainan.")
	gameDuration := 5 * time.Minute
	game.TimeoutID = time.AfterFunc(gameDuration, func() {
		fmt.Println("\nSesi berakhir karena tidak ada aktivitas. Tekan Enter untuk keluar.")
		os.Exit(0)
	})
	fmt.Print("> ")
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			fmt.Print("> ")
			continue
		}
		command := input[0]
		errorIndices := make(map[int]bool)
		switch command {
		case "exit", "keluar":
			fmt.Println("Terima kasih sudah bermain!")
			game.TimeoutID.Stop()
			return
		case "new":
			game.TimeoutID.Stop()
			level := "easy"
			if len(input) > 1 { level = input[1] }
			diff, ok := difficultyMap[level]
			if !ok {
				fmt.Println("Level tidak valid. Gunakan: easy, normal, hard, extreme.")
				break
			}
			game = newGame(diff)
			game.TimeoutID = time.AfterFunc(gameDuration, func() {
				fmt.Println("\nSesi berakhir karena tidak ada aktivitas. Tekan Enter untuk keluar.")
				os.Exit(0)
			})
			fmt.Printf("Game baru (level: %s) telah dibuat! Gambar diperbarui.\n", level)
		case "hint":
			game.TimeoutID.Reset(gameDuration)
			if game.HintsUsed >= 3 {
				fmt.Println("Jatah bantuan sudah habis (3/3).")
				break
			}
			var emptyCells []int
			for i, v := range game.Board {
				if v == 0 { emptyCells = append(emptyCells, i) }
			}
			if len(emptyCells) == 0 {
				fmt.Println("Papan sudah penuh.")
				break
			}
			game.HintsUsed++
			randomIndex := emptyCells[rand.Intn(len(emptyCells))]
			game.Board[randomIndex] = game.Solution[randomIndex]
			fmt.Printf("Bantuan diberikan! Sisa bantuan: %d/3. Gambar diperbarui.\n", 3-game.HintsUsed)
		case "cek":
			game.TimeoutID.Reset(gameDuration)
			isFull, correct := true, true
			for i, v := range game.Board {
				if v == 0 { isFull = false }
				if game.Puzzle[i] == 0 && v != 0 && v != game.Solution[i] {
					errorIndices[i] = true
					correct = false
				}
			}
			if !correct {
				fmt.Printf("Ditemukan %d kesalahan (ditandai warna merah). Gambar diperbarui.\n", len(errorIndices))
			} else if isFull {
				fmt.Println("Luar Biasa! Anda berhasil menyelesaikan puzzle ini!")
				fmt.Print("Ketik 'new' untuk game baru atau 'exit' untuk keluar.\n")
				game.TimeoutID.Stop()
			} else {
				fmt.Println("Sejauh ini semua jawaban benar. Lanjutkan!")
			}
		default:
			if len(input) == 2 {
				game.TimeoutID.Reset(gameDuration)
				index, ok := parseCoord(input[0])
				value, err := strconv.Atoi(input[1])
				if !ok || err != nil || value < 0 || value > 9 {
					fmt.Println("Format gerakan tidak valid. Contoh: a1 5")
					break
				}
				if game.Puzzle[index] != 0 {
					fmt.Println("Kotak ini adalah bagian dari puzzle dan tidak bisa diubah.")
					break
				}
				game.Board[index] = value
				fmt.Printf("Kotak %s diisi dengan %d. Gambar diperbarui.\n", strings.ToUpper(input[0]), value)
				if value != 0 {
					isFull := true
					for _, v := range game.Board {
						if v == 0 {
							isFull = false
							break
						}
					}
					if isFull {
						isSolved := true
						for i := 0; i < 81; i++ {
							if game.Board[i] != game.Solution[i] {
								isSolved = false
								break
							}
						}
						if isSolved {
							fmt.Println("\nLuar Biasa! Anda berhasil menyelesaikan game ini!")
							fmt.Print("Ketik 'new' untuk game baru atau 'exit' untuk keluar.\n")
							game.TimeoutID.Stop()
						}
					}
				}
			} else {
				fmt.Println("Perintah tidak dikenali.")
			}
		}
		generateSudokuImage(game.Puzzle, game.Board, errorIndices)
		fmt.Print("> ")
	}
}