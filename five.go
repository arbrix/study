// Write a program that outputs all possibilities to put + or - or nothing
// between the numbers 1, 2, ..., 9 (in this order) such that the result is always 100.
// For example: 1 + 2 + 34 – 5 + 67 – 8 + 9 = 100
package main

import (
	"fmt"
	"sync/atomic"
	"sync"
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
		atomic.AddInt32(cnt, 1)
		go concurrentSearch(set%10, set/10, sum-num, fmt.Sprintf("%d", num), res, cnt)
	} else {
		atomic.AddInt32(cnt, 1)
		go concurrentSearch(set%10, set/10, sum-num, eq + fmt.Sprintf("%+d", num), res, cnt)
	}
	//minus
	if "" == eq {
		atomic.AddInt32(cnt, 1)
		go concurrentSearch(-set%10, set/10, sum-num, fmt.Sprintf("%d", num), res, cnt)
	} else {
		atomic.AddInt32(cnt, 1)
		go concurrentSearch(-set % 10, set / 10, sum - num, eq + fmt.Sprintf("%+d", num), res, cnt)
	}
	//concat
	concat := num*10
	if concat > 0 {
		concat += set%10
	} else {
		concat -= set%10
	}
	atomic.AddInt32(cnt, 1)
	go concurrentSearch(concat, set/10, sum, eq, res, cnt)
}

//taskCond condition for generation resulting combination
type taskCond struct {
	num, set, sum int
	eq string
}

func NewTask(num, next, subset, sum int, eq string) taskCond {
	if "" == eq {
		return taskCond{next, subset, sum, fmt.Sprintf("%d", num)}
	}
	return taskCond{next, subset, sum, eq + fmt.Sprintf("%+d", num)}
}

func worker(ts chan taskCond, eq chan <- string, finish chan struct{}, cnt *int32, tCnt int32) {
	readyCombination := make(chan taskCond)

	wg := sync.WaitGroup{}
	wg.Add(1) //real wait for first ready combination
	go func() {
		defer wg.Done() //no more data for ready combination
		for {
			select {
			case t, ok := <- ts:
				if !ok {
					return
				}
				if t.set > 0 {
					//plus
					ts <- NewTask(t.num, t.set % 10, t.set / 10, t.sum-t.num, t.eq)
					//minus
					ts <- NewTask(t.num, -t.set % 10, t.set / 10, t.sum-t.num, t.eq)
					//concat
					concat := t.num*10
					if concat > 0 {
						concat += t.set%10
					} else {
						concat -= t.set%10
					}
					ts <- NewTask(t.num, concat, t.set / 10, t.sum, t.eq)
					//задач больше не будет, можем закрывать канал
					if atomic.AddInt32(cnt, 3) >= tCnt {
						close(ts)
					}
				} else {
					readyCombination <- t
				}
			default:
				continue
			}
		}
	}()
	go func() {
		for comb := range readyCombination {
			//подходит ли данная последовательность
			if comb.num == comb.sum {
				//fmt.Printf("%d == %d (%t) cnt: %d tCnt: %d\n", comb.num, comb.sum, comb.num == comb.sum, atomic.LoadInt32(cnt), tCnt)
				comb.eq += fmt.Sprintf("%+d", comb.num)
				eq <- comb.eq
			}
		}
		//each worker send signal when finishing
		finish <- struct{}{}
	}()
	wg.Wait()
	close(readyCombination)
	//fmt.Printf("end worker, cnt: %d tCnt: %d\n", atomic.LoadInt32(cnt), tCnt)
}

//findEquation3(987654321, 100)
func findEquation3(set, sum int,) string {
	var cnt, totalCnt int32
	//считаем количество вариантов = емкости канала с задачами
	//каждый воркер добавляет по 3 задачи и они же их обрабатывают
	//поэтому нужен буфер
	tmp := set
	totalCnt = 1
	for tmp >0 {
		tmp /= 10
		totalCnt *= 3
	}
	totalCnt = (totalCnt - 1) / 2 //sum of geometry progression elements s=b(1-q^n)/(1-q), b=1, q=3, n=len(set)
	equation := make(chan string, 100)
	sTask := make(chan taskCond, totalCnt*10)
	finish := make(chan struct{})
	for i:= 0; i<2; i++ {
		go worker(sTask, equation, finish, &cnt, totalCnt)
	}
	//проверяем завершение работы всех воркеров и закрываем канал решений
	go func() {
		for i:= 0; i<2; i++ {
			<-finish
		}
		close(finish)
		close(equation)
	}()
	//счетчик добавленых задач, влияет на закрытие канала задач
	atomic.AddInt32(&cnt, 1)
	sTask <- taskCond{set%10,set/10,sum,""}
	res := ""
	//после отработки всех задач канал вариантов так же будет закрыт
	for s := range equation {
		res += "\n" + s
	}
	return res
}

func main() {
	fmt.Println(findEquation3(987654321, 100))
}
