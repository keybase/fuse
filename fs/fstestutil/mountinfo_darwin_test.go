package fstestutil

import "testing"

func TestUnescapeDarwinMountInfo(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "FuseTestMarker", want: "FuseTestMarker"},
		{name: "space", in: `FuseTest\ Marker`, want: "FuseTest Marker"},
		{name: "tab", in: "FuseTest\\\tMarker", want: "FuseTest\tMarker"},
		{name: "newline", in: "FuseTest\\\nMarker", want: "FuseTest\nMarker"},
		{name: "backslash", in: `FuseTest\\Marker`, want: `FuseTest\Marker`},
		{name: "double backslash", in: `FuseTest\\\\Marker`, want: `FuseTest\\Marker`},
		{name: "unknown escape preserved", in: `FuseTest\qMarker`, want: `FuseTest\qMarker`},
	}

	for _, tc := range tests {
		if got := unescape(tc.in); got != tc.want {
			t.Fatalf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}
