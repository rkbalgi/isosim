package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"isosim/internal/services/v0/data"
	"time"
)

var timeFormat = "2006-01-02T15"

// DbMessage is an entry of a request/response that will be persisted to
// storage
type DbMessage struct {
	ID     string `json:"id"`
	SpecID int    `json:"spec_id"`
	MsgID  int    `json:"msg_id"`

	RequestTS  int64 `json:"request_ts"`
	ResponseTS int64 `json:"response_ts"`

	RequestMsg        string                  `json:"request_msg"`
	ParsedRequestMsg  []data.JsonFieldDataRep `json:"parsed_request_msg"`
	ResponseMsg       string                  `json:"response_msg"`
	ParsedResponseMsg []data.JsonFieldDataRep `json:"parsed_response_msg"`
}

func Write(dbMsg DbMessage) error {

	var err error

	if dbMsg.MsgID == 0 || dbMsg.SpecID == 0 {
		return errors.New("isosim: Invalid SpecID/MsgID")
	}

	uniqueID, err := uuid.NewUUID()
	if err != nil {
		log.Warn("Failed to generate UUID for DbMessage", err)
	} else {
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

	bkt, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%d_%d", dbMsg.SpecID, dbMsg.MsgID)))
	if err != nil {
		return err
	}

	//hourly buckets
	tBkt, err := bkt.CreateBucketIfNotExists([]byte(time.Now().Format(timeFormat)))
	if err != nil {
		return err
	}
	if err = tBkt.Put([]byte(uniqueID.String()), jsonData); err != nil {
		return err
	}
	log.Println("Wrote ..", jsonData, uniqueID.String(), fmt.Sprintf("%d_%d", dbMsg.SpecID, dbMsg.MsgID))
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

// ReadLast reads last 'n' messages for spec and msg
func ReadLast(specID int, msgID int, n int) ([]string, error) {

	tx, err := bdb.Begin(false)
	if err != nil {
		return nil, err
	}
	bktName := fmt.Sprintf("%d_%d", specID, msgID)
	bkt := tx.Bucket([]byte(bktName))
	if bkt == nil {
		log.Println("No bucket for spec/msg")
		return nil, nil
	}

	res := make([]string, 0)
	retrieved := 0
	now := time.Now()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	for {
		//hourly buckets
		bktName := now.Format(timeFormat)
		log.Println("Reading from .. " + bktName)
		tBkt := bkt.Bucket([]byte(bktName))
		if tBkt != nil {
			c := tBkt.Cursor()
			k, v := c.Last()
			if k == nil || v == nil {
				continue
			}
			for retrieved < n {
				res = append(res, string(v))
				k, v = c.Prev()
				if k == nil || v == nil {
					continue
				} else {
					retrieved++
					if retrieved == n {
						return res, nil
					}
				}

			}

		}
		select {
		case <-ctx.Done():
			return res, nil
		default:
			break
		}
		now = now.Add(-1 * time.Hour)
	}

}
