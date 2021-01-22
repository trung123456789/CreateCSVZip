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
	var numOfSplit, Sum int
	var title string
	var sampleData string
	var fileName, ext string
	fmt.Println("Input file name: ")
	fmt.Scanf("%s\n", &fileName)
	fmt.Println("Input ext: ")
	fmt.Scanf("%s\n", &ext)
	fmt.Println("Input title: ")
	fmt.Scanf("%s\n", &title)
	fmt.Println("Input sample data: ")
	fmt.Scanf("%s\n", &sampleData)

	fmt.Println("Input num of split: ")
	fmt.Scanf("%d\n", &numOfSplit)
	fmt.Println("Input sum line: ")
	fmt.Scanf("%d\n", &Sum)
	if Sum%numOfSplit != 0 {
		fmt.Println("Sum line must be divide by number of split")
		return
	}
	startLine := 1
	endLine := Sum / numOfSplit

	for idx := 0; idx < numOfSplit; idx++ {
		newFolder := fileName + "_" + strconv.Itoa(startLine) + "_" + strconv.Itoa(endLine)
		lastFolder := newFolder + "/" + newFolder
		os.MkdirAll(lastFolder, os.ModePerm)
		tempFileName := lastFolder + "/" + "SampleData" + ext

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
		rangBar := Sum / numOfSplit
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
		endLine += Sum / numOfSplit
	}
}

func ZipWriter(baseFolder string, outFile *os.File) {
	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	//// Add some files to the archive.
	addFiles(w, baseFolder, "")

	//// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
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
			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}
