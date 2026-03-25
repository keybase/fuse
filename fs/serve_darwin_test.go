package fs_test

import (
	"os"
	"strings"
	"testing"

	"bazil.org/fuse"
	"bazil.org/fuse/fs/fstestutil"
	"golang.org/x/sys/unix"
)

type exchangeData struct {
	fstestutil.File
	// this struct cannot be zero size or multiple instances may look identical
	_ int
}

func TestExchangeDataNotSupported(t *testing.T) {
	t.Parallel()
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{&fstestutil.ChildMap{
		"one": &exchangeData{},
		"two": &exchangeData{},
	}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer mnt.Close()

	if err := unix.Exchangedata(mnt.Dir+"/one", mnt.Dir+"/two", 0); err != unix.ENOTSUP {
		t.Fatalf("expected ENOTSUP from exchangedata: %v", err)
	}
}

func isFSKitUnavailable(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "Unsupported macOS version") ||
		strings.Contains(msg, "File system extension requires approval") ||
		strings.Contains(msg, "mount(8) returned 69") ||
		strings.Contains(msg, "exit status 251") ||
		strings.Contains(msg, "exit status 69")
}

func TestFSKitMount(t *testing.T) {
	t.Parallel()

	if _, err := os.Stat(fuse.OSXFUSELocationV4.MountFSKit); err != nil {
		t.Skipf("FSKit helper not installed: %v", err)
	}

	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{&fstestutil.ChildMap{
		"child": fstestutil.File{},
	}}, nil, fuse.FSKitBackend())
	if err != nil {
		if isFSKitUnavailable(err) {
			t.Skipf("FSKit backend unavailable on this host: %v", err)
		}
		t.Fatal(err)
	}
	defer mnt.Close()

	if _, err := os.Stat(mnt.Dir + "/child"); err != nil {
		t.Fatalf("stat child through FSKit mount: %v", err)
	}

	if err := fstestutil.CheckDir(mnt.Dir, map[string]fstestutil.FileInfoCheck{
		"child": nil,
	}); err != nil {
		t.Fatalf("FSKit directory contents mismatch: %v", err)
	}
}
