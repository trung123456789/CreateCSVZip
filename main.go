package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	var numOfSplit, sum int
	var title string
	var sampleData string
	var fileName, ext string

	fileName = TypeInput("Input file name: ")
	ext = TypeInput("Input ext: ")
	title = TypeInput("Input title: ")
	sampleData = TypeInput("Input sample data: ")
	numOfSplit = TypeInputNumber("Input num of file: ")
	sum = TypeInputNumber("Input sum line: ")

	if sum%numOfSplit != 0 {
		fmt.Println("sum line must be divide by number of split")
		return
	}
	startLine := 1
	endLine := sum / numOfSplit

	for idx := 0; idx < numOfSplit; idx++ {
		newFolder := fileName + "_" + strconv.Itoa(startLine) + "_" + strconv.Itoa(endLine)
		lastFolder := newFolder + "/" + newFolder
		os.MkdirAll(lastFolder, os.ModePerm)
		tempFileName := lastFolder + "/" + fileName + "." + ext

		file, err := os.Create(tempFileName)
		file.WriteString(title + "\n")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		zipName := newFolder + ".zip"
		tmpl := `{{ yellow "%s:" }} {{ bar . "[" "=" (cycle . "↖" "↗" "↘" "↙" ) "-" "]"}} {{speed . | rndcolor }} {{percent .}}`
		ltempl := fmt.Sprintf(tmpl, zipName)
		// start bar based on our template
		rangBar := sum / numOfSplit
		bar := pb.ProgressBarTemplate(ltempl).Start64(int64(rangBar))
		// set values for string elements
		bar.Set("my_green_string", "green").
			Set("my_blue_string", "blue")

		//bar := pb.StartNew(100000)
		for i := startLine; i <= endLine; i++ {
			splitStr := strings.Split(sampleData, ",")

			var rs []string
			for j := 0; j < len(splitStr); j++ {
				strOrg := splitStr[j]
				if strings.Contains(strOrg, "~") {
					numZeroDel := len(strconv.Itoa(i))
					if (len(strconv.Itoa(i)) > strings.Count(strOrg, "0")) {
						numZeroDel = strings.Count(strOrg, "0")
					}
					strDel := strings.Repeat("0", numZeroDel)
					temp := strings.ReplaceAll(strOrg, "~", strconv.Itoa(i))
					strOrg = strings.Replace(temp, strDel, "", 1)
				}
				rs = append(rs, strOrg)
			}
			file.WriteString(strings.Join(rs, ",") + "\n")
			bar.Increment()
		}

		bar.Finish()
		outFile, err := os.Create(zipName)
		if err != nil {
			fmt.Println(err)
		}
		file.Close()
		defer outFile.Close()

		ZipWriter(newFolder+"/", outFile)
		dir, err := ioutil.ReadDir(newFolder)
		for _, d := range dir {
			os.RemoveAll(d.Name())
		}
		startLine = endLine + 1
		endLine += sum / numOfSplit
	}
}

func TypeInput(log string) string {
	var str string
    fmt.Println(log)
    for {
		fmt.Scanf("%s\n", &str)
        if str == "" {
            fmt.Println("Enter a string data:")
        } else {
            return str
        }
    }
}

func TypeInputNumber(log string) int {
	var str string
    var num int = 0
    fmt.Println(log)
    for {
        fmt.Scanf("%s\n", &str)
		if str == "" {
            fmt.Println("Enter a number data:")
			continue
        }
        num, err := strconv.Atoi(str)
		fmt.Println(err)
        if err != nil {
            fmt.Println("Enter a valid number:")
        } else {
            return num
        }
    }
	return num
}

func ZipWriter(baseFolder string, outFile *os.File) {
	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	//// Add some files to the archive.
	AddFiles(w, baseFolder, "")

	//// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func AddFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			AddFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}
