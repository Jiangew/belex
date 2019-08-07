package exchange

import (
	"errors"
	"log"
	"reflect"
	"time"
)

/**
  重试 API：一个参数必须是 error
  @retry  重试次数
  @delay  每次重试延迟时间间隔
  @method 调用的函数，比如: api.GetTicker, 注意：不是api.GetTicker(...)
  @params 参数,顺序一定要按照实际调用函数入参顺序一样
  @return 返回
*/
func Retry(retry int, delay time.Duration, method interface{}, params ...interface{}) interface{} {
	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		return errors.New("method not a function")
	}

	var value []reflect.Value = make([]reflect.Value, len(params))
	var i int = 0
	for ; i < len(params); i++ {
		value[i] = reflect.ValueOf(params[i])
	}

	var retV interface{}
	var retryC int = 0

_CALL:
	if retryC > 0 {
		log.Println("sleep ", delay, " after re call")
		time.Sleep(delay)
	}

	retValues := invokeM.Call(value)
	for _, vl := range retValues {
		if vl.Type().String() == "error" {
			if vl.IsNil() {
				continue
			}
			log.Println("[api error]", vl)
			retryC++
			if retryC <= retry-1 {
				log.Printf("Invoke Method[%s] Error , Begin Retry Call [%d] ...", invokeM.String(), retryC)
				goto _CALL
			} else {
				log.Println("Invoke Method Fail ???" + invokeM.String())
				return vl.Interface()
			}
		} else {
			retV = vl.Interface()
		}
	}

	return retV
}

/**
 * call all unfinished orders
 */
func CancelAllUnfinishedOrders(api API, symbol Symbol) int {
	if api == nil {
		log.Println("api instance is nil ??? , please new a api instance")
		return -1
	}

	c := 0
	for {
		ret := Retry(2, 200*time.Millisecond, api.GetActiveOrders, symbol)
		if err, isok := ret.(error); !isok {
			log.Println("[api error]", err)
			break
		}
		if ret == nil {
			break
		}

		orders, isok := ret.([]Order)
		if !isok || len(orders) == 0 {
			break
		}
		for _, ord := range orders {
			_, err := api.CancelOrder(ord.OrderID2, symbol)
			if err != nil {
				log.Println(err)
			} else {
				c++
			}
			time.Sleep(120 * time.Millisecond) // race limit
		}
	}

	return c
}
