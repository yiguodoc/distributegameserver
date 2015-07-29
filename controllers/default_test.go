package controllers

import (
	// "net/http"
	// "net/http/httptest"
	"testing"
	// "runtime"
	// "path/filepath"
	// _ "distributionGame/routers"

	// "github.com/astaxie/beego"
	// . "github.com/smartystreets/goconvey/convey"
	// "fmt"
	"time"
)

func init() {
}

//failed:all goroutines are asleep - deadlock!
func test先同步发消息后接收(t *testing.T) {
	ch := make(chan int)
	ch <- 1
	time.Sleep(1 * time.Second)
	i := <-ch
	t.Logf("result : %d", i)
}

//pass
//对比之前的测试，说明发送消息是个同步操作
func Test先异步发消息后接收(t *testing.T) {
	t.Log("start...")
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	time.Sleep(1 * time.Second)
	i := <-ch
	if i != 1 {
		t.Fail()
	}
	t.Logf("result : %d", i)
	t.Log("end...")
}

func Test先异步发多个消息后接收(t *testing.T) {
	t.Log("start...")
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
	}()
	time.Sleep(1 * time.Second)
	i := <-ch
	if i != 1 {
		t.Fail()
	}
	i = <-ch
	if i != 2 {
		t.Fail()
	}
	t.Logf("result : %d", i)
	t.Log("end...")
}
