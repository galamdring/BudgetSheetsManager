package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_dateEqual(t *testing.T) {
	type args struct {
		date1 string
		date2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "dates are equal",
			args: args{
				date1: "1/6/2006",
				date2: "1/6/2006",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date1, _ := time.Parse("1/2/2006", tt.args.date1)
			date2, _ := time.Parse("1/2/2006", tt.args.date2)
			if got := dateEqual(date1, date2); got != tt.want {
				t.Errorf("dateEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dateGT(t *testing.T) {
	type args struct {
		date1 string
		date2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "dates are equal",
			args: args{
				date1: "1/6/2006",
				date2: "1/6/2006",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date1, _ := time.Parse("1/2/2006", tt.args.date1)
			date2, _ := time.Parse("1/2/2006", tt.args.date2)
			if got := dateGT(date1, date2); got != tt.want {
				t.Errorf("dateGT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dateGTE(t *testing.T) {
	type args struct {
		date1 time.Time
		date2 time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dateGTE(tt.args.date1, tt.args.date2); got != tt.want {
				t.Errorf("dateGTE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMinMaxDate(t *testing.T) {
	type args struct {
		data []Transaction
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getMinMaxDate(tt.args.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMinMaxDate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getMinMaxDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_transactionEqual(t *testing.T) {
	type args struct {
		t1 Transaction
		t2 Transaction
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transactionEqual(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("transactionEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTransactionsForDate(t *testing.T) {
	type args struct {
		data []Transaction
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want []Transaction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTransactionsForDate(tt.args.data, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTransactionsForDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewTransactions(t *testing.T) {
	type args struct {
		loaded     []Transaction
		fromSheets []Transaction
	}
	tests := []struct {
		name string
		args args
		want []Transaction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNewTransactions(tt.args.loaded, tt.args.fromSheets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}
