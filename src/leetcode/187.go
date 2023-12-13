package leetcode

func findRepeatedDnaSequences(s string) []string {
	stringMap := make(map[string]int)
	ans := make([]string, 0)
	for i := 0; i < len(s)-9; i++ {
		if _, ok := stringMap[s[i:i+10]]; ok {
			stringMap[s[i:i+10]]++
		} else {
			stringMap[s[i:i+10]] = 0
		}
	}
	for k, v := range stringMap {
		if v > 0 {
			ans = append(ans, k)
		}
	}
	return ans
}
