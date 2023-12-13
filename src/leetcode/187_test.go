package leetcode

import "testing"

func TestFindRepeatedDnaSequences(t *testing.T) {
	ans := findRepeatedDnaSequences("AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT")
	t.Log("ans:", ans)
}
