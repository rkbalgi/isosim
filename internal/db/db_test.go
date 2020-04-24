package db

import (
	"testing"
	"time"
)

func Test_ReadWriteToBold(t *testing.T) {
	if err := Init("."); err != nil {
		t.Fatal(err)
	}

	dbMsg := DbMessage{
		ID:                "",
		SpecID:            100,
		MsgID:             1,
		RequestTS:         436466364678,
		ResponseTS:        767366436647,
		RequestMsg:        "110100101010010101010",
		ParsedRequestMsg:  nil,
		ResponseMsg:       "11110........",
		ParsedResponseMsg: nil,
	}
	for i := 0; i < 10; i++ {
		if err := Write(dbMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(2 * time.Second)
	}

	entries, err := ReadLast(100, 1, 5)
	if entries == nil {
		t.Fatal("No entries found!")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(entries)
}

func Test_Read(t *testing.T) {

	if err := Init("."); err != nil {
		t.Fatal(err)
	}

	entries, err := ReadLast(100, 1, 20)
	if entries == nil {
		t.Fatal("No entries found!")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(entries)

}
