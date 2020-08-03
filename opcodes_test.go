package main

import (
	"encoding/binary"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	reg := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected uint16
	}{
		{0, []uint16{0, register0, 1, 1}, reg, 2},
		{0, []uint16{0, register0, 1, 0}, reg, 1},
	}

	for _, test := range tests {
		index := add(test.index, &test.memory, &test.reg)

		// takes 3 args so index is incremented by 3
		if index != 3 {
			t.Error("Got:", index, "Expected:", 2)
		}

		result := test.reg.Get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestEq(t *testing.T) {
	reg := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected uint16
	}{
		{0, []uint16{0, register0, 1, 1}, reg, 1},
		{0, []uint16{0, register0, 1, 0}, reg, 0},
	}

	for _, test := range tests {
		index := eq(test.index, &test.memory, &test.reg)

		// takes 3 args so index is incremented by 3
		if index != 3 {
			t.Error("Got:", index, "Expected:", 2)
		}

		result := test.reg.Get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestGetNextValueShiftIndex(t *testing.T) {
	tests := []struct {
		index         int
		memory        []uint16
		reg           registers
		expectedIndex int
		expectedValue uint16
	}{
		{0, []uint16{0, 1}, registers{}, 1, 1},
		{1, []uint16{0, 1, 2}, registers{}, 2, 2},
	}

	for _, test := range tests {
		resultIndex, resultValue := getNextValueShiftIndex(
			test.index, &test.memory, &test.reg,
		)

		if resultIndex != test.expectedIndex {
			t.Error("Got:", resultIndex, "Expected:", test.expectedIndex)
		}
		if resultValue != test.expectedValue {
			t.Error("Got:", resultValue, "Expected:", test.expectedValue)
		}
	}
}

// Jump will get the jump location from the next memory location (like all optimize
// functions do), but returns it decremented (since the VM loop will immediately
// loop to the next iteration).
func TestJump(t *testing.T) {
	// build a memory data set to verify register values are accessed if specified
	fullMem := []uint16{}
	for i := 0; i <= registerEnd; i++ {
		fullMem = append(fullMem, uint16(i))
	}

	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected int
	}{
		{0, []uint16{10, 11, 12}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 10},
		{1, []uint16{10, 11, 12}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 11},
		{0, []uint16{0, registerStart, 12}, registers{200, 1, 2, 3, 4, 5, 6, 7}, 199},
	}

	for _, test := range tests {
		result := jump(test.index, &test.memory, &test.reg)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

// See TestJump docstring for why expected return index is 1 minus the index
// position.
func TestJumpFalse(t *testing.T) {
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected int
	}{
		// i, a, b, a == 0, jump to b, return index of b
		{0, []uint16{0, 0, 3, 4, 5}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 2},
		// i, a, b, a != 0, no jump to b, next index is 2
		{0, []uint16{0, 1, 0, 0, 0}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 2},
	}

	for _, test := range tests {
		result := jumpFalse(test.index, &test.memory, &test.reg)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

// See TestJump docstring for why expected return index is 1 minus the index
// position.
func TestJumpTrue(t *testing.T) {
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected int
	}{
		// i, a, b, a >= 0, jump to b, return index of b
		{0, []uint16{0, 1, 3, 4, 5}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 2},
		// i, a, b, a !> 0, no jump to b -> 2
		{0, []uint16{0, 0, 0, 0, 0}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 2},
	}

	for _, test := range tests {
		result := jumpTrue(test.index, &test.memory, &test.reg)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestNoop(t *testing.T) {
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected int
	}{
		{0, []uint16{0, 1, 2}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 0},
	}

	for _, test := range tests {
		result := noop(test.index, &test.memory, &test.reg)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestOut(t *testing.T) {
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected uint16
	}{
		{0, []uint16{0, 65}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 65},
	}

	for _, test := range tests {
		// Redirect stdout via Pipe()
		r, w, _ := os.Pipe()
		backupStdout := os.Stdout
		os.Stdout = w

		out(test.index, &test.memory, &test.reg)

		buf := make([]byte, 2)
		_, err := r.Read(buf)
		if err != nil {
			t.Error(err)
		}

		// Convert buffer and restore stdout
		result := binary.LittleEndian.Uint16(buf)
		w.Close()
		os.Stdout = backupStdout

		if uint16(result) != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestSet(t *testing.T) {
	reg := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		index    int
		memory   []uint16
		reg      registers
		expected uint16
	}{
		{0, []uint16{0, register0, 42}, reg, 42},
	}

	for _, test := range tests {
		register := test.memory[1]
		index := set(test.index, &test.memory, &test.reg)

		if index != 2 {
			t.Error("Got:", index, "Expected:", 2)
		}

		result := test.reg.Get(register)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Register:", register)
		}
	}
}