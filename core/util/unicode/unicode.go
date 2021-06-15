/**
2 * @Author: Nico
3 * @Date: 2021/5/23 16:19
4 */
package unicode

import (
	"strconv"
)

func IsNumber(s string) bool {
	_, err := strconv.ParseInt(args[1], 10, 64)
	return err == nil
}
