/**
2 * @Author: Nico
3 * @Date: 2021/5/26 9:41
4 */
package cluster

import (
	"fmt"
	"testing"
)

func TestStart30001(t *testing.T) {
	err := Start(30001, nil)
	fmt.Println(err)
	select {}
}

func TestStart30002(t *testing.T) {
	err := Start(30002, []string{"127.0.0.1:30001"})
	fmt.Println(err)
	select {}
}