//----------------------------------------------------------------------------------------
// Copyright 2013 Nam S. Vo
// Test for multigenome2 package
//----------------------------------------------------------------------------------------

package multigenome2

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

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

func TestMultiGenomeBuild(t *testing.T) {
    defer __(o_())

	sequence := fastaRead("data-test/chr1.fasta")
	SNP_array := vcfRead("data-test/vcf_chr_1.vcf")

	genome := buildMultigenome2(SNP_array, sequence)
	SaveSNPLocation("data-test/SNPLocation.txt", SNP_array)
	SaveMulti("data-test/genomestar.txt", genome)
	genome_test := LoadMulti("data-test/genomestar.txt")
	byte_array, flag := LoadSNPLocation("data-test/SNPLocation.txt")

	fmt.Printf("%s \n", genome_test)
	fmt.Println(byte_array)
	fmt.Println(flag)
}
