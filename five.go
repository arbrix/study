package main

import (
	"fmt"
	"sync/atomic"
)

//findEquation([]int{1,2,3,4,5,6,7,8,9}, 100, "")
func findEquation(numbers []int, sum int, resStr string) string {
	num := numbers[0]
	//fmt.Printf("current number: %d, check sum: %d, result string: %s, set: %v\n", num, sum, resStr, numbers)
	if 1 == len(numbers) {
		if num == sum {
			resStr += fmt.Sprintf("%+d", num) + "\n"
			return resStr
		}
		return ""
	}
	var res, str string
	//plus
	plus := numbers[1:]
	if "" == resStr {
		str = findEquation(plus, sum-num, fmt.Sprintf("%d", num))
	} else {
		str = findEquation(plus, sum-num, resStr + fmt.Sprintf("%+d", num))
	}
	if "" != str {
		res += str
	}
	//minus
	minus := numbers[1:]
	minus[0] *= -1
	if "" == resStr {
		str = findEquation(minus, sum-num, fmt.Sprintf("%d",num))
	} else {
		str = findEquation(minus, sum-num, resStr + fmt.Sprintf("%+d",num))
	}
	if "" != str {
		res += str
	}
	//concat
	concat := append([]int(nil), numbers[1:]...)
	if concat[0] < 0 {
		concat[0] *= -1
	}
	if num > 0 {
		concat[0] = num*10+concat[0]
	} else {
		concat[0] = num*10-concat[0]
	}
	//fmt.Printf("concat first two numbers: %v, from set: %v\n", concat, numbers)
	str = findEquation(concat, sum, resStr)
	if "" != str {
		res += str
	}
	return res
}

//findEquation2(1, 98765432, 100, "")
func findEquation2(num, set, sum int, equation string) string {
	if set == 0 {
		if num == sum {
			equation += fmt.Sprintf("%+d", num) + "\n"
			return equation
		}
		return ""
	}
	var s, r string
	//plus
	if "" == equation {
		s = findEquation2(set%10, set/10, sum-num, fmt.Sprintf("%d", num))
	} else {
		s = findEquation2(set%10, set/10, sum-num, equation + fmt.Sprintf("%+d", num))
	}
	if "" != s {
		r += s
	}
	//minus
	if "" == equation {
		s = findEquation2(-set%10, set/10, sum-num, fmt.Sprintf("%d", num))
	} else {
		s = findEquation2(-set % 10, set / 10, sum - num, equation + fmt.Sprintf("%+d", num))
	}
	if "" != s {
		r += s
	}
	//concat
	concat := num*10
	if concat > 0 {
		concat += set%10
	} else {
		concat -= set%10
	}
	s = findEquation2(concat, set/10, sum, equation)
	if "" != s {
		r += s
	}
	return r
}

func concurrentSearch(num, set, sum int, eq string, res chan <- string, cnt *int32) {
	atomic.AddInt32(cnt, 1)
	defer atomic.AddInt32(cnt, -1)
	if set == 0 {
		if num == sum {
			eq += fmt.Sprintf("%+d", num)
			res <- eq
		}
		return
	}
	//plus
	if "" == eq {
		go concurrentSearch(set%10, set/10, sum-num, fmt.Sprintf("%d", num), res, cnt)
	} else {
		go concurrentSearch(set%10, set/10, sum-num, eq + fmt.Sprintf("%+d", num), res, cnt)
	}
	//minus
	if "" == eq {
		go concurrentSearch(-set%10, set/10, sum-num, fmt.Sprintf("%d", num), res, cnt)
	} else {
		go concurrentSearch(-set % 10, set / 10, sum - num, eq + fmt.Sprintf("%+d", num), res, cnt)
	}
	//concat
	concat := num*10
	if concat > 0 {
		concat += set%10
	} else {
		concat -= set%10
	}
	go concurrentSearch(concat, set/10, sum, eq, res, cnt)
}

//findEquation3(987654321, 100)
func findEquation3(set, sum int,) string {
	var cnt int32
	eqChan := make(chan string)
	go concurrentSearch(set%10, set/10, sum, "", eqChan, &cnt)
	res := <- eqChan
	for atomic.LoadInt32(&cnt) > 0 {
		select {
		case s := <- eqChan:
			res += "\n" + s
		default:
			continue
		}
	}
	return res
}

func main() {
	fmt.Println("result: \n", findEquation3(987654321, 100))
}
