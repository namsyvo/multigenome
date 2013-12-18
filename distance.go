//-------------------------------------------------------------------------------------------------
// Copyright 2013 Nam S. Vo

// Distance package provides some functions for distance between reads and "starred" multi-genomes.

// DistanceMulti calculates the distance between reads and "starred" multi-genomes.
//-------------------------------------------------------------------------------------------------

package distance

import (
	"math"
	"fmt"
)

//-------------------------------------------------------------------------------------------------
//Utility functions: max, min
//-------------------------------------------------------------------------------------------------

func minInt(a, b int) int {
	if a < b {
		return a
	} // if
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} // if
	return b
}

func min3Int(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func min3Float32(a, b, c float32) float32 {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

// Minimum of a slice of int arguments
func MinIntS(v []int) (m int) {
        if len(v) > 0 {
                m = v[0]
        }
        for i := 1; i < len(v); i++ {
                if v[i] < m {
                        m = v[i]
                }
        }
        return
}

// Minimum of a variable number of int arguments
func MinIntV(v1 int, vn ...int) (m int) {
        m = v1
        for i := 0; i < len(vn); i++ {
                if vn[i] < m {
                        m = vn[i]
                }
        }
        return
}

//-------------------------------------------------------------------------------------------------
// Functions for distance between reads and multi-genomes.
//-------------------------------------------------------------------------------------------------

//-------------------------------------------------------------------------------------------------
// Constants and global variables.
//-------------------------------------------------------------------------------------------------

// Cost for substitution with error sequencing
var ERR_COST float32 = 1

// Distance threshold for early break
var DIST_THRES int = 1000

// SNP profile
var SNP_PROFILE map[int][][]byte

// SNP profile
var SAME_LEN_SNP map[int]int

// Approx value for Infinity
var INF int = math.MaxInt32

//-------------------------------------------------------------------------------------------------
// Cost functions for computing distance between reads and multi-genomes.
//-------------------------------------------------------------------------------------------------

// Cost for "SNP match"
func Cost(s, t []byte) int {
	mismatch := 0
	for i:= 0; i < len(s); i++ {
		if s[i] != t[i] {
			mismatch++
		}
	}
	return mismatch
}

//-------------------------------------------------------------------------------------------------
// Functions for computing distance between reads and multi-genomes.
//-------------------------------------------------------------------------------------------------

// Initilize constants and global variables
func Init(pERR_COST float32, pDIST_THRES int, pSNP_PROFILE map[int][][]byte, pSAME_LEN_SNP map[int]int) {
	ERR_COST = pERR_COST
	DIST_THRES = pDIST_THRES
	SNP_PROFILE = pSNP_PROFILE
	SAME_LEN_SNP = pSAME_LEN_SNP
}

//-------------------------------------------------------------------------------------------------
// Calculate the distance between reads and multi-genomes.
// The reads include standard bases, the multi-genomes include standard bases and "*" characters.
// 	s is a read.
// 	t is part of a multi-genome.
//-------------------------------------------------------------------------------------------------
func DistanceMulti(s, t []byte, pos int) (float32, int, int, int, [][]int) {

	//fmt.Println(string(s), string(t))
    m, n := len(s), len(t)
    var i, j, k int

    //"Hamming distance" part
    var cost int
    var d, min_d int
    var snp_len int
    var snp_values [][]byte
    var is_snp, is_same_len_snp bool

    d = 0
    for m > 0 && n > 0 {
		snp_values, is_snp = SNP_PROFILE[pos + n - 1]
		snp_len, is_same_len_snp = SAME_LEN_SNP[pos + n - 1]
    	if !is_snp {
        	if s[m-1] != t[n-1] {
        		d++
        	}
    		m--
    		n--
    	} else if is_same_len_snp {
			min_d = INF
    		for i = 0; i < len(snp_values); i++ {
    			cost = Cost(s[m - snp_len: m], snp_values[i])
    			if min_d > cost {
    				min_d = cost
    			}
    		}
    		d += min_d
    		m -= snp_len
    		n--
    	} else {
    		break
    	}
    }
    //fmt.Println(d, m, n)

    //"Edit distance" part

    // Make a 2-D matrix:
    // 	rows correspond to prefixes of source, columns to prefixes of target.
    // 	cells contain edit distances.
    H := make([][]int, m + 1)
    for i = 0; i <= m; i++ {
            H[i] = make([]int, n + 1)
    }
    // Initialize distances (base cases, from/to empty string).
	// Initialize the first row in distance matrix
    H[0][0] = 0
	//Initialize the first column in distance matrix
    for i = 1; i <= m; i++ {
		H[i][0] = INF
    }
	//Initialize the first row in distance matrix
    for j = 1; j <= n; j++ {
		H[0][j] = 0
    }
	//fmt.Println(H)
    // Fill in the remaining cells follow the Levenshtein algorithm.
	//Compute values in all rows in distance matrix
	var temp int
    for i = 1; i <= m; i++ {
		//fmt.Println("i: ", i)
		//Compute values in all cell in row i in distance matrix
        for j = 1; j <= n; j++ {
			snp_values, is_snp = SNP_PROFILE[pos + j - 1]
	    	if !is_snp {
				if s[i-1] != t[j-1] {
					H[i][j] = H[i - 1][j - 1] + 1
				} else {
					H[i][j] = H[i - 1][j - 1]
				}
	   		} else {
				H[i][j] = INF
				for k = 0; k < len(snp_values); k++ {
					if i - len(snp_values[k]) >= 0 {
						if snp_values[k][0] == '-' {
	    						temp = H[i][j - 1]
						} else {
							temp = H[i - len(snp_values[k])][j - 1] + Cost(s[i - len(snp_values[k]) : i], snp_values[k])
						}
	    				if H[i][j] > temp {
	    					H[i][j] = temp
	    				}
						//fmt.Println(">>> ",k, snp_values[k], string(s[i - len(snp_values[k]) : i]), string(snp_values[k]), d)
					}
	   			}
			}
			//fmt.Println("j: ", j, ": ", H[i][j], " ")
		}
		//fmt.Println()
/*		//Early break if ED exceeds given threshold
		flag := true
		for k = 0; k <= n; k++ {
		    if H[i][k] <= DIST_THRES - d {
		        flag = false
		        break
		    }
		}
		if flag {
			return float32(DIST_THRES + 1) * ERR_COST, d, m, n, [][]int{}
		}
*/
	}
	//fmt.Println(H)
    return float32(d + H[m][n]) * ERR_COST, d, m, n, H
}


//-------------------------------------------------------------------------------------------------
// Construct an alignment of s1 and s2 as pair of indexes, traceback from dynamic programming table
//-------------------------------------------------------------------------------------------------
func TraceBack(s1, s2 []byte, d [][]float32, pos int) []int {

	//trace back to find number of insertions / deletions / substitutions
	var i_count, d_count, s_count int = 0, 0, 0

	var i, j int = int(len(s1)), int(len(s2))
	//fmt.Println("(", i, ", ", j, "); ")
	for  i > 0 || j > 0 {

		if i > 0 && j > 0 {
			var diff float32 = d[i][j] - min3Float32(d[i - 1][j - 1], d[i - 1][j], d[i][j - 1]);
		  	if diff == (d[i][j] - d[i - 1][j - 1]) {
				if s1[i-1] == s2[j-1] { // matched, no operation
		  			fmt.Println("(", i - 1, ", ", j - 1, "); ")
		  			i--;
		  			j--;
				} else { // substitution operation
		  			fmt.Println("(", i - 1, ", ", j - 1, "); ")
		  			s_count++;
		  			i--;
		  			j--;
				}
		  	} else if diff == d[i][j] - d[i - 1][j] { // deletion operation
				fmt.Println("(", i - 1, ", ", "-", "); ")
				d_count++;
				i--;
		  	} else if diff == d[i][j] - d[i][j - 1] { // insertion operation
				fmt.Println("(", "-", ", ", j - 1, "); ")
				i_count++;
				j--;
		  	}
		} else if i == 0 { // insertion operation
		  fmt.Println("(","-", ", ", j - 1, "); ")
		  i_count++;
		  j--;	    
		} else if j == 0 { // deletion operation
		  fmt.Println("(", i - 1, ", ", "-", "); ")
		  d_count++;
		  i--;	
		}    
	}
	return []int{i_count, d_count, s_count}
}