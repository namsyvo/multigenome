/*
Distance package provides some functions for Hamming and edit-distance calculation.

EditDistanceStr calculates the standard edit-distance.

EditDistance calculates the edit-distance for multi-genomes,
 and EditDistanceFull returns extra matching infomation.
*/

package Distance

import (
	"math"
    "fmt"
)

var A_to_iupac = map[byte]bool{
    'A': true, 'R': true, 'M': true, 'W': true, 'D': true, 'H': true, 'V': true, 'N': true,
    'a': true, 'r': true, 'm': true, 'w': true, 'd': true, 'h': true, 'v': true, 'n': true,
}

var C_to_iupac = map[byte]bool{
    'C': true, 'Y': true, 'M': true, 'S': true, 'B': true, 'H': true, 'V': true, 'N': true,
    'c': true, 'y': true, 'm': true, 's': true, 'b': true, 'h': true, 'v': true, 'n': true,
}

var G_to_iupac = map[byte]bool{
    'G': true, 'R': true, 'K': true, 'S': true, 'B': true, 'D': true, 'V': true, 'N': true,
    'g': true, 'r': true, 'k': true, 's': true, 'b': true, 'd': true, 'v': true, 'n': true,
}
var T_to_iupac = map[byte]bool{
    'T': true, 'Y': true, 'K': true, 'W': true, 'B': true, 'D': true, 'H': true, 'N': true,
    't': true, 'y': true, 'k': true, 'w': true, 'b': true, 'd': true, 'h': true, 'n': true,
}

var std_to_iupac = map[byte]map[byte]bool{
    'A': A_to_iupac,
    'C': C_to_iupac,
    'G': G_to_iupac,
    'T': T_to_iupac,
}

var snp_pos = map[int]int{
	1:100,
	5:500,
	10:1000,
	15:1500,
}

func PrintMaps(){
	fmt.Println(A_to_iupac)
	fmt.Println(C_to_iupac)
	fmt.Println(G_to_iupac)
	fmt.Println(T_to_iupac)
	fmt.Println(std_to_iupac)
}

var r_e float32 = 0.1
var r_m float32 = 0.01
var r_x float32 = 0.05
var infty float32 = 1000

var snp_sub_cost float32 = float32(math.Log10(float64(1/r_e)) + math.Log10(float64(1/r_m)))
var non_snp_sub_cost float32 = float32(math.Log10(float64(1/r_e)))

var Dist_Threshold int = 100

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

func minFloat32(a, b float32) float32 {
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
            if M1[k] <= Dist_Threshold {
                flag = false
                break
            }
        }
        if flag {
        	return Dist_Threshold + 1;
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
Cost functions for computing edit-distance for multi-genomes.
*/

func CostOfSub(s, t byte, pos int) float32 {
	if std_to_iupac[s][t] {
		return 0
	} else {
		_, ok := snp_pos[pos]
		if ok {
			return snp_sub_cost
		}
		return non_snp_sub_cost
	}
}

func CostOfIns(t byte, pos int) float32 {
	_, ok := snp_pos[pos]
	if ok {
		return 1
	} else {
		return 1 //infty
	}
}

func CostOfDel(s byte) float32 {
	return 1//infty
}

/*
EditDistance returns the edit-distance between reads and (a part of) multi-genomes.
The time complexity is O(mn) where m and n are lengths of a and b, and space complexity is O(n).
*/
func EditDistance(s, t []byte, pos int) float32 {
	ls, lt := len(s), len(t)
	lm := maxInt(ls, lt)

	M0, M1 := make([]float32, lm+1), make([]float32, lm+1)

	//Initialize the first row in ED matrix
	M0[0] = 0
	for j := 1; j <= lt; j++ {
		M0[j] = M0[j-1] + CostOfIns(t[j-1], pos)
	} // for j

	//Compute values in all rows in ED matrix
	for i := 1; i <= ls; i++ {
		//Initialize the first cell in row i in ED matrix
		M1[0] = M0[0]+ CostOfDel(s[i-1])

		//Compute values in all cell in row i in ED matrix
		for j := 1; j <= lt; j++ {
			mn := minFloat32(M0[j] + CostOfDel(s[i-1]), M1[j-1] + CostOfIns(t[j-1], pos)) // deletion & insertion
			mn = minFloat32(mn, M0[j-1] + CostOfSub(s[i-1], t[j-1], pos))                 // match/mismatch
			M1[j] = mn // min of 3 operations
		} // for j

        //Early break if ED exceeds given threshold
        flag := true
        for k := 0; k <= lt; k++ {
            if M1[k] <= float32(Dist_Threshold) {
                flag = false
                break
            }
        }
        if flag {
        	return float32(Dist_Threshold + 1);
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
EditDistanceFull returns the edit-distance and corresponding match indexes defined by Interface.
Each element in matA and matB is the index in the other list, if it is equal to or greater than zero;
 or -1 meaning a deleting or inserting in matA or matB, respectively.

The time and space complexity are all O(mn) where m and n are lengths of a and b.

NOTE if detailed matching information is not necessary, call EditDistance instead because it needs much less memories.
*/
/*
func EditDistanceFull(in Interface) (dist int, matA, matB []int) {
	la, lb := in.LenA(), in.LenB()

	f := make([]int, lb+1)
	ops := make([]byte, la*lb)

	for j := 1; j <= lb; j++ {
		f[j] = f[j-1] + in.CostOfIns(j-1)
	} // for j

	// Matching with dynamic programming
	p := 0
	for i := 0; i < la; i++ {
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0] += in.CostOfDel(i)
		for j := 1; j <= lb; j++ {
			mn, op := f[j]+in.CostOfDel(i), opDEL // delete

			if v := f[j-1] + in.CostOfIns(j-1); v < mn {
				// insert
				mn, op = v, opINS
			} // if

			// change/matched
			if v := fj1 + in.CostOfChange(i, j-1); v < mn {
				// insert
				mn, op = v, opCHANGE
			} // if

			fj1, f[j], ops[p] = f[j], mn, op // save f[j] to fj1(j is about to increase), update f[j] to mn
			p++
		} // for j
	} // for i
	// Reversely find the match info
	matA, matB = matchingFromOps(la, lb, ops)

	return f[lb], matA, matB
}
*/