package glid

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	machineIDBits = 10
	sequenceBits  = 12

	maxMachineID = int64(1<<machineIDBits) - 1
	sequenceMask = int64(1<<sequenceBits) - 1

	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits

	epochMilli = int64(1704067200000)
)

// Generator 生成 Snowflake 风格的递增 ID。
type Generator struct {
	mu            sync.Mutex
	machineID     int64
	lastTimestamp int64
	sequence      int64
}

// NewGenerator 创建指定机器 ID 的生成器，machineID 必须在 0..1023 之间。
func NewGenerator(machineID int64) (*Generator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID must be between 0 and 1023")
	}
	return &Generator{machineID: machineID}, nil
}

// Next 返回下一个 int64 ID。
func (g *Generator) Next() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := currentMilli()
	if now < g.lastTimestamp {
		return 0, errors.New("system clock moved backward")
	}

	if now == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			now = waitNextMilli(g.lastTimestamp)
		}
	} else {
		g.sequence = 0
	}
	g.lastTimestamp = now

	id := ((now - epochMilli) << timestampShift) | (g.machineID << machineIDShift) | g.sequence
	return id, nil
}

// NextString 返回下一个十进制字符串 ID。
func (g *Generator) NextString() (string, error) {
	id, err := g.Next()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func currentMilli() int64 {
	return time.Now().UnixMilli()
}

func waitNextMilli(lastTimestamp int64) int64 {
	now := currentMilli()
	for now <= lastTimestamp {
		now = currentMilli()
	}
	return now
}
