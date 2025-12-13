package wire

import (
	"reflect"
	"testing"
)

func TestDecodeCStrings(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []string
	}{
		{"single string", []byte("one\u0000"), []string{"one"}},
		{"two strings", []byte("one\u0000two\u0000"), []string{"one", "two"}},
		{"last empty", []byte("one\u0000\u0000"), []string{"one", ""}},
		{"first empty", []byte("\u0000two\u0000"), []string{"", "two"}},
		{"all empty", []byte("\u0000\u0000"), []string{"", ""}},
		{"nil in nil out", nil, nil},
		{"empty ok", []byte{}, nil},
		{"missing last null", []byte("one"), []string{"one"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ltt := tt
			t.Parallel()
			if got := DecodeCStrings(ltt.data); !reflect.DeepEqual(got, ltt.want) {
				t.Errorf("DecodeCStrings() = %v, want %v", got, ltt.want)
			}
		})
	}
}

func TestAppendCString(t *testing.T) {
	type args struct {
		dest []byte
		s    string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"append to nil", args{nil, "append"}, []byte("append\u0000")},
		{"append to empty", args{[]byte{}, "append"}, []byte("append\u0000")},
		{"append", args{[]byte("one\u0000"), "append"}, []byte("one\u0000append\u0000")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ltt := tt
			t.Parallel()
			if got := AppendCString(ltt.args.dest, ltt.args.s); !reflect.DeepEqual(got, ltt.want) {
				t.Errorf("AppendCString() = %v, want %v", got, ltt.want)
			}
		})
	}
}

func TestReadCString(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"simple", []byte("simple\u0000"), "simple"},
		{"trailing", []byte("simple\u0000other data"), "simple"},
		{"no null", []byte("simple"), "simple"},
		{"empty", []byte("\u0000"), ""},
		{"nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ltt := tt
			t.Parallel()
			if got := ReadCString(ltt.data); got != ltt.want {
				t.Errorf("ReadCString() = %v, want %v", got, ltt.want)
			}
		})
	}
}

// Benchmark for AppendCString function
func BenchmarkAppendCString(b *testing.B) {
	dest := make([]byte, 0, 100)
	testString := "test string value"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dest = dest[:0]
		_ = AppendCString(dest, testString)
	}
}

// Benchmark for DecodeCStrings function
func BenchmarkDecodeCStrings(b *testing.B) {
	testData := []byte("string1\x00string2\x00string3\x00string4\x00")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DecodeCStrings(testData)
	}
}

// Benchmark for ReadCString function
func BenchmarkReadCString(b *testing.B) {
	testData := []byte("test string\x00remaining data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadCString(testData)
	}
}
