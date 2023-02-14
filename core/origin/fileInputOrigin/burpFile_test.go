package fileInputOrigin

import "testing"

func TestFileInputOrigin_parseData(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				path: "C:\\Users\\Lenovo\\Desktop\\src.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//o := &FileInputOrigin{
			//	path: tt.fields.path,
			//}
			//o.parseData()
		})
	}
}
