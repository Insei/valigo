package zero

import (
	"testing"
)

func TestZero(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 *int
		Field3 map[string]int
		Field4 struct {
			Field5 string
		}
	}

	testCases := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			input: &TestStruct{
				Field1: "hello",
				Field2: nil,
				Field3: nil,
				Field4: struct {
					Field5 string
				}{
					Field5: "world",
				},
			},
			wantErr: false,
		},
		{
			name: "non-addressable input",
			input: TestStruct{
				Field1: "hello",
				Field2: nil,
				Field3: nil,
				Field4: struct {
					Field5 string
				}{
					Field5: "world",
				},
			},
			wantErr: true,
		},
		{
			name:    "non-struct input",
			input:   "hello",
			wantErr: true,
		},
		{
			name: "struct with unaddressable field",
			input: &struct {
				Field1 string
				Field2 func()
			}{
				Field1: "hello",
				Field2: func() {},
			},
			wantErr: false,
		},
		{
			name: "struct with nil field",
			input: &struct {
				Field1 *string
			}{
				Field1: nil,
			},
			wantErr: false,
		},
		{
			name: "struct with empty map field",
			input: &struct {
				Field1 map[string]int
			}{
				Field1: make(map[string]int),
			},
			wantErr: false,
		},
		{
			name: "struct with non-empty map field",
			input: &struct {
				Field1 map[string]int
			}{
				Field1: map[string]int{"key": 1},
			},
			wantErr: false,
		},
		{
			name: "struct with nil map field",
			input: &struct {
				Field1 map[string]int
			}{
				Field1: nil,
			},
			wantErr: false,
		},
		{
			name: "struct with recursive fields",
			input: &struct {
				Field1 *struct {
					Field2 string
				}
			}{
				Field1: &struct {
					Field2 string
				}{
					Field2: "hello",
				},
			},
			wantErr: false,
		},
		{
			name: "struct with multiple levels of nesting",
			input: &struct {
				Field1 string
				Field2 struct {
					Field3 string
					Field4 struct {
						Field5 string
					}
				}
			}{
				Field1: "hello",
				Field2: struct {
					Field3 string
					Field4 struct {
						Field5 string
					}
				}{
					Field3: "world",
					Field4: struct {
						Field5 string
					}{
						Field5: "nested",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "struct with multiple fields of the same type",
			input: &struct {
				Field1 string
				Field2 string
				Field3 string
			}{
				Field1: "hello",
				Field2: "world",
				Field3: "nested",
			},
			wantErr: false,
		},
		{
			name: "struct with a field that is a slice of structs",
			input: &struct {
				Field1 []struct {
					Field2 string
				}
			}{
				Field1: []struct {
					Field2 string
				}{
					{Field2: "hello"},
					{Field2: "world"},
				},
			},
			wantErr: false,
		},
		{
			name: "struct with a field that is a map of structs",
			input: &struct {
				Field1 map[string]struct {
					Field2 string
				}
			}{
				Field1: map[string]struct {
					Field2 string
				}{
					"key1": {Field2: "hello"},
					"key2": {Field2: "world"},
				},
			},
			wantErr: false,
		},
		{
			name: "struct with a field that is a pointer to a struct",
			input: &struct {
				Field1 *struct {
					Field2 string
				}
			}{
				Field1: &struct {
					Field2 string
				}{
					Field2: "hello",
				},
			},
			wantErr: false,
		},
		{
			name: "struct with a field that is a nil pointer to a struct",
			input: &struct {
				Field1 *struct {
					Field2 string
				}
			}{
				Field1: nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Allocate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("Allocate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestMustZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		} else {
			t.Errorf("Must() did not panic")
		}
	}()

	Must("hello")
}

func TestZeroRecursiveTypes(t *testing.T) {
	type node struct {
		Name string
		Next *node
		Meta *struct{ Parent *node }
	}
	obj := &node{}
	if err := Allocate(obj); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}
	if obj.Next != nil || obj.Meta == nil || obj.Meta.Parent != nil {
		t.Errorf("Allocate() must allocate non recursive pointers only, got %+v", obj)
	}
}

func TestZeroErrors(t *testing.T) {
	if err := Allocate(new(int)); err == nil {
		t.Errorf("Allocate() of a non-struct must fail")
	}
}

func TestZeroSkipsUnexportedFields(t *testing.T) {
	type withHidden struct {
		hidden *struct{ X int }
		Shown  *struct{ X int }
	}
	obj := &withHidden{}
	if err := Allocate(obj); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}
	if obj.hidden != nil || obj.Shown == nil {
		t.Errorf("Allocate() must allocate exported fields only, got %+v", obj)
	}
}
