package chandict

import (
	"context"
	"github.com/sargun/goconcurrency/types"
	"runtime"
)

var _ types.ConcurrentDict = (*ChanDict)(nil)

type readRequestResponse struct {
	val    string
	exists bool
}

type readRequest struct {
	readKey      string
	responseChan chan readRequestResponse
}

type writeRequest struct {
	readKey      string
	writeVal     string
	responseChan chan struct{}
}

type casRequest struct {
	setOnNotExists bool
	readKey        string
	oldVal         string
	newVal         string
	responseChan   chan bool
}

type deleteRequest struct {
	deleteKey    string
	responseChan chan struct{}
}

type ChanDict struct {
	dict           map[string]string
	readRequests   chan readRequest
	writeRequests  chan writeRequest
	casRequests    chan casRequest
	deleteRequests chan deleteRequest
}

func NewChanDict() *ChanDict {
	ctx, cancel := context.WithCancel(context.Background())
	readRequests := make(chan readRequest)
	writeRequests := make(chan writeRequest)
	casRequests := make(chan casRequest)
	deleteRequests := make(chan deleteRequest)
	d := &ChanDict{
		readRequests:   readRequests,
		writeRequests:  writeRequests,
		casRequests:    casRequests,
		deleteRequests: deleteRequests,
	}
	// This is a lambda, so we don't have to add members to the struct
	runtime.SetFinalizer(d, func(dict *ChanDict) {
		cancel()
	})
	// We can't have run be a method of ChanDict, because otherwise then the goroutine will keep the reference alive
	go run(ctx, readRequests, writeRequests, casRequests, deleteRequests)
	return d
}

func run(parentCtx context.Context, readRequests <-chan readRequest, writeRequests <-chan writeRequest, casRequests <-chan casRequest, deleteRequests <-chan deleteRequest) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	localDict := make(map[string]string)

	for {
		select {
		case <-ctx.Done():
			return
		case dr := <-deleteRequests:
			delete(localDict, dr.deleteKey)
			close(dr.responseChan)
		case wr := <-writeRequests:
			localDict[wr.readKey] = wr.writeVal
			close(wr.responseChan)
		case rr := <-readRequests:
			val, exists := localDict[rr.readKey]
			rr.responseChan <- readRequestResponse{val, exists}
		case cr := <-casRequests:
			if val, exists := localDict[cr.readKey]; exists && val == cr.oldVal {
				localDict[cr.readKey] = cr.newVal

				cr.responseChan <- true
			} else if !exists && cr.setOnNotExists {
				localDict[cr.readKey] = cr.newVal
				cr.responseChan <- true
			} else {
				cr.responseChan <- false

			}
		}
	}
}

func (dict *ChanDict) SetKey(key, val string) {
	c := make(chan struct{})
	w := writeRequest{readKey: key, writeVal: val, responseChan: c}
	dict.writeRequests <- w
	<-c
}

func (dict *ChanDict) ReadKey(key string) (string, bool) {
	c := make(chan readRequestResponse)
	w := readRequest{readKey: key, responseChan: c}
	dict.readRequests <- w
	resp := <-c
	return resp.val, resp.exists
}

func (dict *ChanDict) CasVal(key, oldVal, newVal string, setOnNotExists bool) bool {
	c := make(chan bool)
	w := casRequest{readKey: key, oldVal: oldVal, newVal: newVal, responseChan: c, setOnNotExists: setOnNotExists}
	dict.casRequests <- w
	return <-c
}

func (dict *ChanDict) DeleteKey(key string) {
	c := make(chan struct{})
	d := deleteRequest{deleteKey: key, responseChan: c}
	dict.deleteRequests <- d
	<-c
}
