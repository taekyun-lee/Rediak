package KangDB

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	getoneQ :="SET key1 value1"
	getoneA :=[]string{"SET", "key1","value1"}

	v,ok := ParseArgs(getoneQ)
	fmt.Println(getoneA,v)
	if ok != nil || !reflect.DeepEqual(getoneA,v){
		t.Errorf("SET 1  parse failed %v %v\n",v,getoneA)
	}
}
