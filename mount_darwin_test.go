package fuse

import "testing"

func TestFSKitBackendOption(t *testing.T) {
	conf := mountConfig{
		options: make(map[string]string),
	}

	if err := FSKitBackend()(&conf); err != nil {
		t.Fatalf("FSKitBackend returned error: %v", err)
	}

	if got, want := conf.osxfuseBackend, "fskit"; got != want {
		t.Fatalf("wrong backend: got %q want %q", got, want)
	}
	if got, want := conf.options["backend"], "fskit"; got != want {
		t.Fatalf("wrong mount option: got %q want %q", got, want)
	}
}

func TestDefaultOSXFUSELocationsForFSKit(t *testing.T) {
	conf := mountConfig{
		options:        make(map[string]string),
		osxfuseBackend: "fskit",
	}

	locations := defaultOSXFUSELocations(&conf)
	if len(locations) != 1 {
		t.Fatalf("wrong number of default locations: got %d want 1", len(locations))
	}
	if got, want := locations[0], OSXFUSELocationV4; got != want {
		t.Fatalf("wrong default location: got %#v want %#v", got, want)
	}
}

func TestMountBinaryUsesFSKitHelper(t *testing.T) {
	conf := mountConfig{
		options:        make(map[string]string),
		osxfuseBackend: "fskit",
	}

	got := mountBinary(OSXFUSELocationV4, &conf)
	if got != OSXFUSELocationV4.MountFSKit {
		t.Fatalf("wrong mount binary: got %q want %q", got, OSXFUSELocationV4.MountFSKit)
	}
}
