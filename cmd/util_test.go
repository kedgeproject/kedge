package cmd

import (
	"reflect"
	"testing"
)

func Test_removeDuplicateFiles(t *testing.T) {

	tests := []struct {
		name       string
		inputFiles []string
		want       []string
	}{

		{
			name:       "Two files with same relative path",
			inputFiles: []string{"../examples/httpd/httpd.yaml", "../examples/httpd/httpd.yaml"},
			want:       []string{"../examples/httpd/httpd.yaml"},
		},
		{
			name:       "Two files with same absolute path",
			inputFiles: []string{"../examples/httpd/httpd.yaml", "../../kedge/examples/httpd/httpd.yaml", "../examples/gitlab"},
			want:       []string{"../examples/httpd/httpd.yaml", "../examples/gitlab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := removeDuplicateFiles(tt.inputFiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test Failed, got = %v, wanted %v", got, tt.want)
			}
		})
	}
}
