package distance_test

import (
	"seq-cmp.2"
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
		d := distance.EditDistanceStr(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestHammingDistance(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"AGCBBHTYHHHWDWD", "", 0},
		{"", "AGCBBHTYHHHWDWD", 0},

		{"AGCCTTCAATATT", "AGCbBHyHHWdWD", 0},
		{"AGCCATCAATATT", "ATCBBHYHHWDWD", 4},
		{"AGCCTTCAATATC", "AGCBBHYHHWDWD", 3},
		{"AGCCGTCAATATT", "GGCBBhYHhWDwD", 1},
		{"AGCCATCATAT"  , "AGCBbHYHHwDWD", 3},

		{"TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC",
		 "TGCGAATGCTGRCATAAGTAACGABAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACHAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAV", 104},

		{"CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT",
		 "AGGCCTCRTTACCAGATTTKCGGTAGCCTTTTATCCABTTTTTAACCGTGTTGTGGCTGGGHATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA", 82},
	}

	for _, cs := range cases {
		d := distance.HammingDistance([]byte(cs[0].(string)), []byte(cs[1].(string)))
		if d != float32(cs[2].(int)) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceArray(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"AGCBBHTYHHHWDWD", "", distance.DIST_THRES + 1},
		{"", "AGCBBHTYHHHWDWD", int(3.2212255e+10)},

		{"AGCCTTCAATATT", "AGCBBHtYhHHWDWD", 2},
		{"AGCCATCAATATT", "AGCBBHtYhHHWDWD", 5},
		{"AGCCTTCAATATC", "AGCBBHtYhHHWDWD", 5},
		{"AGCCGTCAATATT", "AGCBBHtYhHHWDWD", 2},
		{"AGCCATCATATT" , "AGCBBHtYhHHWDWD", int(2.1474836e+09)},

		{"TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC",
		 "TGCGAATGCTGRCATAAGTAACGABAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACHAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAV", 104},

		{"CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT",
		 "AGGCCTCRTTACCAGATTTKCGGTAGCCTTTTATCCABTTTTTAACCGTGTTGTGGCTGGGHATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA", 82},
	}

	fmt.Println("\tSNP_SUB_COST: ", distance.SNP_SUB_COST)
	fmt.Println("\tN_SNP_SUB_COST: ", distance.N_SNP_SUB_COST)
	fmt.Println("\tINF: ", distance.INF)
	fmt.Println("\tDIST_THRES: ", distance.DIST_THRES)

	for _, cs := range cases {
		d := distance.EditDistanceArray([]byte(cs[0].(string)), []byte(cs[1].(string)))
		if d != float32(cs[2].(int)) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceMatrix(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"", "", 0},
		{"AGCBBHTYHHHWDWD", "", distance.DIST_THRES + 1},
		{"", "AGCBBHTYHHHWDWD", int(3.2212255e+10)},

		{"AGCCTTCAATATT", "AGCBBHtYhHHWDWD", 2},
		{"AGCCATCAATATT", "AGCBBHtYhHHWDWD", 5},
		{"AGCCTTCAATATC", "AGCBBHtYhHHWDWD", 5},
		{"AGCCGTCAATATT", "AGCBBHtYhHHWDWD", 2},
		{"AGCCATCATATT" , "AGCBBHtYhHHWDWD", int(2.1474836e+09)},

		{"TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC",
		 "TGCGAATGCTGRCATAAGTAACGABAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACHAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAV", 104},

		{"CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT",
		 "AGGCCTCRTTACCAGATTTKCGGTAGCCTTTTATCCABTTTTTAACCGTGTTGTGGCTGGGHATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA", 82},
	}

	fmt.Println("\tSNP_SUB_COST: ", distance.SNP_SUB_COST)
	fmt.Println("\tN_SNP_SUB_COST: ", distance.N_SNP_SUB_COST)
	fmt.Println("\tINF: ", distance.INF)
	fmt.Println("\tDIST_THRES: ", distance.DIST_THRES)

	for _, cs := range cases {
		d := distance.EditDistanceMatrix([]byte(cs[0].(string)), []byte(cs[1].(string)))
		if d != float32(cs[2].(int)) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}