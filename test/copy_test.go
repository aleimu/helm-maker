package test

import (
	"fmt"
	"testing"
)

type simple struct {
	A string
	B string
	C []byte
}

func TestCopy(t *testing.T) {
	a := new(simple)
	a.A = "A"

	b := a
	b.A = "BA"

	c := *a
	c.A = "CA"

	fmt.Printf("%+v %+v %+v \n", a, b, c)
	fmt.Println(&a.A, &b.A, &c.A)
}

/*
Go语言中所有赋值操作都是值传递，
如果结构中不含指针，则直接赋值就是深度拷贝；
如果结构中含有指针（包括自定义指针，以及切片，map等使用了指针的内置类型），则数据源和拷贝之间对应指针会共同指向同一块内存，这时深度拷贝需要特别处理。

目前，有三种方法，
一是用gob序列化成字节序列再反序列化生成克隆对象；
二是先转换成json字节序列，再解析字节序列生成克隆对象；
三是针对具体情况，定制化拷贝。

前两种方法虽然比较通用但是因为使用了reflex反射，性能比定制化拷贝要低出2个数量级，所以在性能要求较高的情况下应该尽量避免使用前两者。

*/
