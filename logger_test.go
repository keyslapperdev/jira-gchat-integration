package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Accurately logs to one writer", func(t *testing.T) {
		wantedMsg := "Test text to log"

		mLog1 := MockLog1{&[]byte{}}

		logger := StartLogger(mLog1)
		logger.Print(wantedMsg)

		assert.True(t, strings.Contains(string(*mLog1.message), string(wantedMsg)), "Correct message not logged")
	})

	t.Run("Accurately logs to two writers", func(t *testing.T) {
		wantedMsg := "Test text to log"

		mLog1 := MockLog1{&[]byte{}}
		mLog2 := MockLog2{&[]byte{}}

		logger := StartLogger(mLog1, mLog2)
		logger.Print(wantedMsg)

		assert.True(t, strings.Contains(string(*mLog1.message), wantedMsg), "Correct message not logged in first logger")
		assert.True(t, strings.Contains(string(*mLog2.message), wantedMsg), "Correct message not logged in second logger")
	})
}

type MockLog1 struct{ message *[]byte }

func (ml1 MockLog1) Write(p []byte) (int, error) {
	*ml1.message = p
	return len(p), nil
}

type MockLog2 struct{ message *[]byte }

func (ml2 MockLog2) Write(p []byte) (int, error) {
	*ml2.message = p
	return len(p), nil
}
