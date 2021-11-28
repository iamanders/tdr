package main

import (
	"testing"
)

func Test_defaultEmptyString(t *testing.T) {
	type args struct {
		s             string
		defaultString string
	}
	tests := []struct {
		name     string
		expected string
		args     args
	}{
		{name: "Default string", args: args{s: "", defaultString: "Default"}, expected: "Default"},
		{name: "Argument string", args: args{s: "Yepp", defaultString: "Default"}, expected: "Yepp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultEmptyString(&tt.args.s, tt.args.defaultString)
			if tt.args.s != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, tt.args.s)
			}
		})
	}
}

func Test_intOrZero(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "Positive number", args: args{s: "123"}, want: 123},
		{name: "Zero string", args: args{s: "0"}, want: 0},
		{name: "Empty string", args: args{s: ""}, want: 0},
		{name: "Alpha string", args: args{s: "Asd"}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intOrZero(tt.args.s); got != tt.want {
				t.Errorf("Expected %v but got %v", got, tt.want)
			}
		})
	}
}
