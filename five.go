package main

import "fmt"

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

func findEquation2(num, set, sum int, equation string) string {
	if set == 0 {
		if num == sum {
			equation += fmt.Sprintf("%+d", num) + "\n"
			return equation
		}
		return ""
	}
	if num > 123456789 {
		fmt.Println("to big number: ", num)
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

func main() {
	fmt.Print("result: \n", findEquation2(1, 98765432, 100, ""))
}
