//----------------------------------------------------------------------------------------
// Copyright 2013 Nam S. Vo

// Test for distance package
//----------------------------------------------------------------------------------------

package distance_test

import (
	"distance2"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"math"
)

//Initialize constants
var r_e float32 = 0.1
var r_m float32 = 0.01

var ERR_COST float32 = float32(math.Log10(float64(1/r_e))) // = 1

var DIST_THRES int = 1000

//Functions for displaying information
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
	fmt.Println()
}

type type_snpprofile map[int][][]byte
type type_samelensnp map[int]int
type TestCase struct {
	Profile type_snpprofile
	SNPlen type_samelensnp
	genome string
	read string
	d float32
}

// Test for distance between reads and "starred" multi-genomes
func TestDistanceMulti(t *testing.T) {
	defer __(o_())

	var test_cases = []TestCase{
		{ type_snpprofile{}, type_samelensnp{}, "ACG", "G", 0 },
		{ type_snpprofile{}, type_samelensnp{}, "TTACG", "ACT", 1 },

		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "ACCACGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "ACCCCGA", 1 },
		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "CACGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "CGCGT", 1 },
		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "ACCTCGT", 1 },

		{ type_snpprofile{3: {{'A'}, {'C'}, {'-'}} }, 	type_samelensnp{}, "ACC*CGT", "ACCACGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'C'}, {'-'}} }, 	type_samelensnp{}, "ACC*CGT", "ACCCCGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'C'}, {'-'}} }, 	type_samelensnp{}, "ACC*CGT", "ACCCGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'C'}, {'-'}} }, 	type_samelensnp{}, "ACC*CGT", "ACCTCGT", 1 },

		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCACGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCACCGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCATCGT", 1 },

		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCTCGT", 1 },
		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCACCGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'A','C'}} }, type_samelensnp{}, "ACC*CGT", "ACCTCCGT", 1 },

		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCTCGT", 0 },
		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCTTACGT", 0 },
		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCGTACGT", 1 },
		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCTACGT", 2 },
		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCTTCGT", 3 },
		{ type_snpprofile{3: {{'T'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCTTCGG", 4 },

		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "TTTACCACGT", float32(distance.INF) },
		{ type_snpprofile{3: {{'A'}, {'C'}} }, type_samelensnp{3: 1}, "ACC*CGT", "TTTACCACGT", float32(distance.INF) },

		{ type_snpprofile{3: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}} }, type_samelensnp{}, "ACC*CGT", "ACCGTACGT", 1 },
		{ type_snpprofile{3: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{}, "ACC*CGT", "ACCTTACGT", 0 },
		{ type_snpprofile{3: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{}, "ACC*CGT", "ACCGTACGT", 1 },

		{ type_snpprofile{0: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{}, "*ACGT", "TTAACGT", 0 },
		{ type_snpprofile{0: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{}, "*ACGT", "GTAACGT", 1 },
		{ type_snpprofile{0: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{}, "*ACGT", "ATACGT", 2 },

	}
	for i := 0; i < len(test_cases); i++ {
		distance.Init(ERR_COST, DIST_THRES, test_cases[i].Profile, test_cases[i].SNPlen)
		//fmt.Println(">", distance.ERR_COST, distance.INF, distance.DIST_THRES, distance.SNP_PROFILE, distance.SNP_LEN)
		dis, d, m, n, _ := distance.DistanceMulti([]byte(test_cases[i].read), []byte(test_cases[i].genome), 0)
		if dis != test_cases[i].d {
			t.Errorf("Fail alignment (case, calculated distance, true distance):", i, string([]byte(test_cases[i].read)), string([]byte(test_cases[i].genome)), dis, test_cases[i].d, d, m, n)
		} else {
			fmt.Println("Successful alignment (read, genome, distance): ", test_cases[i].read, test_cases[i].genome, test_cases[i].d)
		}
	}
}

/*
// Test for distance between reads and "starred" multi-genomes
func TestDistanceMultiSpeed(t *testing.T) {
	defer __(o_())

	var test_cases = []TestCase{
		{ type_profile{32: {{'A'}, {'T', 'A'}, {'T', 'T', 'A'}, {'-'}} }, type_samelensnp{32: 0}, "ACGTACGTACGTACGTACGTACGTACGTACGT*", "ACGTACGTACGTACGTACGTACGTACGTACGTAT", 2 },
	}
	k:=1000000
	for i := 0; i < len(test_cases); i++ {
		distance.Init(ERR_COST, DIST_THRES, test_cases[i].Profile, test_cases[i].SNPlen)
		//fmt.Println(">", distance.ERR_COST, distance.INF, distance.DIST_THRES, distance.SNP_PROFILE, distance.SNP_LEN)
		for j:=0; j < k; j++ {
			distance.DistanceMulti([]byte(test_cases[i].read), []byte(test_cases[i].genome), 0)
		}
		fmt.Println("Successful alignment for ", k, "reads")
	}
}
*/