package gkgfiler

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
)

//GetGoSrcPath is get $GOPATH/src
func GetGoSrcPath() (string, error) {
	// lookup go path
	gopath := build.Default.GOPATH
	if gopath == "" {
		fmt.Println("GOPATHが設定されていません")
		return "", errors.New("GOPATHが設定されていません")
	}
	//  取得した$GOPATHが:つなぎなどで複数設定されていたら一番先頭を使う。
	srcRoot := filepath.Join(filepath.SplitList(gopath)[0], "src")
	return srcRoot, nil

}

//Exist is check file exist
func Exist(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

// func ReplaceTextInFiles(path, projectPath string) {
// 	paths := getPathsRecurcive(path, []string{})
// 	for _, paths := range paths {
// 		replacePathInFile(file, projectPath)
// 	}
// }

//GetPaths ディレクトリ内の、渡したパターンにマッチしたファイルを全て取得する exsample "*.go","*.yaml"
func GetPaths(dir string, matchingPatterns ...string) (matches []string, e error) {
	for _, match := range matchingPatterns {
		files, e := filepath.Glob(filepath.Join(dir + match))
		if e != nil {
			return nil, e
		}
		matches = append(matches, files...)
	}
	return matches, e
}

//GetPathsRecurcive ディレクトリに含まれるファイルのパスを再帰的に取得
func GetPathsRecurcive(dir string, matchingPatterns ...string) (paths []string, e error) {
	return getPathsRecurciveImpl(dir, []string{}, matchingPatterns...)
}

func getPathsRecurciveImpl(dir string, paths []string, matchingPatterns ...string) (fileNames []string, e error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	for _, file := range files {
		if file.IsDir() {
			paths, e = getPathsRecurciveImpl(filepath.Join(dir, file.Name()), paths, matchingPatterns...)
			if e != nil {
				return nil, e
			}
		}
	}

	f, e := GetPaths(filepath.Join(dir), matchingPatterns...)
	if e != nil {
		return nil, e
	}
	return append(paths, f...), e
}

// // 書き込み処理を行う
// func replaceTextInFile(filename, projectPath string) {
// 	input, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	lines := strings.Split(string(input), "\n")

// 	const origin = "github.com/sasasaiki/gokigen"
// 	for i, line := range lines {
// 		lines[i] = strings.Replace(line, origin, projectPath, -1)
// 	}

// 	output := strings.Join(lines, "\n")
// 	err = ioutil.WriteFile(filename, []byte(output), 0644)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
