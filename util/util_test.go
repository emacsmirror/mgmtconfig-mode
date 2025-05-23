// Mgmt
// Copyright (C) James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Additional permission under GNU GPL version 3 section 7
//
// If you modify this program, or any covered work, by linking or combining it
// with embedded mcl code and modules (and that the embedded mcl code and
// modules which link with this program, contain a copy of their source code in
// the authoritative form) containing parts covered by the terms of any other
// license, the licensors of this program grant you additional permission to
// convey the resulting work. Furthermore, the licensors of this program grant
// the original author, James Shubin, additional permission to update this
// additional permission if he deems it necessary to achieve the goals of this
// additional permission.

//go:build !root

package util

import (
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"
)

func TestNumToAlpha(t *testing.T) {
	var numToAlphaTests = []struct {
		number int
		result string
	}{
		{0, "a"},
		{25, "z"},
		{26, "aa"},
		{27, "ab"},
		{702, "aaa"},
		{703, "aab"},
		{63269, "cool"},
		{-1, ""},
	}

	for _, test := range numToAlphaTests {
		actual := NumToAlpha(test.number)
		if actual != test.result {
			t.Errorf("input: %d, expected: %s, actual: %s", test.number, test.result, actual)
		}
	}
}

func TestUtilT1(t *testing.T) {

	if Dirname("/foo/bar/baz") != "/foo/bar/" {
		t.Errorf("result is incorrect.")
	}

	if Dirname("/foo/bar/baz/") != "/foo/bar/" {
		t.Errorf("result is incorrect.")
	}

	if Dirname("/foo/") != "/" {
		t.Errorf("result is incorrect.")
	}

	if Dirname("/") != "" { // TODO: should this equal "/" or "" ?
		t.Errorf("result is incorrect.")
	}

	if Dirname("foo/bar.conf") != "foo/" {
		t.Errorf("result is incorrect.")
	}

	if Dirname("foo/bar/baz.conf") != "foo/bar/" {
		t.Errorf("result is incorrect.")
	}

	if Dirname("bar.conf") != "" {
		t.Errorf("result is incorrect.")
	}

	if Basename("/foo/bar/baz") != "baz" {
		t.Errorf("result is incorrect.")
	}

	if Basename("/foo/bar/baz/") != "baz/" {
		t.Errorf("result is incorrect.")
	}

	if Basename("/foo/") != "foo/" {
		t.Errorf("result is incorrect.")
	}

	if Basename("/") != "/" { // TODO: should this equal "" or "/" ?
		t.Errorf("result is incorrect.")
	}

	if Basename("") != "" { // TODO: should this equal something different?
		t.Errorf("result is incorrect.")
	}
}

func TestUtilT2(t *testing.T) {

	// TODO: compare the output with the actual list
	p0 := "/"
	r0 := []string{""} // TODO: is this correct?
	if len(PathSplit(p0)) != len(r0) {
		t.Errorf("result should be: %q.", r0)
		t.Errorf("result should have a length of: %v.", len(r0))
	}

	p1 := "/foo/bar/baz"
	r1 := []string{"", "foo", "bar", "baz"}
	if len(PathSplit(p1)) != len(r1) {
		//t.Errorf("result should be: %q.", r1)
		t.Errorf("result should have a length of: %v.", len(r1))
	}

	p2 := "/foo/bar/baz/"
	r2 := []string{"", "foo", "bar", "baz"}
	if len(PathSplit(p2)) != len(r2) {
		t.Errorf("result should have a length of: %v.", len(r2))
	}
}

func TestUtilT3(t *testing.T) {

	if HasPathPrefix("/foo/bar/baz", "/foo/ba") != false {
		t.Errorf("result should be false.")
	}

	if HasPathPrefix("/foo/bar/baz", "/foo/bar") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz", "/foo/bar/") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz/", "/foo/bar") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz/", "/foo/bar/") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz/", "/foo/bar/baz/dude") != false {
		t.Errorf("result should be false.")
	}

	if HasPathPrefix("/foo/bar/baz/boo/", "/foo/") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz", "/foo/bar/baz") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/bar/baz/", "/foo/bar/baz/") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo", "/foo") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/foo/", "/foo/") != true {
		t.Errorf("result should be true.")
	}

	if HasPathPrefix("/", "/") != true {
		t.Errorf("result should be true.")
	}

}

func TestUtilT4(t *testing.T) {

	if PathPrefixDelta("/foo/bar/baz", "/foo/ba") != -1 {
		t.Errorf("result should be -1.")
	}

	if PathPrefixDelta("/foo/bar/baz", "/foo/bar") != 1 {
		t.Errorf("result should be 1.")
	}

	if PathPrefixDelta("/foo/bar/baz", "/foo/bar/") != 1 {
		t.Errorf("result should be 1.")
	}

	if PathPrefixDelta("/foo/bar/baz/", "/foo/bar") != 1 {
		t.Errorf("result should be 1.")
	}

	if PathPrefixDelta("/foo/bar/baz/", "/foo/bar/") != 1 {
		t.Errorf("result should be 1.")
	}

	if PathPrefixDelta("/foo/bar/baz/", "/foo/bar/baz/dude") != -1 {
		t.Errorf("result should be -1.")
	}

	if PathPrefixDelta("/foo/bar/baz/a/b/c/", "/foo/bar/baz") != 3 {
		t.Errorf("result should be 3.")
	}

	if PathPrefixDelta("/foo/bar/baz/", "/foo/bar/baz") != 0 {
		t.Errorf("result should be 0.")
	}
}

func TestUtilT8(t *testing.T) {

	r0 := []string{"/"}
	if fullList0 := PathSplitFullReversed("/"); !reflect.DeepEqual(r0, fullList0) {
		t.Errorf("expected: %v; got: %v.", r0, fullList0)
	}

	r1 := []string{"/foo/bar/baz/file", "/foo/bar/baz/", "/foo/bar/", "/foo/", "/"}
	if fullList1 := PathSplitFullReversed("/foo/bar/baz/file"); !reflect.DeepEqual(r1, fullList1) {
		t.Errorf("expected: %v; got: %v.", r1, fullList1)
	}

	r2 := []string{"/foo/bar/baz/dir/", "/foo/bar/baz/", "/foo/bar/", "/foo/", "/"}
	if fullList2 := PathSplitFullReversed("/foo/bar/baz/dir/"); !reflect.DeepEqual(r2, fullList2) {
		t.Errorf("expected: %v; got: %v.", r2, fullList2)
	}

}

func TestUtilT9(t *testing.T) {
	fileListIn := []string{ // list taken from drbd-utils package
		"/etc/drbd.conf",
		"/etc/drbd.d/global_common.conf",
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/etc/drbd.d",
		"/usr/share/doc/drbd-utils",
		"/var/lib/drbd",
	}
	sort.Strings(fileListIn)

	fileListOut := []string{ // fixed up manually
		"/etc/drbd.conf",
		"/etc/drbd.d/global_common.conf",
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/etc/drbd.d/",               // added trailing slash
		"/usr/share/doc/drbd-utils/", // added trailing slash
		"/var/lib/drbd",              // can't be fixed :(
	}
	sort.Strings(fileListOut)

	dirify := DirifyFileList(fileListIn, false) // TODO: test with true
	sort.Strings(dirify)
	equals := reflect.DeepEqual(fileListOut, dirify)
	if a, b := len(fileListOut), len(dirify); a != b {
		t.Errorf("counts didn't match: %d != %d", a, b)
	} else if !equals {
		t.Errorf("did not match expected!")
		for i := 0; i < len(dirify); i++ {
			if fileListOut[i] != dirify[i] {
				t.Errorf("# %d: %v <> %v", i, fileListOut[i], dirify[i])
			}
		}
	}
}

func TestUtilT10(t *testing.T) {
	fileListIn := []string{ // fake package list
		"/etc/drbd.conf",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/etc/drbd.d",
		"/etc/drbd.d/foo",
		"/var/lib/drbd",
		"/var/somedir/",
	}
	sort.Strings(fileListIn)

	fileListOut := []string{ // fixed up manually
		"/etc/drbd.conf",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/etc/drbd.d/", // added trailing slash
		"/etc/drbd.d/foo",
		"/var/lib/drbd", // can't be fixed :(
		"/var/somedir/", // stays the same
	}
	sort.Strings(fileListOut)

	dirify := DirifyFileList(fileListIn, false) // TODO: test with true
	sort.Strings(dirify)
	equals := reflect.DeepEqual(fileListOut, dirify)
	if a, b := len(fileListOut), len(dirify); a != b {
		t.Errorf("counts didn't match: %d != %d", a, b)
	} else if !equals {
		t.Errorf("did not match expected!")
		for i := 0; i < len(dirify); i++ {
			if fileListOut[i] != dirify[i] {
				t.Errorf("# %d: %v <> %v", i, fileListOut[i], dirify[i])
			}
		}
	}
}

func TestUtilT11(t *testing.T) {
	in1 := []string{"/", "/usr/", "/usr/lib/", "/usr/share/"} // input
	ex1 := []string{"/usr/lib/", "/usr/share/"}               // expected
	sort.Strings(ex1)
	out1 := RemoveCommonFilePrefixes(in1)
	sort.Strings(out1)
	if !reflect.DeepEqual(ex1, out1) {
		t.Errorf("expected: %v; got: %v.", ex1, out1)
	}

	in2 := []string{"/", "/usr/"}
	ex2 := []string{"/usr/"}
	sort.Strings(ex2)
	out2 := RemoveCommonFilePrefixes(in2)
	sort.Strings(out2)
	if !reflect.DeepEqual(ex2, out2) {
		t.Errorf("expected: %v; got: %v.", ex2, out2)
	}

	in3 := []string{"/"}
	ex3 := []string{"/"}
	out3 := RemoveCommonFilePrefixes(in3)
	if !reflect.DeepEqual(ex3, out3) {
		t.Errorf("expected: %v; got: %v.", ex3, out3)
	}

	in4 := []string{"/usr/bin/foo", "/usr/bin/bar", "/usr/lib/", "/usr/share/"}
	ex4 := []string{"/usr/bin/foo", "/usr/bin/bar", "/usr/lib/", "/usr/share/"}
	sort.Strings(ex4)
	out4 := RemoveCommonFilePrefixes(in4)
	sort.Strings(out4)
	if !reflect.DeepEqual(ex4, out4) {
		t.Errorf("expected: %v; got: %v.", ex4, out4)
	}

	in5 := []string{"/usr/bin/foo", "/usr/bin/bar", "/usr/lib/", "/usr/share/", "/usr/bin"}
	ex5 := []string{"/usr/bin/foo", "/usr/bin/bar", "/usr/lib/", "/usr/share/"}
	sort.Strings(ex5)
	out5 := RemoveCommonFilePrefixes(in5)
	sort.Strings(out5)
	if !reflect.DeepEqual(ex5, out5) {
		t.Errorf("expected: %v; got: %v.", ex5, out5)
	}

	in6 := []string{"/etc/drbd.d/", "/lib/drbd/", "/usr/lib/drbd/", "/usr/lib/systemd/system/", "/usr/lib/tmpfiles.d/", "/usr/sbin/", "/usr/share/doc/drbd-utils/", "/usr/share/man/man5/", "/usr/share/man/man8/", "/usr/share/doc/", "/var/lib/"}
	ex6 := []string{"/etc/drbd.d/", "/lib/drbd/", "/usr/lib/drbd/", "/usr/lib/systemd/system/", "/usr/lib/tmpfiles.d/", "/usr/sbin/", "/usr/share/doc/drbd-utils/", "/usr/share/man/man5/", "/usr/share/man/man8/", "/var/lib/"}
	sort.Strings(ex6)
	out6 := RemoveCommonFilePrefixes(in6)
	sort.Strings(out6)
	if !reflect.DeepEqual(ex6, out6) {
		t.Errorf("expected: %v; got: %v.", ex6, out6)
	}

	in7 := []string{"/etc/", "/lib/", "/usr/lib/", "/usr/lib/systemd/", "/usr/", "/usr/share/doc/", "/usr/share/man/", "/var/"}
	ex7 := []string{"/etc/", "/lib/", "/usr/lib/systemd/", "/usr/share/doc/", "/usr/share/man/", "/var/"}
	sort.Strings(ex7)
	out7 := RemoveCommonFilePrefixes(in7)
	sort.Strings(out7)
	if !reflect.DeepEqual(ex7, out7) {
		t.Errorf("expected: %v; got: %v.", ex7, out7)
	}

	in8 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/global_common.conf",
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/etc/drbd.d/",
		"/usr/share/doc/drbd-utils/",
		"/var/lib/drbd",
	}
	ex8 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/global_common.conf",
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/var/lib/drbd",
	}
	sort.Strings(ex8)
	out8 := RemoveCommonFilePrefixes(in8)
	sort.Strings(out8)
	if !reflect.DeepEqual(ex8, out8) {
		t.Errorf("expected: %v; got: %v.", ex8, out8)
	}

	in9 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/",
		"/lib/drbd/drbd",
		"/lib/drbd/",
		"/lib/drbd/",
		"/lib/drbd/",
		"/usr/lib/drbd/",
		"/usr/lib/drbd/",
		"/usr/lib/drbd/",
		"/usr/lib/drbd/",
		"/usr/lib/drbd/",
		"/usr/lib/systemd/system/",
		"/usr/lib/tmpfiles.d/",
		"/usr/sbin/",
		"/usr/sbin/",
		"/usr/share/doc/drbd-utils/",
		"/usr/share/doc/drbd-utils/",
		"/usr/share/man/man5/",
		"/usr/share/man/man5/",
		"/usr/share/man/man8/",
		"/usr/share/man/man8/",
		"/usr/share/man/man8/",
		"/etc/drbd.d/",
		"/usr/share/doc/drbd-utils/",
		"/var/lib/drbd",
	}
	ex9 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/",
		"/lib/drbd/drbd",
		"/usr/lib/drbd/",
		"/usr/lib/systemd/system/",
		"/usr/lib/tmpfiles.d/",
		"/usr/sbin/",
		"/usr/share/doc/drbd-utils/",
		"/usr/share/man/man5/",
		"/usr/share/man/man8/",
		"/var/lib/drbd",
	}
	sort.Strings(ex9)
	out9 := RemoveCommonFilePrefixes(in9)
	sort.Strings(out9)
	if !reflect.DeepEqual(ex9, out9) {
		t.Errorf("expected: %v; got: %v.", ex9, out9)
	}

	in10 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/",                   // watch me, i'm a dir
		"/etc/drbd.d/global_common.conf", // and watch me i'm a file!
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/", // watch me, i'm a dir too
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/var/lib/drbd",
	}
	ex10 := []string{
		"/etc/drbd.conf",
		"/etc/drbd.d/global_common.conf",
		"/lib/drbd/drbd",
		"/lib/drbd/drbdadm-83",
		"/lib/drbd/drbdadm-84",
		"/lib/drbd/drbdsetup-83",
		"/lib/drbd/drbdsetup-84",
		"/usr/lib/drbd/crm-fence-peer.sh",
		"/usr/lib/drbd/crm-unfence-peer.sh",
		"/usr/lib/drbd/notify-emergency-reboot.sh",
		"/usr/lib/drbd/notify-emergency-shutdown.sh",
		"/usr/lib/drbd/notify-io-error.sh",
		"/usr/lib/drbd/notify-out-of-sync.sh",
		"/usr/lib/drbd/notify-pri-lost-after-sb.sh",
		"/usr/lib/drbd/notify-pri-lost.sh",
		"/usr/lib/drbd/notify-pri-on-incon-degr.sh",
		"/usr/lib/drbd/notify-split-brain.sh",
		"/usr/lib/drbd/notify.sh",
		"/usr/lib/drbd/outdate-peer.sh",
		"/usr/lib/drbd/rhcs_fence",
		"/usr/lib/drbd/snapshot-resync-target-lvm.sh",
		"/usr/lib/drbd/stonith_admin-fence-peer.sh",
		"/usr/lib/drbd/unsnapshot-resync-target-lvm.sh",
		"/usr/lib/systemd/system/drbd.service",
		"/usr/lib/tmpfiles.d/drbd.conf",
		"/usr/sbin/drbd-overview",
		"/usr/sbin/drbdadm",
		"/usr/sbin/drbdmeta",
		"/usr/sbin/drbdsetup",
		"/usr/share/doc/drbd-utils/COPYING",
		"/usr/share/doc/drbd-utils/ChangeLog",
		"/usr/share/doc/drbd-utils/README",
		"/usr/share/doc/drbd-utils/drbd.conf.example",
		"/usr/share/man/man5/drbd.conf-8.3.5.gz",
		"/usr/share/man/man5/drbd.conf-8.4.5.gz",
		"/usr/share/man/man5/drbd.conf-9.0.5.gz",
		"/usr/share/man/man5/drbd.conf.5.gz",
		"/usr/share/man/man8/drbd-8.3.8.gz",
		"/usr/share/man/man8/drbd-8.4.8.gz",
		"/usr/share/man/man8/drbd-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview-9.0.8.gz",
		"/usr/share/man/man8/drbd-overview.8.gz",
		"/usr/share/man/man8/drbd.8.gz",
		"/usr/share/man/man8/drbdadm-8.3.8.gz",
		"/usr/share/man/man8/drbdadm-8.4.8.gz",
		"/usr/share/man/man8/drbdadm-9.0.8.gz",
		"/usr/share/man/man8/drbdadm.8.gz",
		"/usr/share/man/man8/drbddisk-8.3.8.gz",
		"/usr/share/man/man8/drbddisk-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-8.3.8.gz",
		"/usr/share/man/man8/drbdmeta-8.4.8.gz",
		"/usr/share/man/man8/drbdmeta-9.0.8.gz",
		"/usr/share/man/man8/drbdmeta.8.gz",
		"/usr/share/man/man8/drbdsetup-8.3.8.gz",
		"/usr/share/man/man8/drbdsetup-8.4.8.gz",
		"/usr/share/man/man8/drbdsetup-9.0.8.gz",
		"/usr/share/man/man8/drbdsetup.8.gz",
		"/var/lib/drbd",
	}
	sort.Strings(ex10)
	out10 := RemoveCommonFilePrefixes(in10)
	sort.Strings(out10)
	if !reflect.DeepEqual(ex10, out10) {
		t.Errorf("expected: %v; got: %v.", ex10, out10)
		for i := 0; i < len(ex10); i++ {
			if ex10[i] != out10[i] {
				t.Errorf("# %d: %v <> %v", i, ex10[i], out10[i])
			}
		}
	}
}

func TestSegmentedPathSplit(t *testing.T) {
	if ex, out := []string{}, SegmentedPathSplit(
		"",
	); !reflect.DeepEqual(out, ex) {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := []string{"/"}, SegmentedPathSplit(
		"/",
	); !reflect.DeepEqual(out, ex) {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := []string{"/", "foo/", "bar/"}, SegmentedPathSplit(
		"/foo/bar/",
	); !reflect.DeepEqual(out, ex) {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := []string{"/", "foo/", "bar"}, SegmentedPathSplit(
		"/foo/bar",
	); !reflect.DeepEqual(out, ex) {
		t.Errorf("expected: %v got: %v", ex, out)
	}
}

func TestCommonPathPrefix1(t *testing.T) {
	if ex, out := "/foo/whatever2/", CommonPathPrefix(
		"/foo/whatever2/",
		"/foo/whatever2/",
		"/foo/whatever2/",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
}

func TestCommonPathPrefix2(t *testing.T) {
	if ex, out := "/whatever1", CommonPathPrefix(
		"/whatever1",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
	if ex, out := "/whatever2", CommonPathPrefix(
		"/whatever2",
		"/whatever2",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
	if ex, out := "/foo/whatever1", CommonPathPrefix(
		"/foo/whatever1",
		"/foo/whatever1",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
	if ex, out := "/foo/whatever2", CommonPathPrefix(
		"/foo/whatever2",
		"/foo/whatever2",
		"/foo/whatever2",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := "/whatever3/", CommonPathPrefix(
		"/whatever3/",
		"/whatever3/",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
	if ex, out := "/foo/whatever3/", CommonPathPrefix(
		"/foo/whatever3/",
		"/foo/whatever3/",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}
	if ex, out := "/foo/whatever4/", CommonPathPrefix(
		"/foo/whatever4/",
		"/foo/whatever4/",
		"/foo/whatever4/",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := "/", CommonPathPrefix(
		"/foo/bar",
		"/bar/baz/",
		"/baz/bing/wow",
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	if ex, out := "/foo/", CommonPathPrefix(
		"/foo/bar/",
		"/foo/bar/dude",
		"/foo/bar", // this is not the same as /foo/bar/ !
	); out != ex {
		t.Errorf("expected: %v got: %v", ex, out)
	}

	// If we want to "safe clean" each path, then this test should be added.
	//if ex, out := "/home/james/tmp/", CommonPathPrefix(
	//	"/home/james/tmp/coverage/test",
	//	"/home/james/tmp/covert/operator",
	//	"/home/james/tmp/coven/members",
	//	"/home//james/tmp/coventry",
	//	"/home/james/././tmp/covertly/foo",
	//	"/home/luser/../james/tmp/coved/bar",
	//); out != ex {
	//	t.Errorf("expected: %v got: %v", ex, out)
	//}
}

func TestUtilFlattenListWithSplit1(t *testing.T) {
	{
		in := []string{} // input
		ex := []string{} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}

	{
		in := []string{"hey"} // input
		ex := []string{"hey"} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}

	{
		in := []string{"a", "b", "c", "d"} // input
		ex := []string{"a", "b", "c", "d"} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}

	{
		in := []string{"a,b,c,d"}          // input
		ex := []string{"a", "b", "c", "d"} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}

	{
		in := []string{"a,b;c d"}          // input (mixed)
		ex := []string{"a", "b", "c", "d"} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}

	{
		in := []string{"a,b,c,d;e,f,g,h;i,j,k,l;m,n,o,p q,r,s,t;u,v,w,x y z"}                                                                            // input (mixed)
		ex := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"} // expected
		out := FlattenListWithSplit(in, []string{",", ";", " "})
		sort.Strings(out)
		sort.Strings(ex)
		if !reflect.DeepEqual(ex, out) {
			t.Errorf("expected: %v; got: %v.", ex, out)
		}
	}
}

func TestRemoveBasePath0(t *testing.T) {
	// expected successes...
	if s, err := RemoveBasePath("/usr/bin/foo", "/usr/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "bin/foo" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := RemoveBasePath("/usr/bin/project/", "/usr/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := RemoveBasePath("/", "/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "" { // TODO: is this correct?
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := RemoveBasePath("/usr/bin/project/", "/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "usr/bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := RemoveBasePath("/usr/bin/project/", "/usr/bin/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := RemoveBasePath("/usr/bin/foo", "/usr/bin/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "foo" {
		t.Errorf("unexpected string, got: %s", s)
	}
	// allow this one, even though it's relative paths
	if s, err := RemoveBasePath("usr/bin/project/", "usr/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}

	// expected errors...
	if s, err := RemoveBasePath("", ""); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := RemoveBasePath("", "/usr/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := RemoveBasePath("usr/bin/project/", ""); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := RemoveBasePath("usr/bin/project/", "/usr/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := RemoveBasePath("/usr/bin/project/", "usr/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	// allow this one, even though it's relative paths
	//if s, err := RemoveBasePath("usr/bin/project/", "usr/"); err == nil {
	//	t.Errorf("expected error, got: %s", s)
	//}
	if s, err := RemoveBasePath("/usr/bin/project/", "/bin/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
}

func TestRebasePath0(t *testing.T) {
	// expected successes...
	if s, err := Rebase("/usr/bin/foo", "/usr/", "/usr/local/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/usr/local/bin/foo" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "/usr/", "/usr/local/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/usr/local/bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/", "/", "/opt/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/opt/" { // TODO: is this correct?
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "/", "/opt/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/opt/usr/bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "/usr/bin/", "/opt/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/opt/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/foo", "/usr/bin/", "/opt/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/opt/foo" {
		t.Errorf("unexpected string, got: %s", s)
	}
	// allow this one, even though it's relative paths
	if s, err := Rebase("usr/bin/project/", "usr/", "/opt/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "/opt/bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	// empty root to build a relative dir path
	if s, err := Rebase("/var/lib/dir/file.conf", "/var/lib/", ""); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "dir/file.conf" {
		t.Errorf("unexpected string, got: %s", s)
	}

	// expected errors...
	if s, err := Rebase("", "", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := Rebase("", "/usr/", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := Rebase("usr/bin/project/", "", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := Rebase("usr/bin/project/", "/usr/", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "usr/", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}
	// allow this one, even though it's relative paths
	//if s, err := Rebase("usr/bin/project/", "usr/", "/opt/"); err == nil {
	//	t.Errorf("expected error, got: %s", s)
	//}
	if s, err := Rebase("/usr/bin/project/", "/bin/", "/opt/"); err == nil {
		t.Errorf("expected error, got: %s", s)
	}

	// formerly a failure:
	//if s, err := Rebase("/usr/bin/project", "/usr/", ""); err == nil {
	//	t.Errorf("expected error, got: %s", s)
	//}
	// replaced with a valid result instead:
	if s, err := Rebase("/usr/bin/project", "/usr/", ""); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "bin/project" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "/usr/", ""); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
	if s, err := Rebase("/usr/bin/project/", "/usr/", "foo/bar/"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if s != "foo/bar/bin/project/" {
		t.Errorf("unexpected string, got: %s", s)
	}
}

func TestRemovePathPrefix0(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{
			in:  "/simple1",
			out: "/",
		},
		{
			in:  "/simple1/foo/bar",
			out: "/foo/bar",
		},
		{
			in:  "/simple1/foo/bar/",
			out: "/foo/bar/",
		},
	}
	for _, test := range testCases {
		out, err := RemovePathPrefix(test.in)
		if err != nil {
			t.Errorf("error: %+v", err)
			continue
		}
		if test.out != out {
			t.Errorf("failed: %s -> %s", test.in, out)
			continue
		}
	}
}

func TestRemovePathSuffix0(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{
			in:  "/simple1/",
			out: "/",
		},
		{
			in:  "/simple1/foo/bar/",
			out: "/simple1/foo/",
		},
		{
			in:  "/simple1/foo/",
			out: "/simple1/",
		},
		// TODO: are these what we want?
		{
			in:  "/simple1/foo",
			out: "/simple1/",
		},
		{
			in:  "/simple1",
			out: "/",
		},
	}
	for _, test := range testCases {
		out, err := RemovePathSuffix(test.in)
		if err != nil {
			t.Errorf("error: %+v", err)
			continue
		}
		if test.out != out {
			t.Errorf("failed: %s -> %s (exp: %s)", test.in, out, test.out)
			continue
		}
	}
}

func TestDirParents0(t *testing.T) {
	tests := []struct {
		in  string
		out []string
	}{
		{
			in:  "",
			out: nil,
		},
		{
			in:  "/",
			out: []string{},
		},
		{
			in: "/tmp/x1/mod1/files/",
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				"/tmp/x1/mod1/",
			},
		},
		{
			in: "/tmp/x1/mod1/files/foo",
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				"/tmp/x1/mod1/",
				"/tmp/x1/mod1/files/",
			},
		},
	}
	for index, tt := range tests {
		result := DirParents(tt.in)
		if a, b := len(result), len(tt.out); a != b {
			t.Errorf("test #%d: expected length differs (%d != %d)", index, a, b)
			t.Errorf("test #%d: actual: %+v", index, result)
			t.Errorf("test #%d: expected: %+v", index, tt.out)
			break
		}
		for i := range result {
			if result[i] != tt.out[i] {
				t.Errorf("test #%d: parents diff: wanted: %s got: %s", index, tt.out[i], result[i])
			}
		}
	}
}

func TestMissingMkdirs0(t *testing.T) {
	tests := []struct {
		in   []string
		out  []string
		fail bool
	}{
		{
			in:  []string{},
			out: []string{},
		},
		{
			in: []string{
				"/",
			},
			out: []string{},
		},
		{
			in: []string{
				"/tmp/x1/metadata.yaml",
				"/tmp/x1/main.mcl",
				"/tmp/x1/files/",
				"/tmp/x1/second.mcl",
				"/tmp/x1/mod1/metadata.yaml",
				"/tmp/x1/mod1/main.mcl",
				"/tmp/x1/mod1/files/",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				"/tmp/x1/mod1/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/main.mcl",
				"/tmp/x1/metadata.yaml",
				"/tmp/x1/mod1/files/",
				"/tmp/x1/mod1/main.mcl",
				"/tmp/x1/mod1/metadata.yaml",
				"/tmp/x1/second.mcl",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				"/tmp/x1/mod1/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/files/a/b/c/",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				//"/tmp/x1/files/", // already exists!
				"/tmp/x1/files/a/",
				"/tmp/x1/files/a/b/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/files/a/b/c/",
				"/tmp/x1/files/a/b/c/",
				"/tmp/x1/files/a/b/c/",
				"/tmp/x1/files/a/b/c/", // duplicates
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				//"/tmp/x1/files/", // already exists!
				"/tmp/x1/files/a/",
				"/tmp/x1/files/a/b/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d2",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				//"/tmp/x1/files/", // already exists!
				"/tmp/x1/files/a/",
				"/tmp/x1/files/a/b/",
				"/tmp/x1/files/a/b/c/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1", // duplicates!
				"/tmp/x1/files/a/b/c/d2",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				//"/tmp/x1/files/", // already exists!
				"/tmp/x1/files/a/",
				"/tmp/x1/files/a/b/",
				"/tmp/x1/files/a/b/c/",
			},
		},
		{
			in: []string{
				"/tmp/x1/files/",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1",
				"/tmp/x1/files/a/b/c/d1", // duplicates!
				"/tmp/x1/files/a/b/",
				"/tmp/x1/files/a/b/c/d2",
			},
			out: []string{
				"/",
				"/tmp/",
				"/tmp/x1/",
				//"/tmp/x1/files/", // already exists!
				"/tmp/x1/files/a/",
				"/tmp/x1/files/a/b/c/",
			},
		},
		// invalid path list, so undefined results
		//{
		//	in: []string{
		//		"/tmp/x1/files/",
		//		"/tmp/x1/files/a/b/c/d",
		//		"/tmp/x1/files/a/b/c/d/", // error: same name as file
		//		"/tmp/x1/files/a/b/c/d1",
		//	},
		//	out: []string{},
		//	fail: true, // TODO: put a specific error?
		//},
		// TODO: add more tests
	}
	for index, tt := range tests {
		result, err := MissingMkdirs(tt.in)

		if !tt.fail && err != nil {
			t.Errorf("test #%d: failed with: %+v", index, err)
			break
		}
		if tt.fail && err == nil {
			t.Errorf("test #%d: passed, expected fail", index)
			break
		}
		if !tt.fail && result == nil {
			t.Errorf("test #%d: output was nil", index)
			break
		}

		if a, b := len(result), len(tt.out); a != b {
			t.Errorf("test #%d: expected length differs (%d != %d)", index, a, b)
			t.Errorf("test #%d: actual: %+v", index, result)
			t.Errorf("test #%d: expected: %+v", index, tt.out)
			break
		}
		for i := range result {
			if result[i] != tt.out[i] {
				t.Errorf("test #%d: missing mkdirs diff: wanted: %s got: %s", index, tt.out[i], result[i])
			}
		}
	}
}

func TestPriorityStrSliceSort0(t *testing.T) {
	in := []string{"foo", "bar", "baz"}
	ex := []string{"bar", "baz", "foo"}

	fn := func(x string) bool {
		return x == "foo"
	}
	out := PriorityStrSliceSort(in, fn)

	if !reflect.DeepEqual(ex, out) {
		t.Errorf("expected: %v; got: %v.", ex, out)
	}
}

func TestPriorityStrSliceSort1(t *testing.T) {
	in := []string{"foo", "bar", "baz"}
	ex := []string{"bar", "foo", "baz"}

	fn := func(x string) bool {
		return x != "bar" // != brings this key to the front
	}
	out := PriorityStrSliceSort(in, fn)

	if !reflect.DeepEqual(ex, out) {
		t.Errorf("expected: %v; got: %v.", ex, out)
	}
}

func TestPriorityStrSliceSort2(t *testing.T) {
	in := []string{"bar", "foo", "bar", "bar", "baz"}
	ex := []string{"foo", "baz", "bar", "bar", "bar"}

	fn := func(x string) bool {
		return x == "bar"
	}
	out := PriorityStrSliceSort(in, fn)

	if !reflect.DeepEqual(ex, out) {
		t.Errorf("expected: %v; got: %v.", ex, out)
	}
}

func TestPriorityStrSliceSort3(t *testing.T) {
	in := []string{"foo", "bar1", "bar2", "bar3", "baz"}
	ex := []string{"bar1", "bar2", "bar3", "foo", "baz"}

	fn := func(x string) bool {
		return !strings.HasPrefix(x, "bar")
	}
	out := PriorityStrSliceSort(in, fn)

	if !reflect.DeepEqual(ex, out) {
		t.Errorf("expected: %v; got: %v.", ex, out)
	}
}

func TestPriorityStrSliceSort4(t *testing.T) {
	in := []string{"foo", "bar1", "bar2", "bar3", "baz"}
	ex := []string{"foo", "baz", "bar1", "bar2", "bar3"}

	fn := func(x string) bool {
		return strings.HasPrefix(x, "bar")
	}
	out := PriorityStrSliceSort(in, fn)

	if !reflect.DeepEqual(ex, out) {
		t.Errorf("expected: %v; got: %v.", ex, out)
	}
}

func TestSortedStrSliceCompare0(t *testing.T) {
	slice0 := []string{"foo", "bar", "baz"}
	slice1 := []string{"bar", "foo", "baz"}

	if err := SortedStrSliceCompare(slice0, slice1); err != nil {
		t.Errorf("slices were not evaluated as equivalent: %v, %v", slice0, slice1)
	}
}

func TestSortedStrSliceCompare1(t *testing.T) {
	slice0 := []string{"foo", "bar", "baz"}
	slice1 := []string{"fi", "fi", "fo"}

	if err := SortedStrSliceCompare(slice0, slice1); err == nil {
		t.Errorf("slices were evaluated as equivalent: %v, %v", slice0, slice1)
	}
}

func TestSortedStrSliceCompare2(t *testing.T) {
	slice0 := []string{"foo", "bar", "baz"}
	slice1 := []string{"foo", "bar"}

	if err := SortedStrSliceCompare(slice0, slice1); err == nil {
		t.Errorf("slices were evaluated as equivalent: %v, %v", slice0, slice1)
	}
}

func TestSortedStrSliceCompare3(t *testing.T) {
	slice0 := []string{"foo", "bar", "baz"}
	slice1 := []string{"zip", "zap", "zop"}

	_ = SortedStrSliceCompare(slice0, slice1)

	if slice0[0] != "foo" || slice0[1] != "bar" || slice0[2] != "baz" {
		t.Errorf("input slice reordered to: %v", slice0)
	}

	if slice1[0] != "zip" || slice1[1] != "zap" || slice1[2] != "zop" {
		t.Errorf("input slice reordered to: %v", slice1)
	}
}

func TestPathSliceSort(t *testing.T) {
	tests := []struct {
		in  []string
		out []string
	}{
		{
			in: []string{
				"/foo/bar/baz",
				"/bing/bang/boom",
				"/1/2/3/",
				"/foo/bar/raz",
				"/bing/buzz/",
				"/foo/",
				"/",
				"/1/",
				"/foo/bar/baz/bam",
				"/bing/bang/",
				"/1/2/",
				"/foo/bar/",
				"/bing/",
			},
			out: []string{
				"/",
				"/1/",
				"/1/2/",
				"/1/2/3/",
				"/bing/",
				"/bing/bang/",
				"/bing/bang/boom",
				"/bing/buzz/",
				"/foo/",
				"/foo/bar/",
				"/foo/bar/baz",
				"/foo/bar/baz/bam",
				"/foo/bar/raz",
			},
		},
	}
	for _, tt := range tests {
		sort.Sort(PathSlice(tt.in))
		for i := range tt.in {
			if tt.in[i] != tt.out[i] {
				t.Errorf("path sort failed: wanted: %s got: %s", tt.out[i], tt.in[i])
			}
		}
	}
}

func TestSortUInt64Slice(t *testing.T) {
	slice0 := []uint64{42, 13, 0}
	sort.Sort(UInt64Slice(slice0))
	if slice0[0] != 0 || slice0[1] != 13 || slice0[2] != 42 {
		t.Errorf("input slice reordered to: %v", slice0)
	}

	slice1 := []uint64{99, 12, 13}
	sort.Sort(UInt64Slice(slice1))
	if slice1[0] != 12 || slice1[1] != 13 || slice1[2] != 99 {
		t.Errorf("input slice reordered to: %v", slice1)
	}
}

func TestSortMapStringValuesByUInt64Keys(t *testing.T) {
	if x := len(SortMapStringValuesByUInt64Keys(nil)); x != 0 {
		t.Errorf("input map of nil caused a: %d", x)
	}

	map0 := map[uint64]string{
		42: "world",
		34: "there",
		13: "hello",
	}
	slice0 := SortMapStringValuesByUInt64Keys(map0)
	if slice0[0] != "hello" || slice0[1] != "there" || slice0[2] != "world" {
		t.Errorf("input slice reordered to: %v", slice0)
	}

	map1 := map[uint64]string{
		99: "a",
		12: "c",
		13: "b",
	}
	slice1 := SortMapStringValuesByUInt64Keys(map1)
	if slice1[0] != "c" || slice1[1] != "b" || slice1[2] != "a" {
		t.Errorf("input slice reordered to: %v", slice1)
	}

	map2 := map[uint64]string{
		12:    "c",
		0:     "d",
		44442: "b",
	}
	slice2 := SortMapStringValuesByUInt64Keys(map2)
	if slice2[0] != "d" || slice2[1] != "c" || slice2[2] != "b" {
		t.Errorf("input slice reordered to: %v", slice2)
	}
}

func TestFirstToUpper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "lowercase word",
			input: "small",
			want:  "Small",
		},
		{
			name:  "capitalized word",
			input: "CAPITAL",
			want:  "CAPITAL",
		},
		{
			name:  "capitalized first letter",
			input: "First",
			want:  "First",
		},
		{
			name:  "lowercase first letter",
			input: "fIRST",
			want:  "FIRST",
		},
		{
			name:  "number",
			input: "0number",
			want:  "0number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstToUpper(tt.input)

			if got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestUint64KeyFromStrInMap(t *testing.T) {
	type input struct {
		needle   string
		haystack map[uint64]string
	}
	type want struct {
		key   uint64
		exist bool
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: `needle "n" in empty haystack`,
			input: input{
				needle:   "n",
				haystack: make(map[uint64]string),
			},
			want: want{
				key:   0,
				exist: false,
			},
		},
		{
			name: `needle "n" in haystack doesn't contain "n"`,
			input: input{
				needle:   "n",
				haystack: map[uint64]string{0: "a", 1: "b"},
			},
			want: want{
				key:   0,
				exist: false,
			},
		},
		{
			name: `needle "n" in haystack contain "n"`,
			input: input{
				needle:   "n",
				haystack: map[uint64]string{0: "a", 1: "b", 2: "n"},
			},
			want: want{
				key:   2,
				exist: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotExist := Uint64KeyFromStrInMap(tt.input.needle, tt.input.haystack)

			if gotKey != tt.want.key {
				t.Errorf("got key: %d, want key: %d", gotKey, tt.want.key)
			}

			if gotExist != tt.want.exist {
				t.Errorf("got exist: %t, want exist: %t", gotExist, tt.want.exist)
			}
		})
	}
}

func TestStrFilterElementsInList(t *testing.T) {
	type input struct {
		filter []string
		list   []string
	}

	tests := []struct {
		name  string
		input input
		want  []string
	}{
		{
			name: "empty filter",
			input: input{
				filter: []string{},
				list:   []string{"first", "second"},
			},
			want: []string{"first", "second"},
		},
		{
			name: "empty list",
			input: input{
				filter: []string{"filter"},
				list:   []string{},
			},
			want: []string{},
		},
		{
			name: "nil",
			input: input{
				filter: nil,
				list:   nil,
			},
			want: []string{},
		},
		{
			name: "filter",
			input: input{
				filter: []string{"filter"},
				list:   []string{"first", "second", "filter"},
			},
			want: []string{"first", "second"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrFilterElementsInList(tt.input.filter, tt.input.list)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestStrListIntersection(t *testing.T) {
	type input struct {
		list1 []string
		list2 []string
	}
	tests := []struct {
		name  string
		input input
		want  []string
	}{
		{
			name: "nil",
			input: input{
				list1: nil,
				list2: nil,
			},
			want: []string{},
		},
		{
			name: "no intersection elements",
			input: input{
				list1: []string{"one", "two"},
				list2: []string{"three", "four"},
			},
			want: []string{},
		},
		{
			name: "contains intersection element",
			input: input{
				list1: []string{"one", "two"},
				list2: []string{"two", "three"},
			},
			want: []string{"two"},
		},
		{
			name: "all intersection elements",
			input: input{
				list1: []string{"one", "two"},
				list2: []string{"one", "two"},
			},
			want: []string{"one", "two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrListIntersection(tt.input.list1, tt.input.list2)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestStrMapKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  []string
	}{
		{
			name:  "nil",
			input: nil,
			want:  []string{},
		},
		{
			name:  "empty map",
			input: map[string]string{},
			want:  []string{},
		},
		{
			name:  "returns sorted keys",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  []string{"key1", "key2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrMapKeys(tt.input)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestStrMapKeysUint64(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]uint64
		want  []string
	}{
		{
			name:  "nil",
			input: nil,
			want:  []string{},
		},
		{
			name:  "empty map",
			input: map[string]uint64{},
			want:  []string{},
		},
		{
			name:  "returns sorted keys",
			input: map[string]uint64{"key1": 1, "key2": 2},
			want:  []string{"key1", "key2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrMapKeysUint64(tt.input)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestBoolMapValues(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]bool
		want  []bool
	}{
		{
			name:  "nil",
			input: nil,
			want:  []bool{},
		},
		{
			name:  "empty map",
			input: map[string]bool{},
			want:  []bool{},
		},
		{
			name:  "return values unordered",
			input: map[string]bool{"key1": true},
			want:  []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolMapValues(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestStrMapValues(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  []string
	}{
		{
			name:  "nil",
			input: nil,
			want:  []string{},
		},
		{
			name:  "empty map",
			input: map[string]string{},
			want:  []string{},
		},
		{
			name:  "return values",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  []string{"value1", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrMapValues(tt.input)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestStrMapValuesUint64(t *testing.T) {
	tests := []struct {
		name  string
		input map[uint64]string
		want  []string
	}{
		{
			name:  "nil",
			input: nil,
			want:  []string{},
		},
		{
			name:  "empty map",
			input: map[uint64]string{},
			want:  []string{},
		},
		{
			name:  "return values",
			input: map[uint64]string{1: "value1", 2: "value2"},
			want:  []string{"value1", "value2"},
		},
	}

	for _, tt := range tests {
		got := StrMapValuesUint64(tt.input)

		if !slices.Equal(got, tt.want) {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}

func TestBoolMapTrue(t *testing.T) {
	tests := []struct {
		name  string
		input []bool
		want  bool
	}{
		{
			name:  "nil",
			input: nil,
			want:  true,
		},
		{
			name:  "empty slice",
			input: []bool{},
			want:  true,
		},
		{
			name:  "all true return true",
			input: []bool{true, true, true},
			want:  true,
		},
		{
			name:  "contain false return false",
			input: []bool{true, false, true},
			want:  false,
		},
		{
			name:  "all false return false",
			input: []bool{false, false, false},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolMapTrue(tt.input)

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestSafePathClean(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  ".",
		},
		{
			name:  "slash",
			input: "/",
			want:  "/",
		},
		{
			name:  "end with slash",
			input: "a//b/",
			want:  "a/b/",
		},
		{
			name:  "end without slash",
			input: "a//b",
			want:  "a/b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafePathClean(tt.input)

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestCommonPathPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "common path",
			input: []string{"/common/uncommon1", "/common/uncommon2"},
			want:  "/common/",
		},
		{
			name:  "empty",
			input: []string{},
			want:  "",
		},
		{
			name:  "single path",
			input: []string{"/path/to"},
			want:  "/path/to",
		},
		// XXX: currently undefined behaviour
		//{
		//	name:  "single path doesn't start with /",
		//	input: []string{"path/to"},
		//	want:  "path/to",
		//},
		{
			name:  "one of the paths doesn't contain /",
			input: []string{"/path/with/slash", "path/without/slash"},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommonPathPrefix(tt.input...)
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestFlattenListWithSplit(t *testing.T) {
	type input struct {
		input []string
		split []string
	}
	tests := []struct {
		name  string
		input input
		want  []string
	}{
		{
			name: "split by spaces and dots",
			input: input{
				input: []string{"a b.c"},
				split: []string{" ", "."},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty split",
			input: input{
				input: []string{"a b.c"},
				split: []string{},
			},
			want: []string{"a b.c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenListWithSplit(tt.input.input, tt.input.split)

			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestRebase(t *testing.T) {
	type input struct {
		path string
		base string
		root string
	}
	type want struct {
		rebasedPath string
		err         error
	}

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "rebased to absolute directory",
			input: input{
				path: "/usr/bin/foo",
				base: "/usr/",
				root: "/usr/local/",
			},
			want: want{
				rebasedPath: "/usr/local/bin/foo",
				err:         nil,
			},
		},
		{
			name: "rebased to relative directory",
			input: input{
				path: "/var/lib/dir/file.conf",
				base: "/var/lib/",
				root: "",
			},
			want: want{
				rebasedPath: "dir/file.conf",
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		gotStr, gotErr := Rebase(tt.input.path, tt.input.base, tt.input.root)

		if gotStr != tt.want.rebasedPath {
			t.Errorf("got rebased path: %v, want rebased path: %v", gotStr, tt.want.rebasedPath)
		}
		if gotErr != tt.want.err {
			t.Errorf("got error: %v, want error to be: %v", gotErr, tt.want.err)
		}
	}

	t.Run("root doesn't end with /", func(t *testing.T) {
		// XXX: in Rebase function return a predefined error
		// e.g. var ErrRootNotDirectory = errors.New("root is not a directory")
		// so it would be clearer and easier to test
		gotStr, gotErr := Rebase("/usr/bin/foo", "/user/", "/usr/local")

		if gotStr != "" {
			t.Errorf("rebased path should be empty")
		}

		if gotErr.Error() != "root is not a directory" {
			t.Errorf(`should receive error: "root is not a directory"`)
		}
	})
}

func TestRemovePathPrefix(t *testing.T) {
	t.Run("removes path prefix", func(t *testing.T) {
		gotStr, gotErr := RemovePathPrefix("/removed/path/to")

		if gotStr != "/path/to" {
			t.Errorf("got: %v, want: %v", gotStr, gotErr)
		}

		if gotErr != nil {
			t.Errorf("got error: %v, want nil error", gotErr)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		gotStr, gotErr := RemovePathPrefix("path/to")

		if gotStr != "" {
			t.Errorf("got: %v, want empty string", gotStr)
		}

		if gotErr.Error() != "must be absolute" {
			t.Errorf(`got error: %v, want error "must be absolute"`, gotErr.Error())
		}
	})

	// XXX: edge cases currently panic, handle edge cases. "/", ""
}

func TestRemovePathSuffix(t *testing.T) {
	t.Run("removes path prefix", func(t *testing.T) {
		gotStr, gotErr := RemovePathSuffix("/path/to/removed")

		if gotStr != "/path/to/" {
			t.Errorf("got: %v, want: %v", gotStr, "/path/to/")
		}

		if gotErr != nil {
			t.Errorf("got error: %v, want nil error", gotErr)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		gotStr, gotErr := RemovePathSuffix("path/to")

		if gotStr != "" {
			t.Errorf("got: %v, want empty string", gotStr)
		}

		if gotErr.Error() != "must be absolute" {
			t.Errorf(`got error: %v, want error "must be absolute"`, gotErr.Error())
		}
	})

	t.Run("/", func(t *testing.T) {
		gotStr, gotErr := RemovePathSuffix("/")

		if gotStr != "" {
			t.Errorf("got: %v, want empty string", gotStr)
		}

		if gotErr.Error() != "input is /" {
			t.Errorf(`got error: %v, want error "input is /"`, gotErr.Error())
		}
	})
	// XXX: double check desired behavior for edge cases. "/", ""
}

func TestSystemBusPrivateUsable(t *testing.T) {
	t.Run("return conn", func(t *testing.T) {
		conn, err := SystemBusPrivateUsable()

		if conn == nil {
			t.Errorf("got conn %v", conn)
		}

		if err != nil {
			t.Errorf("got error %v", err)
		}
	})
	// XXX: testing other cases require refactoring(dependency injection, mock provider)
}

func TestSessionBusPrivateUsable(t *testing.T) {
	t.Run("return conn", func(t *testing.T) {
		conn, err := SessionBusPrivateUsable()

		if conn == nil {
			t.Errorf("got conn %v", conn)
		}

		if err != nil {
			t.Errorf("got error %v", err)
		}
	})
	// XXX: testing other cases require refactoring(dependency injection, mock provider)
}

func TestPathSliceSortMethod(t *testing.T) {
	s := PathSlice{"/c", "/b", "/a"}
	s.Sort()

	if s[0] != "/a" || s[1] != "/b" || s[2] != "/c" {
		t.Errorf("function PathSlice.Sort did not sort correctly, got: %v", s)
	}
}

func TestUInt64SliceSortMethod(t *testing.T) {
	s := UInt64Slice{3, 2, 1}
	s.Sort()

	if s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Errorf("function UInt64Slice.Sort did not sort correctly, got: %v", s)
	}
}

func TestPathSliceLessMethod(t *testing.T) {
	type input struct {
		s PathSlice
		i int
		j int
	}
	tests := []struct {
		name  string
		input input
		want  bool
	}{
		{
			name: "less",
			input: input{
				s: PathSlice{"/a", "/b"},
				i: 0,
				j: 1,
			},
			want: true,
		},
		{
			name: "not less",
			input: input{
				s: PathSlice{"/b", "/a"},
				i: 0,
				j: 1,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.s.Less(tt.input.i, tt.input.j)

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
	// XXX: check other test cases (unable to reach the code)
}

func TestValueToB64(t *testing.T) {
	t.Run("value to b64", func(t *testing.T) {
		v := "value"
		str, err := ValueToB64(v)

		valueInB64 := "EhAABnN0cmluZwwHAAV2YWx1ZQ=="
		if str != valueInB64 {
			t.Errorf("got: %v, want: %v", str, valueInB64)
		}

		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
		}
	})

	t.Run("passing function", func(t *testing.T) {
		str, err := ValueToB64(func() {})

		if str != "" {
			t.Errorf("wanted empty string, got: %v", str)
		}

		if !strings.Contains(err.Error(), "gob failed to encode") {
			t.Errorf(`expected error to contain "gob failed to encode", got: %v`, err.Error())
		}
	})
}

func TestB64ToValue(t *testing.T) {
	t.Run("b64 to value", func(t *testing.T) {
		b64 := "EhAABnN0cmluZwwHAAV2YWx1ZQ=="
		str, err := B64ToValue(b64)

		value := "value"
		if str != value {
			t.Errorf("got: %v, want: %v", str, value)
		}

		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
		}
	})

	t.Run("invalid b64", func(t *testing.T) {
		i, err := B64ToValue("invalid value")

		if i != nil {
			t.Errorf("wanted empty string, got: %v", i)
		}

		if !strings.Contains(err.Error(), "base64 failed to decode") {
			t.Errorf(`expected error to contain "base64 failed to decode", got: %v`, err.Error())
		}
	})

	t.Run("invalid gob", func(t *testing.T) {
		i, err := B64ToValue("dmFsdWU=")

		if i != nil {
			t.Errorf("wanted empty string, got: %v", i)
		}

		if !strings.Contains(err.Error(), "gob failed to decode") {
			t.Errorf(`expected error to contain "gob failed to decode", got: %v`, err.Error())
		}
	})

	// XXX: check unreachable case: "output `%v` is not a value"
}
