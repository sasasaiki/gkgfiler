package gkgfiler

import (
	"fmt"
	"go/build"
	"reflect"
	"testing"
)

func TestGetGoSrc(t *testing.T) {
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

			got, e := GetGoSrc()
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

func TestGetPaths(t *testing.T) {
	type args struct {
		dir              string
		matchingPatterns []string
	}
	tests := []struct {
		name        string
		args        args
		wantMatches []string
		wantErr     bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, err := GetPaths(tt.args.dir, tt.args.matchingPatterns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("GetPaths() = %v, want %v", gotMatches, tt.wantMatches)
			}
		})
	}
}

func TestGetPathsRecurcive(t *testing.T) {
	type args struct {
		dir              string
		matchingPatterns []string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaths, err := GetPathsRecurcive(tt.args.dir, tt.args.matchingPatterns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPathsRecurcive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("GetPathsRecurcive() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
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
			gotFileNames, err := getPathsRecurciveImpl(tt.args.dir, tt.args.paths, tt.args.matchingPatterns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPathsRecurciveImpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFileNames, tt.wantFileNames) {
				t.Errorf("getPathsRecurciveImpl() = %v, want %v", gotFileNames, tt.wantFileNames)
			}
		})
	}
}
