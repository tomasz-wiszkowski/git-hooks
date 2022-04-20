package hooks

import (
	"reflect"
	"testing"
)

func Test_substituteCommandLine(t *testing.T) {
	type args struct {
		inputCmdLine   []string
		substituteArgs map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "No arguments",
			args: args{[]string{}, map[string]interface{}{}},
			want: []string{},
		},
		{
			name: "No substitutions",
			args: args{[]string{"one", "two", "three"}, map[string]interface{}{}},
			want: []string{"one", "two", "three"},
		},
		{
			name: "No active substitutions",
			args: args{[]string{"one", "two", "three"}, map[string]interface{}{"four": "test"}},
			want: []string{"one", "two", "three"},
		},
		{
			name: "Single active substitution with string",
			args: args{[]string{"one", "two", "three"}, map[string]interface{}{"two": "test"}},
			want: []string{"one", "test", "three"},
		},
		{
			name: "Single active substitution with array",
			args: args{[]string{"one", "two", "three"}, map[string]interface{}{"two": []string{"test1", "test2"}}},
			want: []string{"one", "test1", "test2", "three"},
		},
		{
			name: "Repeated active substitutions with string",
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": "test"}},
			want: []string{"one", "test", "three", "test"},
		},
		{
			name: "Repeated active substitutions with array",
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": []string{"test1", "test2"}}},
			want: []string{"one", "test1", "test2", "three", "test1", "test2"},
		},
		{
			name: "Multiple active substitutions with string",
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": "test", "three": "other test"}},
			want: []string{"one", "test", "other test", "test"},
		},
		{
			name: "Multiple active substitutions with array",
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": []string{"test1", "test2"}, "three": []string{"test3", "test4"}}},
			want: []string{"one", "test1", "test2", "test3", "test4", "test1", "test2"},
		},
		{
			name: "String substitutions are nonrecursive",
			// Check that "two" does not become "four" by going "two" -> "three" -> "four"
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": "three", "three": "four"}},
			want: []string{"one", "three", "four", "three"},
		},
		{
			name: "Array substitutions are nonrecursive",
			// Check that "two" does not become "three", "four" by going "two" -> "three" -> "four"
			args: args{[]string{"one", "two", "three", "two"}, map[string]interface{}{"two": []string{"three", "two"}, "three": "four"}},
			want: []string{"one", "three", "two", "four", "three", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substituteCommandLine(tt.args.inputCmdLine, tt.args.substituteArgs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("substituteCommandLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
