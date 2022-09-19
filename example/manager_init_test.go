package example

import (
	"context"
	es "eventsourcing"
	"eventsourcing/example/currency"
	"github.com/aws/smithy-go/ptr"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type EventTypeAndRequest struct {
	Et  *es.EventType
	Req *currency.Request
}

func (e *EventTypeAndRequest) String() string {
	return es.JsonString(e)
}

type MakeRequestFunc func() *EventTypeAndRequest

// event request 를 제네레이팅 해주는 func 정의
var (
	makeCreateRequest = func() *EventTypeAndRequest {
		return &EventTypeAndRequest{
			Et:  &currency.CreateAmountStateEvent,
			Req: nil, // request 가 없다
		}
	}
	makeAddAmountRequest = func() *EventTypeAndRequest {
		rand.Seed(time.Now().UnixNano())
		return &EventTypeAndRequest{
			Et: &currency.AddAmountEvent,
			Req: &currency.Request{
				Amount: rand.Intn(1000),
			},
		}
	}
	makeMinusAmountRequest = func() *EventTypeAndRequest {
		rand.Seed(time.Now().UnixNano())
		return &EventTypeAndRequest{
			Et: &currency.MinusAmountEvent,
			Req: &currency.Request{
				Amount: rand.Intn(500),
			},
		}
	}
	makeChangeStatusRequest = func() *EventTypeAndRequest {
		rand.Seed(time.Now().UnixNano())
		return &EventTypeAndRequest{
			Et: &currency.ChangeStatusEvent,
			Req: &currency.Request{
				Status: (*currency.Status)(ptr.Int(currency.IDLE)),
			},
		}
	}
	makeChangeValueRequest = func() *EventTypeAndRequest {
		rand.Seed(time.Now().UnixNano())
		return &EventTypeAndRequest{
			Et: &currency.ChangeValueEvent,
			Req: &currency.Request{
				Status: (*currency.Status)(ptr.Int(currency.CLAIM)),
				Value:  ptr.String(strconv.Itoa(rand.Int())),
			},
		}
	}
	makeChangeValueV2Request = func() *EventTypeAndRequest {
		rand.Seed(time.Now().UnixNano())
		return &EventTypeAndRequest{
			Et: &currency.ChangeValueV2Event,
			Req: &currency.Request{
				Status: (*currency.Status)(ptr.Int(currency.IDLE)),
				Value:  ptr.String(strconv.Itoa(rand.Int())),
			},
		}
	}
	makeBurnRequest = func() *EventTypeAndRequest {
		return &EventTypeAndRequest{
			Et:  &currency.BurnEvent,
			Req: nil, // request 가 없다
		}
	}
)

func TestCurrencyManager(t *testing.T) {
	pk := es.PartitionKey("test_pk")

	ch := make(chan *EventTypeAndRequest, 10)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		choiceMakeRequestFuncList := []MakeRequestFunc{
			makeAddAmountRequest,
			makeMinusAmountRequest,
			makeChangeStatusRequest,
			makeChangeValueRequest,
			makeChangeValueV2Request,
		}

		// 첫 시작은 create 부터
		ch <- makeCreateRequest()
		rnd := rand.New(rand.NewSource(time.Now().Unix()))

		for {
			rnd.Seed(time.Now().UnixNano())
			f := choiceMakeRequestFuncList[rnd.Intn(len(choiceMakeRequestFuncList))]
			ch <- f()

			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for eventTypeRequest := range ch {
			log.Println("receive", eventTypeRequest)

			err := CurrencyEsManager.Validate(pk, eventTypeRequest.Et)
			if err != nil {
				t.Errorf("validation failed. %s - %s", eventTypeRequest.Et.String(), err)
				continue
			}
			err = CurrencyEsManager.Put(pk, eventTypeRequest.Et, eventTypeRequest.Req)
			if err != nil {
				t.Errorf("put failed. %s - %s", eventTypeRequest.Et.String(), err)
				continue
			}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			err := CurrencyEsManager.ApplyEvents(pk)
			if err != nil {
				t.Logf("update state snapshot failed. %s - %s", pk, err)
			}
			snapshot, _ := CurrencyEsManager.GetStateSnapshot(pk)
			log.Println("current state", snapshot)
		}
	}()

	<-ctx.Done()

	log.Println("[events]")
	log.Println(CurrencyEsManager.GetEvents(pk, 0))
	log.Println("[snapshot + events replayed]")
	log.Println(CurrencyEsManager.GetLatestState(pk))
}
