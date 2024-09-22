package num

//type testStruct struct {
//	Int         int
//	Int16       int16
//	IntPtr      *int
//	Int16Ptr    *int16
//	IntSlice    []int
//	IntSlicePtr *[]int
//}
//
//type helperImpl struct{}
//
//func (h *helperImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
//	return shared.Error{
//		Message: fmt.Sprintf(localeKey, value),
//	}
//}
//
//func TestNewIntBundle(t *testing.T) {
//	// Test case 1: Valid input
//	tt := testStruct{
//		Int: 5,
//	}
//	storage, _ := fmap.GetFrom(tt)
//	helper := helperImpl{}
//	deps := shared.BundleDependencies{
//		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
//		Fields:   storage,
//		Object:   struct{}{},
//		Helper:   &helper,
//	}
//	sb := NewNumBundle(deps)
//
//	if sb.storage != deps.Fields {
//		t.Errorf("storage not set correctly")
//	}
//	if sb.obj != deps.Object {
//		t.Errorf("obj not set correctly")
//	}
//	if sb.h != deps.Helper {
//		t.Errorf("h not set correctly")
//	}
//
//	// Test case 2: Nil storage
//	deps.Fields = nil
//	sb = NewNumBundle(deps)
//	if sb.storage != nil {
//		t.Errorf("storage should be nil")
//	}
//
//	// Test case 3: Nil helper
//	deps.Fields = storage
//	deps.Helper = nil
//	sb = NewNumBundle(deps)
//	if sb.h != nil {
//		t.Errorf("helper should be nil")
//	}
//
//	// Test case 4: Invalid AppendFn
//	deps.Helper = &helper
//	deps.AppendFn = func(field fmap.Field, fn shared.FieldValidationFn) { panic("invalid AppendFn") }
//	defer func() {
//		if r := recover(); r == nil {
//			t.Errorf("expected panic")
//		}
//	}()
//	sb = NewNumBundle(deps)
//	sb.appendFn(nil, nil)
//
//}
//
//func TestIntBundle(t *testing.T) {
//	tt := testStruct{
//		Int: 5,
//	}
//	storage, err := fmap.GetFrom(tt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	helper := helperImpl{}
//	deps := shared.BundleDependencies{
//		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
//		Fields:   storage,
//		Object:   &tt,
//		Helper:   &helper,
//	}
//
//	testCases := []struct {
//		name          string
//		field         *int
//		expectedError bool
//	}{
//		{
//			name:          "Valid field",
//			field:         &tt.Int,
//			expectedError: false,
//		},
//		{
//			name:          "Invalid field",
//			field:         (*int)(nil),
//			expectedError: true,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			defer func() {
//				if r := recover(); r != nil {
//					if tc.expectedError {
//						t.Logf("Expected error: %v", r)
//					} else {
//						t.Errorf("Unexpected error: %v", r)
//					}
//				}
//			}()
//
//			sb := NewNumBundle(deps)
//			sbStr := sb.Int(tc.field)
//			if sbStr == nil && !tc.expectedError {
//				t.Errorf("sbStr should not be nil")
//			}
//		})
//	}
//}
//
//func TestIntBundleIntPtr(t *testing.T) {
//	temp := 5
//	tt := testStruct{
//		IntPtr: &temp,
//	}
//	storage, err := fmap.GetFrom(tt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	helper := helperImpl{}
//	deps := shared.BundleDependencies{
//		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
//		Fields:   storage,
//		Object:   &tt,
//		Helper:   &helper,
//	}
//	sb := NewNumBundle(deps)
//	sbStr := sb.IntPtr(&tt.IntPtr)
//	if sbStr == nil {
//		t.Errorf("field not set correctly")
//	}
//
//	// Test case with an invalid field pointer
//	invalidField := (**int)(nil)
//	defer func() {
//		if r := recover(); r != nil {
//			if err, ok := r.(error); !ok {
//				t.Errorf("expected error, got %v", r)
//			} else if err.Error() != "field not found" {
//				t.Errorf("expected error 'field not found', got %v", err)
//			}
//		} else {
//			t.Errorf("expected panic, got none")
//		}
//	}()
//	sb.IntPtr(invalidField)
//}
//
//func TestIntBundleIntSlice(t *testing.T) {
//	tt := testStruct{
//		IntSlice: []int{10, 20},
//	}
//	storage, err := fmap.GetFrom(tt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	helper := helperImpl{}
//	deps := shared.BundleDependencies{
//		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
//		Fields:   storage,
//		Object:   &tt,
//		Helper:   &helper,
//	}
//	sb := NewNumBundle(deps)
//	modelsPtr := &tt.IntSlice
//	sbStr := sb.IntSlice(modelsPtr)
//	if sbStr == nil {
//		t.Errorf("field not set correctly")
//	}
//
//	// Test case with an invalid field pointer
//	invalidField := (*[]int)(nil)
//	defer func() {
//		if r := recover(); r != nil {
//			if err, ok := r.(error); !ok {
//				t.Errorf("expected error, got %v", r)
//			} else if err.Error() != "field not found" {
//				t.Errorf("expected error 'field not found', got %v", err)
//			}
//		} else {
//			t.Errorf("expected panic, got none")
//		}
//	}()
//	sb.IntSlice(invalidField)
//}
//
//func TestIntBundleIntSlicePtr(t *testing.T) {
//	temp := []int{10, 20}
//	tt := testStruct{
//		IntSlicePtr: &temp,
//	}
//	storage, err := fmap.GetFrom(tt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	helper := helperImpl{}
//	deps := shared.BundleDependencies{
//		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
//		Fields:   storage,
//		Object:   &tt,
//		Helper:   &helper,
//	}
//	sb := NewNumBundle(deps)
//	sbStr := sb.IntSlicePtr(&tt.IntSlicePtr)
//	if sbStr == nil {
//		t.Errorf("field not set correctly")
//	}
//
//	// Test case with an invalid field pointer
//	invalidField := (**[]int)(nil)
//	defer func() {
//		if r := recover(); r != nil {
//			if err, ok := r.(error); !ok {
//				t.Errorf("expected error, got %v", r)
//			} else if err.Error() != "field not found" {
//				t.Errorf("expected error 'field not found', got %v", err)
//			}
//		} else {
//			t.Errorf("expected panic, got none")
//		}
//	}()
//	sb.IntSlicePtr(invalidField)
//}
