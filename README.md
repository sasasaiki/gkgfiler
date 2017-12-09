[日本語はこちら](http://3jigen.net/2017/12/post-598/ "日本語はこちら")
# gkgfiler

gkgfiler(gokigen-filer) make you happy about file operations.

Sometimes it is troublesome to do something sometimes, such as recursively searching for files, replacing parts of files,adding it to the end of the file, checking if the file contains a string, etc. 

So compiled those into a function.

## package gkgfiler
```
    import "github.com/sasasaiki/gkgfiler"
```

## FUNCTIONS

```
func AppendText(path, appendStr string, perm os.FileMode) error
    AppendText apeend appendStr at the end of the file
    perm is permittion. exsample 0777.
```
```
func Contains(filename, findStr string) (bool, error)
    Contains return whether file contains findStr
```
```
func Exist(path string) bool
    Exist is check file exist
```
```
func GetGoSrcPath() (string, error)
    GetGoSrcPath is get $GOPATH/src
```
```
func GetPaths(dir string, includeDir bool, matchingPatterns ...string) (matches []string, e error)
    GetPaths get files (and directory if includeDir=true) that match
    patterns

	matchingPatterns exsample "*.go","*.yaml"
	if you want all , "*"
```
```
func GetPathsRecurcive(dir string, includeDir bool, matchingPatterns ...string) (paths []string, e error)
    GetPathsRecurcive recursively find and get files (and directory if
    includeDir=true) that match patterns recurcive.

	matchingPatterns exsample "*.go","*.yaml"
	if you want all , "*"
```
```
func IsDir(path string) (bool, error)
    IsDir return whether it is directory or notDirectory
```
```
func ReplaceText(filename, origin, replace string, perm os.FileMode) error
    ReplaceText replace originStr to replaceStr
    perm is permittion. exsample 0777.
```
```
func WriteText(path, str string, createIfNothing bool, perm os.FileMode) error
    WriteText write text to file
    When createIfNothing is set to true, a case where it does not exist is created and written
    perm is permittion. exsample 0777.
```
