package queue

import (
	"AS/commonlib-go"
	"fmt"
	"sync"
	"time"
)

type Qbus struct {
	Name	string
	Success []*QbusData
	client  *Client
	dm 		sync.Mutex
}

type QbusResponse struct {
	Code 		int 			`json:"code"`
	Message 	string 			`json:"message"`
	Data		[]*QbusData 	`json:"data"`
	RequestId 	string   		`json:"requestId"`

}

type QbusData struct {
	DeQueueCount		int			`json:"dequeueCount"`
	EnqueueTime			int			`json:"enqueueTime"`
	FirstDeQueueTime 	int			`json:"firstDequeueTime"`
	MsgBody 			string		`json:"msgBody"`
	MsgId 				string		`json:"msgId"`
	MsgTag 				[]string 	`json:"msgTag"`
	NextVisibleTime 	int			`json:"nextVisibleTime"`
	ReceiptHandle   	string		`json:"receiptHandle"`
}

const MaxNumber = 16
type Run interface {
	GetQueueName() string
	MsgHandler(q *Qbus, resp *QbusResponse) bool
}

func NewQueue(name string, config ClientConfig) *Qbus{
	Ins := &Qbus{
		name,
		nil,
		NewClient(config),
		sync.Mutex{},
	}
	return Ins
}

// 批量发送数据
func (q *Qbus) sendBatchMessage(params interface{}, delay int) interface{} {
	resp := q.client.sendMessage(q.Name, params, delay)
	return resp
}

func (q *Qbus) EnQueue(message interface{}, delay int) interface{}{
	resp := q.sendBatchMessage(message, delay)
	return resp
}

func (q *Qbus) Exec(run Run, readNum uint8, duration int, pollingSecond int) {
	if readNum > MaxNumber {
		readNum = MaxNumber
	}

	timer := time.NewTimer(time.Duration(duration) * time.Second)
	for {
		select {
		case <- timer.C:
			_ = q.deleteBatchMessage()
			return // 暂时依赖msgId 做业务幂等
		default:
			q.DeQueue(run, readNum, pollingSecond)
		}
	}
}


// 拉取数据
func (q *Qbus) DeQueue(run Run, readNum uint8, pollingSecond int) bool{
	lib.Trace(fmt.Sprintf("receive queue[%s]batch message\n", run.GetQueueName()))
	resp, err := q.receiveBatchMessage(readNum, pollingSecond)
	if err != nil {
		lib.Warning("receive batch message error", err, "qbus")
		panic(err)
	}

	lib.Trace(fmt.Sprintf("exec msg queue[%s] handler", run.GetQueueName()))
	// 执行方法 执行成功的message 推入success 当中， 删除对列消息
	if ok := run.MsgHandler(q, resp);ok {
		_ = q.deleteBatchMessage()
	}

	return true
}

func (q *Qbus) receiveBatchMessage(readNum uint8, pollingSecond int) (*QbusResponse, error) {
	return q.client.receiveMessage(q.Name, readNum, pollingSecond)
}

// 批量删除数据
func (q *Qbus) deleteBatchMessage() error{
	q.dm.Lock()
	defer q.dm.Unlock()
	_, err := q.client.deleteMessage(q.Name, q.Success)
	if err != nil {
		return err
	}

	q.Success = []*QbusData{}
	return nil
}