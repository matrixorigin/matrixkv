package kv

import (
	"fmt"
	"log"

	"github.com/matrixorigin/matrixcube/storage"
	"github.com/matrixorigin/matrixcube/util"
	"github.com/matrixorigin/tinykv/pkg/metadata"
)

var (
	// OK response
	OK = []byte("OK")
)

type simpleKVExecutor struct {
	kv storage.KVStorage
}

var _ storage.Executor = (*simpleKVExecutor)(nil)

// NewSimpleKVExecutor returns a simple kv executor that supports set/get/delete
// commands.
func NewSimpleKVExecutor(kv storage.KVStorage) storage.Executor {
	return &simpleKVExecutor{kv: kv}
}

func (ce *simpleKVExecutor) UpdateWriteBatch(ctx storage.WriteContext) error {
	writtenBytes := uint64(0)
	r := ctx.WriteBatch()
	wb := r.(util.WriteBatch)
	batch := ctx.Batch()
	requests := batch.Requests
	for j := range requests {
		switch requests[j].CmdType {
		case metadata.SetType:
			wb.Set(requests[j].Key, requests[j].Cmd)
			log.Printf("write %s, %+v", requests[j].Key, requests[j].Cmd)
			writtenBytes += uint64(len(requests[j].Key) + len(requests[j].Cmd))
			ctx.AppendResponse(OK)
		case metadata.DeleteType:
			wb.Delete(requests[j].Key)
			writtenBytes += uint64(len(requests[j].Key))
			ctx.AppendResponse(OK)
		default:
			panic(fmt.Errorf("invalid write cmd %d", requests[j].CmdType))
		}
	}

	writtenBytes += uint64(16)
	ctx.SetDiffBytes(int64(writtenBytes))
	ctx.SetWrittenBytes(writtenBytes)
	return nil
}

func (ce *simpleKVExecutor) ApplyWriteBatch(r storage.Resetable) error {
	wb := r.(util.WriteBatch)
	return ce.kv.Write(wb, false)
}

func (ce *simpleKVExecutor) Read(ctx storage.ReadContext) ([]byte, error) {
	request := ctx.Request()
	switch request.CmdType {
	case metadata.GetType:
		v, err := ce.kv.Get(request.Key)
		if err != nil {
			return nil, err
		}
		ctx.SetReadBytes(uint64(len(v)))
		return v, nil
	default:
		panic(fmt.Errorf("invalid read cmd %d", request.CmdType))
	}
}
