package dsl

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"sort"
	"testing"

	"github.com/kr/pretty"
)

func Test_isModelBase(t *testing.T) {
	fset := token.NewFileSet()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code string
		want []bool
	}{
		{
			name: "input1",
			code: input1,
			want: []bool{true},
		},
		{
			name: "input2",
			code: input2,
			want: []bool{true},
		},
		{
			name: "input3",
			code: input3,
			want: []bool{true, true},
		},
		{
			name: "input4",
			code: input4,
			want: []bool{false},
		},
		{
			name: "input5",
			code: input5,
			want: []bool{true},
		},
		{
			name: "input6",
			code: input6,
			want: []bool{false, false},
		},
		{
			name: "input7",
			code: input7,
			want: []bool{false, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			modelBases := []bool{}
			for _, decl := range f.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl == nil || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec == nil {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok || structType == nil || structType.Fields == nil {
						continue
					}
					var hasModelBase bool
					for _, field := range structType.Fields.List {
						if isModelBase(f, field) {
							hasModelBase = true
							break
						}
					}
					if hasModelBase {
						modelBases = append(modelBases, true)
					} else {
						modelBases = append(modelBases, false)
					}
				}

			}
			if !reflect.DeepEqual(modelBases, tt.want) {
				t.Errorf("isModelBase() = %v, want %v", modelBases, tt.want)
			}
		})
	}
}

func Test_isModelEmpty(t *testing.T) {
	fset := token.NewFileSet()

	tests := []struct {
		name string // description of this test case
		code string
		want []bool
	}{
		{
			name: "input1",
			code: input1,
			want: []bool{false},
		},
		{
			name: "input2",
			code: input2,
			want: []bool{false},
		},
		{
			name: "input3",
			code: input3,
			want: []bool{false, false},
		},
		{
			name: "input4",
			code: input4,
			want: []bool{false},
		},
		{
			name: "input5",
			code: input5,
			want: []bool{false},
		},
		{
			name: "input6",
			code: input6,
			want: []bool{true, false},
		},
		{
			name: "input7",
			code: input7,
			want: []bool{false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			modelEmptys := []bool{}
			for _, decl := range f.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl == nil || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec == nil {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok || structType == nil || structType.Fields == nil {
						continue
					}
					var hasModelEmpty bool
					for _, field := range structType.Fields.List {
						if isModelEmpty(f, field) {
							hasModelEmpty = true
							break
						}
					}
					if hasModelEmpty {
						modelEmptys = append(modelEmptys, true)
					} else {
						modelEmptys = append(modelEmptys, false)
					}
				}

			}
			if !reflect.DeepEqual(modelEmptys, tt.want) {
				t.Errorf("isModelBase() = %v, want %v", modelEmptys, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code  string
		want1 map[string]struct{}
		want2 map[string]struct{}
	}{
		{
			name:  "input1",
			code:  input1,
			want1: map[string]struct{}{"User": {}},
			want2: map[string]struct{}{},
		},
		{
			name:  "input2",
			code:  input2,
			want1: map[string]struct{}{"User2": {}},
			want2: map[string]struct{}{},
		},
		{
			name:  "input3",
			code:  input3,
			want1: map[string]struct{}{"User3": {}, "User4": {}},
			want2: map[string]struct{}{},
		},
		{
			name:  "input4",
			code:  input4,
			want1: map[string]struct{}{},
			want2: map[string]struct{}{},
		},
		{
			name:  "input5",
			code:  input5,
			want1: map[string]struct{}{"User5": {}},
			want2: map[string]struct{}{},
		},
		{
			name:  "input6",
			code:  input6,
			want1: map[string]struct{}{},
			want2: map[string]struct{}{"User6": {}},
		},
		{
			name:  "input7",
			code:  input7,
			want1: map[string]struct{}{},
			want2: map[string]struct{}{"User8": {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			res1, res2 := parse(f)
			got1 := make(map[string]struct{})
			for k := range res1 {
				got1[k] = struct{}{}
			}
			got2 := make(map[string]struct{})
			for k := range res2 {
				got2[k] = struct{}{}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parse() return 1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parse() return 2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_findAllModelBase(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code string
		want []string
	}{
		{
			name: "input1",
			code: input1,
			want: []string{"User"},
		},
		{
			name: "input2",
			code: input2,
			want: []string{"User2"},
		},
		{
			name: "input3",
			code: input3,
			want: []string{"User3", "User4"},
		},
		{
			name: "input4",
			code: input4,
			want: []string{},
		},
		{
			name: "input5",
			code: input5,
			want: []string{"User5"},
		},
		{
			name: "input6",
			code: input6,
			want: []string{},
		},
		{
			name: "input7",
			code: input7,
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			got := findAllModelBase(f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllModelBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findAllModelEmpty(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code string
		want []string
	}{
		{
			name: "input1",
			code: input1,
			want: []string{},
		},
		{
			name: "input2",
			code: input2,
			want: []string{},
		},
		{
			name: "input3",
			code: input3,
			want: []string{},
		},
		{
			name: "input4",
			code: input4,
			want: []string{},
		},
		{
			name: "input5",
			code: input5,
			want: []string{},
		},
		{
			name: "input6",
			code: input6,
			want: []string{"User6"},
		},
		{
			name: "input7",
			code: input7,
			want: []string{"User8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			got := findAllModelEmpty(f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllModelEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code     string
		endpoint string
		want     map[string]*Design
	}{
		{
			name:     "input1",
			code:     input1,
			endpoint: "",
			want: map[string]*Design{
				"User": {
					Enabled:    true,
					Endpoint:   "user2",
					Migrate:    true,
					Create:     &Action{Enabled: true, Service: false, Payload: "User", Result: "*User"},
					Delete:     &Action{Enabled: true, Service: true, Payload: "*User", Result: "*User"},
					Update:     &Action{Enabled: false, Service: true, Payload: "*User", Result: "User"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User", Result: "*User"},
				},
			},
		},
		{
			name:     "input2",
			code:     input2,
			endpoint: "",
			want: map[string]*Design{
				"User2": {
					Enabled:    false,
					Endpoint:   "user2",
					Migrate:    false,
					Create:     &Action{Enabled: false, Service: true, Payload: "User2", Result: "*User3"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					Update:     &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					Patch:      &Action{Enabled: true, Service: true, Payload: "*User", Result: "User"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User2", Result: "*User2"},
				},
			},
		},
		{
			name:     "input3",
			code:     input3,
			endpoint: "",
			want: map[string]*Design{
				"User3": {
					Enabled:    true,
					Endpoint:   "user",
					Migrate:    false,
					Create:     &Action{Enabled: false, Service: true, Payload: "User", Result: "*User"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					Update:     &Action{Enabled: true, Service: true, Payload: "*User", Result: "User"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User3", Result: "*User3"},
				},
				"User4": {
					Enabled:    true,
					Endpoint:   "user4",
					Migrate:    false,
					Create:     &Action{Enabled: true, Service: true, Payload: "User", Result: "*User"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					Update:     &Action{Enabled: false, Service: true, Payload: "*User", Result: "User"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User4", Result: "*User4"},
				},
			},
		},
		{
			name:     "input4",
			code:     input4,
			endpoint: "",
			want:     map[string]*Design{},
		},
		{
			name:     "input5",
			code:     input5,
			endpoint: "",
			want: map[string]*Design{
				"User5": {
					Enabled:    true,
					Endpoint:   "user5",
					Migrate:    false,
					Create:     &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Update:     &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User5", Result: "*User5"},
				},
			},
		},
		{
			name:     "input6",
			code:     input6,
			endpoint: "",
			want: map[string]*Design{
				"User6": {
					Enabled:    true,
					Endpoint:   "user6",
					Migrate:    false,
					Create:     &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Update:     &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User6", Result: "*User6"},
				},
			},
		},
		{
			name:     "input7",
			code:     input7,
			endpoint: "",
			want: map[string]*Design{
				"User8": {
					Enabled:    true,
					Endpoint:   "user8",
					Migrate:    false,
					Create:     &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Delete:     &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Update:     &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Patch:      &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					List:       &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Get:        &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					CreateMany: &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					DeleteMany: &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					UpdateMany: &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					PatchMany:  &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Import:     &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
					Export:     &Action{Enabled: false, Service: false, Payload: "*User8", Result: "*User8"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			got := Parse(f, tt.endpoint)
			if len(got) != len(tt.want) {
				t.Fatalf("Parse() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
			var keys1 []string
			var keys2 []string
			for k := range got {
				keys1 = append(keys1, k)
			}
			for k := range tt.want {
				keys2 = append(keys2, k)
			}
			sort.Strings(keys1)
			sort.Strings(keys2)
			if !reflect.DeepEqual(keys1, keys2) {
				t.Fatalf("Parse() = %v, want %v", got, tt.want)
			}
			for _, k := range keys1 {
				if !reflect.DeepEqual(got[k], tt.want[k]) {
					t.Fatalf("Parse() = \n%v\nwant \n%v\ndiff: \n%v\n",
						pretty.Sprintf("% #v", got[k]),
						pretty.Sprintf("% #v", tt.want[k]),
						pretty.Diff(got[k], tt.want[k]))
				}
			}
		})
	}
}
