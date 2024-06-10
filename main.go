package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// コマンドライン引数のチェック
	if len(os.Args) != 4 {
		fmt.Println("Usage: HexCharFilter <inputFilePath> <outputFilePath> <searchCharsFilePath>")
		return
	}

	inputFilePath := os.Args[1]
	outputFilePath := os.Args[2]
	searchCharsFilePath := os.Args[3]

	// 検索対象文字の読み込み
	searchChars, err := readSearchCharsHex(searchCharsFilePath)
	if err != nil {
		fmt.Println("Error reading search characters:", err)
		return
	}

	// 入力ファイルの処理
	err = processInputFile(inputFilePath, outputFilePath, searchChars)
	if err != nil {
		fmt.Println("Error processing input file:", err)
	}
}

// 16進数形式で記載された検索対象文字をファイルから読み込む
func readSearchCharsHex(filePath string) (map[rune]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	searchChars := make(map[rune]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// #で始まる場合はコメントとみなし、読み飛ばす
		if strings.HasPrefix(line, "#") {
			continue
		}

		if len(line) > 0 {
			r, err := strconv.ParseInt(line, 16, 32)
			if err != nil {
				return nil, err
			}
			searchChars[rune(r)] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return searchChars, nil
}

// 入力ファイルを処理し、出力ファイルに書き出す
func processInputFile(inputFilePath, outputFilePath string, searchChars map[rune]bool) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// ファイルサイズ
	fi, err := inputFile.Stat()
	if err != nil {
		return err
	}
	filesize := fi.Size()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	r := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	var c int64 = 0
	var oldP int64 = 0
	for {
		line, err := r.ReadString('\n') // LF(\n)まで読み込むので、CRLF(\r\n)でも問題なし
		if err != nil && err != io.EOF {
			return err
		}
		// 最終行に改行がない場合を考慮し、len(row) == 0 を入れる
		if err == io.EOF && len(line) == 0 {
			break
		}
		if containsSearchChar(line, searchChars) {
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
		}

		c = c + int64(len(line))
		p := c / (filesize / 100)
		if p != oldP {
			fmt.Printf("\rReading: %2d%%", p)
			// fmt.Printf("Reading: %2d%%\n", p)
			oldP = p
		}
	}

	fmt.Printf("\nfile size=%d, read size=%d", filesize, c)

	return nil
}

// 行に検索対象文字が含まれているかチェックする
func containsSearchChar(line string, searchChars map[rune]bool) bool {
	for _, char := range line {
		if searchChars[char] {
			return true
		}
	}
	return false
}
