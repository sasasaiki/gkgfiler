package gkgfiler

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestGetGoSrcPath(t *testing.T) {
	gopath := build.Default.GOPATH
	fmt.Println(gopath)

	tests := []struct {
		name    string
		gopath  string
		want    string
		wantErr bool
	}{
		{
			name:    "GOPATHが単体で設定されていたら+/src",
			gopath:  gopath,
			want:    gopath + "/src",
			wantErr: false,
		},
		{
			name:    "GOPATHが複数で設定されていたら一個め+/src",
			gopath:  "/hoge" + ":" + gopath,
			want:    "/hoge" + "/src",
			wantErr: false,
		},
		{
			name:    "GOPATHが空（設定されていない）ならerror",
			gopath:  "",
			want:    "" + "/src",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			build.Default.GOPATH = tt.gopath

			got, e := GetGoSrcPath()
			if (e != nil) != tt.wantErr {
				t.Errorf("GetGoSrc() = %v, want %v", e != nil, tt.wantErr)
			}

			if e != nil {
				return
			}

			if got != tt.want {
				t.Errorf("GetGoSrc() = %v, want %v", got, tt.want)
			}
		})
	}

	build.Default.GOPATH = gopath
}

func TestExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "pathを渡したファイルが存在すればtrue",
			args: args{
				path: "./gkgfiler.go",
			},
			want: true,
		},
		{
			name: "pathを渡したフォルダが存在すればtrue",
			args: args{
				path: "./vendor",
			},
			want: true,
		},
		{
			name: "pathを渡したファイルが存在しなければfalse",
			args: args{
				path: "./nothing",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exist(tt.args.path); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		existErr bool
	}{
		{
			name: "pathがディレクトリのものならtrue",
			args: args{
				path: "./testDir0",
			},
			want:     true,
			existErr: false,
		},
		{
			name: "pathがディレクトリのものならfalse",
			args: args{
				path: "./gkgfiler.go",
			},
			want:     false,
			existErr: false,
		},
		{
			name: "pathを渡したファイルが存在しなければfalseとエラー",
			args: args{
				path: "./nothing",
			},
			want:     false,
			existErr: true,
		},
	}
	createTestDirsAndFiles()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := IsDir(tt.args.path)
			if got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
			if (e != nil) != tt.existErr {
				t.Errorf("IsDir() = %v, want %v", e, tt.existErr)
			}
		})
	}
	deleteTestDir()
}

func TestGetPaths(t *testing.T) {
	type args struct {
		dir              string
		includeDir       bool
		matchingPatterns []string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
		wantErr   bool
	}{
		{
			name: "testDir0と*とtrueを渡すことでtestDir0の全てのファイルとフォルダを取得できる",
			args: args{
				dir:              "testDir0",
				includeDir:       true,
				matchingPatterns: []string{"*"},
			},
			wantPaths: []string{
				"testDir0/test.yaml",
				"testDir0/test.text",
				"testDir0/testDir1",
				"testDir0/testDir2",
			},
			wantErr: false,
		},
		{
			name: "testDir0と*とfalseを渡すことでtestDir0の全てのファイルだけを取得できる",
			args: args{
				dir:              "testDir0",
				includeDir:       false,
				matchingPatterns: []string{"*"},
			},
			wantPaths: []string{
				"testDir0/test.yaml",
				"testDir0/test.text",
			},
			wantErr: false,
		},
		{
			name: "testDir0/testDir1と*とtrueを渡すことでtestDir1の全てのファイルとフォルダを取得できる",
			args: args{
				dir:              "testDir0/testDir1",
				includeDir:       true,
				matchingPatterns: []string{"*"},
			},
			wantPaths: []string{
				"testDir0/testDir1/test.go",
				"testDir0/testDir1/test1.go",
				"testDir0/testDir1/test.text",
				"testDir0/testDir1/testDir3",
			},
			wantErr: false,
		},
		{
			name: "testDir0/testDir1と*.goとfalseを渡すことでtestDir0の全ての.goファイルを取得できる",
			args: args{
				dir:              "testDir0/testDir1",
				includeDir:       false,
				matchingPatterns: []string{"*.go"},
			},
			wantPaths: []string{
				"testDir0/testDir1/test.go",
				"testDir0/testDir1/test1.go",
			},
			wantErr: false,
		},
		{
			name: "testDir0/testDir1/testDir3と{*.go,*.yaml}とfalseを渡すことでtestDir0の全ての.goと.yamlファイルを取得できる",
			args: args{
				dir:              "testDir0/testDir1/testDir3",
				includeDir:       false,
				matchingPatterns: []string{"*.go", "*.yaml"},
			},
			wantPaths: []string{
				"testDir0/testDir1/testDir3/test.go",
				"testDir0/testDir1/testDir3/test.yaml",
			},
			wantErr: false,
		},
		{
			name: "存在しないpathを渡すとからのスライス",
			args: args{
				dir:              "nothing",
				includeDir:       true,
				matchingPatterns: []string{"*"},
			},
			wantPaths: []string{},
			wantErr:   false,
		},
	}
	createTestDirsAndFiles()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, err := GetPaths(tt.args.dir, tt.args.includeDir, tt.args.matchingPatterns...)
			sort.Strings(gotMatches)
			sort.Strings(tt.wantPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMatches, tt.wantPaths) {
				t.Errorf("GetPaths() = %v, want %v", gotMatches, tt.wantPaths)
			}
		})
	}
	deleteTestDir()
}

func TestGetPathsRecurcive(t *testing.T) {
	type args struct {
		dir              string
		includeDir       bool
		matchingPatterns []string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
		wantErr   bool
	}{
		{
			name: "testDir0と*を渡すことでtestDir0以下の全てのファイルを取得できる",
			args: args{
				dir:              "testDir0",
				matchingPatterns: []string{"*"},
			},
			wantPaths: []string{
				"testDir0/test.yaml",
				"testDir0/test.text",
				"testDir0/testDir1/test.go",
				"testDir0/testDir1/test1.go",
				"testDir0/testDir1/test.text",
				"testDir0/testDir2/test.go",
				"testDir0/testDir2/test.text",
				"testDir0/testDir1/testDir3/test.yaml",
				"testDir0/testDir1/testDir3/test.go",
				"testDir0/testDir1/testDir3/test.text",
			},
			wantErr: false,
		},
		{
			name: "testDir0と*.goを渡すことでtestDir0以下の全ての.goファイルを取得できる",
			args: args{
				dir:              "testDir0",
				matchingPatterns: []string{"*.go"},
			},
			wantPaths: []string{
				"testDir0/testDir1/test.go",
				"testDir0/testDir1/test1.go",
				"testDir0/testDir2/test.go",
				"testDir0/testDir1/testDir3/test.go",
			},
			wantErr: false,
		},
	}

	createTestDirsAndFiles()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaths, err := GetPathsRecurcive(tt.args.dir, tt.args.includeDir, tt.args.matchingPatterns...)
			sort.Strings(gotPaths)
			sort.Strings(tt.wantPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPathsRecurcive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("GetPathsRecurcive() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}

	deleteTestDir()
}

}

func Test_getPathsRecurciveImpl(t *testing.T) {
	type args struct {
		dir              string
		paths            []string
		matchingPatterns []string
	}
	tests := []struct {
		name          string
		args          args
		wantFileNames []string
		wantErr       bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
func TestContains(t *testing.T) {
	type args struct {
		filename string
		findStr  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "渡した文字列が含まれていればtrue",
			args: args{
				filename: "testDir0/test.text",
				findStr:  testText,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "渡した単語が含まれてなければfalse",
			args: args{
				filename: "testDir0/test.text",
				findStr:  "hoge",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "存在しないファイルを渡すとerror",
			args: args{
				filename: "nothing",
				findStr:  "hoge",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "ディレクトリを渡すとerror",
			args: args{
				filename: "testDir0",
				findStr:  "hoge",
			},
			want:    false,
			wantErr: true,
		},
	}
	createTestDirsAndFiles()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.args.filename, tt.args.findStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
	deleteTestDir()
}

func TestReplaceText(t *testing.T) {
	type args struct {
		filename string
		origin   string
		replace  string
		perm     os.FileMode
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		containOrig bool
		containNew  bool
	}{
		{
			name: "originがreplaceに置き換えられる",
			args: args{
				filename: "testDir0/test.text",
				origin:   "test",
				replace:  "replaced",
				perm:     0777,
			},
			wantErr:     false,
			containOrig: false,
			containNew:  true,
		},
		{
			name: "originが存在しなければreplaceに置き換えらない",
			args: args{
				filename: "testDir0/test.text",
				origin:   "nothing",
				replace:  "replaced",
				perm:     0777,
			},
			wantErr:     false,
			containOrig: true,
			containNew:  false,
		},
		{
			name: "存在しないパスを渡すとエラー",
			args: args{
				filename: "nothing",
				origin:   "test",
				replace:  "replaced",
				perm:     0777,
			},
			wantErr:     true,
			containOrig: false,
			containNew:  false,
		},
		{
			name: "ディレクトリのパスを渡すとエラー",
			args: args{
				filename: "testDir0",
				origin:   "test",
				replace:  "replaced",
				perm:     0777,
			},
			wantErr:     true,
			containOrig: false,
			containNew:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTestDirsAndFiles()
			defer deleteTestDir()

			err := ReplaceText(tt.args.filename, tt.args.origin, tt.args.replace, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			gotContainOrig, e := Contains("testDir0/test.text", testText)
			if e != nil {
				return
			}
			if gotContainOrig != tt.containOrig {
				t.Errorf("ReplaceText() containOrig = %v, wantContainOrig %v", gotContainOrig, tt.containOrig)
			}

			gotContainNew, e := Contains("testDir0/test.text", "replaced")
			if e != nil {
				return
			}
			if gotContainNew != tt.containNew {
				t.Errorf("ReplaceText() containNew = %v, wantContainNew %v", gotContainNew, tt.containNew)
			}

		})
	}
}

func createTestDirsAndFiles() {
	os.MkdirAll("./testDir0/testDir1/testDir3", 0777)
	os.Mkdir("./testDir0/testDir2", 0777)
	os.Create("./testDir0/test.yaml")
	os.Create("./testDir0/test.text")
	os.Create("./testDir0/testDir1/test.go")
	os.Create("./testDir0/testDir1/test1.go")
	os.Create("./testDir0/testDir1/test.text")
	os.Create("./testDir0/testDir2/test.go")
	os.Create("./testDir0/testDir2/test.text")
	os.Create("./testDir0/testDir1/testDir3/test.yaml")
	os.Create("./testDir0/testDir1/testDir3/test.go")
	os.Create("./testDir0/testDir1/testDir3/test.text")
	ioutil.WriteFile("./testDir0/test.text", []byte(testText), 0777)
}

const testText = "this is test text"

func deleteTestDir() {
	os.RemoveAll("./testDir0")
}
