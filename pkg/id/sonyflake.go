package id

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/sonyflake/v2"
)

type Sonyflake struct {
	ops   SonyflakeOptions
	sf    *sonyflake.Sonyflake
	Error error
}

// NewSonyflake can get a unique code by id(You need to ensure that id is unique).
func NewSonyflake(options ...func(*SonyflakeOptions)) *Sonyflake {
	ops := getSonyflakeOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	sf := &Sonyflake{
		ops: *ops,
	}
	st := sonyflake.Settings{
		StartTime: ops.startTime,
	}
	if ops.machineId > 0 {
		st.MachineID = func() (int, error) {
			return ops.machineId, nil
		}
	}

	ins, err := sonyflake.New(st)
	if ins == nil || err != nil {
		sf.Error = fmt.Errorf("create snoyflake failed")
	}

	if _, err = ins.NextID(); err != nil {
		sf.Error = fmt.Errorf("invalid start time: %w", err)
	}

	sf.sf = ins

	return sf
}

func (s *Sonyflake) NextID() (uint64, error) {
	if s.Error != nil {
		return 0, s.Error
	}
	id, err := s.sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("sonyflake get id failed: %w", err)
	}
	return uint64(id), nil
}

func (s *Sonyflake) Id(ctx context.Context) (id uint64) {
	if s.Error != nil {
		return 0
	}

	sleep := 1
	for {
		id, err := s.NextID()
		if err == nil {
			return id
		}

		select {
		case <-ctx.Done():
			return 0
		case <-time.After(time.Duration(sleep) * time.Millisecond):
			sleep *= 2
		}
	}
}
