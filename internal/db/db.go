package db

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"isosim/internal/services/data"
	"time"
)

const timeFormat = "2006-01-02T15"

// DbMessage is an entry of a request/response that will be persisted to
// storage
type DbMessage struct {
	ID       string `json:"id"`
	SpecID   int    `json:"spec_id"`
	MsgID    int    `json:"msg_id"`
	HostAddr string `json:"host_addr"`

	LogTS      string `json:"log_ts"`
	RequestTS  int64  `json:"request_ts"`
	ResponseTS int64  `json:"response_ts"`

	RequestMsg        string                  `json:"request_msg"`
	ParsedRequestMsg  []data.JsonFieldDataRep `json:"parsed_request_msg"`
	ResponseMsg       string                  `json:"response_msg"`
	ParsedResponseMsg []data.JsonFieldDataRep `json:"parsed_response_msg"`
}

// Write writes a message into bolt (into a hourly bucket)
func Write(dbMsg DbMessage) error {

	var err error

	if dbMsg.MsgID == 0 || dbMsg.SpecID == 0 {
		return errors.New("isosim: Invalid SpecID/MsgID")
	}

	uniqueID, err := uuid.NewUUID()
	if err != nil {
		log.Warn("Failed to generate UUID for DbMessage", err)
	} else {
		dbMsg.LogTS = time.Now().Format(time.RFC3339)
		dbMsg.ID = uniqueID.String()
	}
	var jsonData []byte

	if jsonData, err = json.Marshal(dbMsg); err != nil {
		return err
	}

	tx, err := bdb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%d_%d", dbMsg.SpecID, dbMsg.MsgID)))
	if err != nil {
		return err
	}

	//hourly buckets
	tBkt, err := bkt.CreateBucketIfNotExists([]byte(time.Now().Format(timeFormat)))
	if err != nil {
		return err
	}
	bSeq := make([]byte, 8)
	binary.BigEndian.PutUint64(bSeq, tBkt.Sequence())
	if err = tBkt.Put(bSeq, jsonData); err != nil {
		return err
	}
	_, err = tBkt.NextSequence()
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// ReadLast reads last 'n' messages for spec and msg
func ReadLast(specID int, msgID int, n int) ([]string, error) {

	res := make([]string, 0)

	err := bdb.View(func(tx *bolt.Tx) error {

		bktName := fmt.Sprintf("%d_%d", specID, msgID)
		bkt := tx.Bucket([]byte(bktName))
		if bkt == nil {
			log.Debugf("No bucket for spec/msg - %d:%d", specID, msgID)
			return nil
		}

		retrieved := 0
		now := time.Now()
		ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFunc()
		for {
			//hourly buckets
			bktName := now.Format(timeFormat)
			tBkt := bkt.Bucket([]byte(bktName))
			if tBkt != nil {
				//start from the last on the latest bucket
				c := tBkt.Cursor()
				k, v := c.Last()

				if k == nil || v == nil {
					now = now.Add(-1 * time.Hour)
					continue
				}
				for len(res) < n {
					res = append(res, string(v))
					retrieved++
					if len(res) == n {
						return nil
					}
					k, v = c.Prev()
					if k == nil || v == nil {
						// nothing more in this hour,
						// break out of this loop
						goto PREV_HOUR
					}
				}

			}
		PREV_HOUR:
			// we cannot keep looking endlessly
			select {
			case <-ctx.Done():
				return nil
			default:
				break
			}
			now = now.Add(-1 * time.Hour)
		}
	})

	return res, err

}
