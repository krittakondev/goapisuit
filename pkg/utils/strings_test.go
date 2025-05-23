package utils

import (
	"testing"
)

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
			args: "test-case-some",
			want: "TestCaseSome",
		},
		{
			args: "Test-Case",
			want: "TestCase",
		},
		{
			args: "Test-Case-some",
			want: "TestCaseSome",
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

func TestCamelToKebab(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args: args{s: "testcase"},
			want: "testcase",
		},
		{
			args: args{s: "TestCase"},
			want: "test-case",
		},
		{
			args: args{s: "TestCaseSome"},
			want: "test-case-some",
		},
		{
			args: args{s: "TestCaseSome"},
			want: "test-case-some",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelToKebab(tt.args.s); got != tt.want {
				t.Errorf("CamelToKebab() = %v, want %v", got, tt.want)
			}
		})
	}
}
