package defaultprops_test

import (
	"errors"
	"testing"

	"github.com/jecepeda/defaultprops"
	"github.com/stretchr/testify/assert"
)

func intPointer(i int) *int {
	return &i
}

func uintPointer(i uint) *uint {
	return &i
}

func stringPointer(s string) *string {
	return &s
}

func float32Pointer(f float32) *float32 {
	return &f
}

func float64Pointer(f float64) *float64 {
	return &f
}

func boolPointer(b bool) *bool {
	return &b
}

func complexPointer(c complex128) *complex128 {
	return &c
}

func chanPointer(c chan bool) *chan bool {
	if c != nil {
		return &c
	}
	var c1 chan bool
	return &c1
}

func slicePointer(sl []string) *[]string {
	return &sl
}

func arrayPointer(arr [2]string) *[2]string {
	return &arr
}

func mapPointer(m map[string]string) *map[string]string {
	return &m
}

var (
	emptyChan  chan bool
	filledChan = make(chan bool, 100)
	fakeFunc   = func() error {
		return nil
	}
	notSoFakeFunc = func() error {
		return errors.New("yeah boi")
	}
)

type testStruct struct {
	Int            int
	String         string
	Float          float64
	Bool           bool
	InnerPtrStruct *substruct
}

type substruct struct {
	PtrInt *int
}

func TestSubstitute(t *testing.T) {
	type args struct {
		orig   interface{}
		dest   interface{}
		config defaultprops.Config
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected interface{}
	}{
		{
			name: "not pointers",
			args: args{
				orig: testStruct{},
				dest: testStruct{},
			},
			wantErr: true,
		},
		{
			name: "different types",
			args: args{
				orig: intPointer(10),
				dest: &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "string set",
			args: args{
				orig: stringPointer("123"),
				dest: stringPointer("1234"),
			},
			wantErr:  false,
			expected: stringPointer("123"),
		},
		{
			name: "string not set",
			args: args{
				orig: stringPointer(""),
				dest: stringPointer("1234"),
			},
			wantErr:  false,
			expected: stringPointer("1234"),
		},
		{
			name: "string skip non zero",
			args: args{
				orig: stringPointer("123"),
				dest: stringPointer("1234"),
				config: defaultprops.Config{
					SkipIfNonZeroValue: true,
				},
			},
			wantErr:  false,
			expected: stringPointer("1234"),
		},
		{
			name: "float32 set",
			args: args{
				orig: float32Pointer(3),
				dest: float32Pointer(4),
			},
			expected: float32Pointer(3),
		},
		{
			name: "float64 set",
			args: args{
				orig: float64Pointer(3),
				dest: float64Pointer(4),
			},
			expected: float64Pointer(3),
		},
		{
			name: "float64 not set",
			args: args{
				orig: float64Pointer(0),
				dest: float64Pointer(4),
			},
			expected: float64Pointer(4),
		},
		{
			name: "uint set",
			args: args{
				orig: uintPointer(12),
				dest: uintPointer(23),
			},
			expected: uintPointer(12),
		},
		{
			name: "uint not set",
			args: args{
				orig: uintPointer(0),
				dest: uintPointer(23),
			},
			expected: uintPointer(23),
		},
		{
			name: "int set",
			args: args{
				orig: intPointer(12),
				dest: intPointer(23),
			},
			expected: intPointer(12),
		},
		{
			name: "int not set",
			args: args{
				orig: intPointer(0),
				dest: intPointer(23),
			},
			expected: intPointer(23),
		},
		{
			name: "complex set",
			args: args{
				orig: complexPointer(128),
				dest: complexPointer(20),
			},
			expected: complexPointer(128),
		},
		{
			name: "complex not set",
			args: args{
				orig: complexPointer(0),
				dest: complexPointer(23),
			},
			expected: complexPointer(23),
		},
		{
			name: "bool set",
			args: args{
				orig: boolPointer(true),
				dest: boolPointer(false),
			},
			expected: boolPointer(true),
		},
		{
			name: "bool with config: set false bools",
			args: args{
				orig: boolPointer(false),
				dest: boolPointer(true),
				config: defaultprops.Config{
					SetFalseBools: true,
				},
			},
			expected: boolPointer(false),
		},
		{
			name: "chan set",
			args: args{
				orig: chanPointer(filledChan),
				dest: chanPointer(emptyChan),
			},
			expected: chanPointer(filledChan),
		},
		{
			name: "chan not set",
			args: args{
				orig: chanPointer(nil),
				dest: chanPointer(filledChan),
			},
			expected: chanPointer(filledChan),
		},
		{
			name: "slices set",
			args: args{
				orig: slicePointer([]string{"1", "2", "3", "4"}),
				dest: slicePointer([]string{}),
			},
			expected: slicePointer([]string{"1", "2", "3", "4"}),
		},
		{
			name: "slices empty",
			args: args{
				orig: slicePointer([]string{}),
				dest: slicePointer([]string{"1", "2", "3", "4"}),
			},
			expected: slicePointer([]string{"1", "2", "3", "4"}),
		},
		{
			name: "array set",
			args: args{
				orig: arrayPointer([2]string{"1", "2"}),
				dest: arrayPointer([2]string{}),
			},
			expected: arrayPointer([2]string{"1", "2"}),
		},
		{
			name: "array empty",
			args: args{
				orig: arrayPointer([2]string{}),
				dest: arrayPointer([2]string{"1", "2"}),
			},
			expected: arrayPointer([2]string{"1", "2"}),
		},
		{
			name: "map merge",
			args: args{
				orig: mapPointer(map[string]string{"1": "2"}),
				dest: mapPointer(map[string]string{"2": "3"}),
			},
			expected: mapPointer(map[string]string{"1": "2", "2": "3"}),
		},
		{
			name: "map empty",
			args: args{
				orig: mapPointer(map[string]string{}),
				dest: mapPointer(map[string]string{"1": "2"}),
			},
			expected: mapPointer(map[string]string{"1": "2"}),
		},
		{
			name: "map with config does not merge",
			args: args{
				orig: mapPointer(map[string]string{"1": "234"}),
				dest: mapPointer(map[string]string{"2": "2"}),
				config: defaultprops.Config{
					ReplaceMaps: true,
				},
			},
			expected: mapPointer(map[string]string{"1": "234"}),
		},
		{
			name: "function set: ignore",
			args: args{
				orig: &notSoFakeFunc,
				dest: &fakeFunc,
			},
			expected: &fakeFunc,
		},
		{
			name: "struct set",
			args: args{
				orig: &testStruct{
					Int:    2,
					String: "2",
					Bool:   true,
					Float:  3.4,
					InnerPtrStruct: &substruct{
						PtrInt: intPointer(20),
					},
				},
				dest: &testStruct{},
			},
			expected: &testStruct{
				Int:    2,
				String: "2",
				Bool:   true,
				Float:  3.4,
				InnerPtrStruct: &substruct{
					PtrInt: intPointer(20),
				},
			},
		},
		{
			name: "struct set",
			args: args{
				orig: &testStruct{
					Int:    2,
					String: "2",
					Bool:   true,
					Float:  3.4,
					InnerPtrStruct: &substruct{
						PtrInt: intPointer(20),
					},
				},
				dest: &testStruct{
					InnerPtrStruct: &substruct{
						PtrInt: intPointer(30),
					},
				},
			},
			expected: &testStruct{
				Int:    2,
				String: "2",
				Bool:   true,
				Float:  3.4,
				InnerPtrStruct: &substruct{
					PtrInt: intPointer(20),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := defaultprops.Substitute(tt.args.orig, tt.args.dest, tt.args.config)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.args.dest)
			}
		})
	}
}

func TestSubstituteNonConfig(t *testing.T) {
	type args struct {
		orig interface{}
		dest interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected interface{}
	}{
		{
			name: "struct with no config",
			args: args{
				orig: &testStruct{
					Int:    2,
					String: "2",
					Bool:   true,
					Float:  3.4,
					InnerPtrStruct: &substruct{
						PtrInt: intPointer(20),
					},
				},
				dest: &testStruct{
					InnerPtrStruct: &substruct{
						PtrInt: intPointer(30),
					},
				},
			},
			expected: &testStruct{
				Int:    2,
				String: "2",
				Bool:   true,
				Float:  3.4,
				InnerPtrStruct: &substruct{
					PtrInt: intPointer(20),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := defaultprops.SubstituteNonConfig(tt.args.orig, tt.args.dest)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.args.dest)
			}
		})
	}
}
