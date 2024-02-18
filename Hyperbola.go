package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var symbolsToCompare = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰჱჲჳჴჵჹჺჶɑɐɒæɓʙβɔɕçɗɖðʤəɘɚɛɜɝɞɟʄɡɠɢʛɦɧħɥʜɨɪʝɭɬɫɮʟɱɯɰŋɳɲɴøɵɸθœɶʘɹɺɾɻʀʁɽʂʃʈʧʉʊʋⱱʌɣɤʍχʎʏʑʐʒʔʡʕʢǀǁǂ"

func processFile(filePath string, counts map[rune]int) {
	// ფაილის გახსნა წარმოადგენს
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("ფაილი %s-ს გახსნის შეცდომა: %s\n", filePath, err)
		return
	}
	defer file.Close()

	// ფაილის შიგაადგენს წაკითხვა
	content := make([]byte, 1024) // წაკითხვის ბუფერი
	for {
		n, err := file.Read(content)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Printf("ფაილი %s-ში შიგაადგენს წაკითხვის შეცდომა: %s\n", filePath, err)
			return
		}
		// ყველა სიმბოლოს თვლა ფაილში
		for _, char := range string(content[:n]) {
			if strings.ContainsRune(symbolsToCompare, char) {
				counts[char]++
			}
		}
	}
}

func visitFile(path string, info os.FileInfo, err error, counts map[rune]int) error {
	if err != nil {
		fmt.Printf("%q მისამართზე შეცდომის მოხერხება: %v\n", path, err)
		return err
	}
	if !info.IsDir() {
		processFile(path, counts)
	}
	return nil
}

func main() {
	folderPath := "./paste"
	counts := make(map[rune]int)

	// ფოლდერში (და მის ქვეფოლდერებში) ფაილების გადასვლა
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		return visitFile(path, info, err, counts)
	})
	if err != nil {
		fmt.Printf("ფოლდერების გადასვლის შეცდომა: %v\n", err)
		return
	}

	// CSV ფაილში შედეგების ჩაწერა
	csvFile, err := os.Create("result.csv")
	if err != nil {
		fmt.Printf("CSV ფაილის შექმნაში შეცდომა: %v\n", err)
		return
	}
	defer csvFile.Close()

	// CSV მწკრივის მეწარმე შექმნა
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// CSV-ში შედეგების ჩაწერა
	for char, count := range counts {
		if err := writer.Write([]string{string(char), fmt.Sprintf("%d", count)}); err != nil {
			fmt.Printf("CSV ფაილში ჩაწერის შეცდომა: %v\n", err)
			return
		}
	}

	fmt.Println("შედეგები ჩაწერილია result.csv ფაილში")
}