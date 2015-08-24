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
	"reflect"
	"time"
)

func init() {
}

func Compose(o interface{}, fns ...func(interface{}) interface{}) interface{} {
	if len(fns) <= 0 {
		return o
	} else {
		cnt := len(fns)
		last := fns[cnt-1]
		return Compose(last(o), fns[:cnt-1]...)
	}
}

type stru struct {
	name string
}

func Test函数组合(t *testing.T) {
	strus := []*stru{
		&stru{"1"},
		&stru{"2"},
		&stru{"3"},
	}
	f1 := func(o interface{}) interface{} {
		l := o.([]*stru)
		for _, item := range l {
			item.name = item.name + "01"
		}
		return l
	}
	f2 := func(o interface{}) interface{} {
		l := o.([]*stru)
		for _, item := range l {
			item.name = item.name + "02"
		}
		return l
	}
	strus2 := Compose(strus, f1, f2).([]*stru)
	for i, item := range strus2 {
		t.Logf(" %d  => %s", i, item.name)
	}
	t.Fail()
}

func Test反射(t *testing.T) {
	struType := reflect.TypeOf(stru{})
	t.Fail()
	t.Log(struType)
}

func (s *stru) setName(n string) {
	s.name = n
}
func tryFunc(f func(string)) {
	f("111")
}

//将结构的方法作为函数传递时，必须通过结构调用，并且结构本身是和方法一起传递的，类似于多了一个参数的函数
func Test函数签名调用(t *testing.T) {
	var s stru
	tryFunc(s.setName)
	// t.Fail()
	t.Log(s)
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
