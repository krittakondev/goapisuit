package utils

import "testing"

func Test_KebabToCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{
			args: "test-case",
			want: "TestCase",
		},
		{
			args: "Test-Case",
			want: "TestCase",
		},
		{
			args: "TEST-CASE",
			want: "TestCase",
		},
		{
			args: "TESTCASE",
			want: "Testcase",
		},
		{
			args: "testcase",
			want: "Testcase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KebabToCamel(tt.args); got != tt.want {
				t.Errorf("kebabToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
