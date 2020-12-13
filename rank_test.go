package rank

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	length = 7

	charMaxLength = 128

	nextPrevRequirement = 100000
)

var ANSWER_SET_7RANKS = []string{"4DI", "8R", "D4I", "HI", "LVI", "Q9", "UMI"}
var ANSWER_SET_7RANKS_SUB = []string{"1", "2", "3", "4", "5", "6", "7"}
var ANSWER_SET_EQUAL = []bool{true, false, false}
var ANSWER_SET_LESS = []bool{false, true, false}
var ANSWER_SET_GREATER = []bool{false, false, true}

func TestNewRanks(t *testing.T) {
	assert := assert.New(t)

	r, _ := NewRank(10)

	ranksA := r.NewRanks(length)
	ranksB := r.NewRanks(length)
	assert.Equal(length, len(ranksA))
	assert.Equal(length, len(ranksB))
	assert.Equal(ranksA, ranksB)

	if length == 7 {
		assert.Equal(ANSWER_SET_7RANKS, ranksA)
	}

	ranks := r.NewRanksBetween("0", "8", 7)
	assert.Equal(ANSWER_SET_7RANKS_SUB, ranks)
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	r, _ := NewRank(10)

	assert.Equal("HI", r.Insert(RANK_MIN, RANK_MAX))
	assert.Equal("B", r.Insert("A", "C"))

	assert.Equal("00I", r.Insert("0", "01"))
	assert.Equal("01", r.Insert("0", "02"))
	assert.Equal("010I", r.Insert("011", "010"))
	assert.Equal("01Z", r.Insert("02", "01Y"))

	assert.Equal("9", r.Insert("7", "B"))
	assert.Equal("A", r.Insert("9I", "AI"))
	assert.Equal("A", r.Insert("9", "B"))

	assert.Equal("B", r.Insert("B", "B"))
	assert.Equal("B", r.Insert("B", "B00"))
	assert.Equal("B", r.Insert("B000", "B"))
}

func TestNext(t *testing.T) {
	assert := assert.New(t)
	r, _ := NewRank(10)

	cases := []struct {
		scenario string
		input    string
		output   string
		err      error
	}{
		{
			scenario: "begin",
			input:    "A",
			output:   "A000000001",
			err:      nil,
		},
		{
			scenario: "begin next",
			input:    "A1",
			output:   "A100000001",
			err:      nil,
		},
		{
			scenario: "end with Z but not reach limit",
			input:    "AZ",
			output:   "AZ00000001",
			err:      nil,
		},
		{
			scenario: "end with Z and reach limit",
			input:    "BZAZZZZZZZ",
			output:   "BZB",
			err:      nil,
		},
		{
			scenario: "end with Z and reach limit. (continuous Z)",
			input:  "BZZZZZZZZZ",
			output: "C",
			err:    nil,
		},
		{
			scenario: "longer than limit end with z",
			input:    "BZAZZZZZZZAZZZ",
			output:   "BZAZZZZZZZB",
			err:      nil,
		},
		{
			scenario: "invalid input (invalid char)",
			input:    "aa",
			output:   "",
			err:      InvalidRankError,
		},
		{
			scenario: "invalid input (overflow)",
			input:    "ZZ",
			output:   "",
			err:      OverflowError,
		},
		{
			scenario: "invalid input (underflow)",
			input:    "",
			output:   "",
			err:      UnderflowError,
		},
	}

	for _, c := range cases {
		rank, err := r.Next(c.input)
		assert.Equal(c.output, rank, c.scenario)
		assert.Equal(c.err, err, c.scenario)
	}

	tail := r.NewRanks(length)[length-1]
	pre := tail
	var (
		maxStep int
		err     error
	)

	for maxStep = 0; len(tail) <= charMaxLength && maxStep <= nextPrevRequirement; maxStep++ {
		tail, err = r.Next(tail)
		if err != nil {
			t.Errorf("NEXT stress test fail: err: %v", err)
			break
		}
		if r.Greater(pre, tail) {
			t.Errorf("pre: %s greater than next: %s", pre, tail)
		}
		pre = tail
	}
	if maxStep <= nextPrevRequirement {
		t.Errorf("NEXT stress test fail: maxStep: %v", maxStep)
	}
}

func TestPrev(t *testing.T) {
	assert := assert.New(t)
	r, _ := NewRank(10)
	cases := []struct {
		scenario string
		input    string
		output   string
		err      error
	}{
		{
			scenario: "normal",
			input:    "CAD",
			output:   "CACZZZZZZZ",
			err:      nil,
		},
		{
			scenario: "normal next",
			input:    "CAC",
			output:   "CABZZZZZZZ",
			err:      nil,
		},
		{
			scenario: "only 1 letter",
			input:    "A",
			output:   "9ZZZZZZZZZ",
			err:      nil,
		},
		{
			scenario: "end with Z and longer than limit. (continuous Z)",
			input:  "9ZZZZZZZZZZZZZ",
			output: "9ZZZZZZZZZZZZY",
			err:    nil,
		},
		{
			scenario: "longer than limit",
			input:    "9ZZZZZZZZZZZZZ",
			output:   "9ZZZZZZZZZZZZY",
			err:      nil,
		},
		{
			scenario: "longer than limit, getting shorter",
			input:    "9ZZZZZZZZZZZZ1",
			output:   "9ZZZZZZZZZZZZ",
			err:      nil,
		},
		{
			scenario: "a lot of 0 and with only 1 letter in the end",
			input:    "90000000000001",
			output:   "9",
			err:      nil,
		},
		{
			scenario: "a lot of 0 and with only 1 letter in the end(shorter than limit)",
			input:    "900001",
			output:   "900000ZZZZ",
			err:      nil,
		},
		{
			scenario: "invalid input (invalid rank)",
			input:    "aa",
			output:   "",
			err:      InvalidRankError,
		},
		{
			scenario: "invalid input (overflow)",
			input:    "ZZ",
			output:   "",
			err:      OverflowError,
		},
		{
			scenario: "invalid input (underflow)",
			input:    "",
			output:   "",
			err:      UnderflowError,
		},
	}

	for _, c := range cases {
		rank, err := r.Prev(c.input)
		assert.Equal(c.output, rank, c.scenario)
		assert.Equal(c.err, err, c.scenario)
	}

	head := r.NewRanks(length)[0]
	pre := head
	var (
		maxStep int
		err     error
	)

	for maxStep = 0; len(head) <= charMaxLength && maxStep <= nextPrevRequirement; maxStep++ {
		head, err = r.Prev(head)
		if err != nil {
			t.Errorf("NEXT stress test fail: err: %v", err)
			break
		}
		if r.Less(pre, head) {
			t.Errorf("pre: %s less than prev: %s", pre, head)
		}
		pre = head
	}
	if maxStep <= nextPrevRequirement {
		t.Errorf("Prev stress test fail: maxStep: %v", maxStep)
	}
}

func TestRankComparison(t *testing.T) {
	assert := assert.New(t)
	r, _ := NewRank(10)
	//ranks := NewRanks(length)
	ranks := ANSWER_SET_7RANKS

	for i := range ranks {
		for j := range ranks {
			result := ANSWER_SET_EQUAL
			if i > j {
				result = ANSWER_SET_GREATER
			} else if i < j {
				result = ANSWER_SET_LESS
			}
			assert.Equal(r.Equal(ranks[i], ranks[j]), result[0], fmt.Sprintf("%s not equal to %s", ranks[i], ranks[j]))
			assert.Equal(r.Less(ranks[i], ranks[j]), result[1], fmt.Sprintf("%s not less than %s", ranks[i], ranks[j]))
			assert.Equal(r.Greater(ranks[i], ranks[j]), result[2], fmt.Sprintf("%s not greater than %s", ranks[i], ranks[j]))
		}
	}
}
