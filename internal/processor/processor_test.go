package processor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nasjp/tiny-gotests/internal/processor"
	"golang.org/x/tools/imports"
)

func TestRun(t *testing.T) {
	t.Parallel()

	type args struct {
		filePath string
		funcName string
	}

	tests := []struct {
		name           string
		args           args
		goldenFilePath string
	}{
		{"Minimum", args{filePath: "testdata/sample.go", funcName: "Minimum"}, "testdata/minimum_test.go.golden"},
		{"ParamOne", args{filePath: "testdata/sample.go", funcName: "ParamOne"}, "testdata/param_one_test.go.golden"},
		{"ParamTwo", args{filePath: "testdata/sample.go", funcName: "ParamTwo"}, "testdata/param_two_test.go.golden"},
		{"ParamInternalStruct", args{filePath: "testdata/sample.go", funcName: "ParamInternalStruct"}, "testdata/param_internal_struct_test.go.golden"},
		{"ParamInternalPointer", args{filePath: "testdata/sample.go", funcName: "ParamInternalPointer"}, "testdata/param_internal_pointer_test.go.golden"},
		{"ParamExternalStruct", args{filePath: "testdata/sample.go", funcName: "ParamExternalStruct"}, "testdata/param_external_struct_test.go.golden"},
		{"ParamExternalPointer", args{filePath: "testdata/sample.go", funcName: "ParamExternalPointer"}, "testdata/param_external_pointer_test.go.golden"},
		{"ResultOne", args{filePath: "testdata/sample.go", funcName: "ResultOne"}, "testdata/result_one_test.go.golden"},
		{"ResultTwo", args{filePath: "testdata/sample.go", funcName: "ResultTwo"}, "testdata/result_two_test.go.golden"},
		{"ResultErr", args{filePath: "testdata/sample.go", funcName: "ResultErr"}, "testdata/result_err_test.go.golden"},
		{"ResultOneAndErr", args{filePath: "testdata/sample.go", funcName: "ResultOneAndErr"}, "testdata/result_one_and_err_test.go.golden"},
		{"ResultInternalStruct", args{filePath: "testdata/sample.go", funcName: "ResultInternalStruct"}, "testdata/result_internal_struct_test.go.golden"},
		{"ResultInternalPointer", args{filePath: "testdata/sample.go", funcName: "ResultInternalPointer"}, "testdata/result_internal_pointer_test.go.golden"},
		{"ResultExternalStruct", args{filePath: "testdata/sample.go", funcName: "ResultExternalStruct"}, "testdata/result_external_struct_test.go.golden"},
		{"ResultExternalPointer", args{filePath: "testdata/sample.go", funcName: "ResultExternalPointer"}, "testdata/result_external_pointer_test.go.golden"},
		{"ReceiverStruct", args{filePath: "testdata/sample.go", funcName: "ReceiverStruct"}, "testdata/receiver_struct_test.go.golden"},
		{"ReceiverPointer", args{filePath: "testdata/sample.go", funcName: "ReceiverPointer"}, "testdata/receiver_pointer_test.go.golden"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := processor.Run(tt.args.filePath, tt.args.funcName)
			if err != nil {
				t.Errorf("Unexpected Error: %v", err)
				return
			}

			if diff := cmp.Diff(getGolden(t, tt.goldenFilePath), got); diff != "" {
				t.Errorf("processor.Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func getGolden(t *testing.T, filename string) []byte {
	imported, err := imports.Process(filename, nil, nil)
	if err != nil {
		t.Fatalf("reading and formatting file: %v", err)
	}

	return imported
}
