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

//GetPaths ディレクトリ内の、渡したパターンにマッチしたファイル,ディレクトリを全て取得する exsample "*.go","*.yaml"
func GetPaths(dir string, includeDir bool, matchingPatterns ...string) (matches []string, e error) {
	matches = []string{}
	for _, match := range matchingPatterns {
		paths, e := filepath.Glob(filepath.Join(dir, match))
		if e != nil {
			return nil, e
		}

		if !includeDir {
			nonDir := make([]string, 0, len(paths))

			for _, p := range paths {
				isDir, e := IsDir(p)
				if e != nil || isDir {
					continue
				}

				nonDir = append(nonDir, p)
			}
			paths = nonDir
		}

		matches = append(matches, paths...)
	}
	return matches, e
}

//IsDir return directory or notDirectory
func IsDir(path string) (bool, error) {
	fInfo, e := os.Stat(path)
	if e != nil {
		return false, e
	}
	return fInfo.IsDir(), nil
}

//GetPathsRecurcive ディレクトリに含まれるファイルのパスを再帰的に取得
func GetPathsRecurcive(dir string, includeDir bool, matchingPatterns ...string) (paths []string, e error) {
	return getPathsRecurciveImpl(dir, []string{}, includeDir, matchingPatterns...)
}

func getPathsRecurciveImpl(dir string, paths []string, includeDir bool, matchingPatterns ...string) (fileNames []string, e error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	for _, file := range files {
		if file.IsDir() {
			paths, e = getPathsRecurciveImpl(filepath.Join(dir, file.Name()), paths, includeDir, matchingPatterns...)
			if e != nil {
				return nil, e
			}
		}
	}

	f, e := GetPaths(filepath.Join(dir), includeDir, matchingPatterns...)
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
