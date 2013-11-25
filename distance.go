/*
Distance package provides some functions for Hamming and edit-distance calculation.

EditDistanceStr calculates the standard edit-distance.

EditDistance calculates the edit-distance for multi-genomes,
 and EditDistanceFull returns extra matching INFomation.
*/

package distance

import (
	"math"
    "fmt"
)

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

func mINFloat32(a, b float32) float32 {
	if a < b {
		return a
	} // if
	return b
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	} // if
	return b
}

/*
EditDistanceStr calculates the standard edit-distance between two strings.
Memory-improved version.
The time complexity is O(mn) where m and n are lengths of s and t, and space complexity is O(n).
*/

func EditDistanceStr(s, t string) int {
	ls, lt := len(s), len(t)
	lm := maxInt(ls, lt)

	M0, M1 := make([]int, lm+1), make([]int, lm+1)

	//Initialize the first row in ED matrix
	for j := 0; j <= lt; j++ {
		M0[j] = j
	} // for j
	//Compute values in all rows in ED matrix
	for i := 1; i <= ls; i++ {
		//Initialize the first cell in row i in ED matrix
		M1[0] = i

		//Compute values in all cell in row i in ED matrix
		for j := 1; j <= lt; j++ {
			mn := minInt(M0[j] + 1, M1[j-1] + 1) // deletion & insertion
			if s[i-1] == t[j-1] {
				mn = minInt(mn, M0[j-1])         // match
			} else {
				mn = minInt(mn, M0[j-1] + 1)     // substitution (mismatch)
			}
			M1[j] = mn // min of 3 operations
		} // for j

        //Early break if ED exceeds given threshold
        flag := true
        for k := 0; k <= lt; k++ {
            if M1[k] <= DIST_THRES {
                flag = false
                break
            }
        }
        if flag {
        	return DIST_THRES + 1;
        }

		//Swap M0 and M1
		M := make([]int, lm+1)
		copy(M, M0)
		copy(M0, M1)
		copy(M1, M)

	} // for i

	return M0[lt]
}

/*
Map of standard bases (A, C, G, T) and extended IUPAC bases.
A standard base is mapped into any extended IUPAC base that it belongs to.
An extended IUPAC base is an IUPAC base or a gap.
*/

var std_to_iupac = map[byte]map[byte]bool {
    'A': map[byte]bool{
	    'A': true, 'R': true, 'M': true, 'W': true, 'D': true, 'H': true, 'V': true, 'N': true,
    	'a': true, 'r': true, 'm': true, 'w': true, 'd': true, 'h': true, 'v': true, 'n': true,
    	},
    'C': map[byte]bool{
    	'C': true, 'Y': true, 'M': true, 'S': true, 'B': true, 'H': true, 'V': true, 'N': true,
    	'c': true, 'y': true, 'm': true, 's': true, 'b': true, 'h': true, 'v': true, 'n': true,
    	},
    'G': map[byte]bool{
    	'G': true, 'R': true, 'K': true, 'S': true, 'B': true, 'D': true, 'V': true, 'N': true,
    	'g': true, 'r': true, 'k': true, 's': true, 'b': true, 'd': true, 'v': true, 'n': true,
		},
    'T': map[byte]bool{
	    'T': true, 'Y': true, 'K': true, 'W': true, 'B': true, 'D': true, 'H': true, 'N': true,
    	't': true, 'y': true, 'k': true, 'w': true, 'b': true, 'd': true, 'h': true, 'n': true,
		},
}

//Print Std-IUPAC map
func PrintMaps(){
	fmt.Println(std_to_iupac)
}

/*
Map of SNP region and multi-genomes.
A position on SNP region is mapped into an actual position on the multi-genome.
*/
var snp_pos = map[int]int{
	3:100,
	4:100,
	5:100,
	6:100,
	7:100,
	8:100,
	9:100,
	10:100,
	11:100,
	12:100,
	13:100,
	14:100,
}

//Const for computing cost for edit distance
var r_e float32 = 0.1
var r_m float32 = 0.01
var r_x float32 = 0.05

var SNP_SUB_COST float32 = float32(math.Log10(float64(1/r_e)) + math.Log10(float64(1/r_m))) // = 3
var N_SNP_SUB_COST float32 = float32(math.Log10(float64(1/r_e))) // = 1

var INF float32 = float32(math.MaxInt32)
var DIST_THRES int = 1000


//Cost functions for computing edit-distance for multi-genomes.

// Cost for aligning (x_i, -)
var CostOfDel = INF

// Cost for aligning (-, y_j)
// Go should called this as an inline function
func CostOfIns(t byte) float32 {
	//If t is lowercase (that is, an IUPAC base or a gap)
	if 97 <= t && t <= 122 {
		return 1
	} else {
		return INF
	}
}

// Cost for aligning (x_i, y_j)
func CostOfSub(s, t byte, pos int) float32 {
	if s == t {
		return 0
	}
	//Check for SNP position
	_, is_snp := snp_pos[pos]
	if !is_snp {
		return N_SNP_SUB_COST
	} else {
		//Check for match between a standard base and an extended IUPAC base
		_, matched := std_to_iupac[s][t]
		if matched {
			return 0
		}
		return SNP_SUB_COST
	}
}


/*
HammingDistance returns the Hamming distance between reads (s) and (a part of) multi-genome (t),
s is a read
t is part of the SNP region
The time complexity is O(min(m, n)) where m and n are lengths of s and t.
*/

func HammingDistance(s, t []byte) float32 {
	lm := minInt(len(s), len(t))
	var dis float32 = 0;
	//Compute values in all rows in ED matrix
	for i := 0; i < lm; i++ {
        dis = dis + CostOfSub(s[i], t[i], i)

        //Early break if ED exceeds given threshold
        if dis > float32(DIST_THRES) {
            break
        }
    }
	return dis
}

/*
EditDistance returns the edit-distance between reads (s) and (a part of) multi-genome (t),
s is a read
t is part of the SNP region
using the Levenshtein algorithm with a memory-improved version using two 1-D arrays.
The time complexity is O(mn) where m and n are lengths of s and t, and space complexity is O(m).
*/

func EditDistanceArray(s, t []byte) float32 {
	ls, lt := len(s), len(t)
	lm := maxInt(ls, lt)

	// Two 1-D array storing dynamic programing values
	M0, M1 := make([]float32, lm+1), make([]float32, lm+1)

	//Initialize the first row in ED matrix
	M0[0] = 0
	for j := 1; j <= lt; j++ {
		M0[j] = M0[j-1] + CostOfIns(t[j-1])
	} // for j

	//Compute values in all rows in ED matrix
	for i := 1; i <= ls; i++ {
		//Initialize the first cell in row i in ED matrix
		M1[0] = M0[0]+ CostOfDel
		//Compute values in all cell in row i in ED matrix
		for j := 1; j <= lt; j++ {
            insCost := M1[j-1] + CostOfIns(t[j-1])
            delCost := M0[j] + CostOfDel
            matchSubCost := M0[j-1] + CostOfSub(s[i-1], t[j-1], j-1)
			M1[j] = mINFloat32(delCost, mINFloat32(matchSubCost, insCost))
		} // for j

        //Early break if ED exceeds given threshold
        flag := true
        for k := 0; k <= lt; k++ {
            if M1[k] <= float32(DIST_THRES) {
                flag = false
                break
            }
        }
        if flag {
        	return float32(DIST_THRES + 1);
        }

		//Swap M0 and M1
		M := make([]float32, lm+1)
		copy(M, M0)
		copy(M0, M1)
		copy(M1, M)

	} // for i

	return M0[lt]
}

/*
EditDistance returns the edit-distance between reads (s) and (a part of) multi-genome (t),
s is a read
t is part of the SNP region
using the Levenshtein algorithm with dynamic programming table as 2-D matrix.
The time complexity is O(mn) where m and n are lengths of a and b, and space complexity is O(mn).
*/

func EditDistanceMatrix(s, t []byte) float32 {
    // Make a 2-D matrix: rows correspond to prefixes of source, columns to prefixes of target.
    // Cells will contain edit distances.
    ls, lt := len(s), len(t)
    H := make([][]float32, ls + 1)

    // Initialize distances (base cases, from/to empty string).
    // Fill the left column and the top row with row/column indices.
    for i := 0; i <= ls; i++ {
            H[i] = make([]float32, lt + 1)
    }

	//Initialize the first row in ED matrix
    H[0][0] = 0
    for i := 1; i <= ls; i++ {
            H[i][0] = H[i-1][0] + CostOfDel
    }
	//Initialize the first column in ED matrix
    for j := 1; j <= lt; j++ {
            H[0][j] = H[0][j-1] + CostOfIns(t[j-1])
    }

    // Fill in the remaining cells follow the Levenshtein algorithm.
	//Compute values in all rows in ED matrix
    for i := 1; i <= ls; i++ {
		//Compute values in all cell in row i in ED matrix
        for j := 1; j <= lt; j++ {
            insCost := H[i][j - 1] + CostOfIns(t[j-1])
            delCost := H[i - 1][j] + CostOfDel
            matchSubCost := H[i - 1][j - 1] + CostOfSub(s[i-1], t[j-1], j-1)
            H[i][j] = mINFloat32(delCost, mINFloat32(matchSubCost, insCost))
       	}
		//Early break if ED exceeds given threshold
		flag := true
		for k := 0; k <= lt; k++ {
		    if H[i][k] <= float32(DIST_THRES) {
		        flag = false
		        break
		    }
		}
		if flag {
			return float32(DIST_THRES + 1);
		}
    }
    //LogMatrix(source, target, matrix)
    return H[ls][lt]
}