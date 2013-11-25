package Distance

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func o_() string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	if p := strings.LastIndexAny(name, `./\`); p >= 0 {
		name = name[p+1:]
	} // if
	fmt.Println("== BEGIN", name, "===")
	return name
}

func __(name string) {
	fmt.Println("== END", name, "===")
}

func TestEditDistanceStr(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"abcde", "", 5},
		{"", "abcde", 5},
		{"abcd", "bcde", 2},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 2},
		{"abcde", "abfde", 1},
	}

	for _, cs := range cases {
		d := EditDistanceStr(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
	} // for case

}

func TestEditDistance(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"ACGTACGT", "", 8},
		{"", "ACGTACGT", 8},
		{"ACGT", "ACGT", 0},
		{"CGTA", "ACGT", 2},
		{"TACGT", "ACGAT", 2},
		{"GTTAC", "GTCAC", 1},
		{"TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC",
		 "TGCGAATGCTGACATAAGTAACGATAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACCAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAC", 52},
		{"CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT",
		 "AGGCCTCATTACCAGATTTGCGGTAGCCTTTTATCCAGTTTTTAACCGTGTTGTGGCTGGGAATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA", 56},
	}

	for _, cs := range cases {
		d := EditDistance([]byte(cs[0].(string)), []byte(cs[1].(string)), 0)
		if d != float32(cs[2].(int)) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceIUPAC(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"AGCBBHTYHHHWDWD", "", 15},
		{"", "AGCBBHTYHHHWDWD", 15},

		{"AGCCTTCAATATT", "AGCBBHTYHHHWDWD", 2},
		{"AGCCATCAATATT", "AGCBBHTYHHHWDWD", 2},
		{"AGCCTTCAATATC", "AGCBBHTYHHHWDWD", 3},
		{"AGCCGTCAATATT", "AGCBBHTYHHHWDWD", 2},
		{"AGCCATCATATT" , "AGCBBHTYHHHWDWD", 3},

		{"TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC",
		 "TGCGAATGCTGRCATAAGTAACGABAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACHAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAV", 52},

		{"CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT",
		 "AGGCCTCRTTACCAGATTTKCGGTAGCCTTTTATCCABTTTTTAACCGTGTTGTGGCTGGGHATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA", 56},
	}
	for _, cs := range cases {
		d := EditDistance([]byte(cs[0].(string)), []byte(cs[1].(string)), 0)
		if d != float32(cs[2].(int)) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

/*
func TestEditDistanceFull(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2, []int{-1, 0, 1, 2}, []int{1, 2, 3, -1}},
		{"abcde", "", 5, []int{-1, -1, -1, -1, -1}, []int{}},
		{"", "abcde", 5, []int{}, []int{-1, -1, -1, -1, -1}},
		{"", "", 0, []int{}, []int{}},
		{"abcde", "abcde", 0, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"abcde", "dabce", 2, []int{1, 2, 3, -1, 4}, []int{-1, 0, 1, 2, 4}},
		{"abcde", "abfde", 1, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}}}

	for _, cs := range cases {
		d, matA, matB := EditDistanceFull(&stringInterface{[]rune(cs[0].(string)), []rune(cs[1].(string))})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
		if fmt.Sprint(matA) != fmt.Sprint(cs[3]) {
			t.Errorf("matA for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[3], matA)
		} // if
		if fmt.Sprint(matB) != fmt.Sprint(cs[4]) {
			t.Errorf("matB for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[4], matB)
		} // if
	} // for case
}
*/