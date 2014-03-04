//-------------------------------------------------------------------------------------------------
// Copyright 2013 Nam S. Vo

// Distance package provides some functions for distance between reads and "starred" multi-genomes.

// DistanceMulti calculates the distance between reads and "starred" multi-genomes.
//-------------------------------------------------------------------------------------------------

package distance

import (
	"math"
	//"fmt"
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

// SNP profile
var SNP_PROFILE map[int][][]byte

// Same length SNP
var SAME_LEN_SNP map[int]int

// Distance threshold for early break
var DIST_THRES int = math.MaxInt16

// Value for Infinity
var INF int = math.MaxInt16

//-------------------------------------------------------------------------------------------------
// Cost functions for computing distance between reads and multi-genomes.
//-------------------------------------------------------------------------------------------------

// Cost for "SNP match"
// Input slices should have same length
func Cost(s, t []byte) int {
	for i:= 0; i < len(s); i++ {
		if s[i] != t[i] {
			return INF
		}
	}
	return 0
}

//-------------------------------------------------------------------------------------------------
// Functions for computing distance between reads and multi-genomes.
//-------------------------------------------------------------------------------------------------

// Initilize constants and global variables
func Init(pDIST_THRES int, pSNP_PROFILE map[int][][]byte, pSAME_LEN_SNP map[int]int) {
	SNP_PROFILE = pSNP_PROFILE
	SAME_LEN_SNP = pSAME_LEN_SNP
	DIST_THRES = pDIST_THRES
}

//-------------------------------------------------------------------------------------------------
// Calculate the distance between s and t in backward direction.
// 	s is a read.
// 	t is part of a multi-genome.
// The reads include standard bases, the multi-genomes include standard bases and "*" characters.
//-------------------------------------------------------------------------------------------------
func BackwardDistanceMulti(s, t []byte, pos int) (int, int, int, int, map[int][]byte, map[int]int) {

    var cost int
    var d, min_d int
    var snp_len int
    var snp_values [][]byte
    var is_snp, is_same_len_snp bool
	var S = make(map[int][]byte)

    var i, j, k int

    d = 0
    m, n := len(s), len(t)
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
			min_d = 1000*INF
    		for i = 0; i < len(snp_values); i++ {
    			cost = Cost(s[m - snp_len: m], snp_values[i])
    			if min_d > cost {
    				min_d = cost
    			}
    		}
	  		if min_d == INF {
	  			return INF, 0, m, n, make(map[int][]byte), make(map[int]int)
	  		}
	  		S[pos + n - 1] = s[m - snp_len: m]
    		d += min_d
    		m -= snp_len
    		n--
    	} else {
    		break
    	}
    }

    D := make([][]int, m + 1)
    T := make(map[int]int)
    for i = 0; i <= m; i++ {
		D[i] = make([]int, n + 1)
    }
    D[0][0] = 0
    for i = 1; i <= m; i++ {
		D[i][0] = INF
    }
    for j = 1; j <= n; j++ {
		D[0][j] = 0
    }
	var temp_dis, min_index int
    for i = 1; i <= m; i++ {
        for j = 1; j <= n; j++ {
			snp_values, is_snp = SNP_PROFILE[pos + j - 1]
	    	if !is_snp {
				if s[i-1] != t[j-1] {
					D[i][j] = D[i - 1][j - 1] + 1
				} else {
					D[i][j] = D[i - 1][j - 1]
				}
	   		} else {
				D[i][j] = 1000*INF //1000*INF is a value for testing, will change to a better solution later
				min_index = 0
				for k = 0; k < len(snp_values); k++ {
					snp_len = len(snp_values[k])
					//One possible case: i - snp_len < 0 for all k
					if i - snp_len >= 0 {
						if snp_values[k][0] != '.' {
							temp_dis = D[i - snp_len][j - 1] + Cost(s[i - snp_len : i], snp_values[k])
						} else {
    						temp_dis = D[i][j - 1]
						}
	    				if D[i][j] > temp_dis {
	    					D[i][j] = temp_dis
							min_index = k
	    				}
					}
	   			}
				if (D[i][j] < INF) {
					T[pos + j - 1] = min_index
				}
			}
		}
	}
	if D[m][n] >= INF {
		return 0, INF, m, n, S, make(map[int]int)
	}
    return d, D[m][n], m, n, S, T
}


//-------------------------------------------------------------------------------------------------
// BackwardTraceBack constructs alignment between s and t based on the results from BackwardDistanceMulti.
// 	s is a read.
// 	t is part of a multi-genome.
// The reads include standard bases, the multi-genomes include standard bases and "*" characters.
//-------------------------------------------------------------------------------------------------
func BackwardTraceBack(s, t []byte, m, n int, S map[int][]byte, T map[int]int, pos int) map[int][]byte {

	var snp_values [][]byte
	var is_snp bool
	var snp_len int
	var snp_calling = make(map[int][]byte)

	for k, v := range S {
		snp_calling[k] = v
	}
	var i, j int = m, n
	for  i > 0 || j > 0 {
		snp_values, is_snp = SNP_PROFILE[pos + j - 1]
		if i > 0 && j > 0 {
		  	if !is_snp {
		  		i, j = i - 1, j - 1
		  	} else {
				if snp_values[T[pos + j - 1]][0] != '.' {
			  		snp_len = len(snp_values[T[pos + j - 1]])
			  	} else {
			  		snp_len = 0
			  	}
		  		snp_calling[pos + j - 1] = s[i - snp_len : i]
		  		i, j = i - snp_len, j - 1
		  	}
		} else if i == 0 {
			j = j - 1;
		} else if j == 0 {
			i = i - 1;
		}
	}
	return snp_calling
}


//-------------------------------------------------------------------------------------------------
// Calculate the distance between s and t in forward direction.
// 	s is a read.
// 	t is part of a multi-genome.
// The reads include standard bases, the multi-genomes include standard bases and "*" characters.
//-------------------------------------------------------------------------------------------------
func ForwardDistanceMulti(s, t []byte, pos int) (int, int, int, int, map[int][]byte, map[int]int) {

    M, N := len(s), len(t)

    var cost int
    var d, min_d int
    var snp_len int
    var snp_values [][]byte
    var is_snp, is_same_len_snp bool
	var S = make(map[int][]byte)

    var i, j, k int

    d = 0
    m, n := M, N
    for m > 0 && n > 0 {
		snp_values, is_snp = SNP_PROFILE[pos + (N - 1) - (n - 1)]
		snp_len, is_same_len_snp = SAME_LEN_SNP[pos + (N - 1) - (n - 1)]
    	if !is_snp {
        	if s[(M - 1) - (m - 1)] != t[(N - 1) - (n - 1)] {
        		d++
        	}
    		m--
    		n--
    	} else if is_same_len_snp {
			min_d = 1000*INF
    		for i = 0; i < len(snp_values); i++ {
    			cost = Cost(s[M - m: M - (m - snp_len)], snp_values[i])
    			if min_d > cost {
    				min_d = cost
    			}
    		}
	  		if min_d == INF {
	  			return INF, 0, m, n, make(map[int][]byte), make(map[int]int)
	  		}
	  		S[pos + (N - 1) - (n - 1)] = s[M - m: M - (m - snp_len)]
    		d += min_d
    		m -= snp_len
    		n--
    	} else {
    		break
    	}
    }

    D := make([][]int, m + 1)
    T := make(map[int]int)
    for i = 0; i <= m; i++ {
		D[i] = make([]int, n + 1)
    }
    D[0][0] = 0
    for i = 1; i <= m; i++ {
		D[i][0] = INF
    }
    for j = 1; j <= n; j++ {
		D[0][j] = 0
    }
	var temp_dis, min_index int
    for i = 1; i <= m; i++ {
        for j = 1; j <= n; j++ {
			snp_values, is_snp = SNP_PROFILE[pos + (N - 1) - (j - 1)]
		    if !is_snp {
				if s[(M - 1) - (i - 1)] != t[(N - 1) - (j - 1)] {
						D[i][j] = D[i - 1][j - 1] + 1
				} else {
						D[i][j] = D[i - 1][j - 1]
				}
		   	} else {
				D[i][j] = 1000*INF //1000*INF is a value for testing, will change to a better solution later
				min_index = 0
				for k = 0; k < len(snp_values); k++ {
					snp_len = len(snp_values[k])
					//One possible case: i - snp_len < 0 for all k
					if i - snp_len >= 0 {
						if snp_values[k][0] != '.' {
							temp_dis = D[i - snp_len][j - 1] + Cost(s[M - i : M - (i - snp_len)], snp_values[k])
						} else {
	    						temp_dis = D[i][j - 1]
						}
	    				if D[i][j] > temp_dis {
	    					D[i][j] = temp_dis
							min_index = k
		   				}
					}
		   		}
				if (D[i][j] < INF) {
					T[pos + (N - 1) - (j - 1)] = min_index
				}
			}
		}
    }
	if D[m][n] >= INF {
		return 0, INF, m, n, S, make(map[int]int)
	}
    return d, D[m][n], m, n, S, T
}


//-------------------------------------------------------------------------------------------------
// ForwardTraceBack constructs alignment between s and t based on the results from ForwardDistanceMulti.
// 	s is a read.
// 	t is part of a multi-genome.
// The reads include standard bases, the multi-genomes include standard bases and "*" characters.
//-------------------------------------------------------------------------------------------------
func ForwardTraceBack(s, t []byte, m, n int, S map[int][]byte, T map[int]int, pos int) map[int][]byte {

	var is_snp bool
	var snp_values [][]byte
	var snp_len int
	var snp_calling = make(map[int][]byte)

	for k, v := range S {
		snp_calling[k] = v
	}
	var i, j int = m, n
	var M, N int = len(s), len(t)
	for  i > 0 || j > 0 {
		snp_values, is_snp = SNP_PROFILE[pos + (N - 1) - (j - 1)]
		if i > 0 && j > 0 {
		  	if !is_snp {
		  		i, j = i - 1, j - 1
		  	} else {
				if snp_values[T[pos + (N - 1) - (j - 1)]][0] != '.' {
			  		snp_len = len(snp_values[T[pos + (N - 1) - (j - 1)]])
			  	} else {
			  		snp_len = 0
			  	}
		  		snp_calling[pos + (N - 1) - (j - 1)] = s[M - i : M - (i - snp_len)]
		  		i, j = i - snp_len, j - 1
		  	}
		} else if i == 0 {
			j = j - 1;
		} else if j == 0 {
			i = i - 1;
		}
	}
	return snp_calling
}
