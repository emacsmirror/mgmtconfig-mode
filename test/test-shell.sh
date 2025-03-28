#!/usr/bin/env bash
# simple test harness for testing mgmt
# NOTE: this will rm -rf /tmp/mgmt/

echo running "$0"
set -o errexit
set -o pipefail

#ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"	# dir!
ROOT=$(dirname "${BASH_SOURCE}")/..
cd "${ROOT}"
. test/util.sh
cd - >/dev/null

if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
	echo -e "usage: ./"`basename $0`" [[--help] | <test>]"
	echo -e "where: <test> is empty to run all tests, or <file>.sh from shell/ dir"
	exit 1
fi

COLS="$(tput cols 2>/dev/null || echo 80)"	# github-actions has no $TERM
LINE=$(printf '=%.0s' `seq -s ' ' $COLS`)	# a terminal width string
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"	# dir!
cd "$DIR" >/dev/null	# work from main mgmt directory
failures=""
count=0

# loop through tests
for test_script in $DIR/test/shell/*.sh; do
	[ -x "$test_script" ] || continue	# file must be executable
	test_name=`basename "$test_script"`	# short name
	# if ARGV has test names, only execute those!
	if [ "$1" != '' ]; then
		[ "$test_name" != "$1" ] && continue
	fi
	cd $DIR/test/shell/ >/dev/null	# shush the cd operation
	mkdir -p '/tmp/mgmt/'	# directory for mgmt to put files in
	#echo "Running: $test_name"
	export MGMT_TMPDIR='/tmp/mgmt/'	# we can add to env like this
	count=`expr $count + 1`
	set +o errexit	# don't kill script on test failure
	out=$($test_script 2>&1)	# run and capture stdout & stderr
	e=$?	# save exit code
	set -o errexit	# re-enable killing on script failure
	cd - >/dev/null
	if [ -L '/tmp/mgmt/' ]; then	# this was once a symlink :(
		fail_test "Can't remove symlink in /tmp/mgmt/"
	fi
	rm -rf '/tmp/mgmt/'	# clean up after test
	if [ $e -ne 0 ]; then
		echo -e "FAIL\t$test_name"	# fail
		# store failures...
		failures=$(
			# prepend previous failures if any
			[ -n "${failures}" ] && echo "$failures" && echo "$LINE"
			echo "Script: $test_name"
			# if we see 124, it might be the exit value of timeout!
			[ $e -eq 124 ] && echo "Exited: $e (timeout?)" || echo "Exited: $e"
			if [ "$out" = "" ]; then
				echo "Output: (empty!)"
			else
				echo "Output:"
				echo "$out"
			fi
		)
	else
		echo -e "ok\t$test_name"	# pass
	fi
done

if [ "$count" = '0' ]; then
	fail_test 'No tests were run!'
fi

# display errors
if [[ -n "${failures}" ]]; then
	echo 'FAIL'
	echo 'The following tests failed:'
	echo "${failures}"
	exit 1
fi
echo 'PASS'
