package chandict

import (
	"context"
	"github.com/sargun/goconcurrency/types"
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

func NewChanDict(ctx context.Context) *ChanDict {
	d := &ChanDict{
		dict:           make(map[string]string),
		readRequests:   make(chan readRequest),
		writeRequests:  make(chan writeRequest),
		casRequests:    make(chan casRequest),
		deleteRequests: make(chan deleteRequest),
	}
	go d.run(ctx)
	return d
}

func (dict *ChanDict) run(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case dr := <-dict.deleteRequests:
			delete(dict.dict, dr.deleteKey)
			close(dr.responseChan)
		case wr := <-dict.writeRequests:
			dict.dict[wr.readKey] = wr.writeVal
			close(wr.responseChan)
		case rr := <-dict.readRequests:
			val, exists := dict.dict[rr.readKey]
			rr.responseChan <- readRequestResponse{val, exists}
		case cr := <-dict.casRequests:
			if val, exists := dict.dict[cr.readKey]; exists && val == cr.oldVal {
				dict.dict[cr.readKey] = cr.newVal

				cr.responseChan <- true
			} else if !exists && cr.setOnNotExists {
				dict.dict[cr.readKey] = cr.newVal
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
