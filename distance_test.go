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

type unitStr struct {
	Base
	a, b []byte
}

func (us *unitStr) CostOfChange(iA, iB int) int {
	if us.a[iA] == us.b[iB] {
		return 0
	} // if

	return 1
}

func TestEditDistanceStr(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2},
		{"abcde", "", 5},
		{"", "abcde", 5},
		{"", "", 0},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 2},
		{"abcde", "abfde", 1}}

	for _, cs := range cases {
		d := EditDistanceStr(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
	} // for case

	for _, cs := range cases {
		a, b := []byte(cs[0].(string)), []byte(cs[1].(string))
		d := EditDistance(&unitStr{Base{len(a), len(b), 1}, a, b})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
	} // for case
}

func TestEditDistanceStrFull(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2, "bcd"},
		{"abcde", "", 5, ""},
		{"", "abcde", 5, ""},
		{"", "", 0, ""},
		{"abcde", "abcde", 0, "abcde"},
		{"abcde", "dabce", 2, "abce"},
		{"abcde", "abfde", 1, "abde"}}

	for _, cs := range cases {
		d, lcs := EditDistanceStrFull(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
		if lcs != cs[3].(string) {
			t.Errorf("Longest-common-string between %s and %s is expected to be %s, but %s got!", cs[0], cs[1], cs[3], lcs)
		} // if
	} // for case
}

func ExampleString() {
	fmt.Println(EditDistanceStr("abcde", "bfdeg"))

	fmt.Println(EditDistanceStrFull("abcde", "bfdeg"))
	/* Output:
	3
3 bde
	*/
}

type stringInterface struct {
	a, b []byte
}

func (in *stringInterface) LenA() int {
	return len(in.a)
}

func (in *stringInterface) LenB() int {
	return len(in.b)
}

func (in *stringInterface) CostOfChange(iA, iB int) int {
	if base_cmp(in.a[iA], in.b[iB]) {
		return 0
	} // if

	return 1
}

func (in *stringInterface) CostOfDel(iA int) int {
	return 1
}

func (in *stringInterface) CostOfIns(iB int) int {
	return 1
}

func TestEditDistance(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"ACGTACGT", "", 8},
		{"", "ACGTACGT", 8},
		{"ACGT", "ACGT", 0},
		{"ACGT", "CGTA", 2},
		{"ACGAT", "TACGT", 2},
		{"GTCAC", "GTTAC", 1},
		{"TGCGAATGCTGACATAAGTAACGATAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACCAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAC",
		 "TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC", 52},
		{"AGGCCTCATTACCAGATTTGCGGTAGCCTTTTATCCAGTTTTTAACCGTGTTGTGGCTGGGAATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA",
		 "CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT", 56},
		{"AGCBBHTYHHHWDWD", "", 15},
		{"", "AGCBBHTYHHHWDWD", 15},
		{"AGCBBHTYHHHWDWD", "AGCCTTCAATATT", 2},
		{"AGCBBHTYHHHWDWD", "AGCCATCAATATT", 2},
		{"AGCBBHTYHHHWDWD", "AGCCTTCAATATC", 3},
		{"AGCBBHTYHHHWDWD", "AGCCGTCAATATT", 2},
		{"AGCBBHTYHHHWDWD", "AGCCATCATATT", 3},
		{"TGCGAATGCTGRCATAAGTAACGABAAAGCGGGTGAAAAGCCCGCTCGCCGGAAGACHAAGGGTTCCTGTCCAACGTTAATCGGGGCAGGGGGAGTCGAV",
		 "TCGACTCACCCTGCCCCGATTAACGTTGGACAGGAACCCTTGGTCTTCCGGCGAGCGGGCTTTTCACCCGCTTTATCGTTACTTATGTCAGCATTCGCAC", 52},
		{"AGGCCTCRTTACCAGATTTKCGGTAGCCTTTTATCCABTTTTTAACCGTGTTGTGGCTGGGHATGTTAAACCGGGCAGCTACCCGGGGTAGAGACTCAGA",
		 "CTGAGTCTCTACCCCGGGTAGCTGCCCGGTTTAACATTCCCAGCCACAACACGGTTAAAAACTGGATAAAAGGCTACCGCAAATCTGGTAATGAGGCCTT", 56},
		{"", "", 0}}

	for _, cs := range cases {
		d := EditDistance(&stringInterface{[]byte(cs[0].(string)), []byte(cs[1].(string))})
		if d != cs[2].(int) {
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
/*
func TestEditDistanceFunc(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2},
		{"abcde", "", 5},
		{"", "abcde", 5},
		{"", "", 0},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 2},
		{"abcde", "abfde", 1}}

	for _, cs := range cases {
		a, b := []rune(cs[0].(string)), []rune(cs[1].(string))
		d := EditDistanceFunc(len(a), len(b),
			func(iA, iB int) int {
				if a[iA] == b[iB] {
					return 0
				} // if

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceFuncFull(t *testing.T) {
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
		a, b := []rune(cs[0].(string)), []rune(cs[1].(string))
		d, matA, matB := EditDistanceFuncFull(len(a), len(b),
			func(iA, iB int) int {
				if a[iA] == b[iB] {
					return 0
				} // if

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
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

func ExampleEditDistanceFunc() {
	a, b := "abcd", "bcde"
	d := EditDistanceFunc(len(a), len(b), func(iA, iB int) int {
		return Ternary(a[iA] == b[iB], 0, 1)
	}, ConstCost(1), ConstCost(1))
	fmt.Println(a, b, d)
	// Output:
	// abcd bcde 2
}
*/