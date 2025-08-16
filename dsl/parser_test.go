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
			want: []bool{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			hasModelBases := []bool{}
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
					for _, field := range structType.Fields.List {
						if isModelBase(f, field) {
							hasModelBases = append(hasModelBases, true)
							continue
						}
					}
				}

			}
			if !reflect.DeepEqual(hasModelBases, tt.want) {
				t.Errorf("isModelBase() = %v, want %v", hasModelBases, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code string
		want map[string]struct{}
	}{
		{
			name: "input1",
			code: input1,
			want: map[string]struct{}{"User": {}},
		},
		{
			name: "input2",
			code: input2,
			want: map[string]struct{}{"User2": {}},
		},
		{
			name: "input3",
			code: input3,
			want: map[string]struct{}{"User3": {}, "User4": {}},
		},
		{
			name: "input4",
			code: input4,
			want: map[string]struct{}{},
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
			res := parse(f)
			got := make(map[string]struct{})
			for k := range res {
				got[k] = struct{}{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findAllModelNames(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			got := findAllModelNames(f)
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllModelNames() = %v, want %v", got, tt.want)
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
