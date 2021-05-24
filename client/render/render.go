/**
2 * @Author: Nico
3 * @Date: 2021/5/23 20:28
4 */
package render

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

var (
	in *bufio.Reader = bufio.NewReader(os.Stdin)
	out *bufio.Writer = bufio.NewWriter(os.Stdout)
	buffer = bytes.Buffer{}
)

func Readline() ([]byte, error){
	for {
		b, err := in.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == '\n'{
			lines := buffer.Bytes()
			buffer.Reset()
			return lines, nil
		}
		buffer.WriteByte(b)
	}
}

func ReadBuffer() []byte{
	return buffer.Bytes()
}

func EraseLine() {
	fmt.Printf("\033[1A")
	fmt.Printf("\r\r")
}


