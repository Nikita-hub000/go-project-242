package sizefmt

import "testing"

func TestBytesIEC(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{name: "below kb threshold stays bytes", size: 1023, want: "1023B"},
		{name: "exactly 1.0KB", size: 1024, want: "1.0KB"},
		{name: "fractional kilobytes", size: 1536, want: "1.5KB"},
		{name: "exactly 2.0MB", size: 2 * 1024 * 1024, want: "2.0MB"},
		{name: "exactly 3.0GB", size: 3 * 1024 * 1024 * 1024, want: "3.0GB"},
		{name: "exactly 1.0TB", size: 1024 * 1024 * 1024 * 1024, want: "1.0TB"},
		{name: "exactly 1.0PB", size: 1024 * 1024 * 1024 * 1024 * 1024, want: "1.0PB"},
		{name: "exactly 1.0EB", size: 1024 * 1024 * 1024 * 1024 * 1024 * 1024, want: "1.0EB"},
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
