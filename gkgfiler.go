package gkgfiler

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

//ReplaceText ファイル内のorigin文字列をreplaceに変更する
func ReplaceText(filename, origin, replace string, perm os.FileMode) error {
	input, e := ioutil.ReadFile(filename)
	if e != nil {
		return e
	}

	output := strings.Replace(string(input), origin, filename, -1)

	e = ioutil.WriteFile(filename, []byte(output), perm)
	if e != nil {
		return e
	}

	return nil
}

//Contains ファイル内にfindStrが含まれるかどうかを返す
func Contains(filename, findStr string, perm os.FileMode) (bool, error) {
	input, e := ioutil.ReadFile(filename)
	if e != nil {
		return false, e
	}
	return strings.Contains(string(input), findStr), nil
}

//AppendText ファイルの末尾にテキストを追加
func AppendText(path, appendStr string, perm os.FileMode) error {
	f, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, perm)
	if e != nil {
		return e
	}
	defer f.Close()

	fmt.Fprintln(f, appendStr)

	return nil
}

//WriteText ファイルにテキストを上書き
func WriteText(path, str string, perm os.FileMode) error {
	e := ioutil.WriteFile(path, []byte(str), perm)
	if e != nil {
		return e
	}
	return nil
}
