package main

import (
	"reflect"
	"testing"
)

func TestPathMapper_mapUrlPathToFilesystem(t *testing.T) {
	type fields struct {
		rootDir string
	}
	type args struct {
		urlPath string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		calculatedPath string
		fileExtension  string
	}{
		{
			name:           "combines request-url and root-dir",
			fields:         fields{"/video"},
			args:           args{"/Documentations/Nature"},
			calculatedPath: "/video/Documentations/Nature",
			fileExtension:  "",
		}, {
			name:           "cleans result-directory",
			fields:         fields{"/video/./"},
			args:           args{"/Documentations/../Science Fiction/./StarTrek/"},
			calculatedPath: "/video/Science Fiction/StarTrek",
			fileExtension:  "",
		}, {
			name:           "avoids directory traversal",
			fields:         fields{"/video"},
			args:           args{"/Documentations/../../../../../"},
			calculatedPath: "/video",
			fileExtension:  "",
		}, {
			name:           "handles minimal number of slashes",
			fields:         fields{"/video"},
			args:           args{"Documentations/Nature"},
			calculatedPath: "/video/Documentations/Nature",
			fileExtension:  "",
		}, {
			name:           "handles multiple slashes",
			fields:         fields{"/video/"},
			args:           args{"/Documentations//Nature/BBC//"},
			calculatedPath: "/video/Documentations/Nature/BBC",
			fileExtension:  "",
		}, {
			name:           "handles filenames",
			fields:         fields{"/video/"},
			args:           args{"/Documentations/Nature/BBC/Apes.avi"},
			calculatedPath: "/video/Documentations/Nature/BBC/Apes.avi",
			fileExtension:  "avi",
		}, {
			name:           "handles multiple file-extensions",
			fields:         fields{"/video/"},
			args:           args{"/Documentations/Nature/BBC/Apes.old.flv"},
			calculatedPath: "/video/Documentations/Nature/BBC/Apes.old.flv",
			fileExtension:  "flv",
		}, {
			name:           "handles uppercase file-extensions",
			fields:         fields{"/video/"},
			args:           args{"/Documentations/Nature/BBC/Apes.FLV"},
			calculatedPath: "/video/Documentations/Nature/BBC/Apes.FLV",
			fileExtension:  "flv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directoryMapper := PathMapper{
				rootDir: tt.fields.rootDir,
			}
			got := directoryMapper.MapUrlPathToFilesystem(tt.args.urlPath);
			if !reflect.DeepEqual(got.CalculatedPath, tt.calculatedPath) {
				t.Errorf("MapUrlPathToFilesystem().CalculatedPath = %v, want %v", got, tt.calculatedPath)
			}
			if !reflect.DeepEqual(got.FileExtension, tt.fileExtension) {
				t.Errorf("MapUrlPathToFilesystem().FileExtension = %v, want %v", got, tt.fileExtension)
			}
		})
	}
}
