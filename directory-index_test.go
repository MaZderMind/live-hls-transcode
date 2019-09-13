package main

import "testing"

func TestDirectoryIndex_calculatePath(t *testing.T) {
	type fields struct {
		rootDir string
	}
	type args struct {
		urlPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "combines request-url and root-dir",
			fields: fields{"/video"},
			args:   args{"/Documentations/Nature"},
			want:   "/video/Documentations/Nature",
		}, {
			name:   "cleans result-directory",
			fields: fields{"/video/./"},
			args:   args{"/Documentations/../Science Fiction/./StarTrek/"},
			want:   "/video/Science Fiction/StarTrek",
		}, {
			name:   "avoids directory traversal",
			fields: fields{"/video"},
			args:   args{"/Documentations/../../../../../"},
			want:   "/video",
		}, {
			name:   "handles minimal number of slashes",
			fields: fields{"/video"},
			args:   args{"Documentations/Nature"},
			want:   "/video/Documentations/Nature",
		}, {
			name:   "handles maximal number of slashes",
			fields: fields{"/video/"},
			args:   args{"/Documentations//Nature/"},
			want:   "/video/Documentations/Nature",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := DirectoryIndex{
				rootDir: tt.fields.rootDir,
			}
			if got := i.calculatePath(tt.args.urlPath); got != tt.want {
				t.Errorf("calculatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
