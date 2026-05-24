package sizefmt

import "testing"

func TestBytesIEC(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{name: "bytes", size: 512, want: "512B"},
		{name: "kilobytes", size: 1536, want: "1.5KB"},
		{name: "megabytes", size: 2 * 1024 * 1024, want: "2.0MB"},
		{name: "gigabytes", size: 3 * 1024 * 1024 * 1024, want: "3.0GB"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := BytesIEC(tt.size)
			if got != tt.want {
				t.Fatalf("BytesIEC() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBytesRaw(t *testing.T) {
	if got := BytesRaw(42); got != "42B" {
		t.Fatalf("BytesRaw() = %q, want %q", got, "42B")
	}
}
