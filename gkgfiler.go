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

// Exist is check file exist
func Exist(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

//GetPaths get files (and directory if includeDir=true) that match patterns
/*
	matchingPatterns exsample "*.go","*.yaml"
	if you want all , "*"
*/
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

//IsDir return whether it is directory or notDirectory
func IsDir(path string) (bool, error) {
	fInfo, e := os.Stat(path)
	if e != nil {
		return false, e
	}
	return fInfo.IsDir(), nil
}

//GetPathsRecurcive recursively find and get files (and directory if includeDir=true) that match patterns recurcive.
/*
	matchingPatterns exsample "*.go","*.yaml"
	if you want all , "*"
*/
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

//ReplaceText replace originStr to replaceStr
//perm is permittion. exsample 0777.
func ReplaceText(filename, origin, replace string, perm os.FileMode) error {
	input, e := ioutil.ReadFile(filename)
	if e != nil {
		return e
	}

	output := strings.Replace(string(input), origin, replace, -1)

	e = ioutil.WriteFile(filename, []byte(output), perm)
	if e != nil {
		return e
	}

	return nil
}

//Contains return whether file contains findStr
func Contains(filename, findStr string) (bool, error) {
	input, e := ioutil.ReadFile(filename)
	if e != nil {
		return false, e
	}
	return strings.Contains(string(input), findStr), nil
}

//AppendText apeend appendStr at the end of the file
//perm is permittion. exsample 0777.
func AppendText(path, appendStr string, perm os.FileMode) error {
	f, e := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, perm)
	if e != nil {
		return e
	}
	defer f.Close()

	fmt.Fprintln(f, appendStr)

	return nil
}

//WriteText write text to file
//perm is permittion. exsample 0777.
//When createIfNothing is set to true, a case where it does not exist is created and written
func WriteText(path, str string, createIfNothing bool, perm os.FileMode) error {
	if !createIfNothing && !Exist(path) {
		return errors.New("file not found")
	}
	e := ioutil.WriteFile(path, []byte(str), perm)
	if e != nil {
		return e
	}
	return nil
}
