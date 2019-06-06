package test

import (
	"fmt"
	"testing"
)

type Parent struct {
	Name string
	Children
}

type Children struct {
	age string
}

func (c *Children) GetName() {
	fmt.Println(c.GetName()
}

func Test_sdf(t testing.T)  {
	test : &Parent{
		Name: "fdsf"
	}
}
