package rank

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// InitialWithBadLimit is a error to inidicate bad initialization
	InitialWithBadLimit = errors.New("bad inititalization param. limit <= 0")
	// OverflowError is an error to indicate overflow
	OverflowError = fmt.Errorf("Rank overflow. Maximun value: %s", RANK_MAX)
	// UnderflowError is an error to indicate underflow
	UnderflowError = fmt.Errorf("Rank underflow. Minimum value: %s", RANK_MIN)
	// InvalidRankError is an error to indicate invalid rank
	InvalidRankError = fmt.Errorf("Invalid digit input. Allowed digigts: %s", DIGITS)
	// InvalidRatioError is an error to indicate invalid ratio
	InvalidRatioError = errors.New("invalid ratio. ratio should within -100 to 100")
)

// Rank is an interface for rank module
type Rank interface {
	NewRanks(length uint64) []string
	NewRanksBetween(start, end string, length uint64) []string
	Prev(curr string) (string, error)
	Next(curr string) (string, error)
	Insert(prev, next string) string
	Equal(a, b string) bool
	Less(a, b string) bool
	Greater(a, b string) bool
}

const (
	// using 0~9 + A~Z as digit of rank

	// RANK_MIN is a minimal rank
	RANK_MIN = "0"
	// RANK_MAX is a maximal rank
	RANK_MAX = "Z"
	// DIGITS includes all possible digit
	DIGITS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type rankImpl struct {
	Limit int
}

// NewRank initialize a rank impl
func NewRank(limit int) (*rankImpl, error) {
	if limit == 0 {
		return nil, InitialWithBadLimit
	}
	return &rankImpl{Limit: limit}, nil
}

// NewRanks return a list of rank
func (r *rankImpl) NewRanks(length uint64) []string {
	return r.NewRanksBetween(RANK_MIN, RANK_MAX, length)
}

// NewRanksBetween get a list of new rank between 2 string
func (r *rankImpl) NewRanksBetween(start, end string, length uint64) []string {
	return r.recursiveInsertRankBetween(0, length, start, end)
}

// Prev is a method to get a rank just one step smaller than input string
func (r *rankImpl) Prev(curr string) (string, error) {
	if !r.isValidRank(curr) {
		return "", InvalidRankError
	}
	if r.Greater(curr, RANK_MAX) {
		return "", OverflowError
	}
	if r.Less(curr, RANK_MIN) || curr == "" {
		return "", UnderflowError
	}
	curr = strings.TrimRight(curr, RANK_MIN)
	if len(curr) < r.Limit {
		newIdx := strings.Index(DIGITS, string(curr[len(curr)-1])) - 1
		ret := curr[:len(curr)-1] + string(DIGITS[newIdx])
		for i := 0; i < r.Limit-len(curr); i++ {
			ret += RANK_MAX
		}
		return strings.TrimRight(ret, RANK_MIN), nil
	}
	lastNewIdx := strings.Index(DIGITS, string(curr[len(curr)-1])) - 1
	ret := strings.TrimRight(curr[:len(curr)-1]+string(DIGITS[lastNewIdx]), RANK_MIN)
	return ret, nil
}

// Next is a method to get a rank just one step bigger than input string
func (r *rankImpl) Next(curr string) (string, error) {
	if !r.isValidRank(curr) {
		return "", InvalidRankError
	}
	if r.Greater(curr, RANK_MAX) {
		return "", OverflowError
	}
	if r.Less(curr, RANK_MIN) || curr == "" {
		return "", UnderflowError
	}
	curr = strings.TrimRight(curr, RANK_MIN)
	if len(curr) < r.Limit {
		ret := curr
		for i := 0; i < r.Limit-len(curr)-1; i++ {
			ret += RANK_MIN
		}
		return ret + string(DIGITS[1]), nil
	}
	carry := 1
	ret := ""
	for i := len(curr) - 1; i >= 0; i-- {
		newIdx := strings.Index(DIGITS, string(curr[i])) + carry
		carry = newIdx / len(DIGITS)
		newIdx = newIdx % len(DIGITS)
		ret = string(DIGITS[newIdx]) + ret
	}
	return strings.TrimRight(ret, RANK_MIN), nil
}

// Insert get a new rank between 2 string
func (r *rankImpl) Insert(prev, next string) string {
	if r.Equal(prev, next) {
		if len(next) < len(prev) {
			return next
		}
		return prev
	}
	if r.Greater(prev, next) {
		tmp := prev
		prev = next
		next = tmp
	}
	// rank := (prev + next)/2
	lenP := len(prev)
	lenN := len(next)
	rank := ""
	for i := 0; i < lenP && i < lenN; i++ {
		if prev[i] != next[i] {
			digit := r.average(prev[i:], next[i:])
			rank = prev[0:i] + digit
			break
		}
	}
	if rank != "" {
		return rank
	}
	// prev is a substring of next (opposite case should not happen here)
	digit := r.average(RANK_MIN, next[lenP:])
	rank = prev + digit
	return strings.TrimRight(rank, RANK_MIN)
}

// Equal is a method to check if 2 rank is equal
func (r *rankImpl) Equal(a, b string) bool {
	return strings.TrimRight(a, RANK_MIN) == strings.TrimRight(b, RANK_MIN)
}

// Less is a method to check if a < b
func (r *rankImpl) Less(a, b string) bool {
	if r.Equal(a, b) {
		return false
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return true
		} else if b[i] < a[i] {
			return false
		}
	}
	if len(a) < len(b) {
		return true
	}
	return false
}

// Greater is a method to check if a > b
func (r *rankImpl) Greater(a, b string) bool {
	return !r.Equal(a, b) && !r.Less(a, b)
}

func (r *rankImpl) recursiveInsertRankBetween(leftIndex, rightIndex uint64, prevRank, nextRank string) []string {
	if leftIndex == rightIndex {
		return []string{}
	}

	index := (leftIndex + rightIndex) / 2
	rank := r.Insert(prevRank, nextRank)

	left := r.recursiveInsertRankBetween(leftIndex, index, prevRank, rank)
	right := r.recursiveInsertRankBetween(index+1, rightIndex, rank, nextRank)

	left = append(left, rank)
	return append(left, right...)
}

func (r *rankImpl) average(prev, next string) string {
	lengthMax := len(prev)
	if len(next) > lengthMax {
		lengthMax = len(next)
	}

	avg := ""
	base := len(DIGITS)
	remain := 0
	for i := 0; i < lengthMax || remain != 0; i++ {
		p := 0
		n := 0
		if i < len(prev) {
			p = strings.Index(DIGITS, string(prev[i]))
		}
		if i < len(next) {
			n = strings.Index(DIGITS, string(next[i]))
		}

		sum := p + n + remain
		curr := sum / 2
		if curr >= base {
			plusback := 1 // can only be 1. (p + n)/2 <= base and remain <= 1  -> (p + n) + remain / 2 < 2 * based
			firstLetterInBetween := ""
			for i := 1; len(avg)-i >= 0; i++ {
				lastBitIndex := strings.Index(DIGITS, string(avg[len(avg)-i]))
				if lastBitIndex+plusback >= base {
					firstLetterInBetween += string(DIGITS[0])
					continue
				}
				lastBit := DIGITS[lastBitIndex+plusback]
				avg = avg[0:len(avg)-i] + string(lastBit) + firstLetterInBetween + avg[len(avg)-i+1:]
				break
			}
			// impossible to have plusback > 0 after the loop.
			curr = curr % base
		}
		avg += string(DIGITS[curr])
		remain = (sum % 2) * base
	}
	return strings.TrimRight(avg, RANK_MIN)
}

func (r *rankImpl) isValidRank(ranks ...string) bool {
	for _, s := range ranks {
		for i := 0; i < len(s); i++ {
			if idx := strings.Index(DIGITS, string(s[i])); idx == -1 {
				return false
			}
		}
	}
	return true
}
