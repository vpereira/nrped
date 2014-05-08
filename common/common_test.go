package common

import (
    "testing"
)

func TestCommonDoCRC32(t *testing.T) {
    if crc32, _ := DoCRC32("hello"); crc32 != 907060870 {
        t.Error("DoCRC32 failed")
    }
}
