package fstestutil

import (
	"syscall"
)

// unescape removes the backslash-escaping used by Darwin mount info for the
// handful of characters we care about in tests.
func unescape(s string) string {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			buf = append(buf, s[i])
			continue
		}
		switch s[i+1] {
		case '\\', ' ', '\t', '\n':
			buf = append(buf, s[i+1])
			i++
		default:
			buf = append(buf, s[i])
		}
	}
	return string(buf)
}

func getMountInfo(mnt string) (*MountInfo, error) {
	var st syscall.Statfs_t
	err := syscall.Statfs(mnt, &st)
	if err != nil {
		return nil, err
	}
	i := &MountInfo{
		// osx getmntent(3) fails to un-escape the data, so we do it..
		// this might lead to double-unescaping in the future. fun.
		// TestMountOptionFSNameEvilBackslashDouble checks for that.
		FSName: unescape(cstr(st.Mntfromname[:])),
	}
	return i, nil
}
