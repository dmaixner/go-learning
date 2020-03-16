//http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func arraysAreValues() {
	x := [3]int{1, 2, 3}

	func(arr *[3]int) {
		(*arr)[0] = 7
		fmt.Println(arr)  //prints &[7 2 3]
		fmt.Println(*arr) //prints [7 2 3]
	}(&x)

	fmt.Println(x) //prints [7 2 3]
}

func slicesArePointers() {
	x := []int{1, 2, 3}

	func(arr []int) {
		arr[0] = 7
		fmt.Println(arr) //prints [7 2 3]
	}(x)

	fmt.Println(x) //prints [7 2 3]
}

type MyData struct {
	One int // this one is exported
	two string
}

func waitgroup() {
	var wg sync.WaitGroup
	done := make(chan struct{})
	wq := make(chan interface{})
	workerCount := 2

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go doit(i, wq, done, &wg)
	}

	for i := 0; i < workerCount; i++ {
		wq <- i
	}

	fmt.Println("closing done")
	close(done)
	wg.Wait()
	fmt.Println("all done!")
}

func doit(workerId int, wq <-chan interface{}, done <-chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("[%v] is running\n", workerId)
	defer wg.Done()
	for {
		select {
		case m := <-wq:
			fmt.Printf("[%v] m => %v\n", workerId, m)
			time.Sleep(1 * time.Second)
		case <-done:
			fmt.Printf("[%v] is done\n", workerId)
			return
		}
	}
}

type data struct {
	num   int
	key   *string
	items map[string]bool
}

func (this *data) pmethod() {
	this.num = 7
}

func (this data) vmethod() {
	this.num = 8
	*this.key = "v.key"
	this.items["vmethod"] = true
}

func value_pointer() {
	key := "key.1"
	d := data{1, &key, make(map[string]bool)}

	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
	//prints num=1 key=key.1 items=map[]

	d.pmethod()
	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
	//prints num=7 key=key.1 items=map[]

	d.vmethod()
	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
	//prints num=7 key=v.key items=map[vmethod:true]
}

func json_gotchas() {
	data := "x < y"

	raw, _ := json.Marshal(data)
	fmt.Println(string(raw))
	//prints: "x \u003c y" <- probably not what you expected

	var b1 bytes.Buffer
	json.NewEncoder(&b1).Encode(data)
	fmt.Println(b1.String())
	//prints: "x \u003c y"\n <- probably not what you expected

	var b2 bytes.Buffer
	enc := json.NewEncoder(&b2)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	fmt.Println(b2.String())
	//prints: "x < y" <- looks better

	var data1 = []byte(`{"status": 200}`)
	var result struct {
		Status uint64 `json:"status"`
	}
	if err := json.NewDecoder(bytes.NewReader(data1)).Decode(&result); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("result => %+v\n", result)
	//prints: result => {Status:200}
}

func deferf() {
	var i int = 1
	defer fmt.Println("result =>", func() int { return i * 2 }())
	i++
	//prints: result => 2 (not ok if you expected 4)
	i = 1
	defer func(in *int) { fmt.Println("result =>", *in) }(&i)
	i = 2
	//prints: result => 2
}

func interfacesnil() {
	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) //prints: <nil> true
	fmt.Println(in, in == nil)     //prints: <nil> true

	in = data
	fmt.Println(in, in == nil) //prints: <nil> false
	//'data' is 'nil', but 'in' is not 'nil'

	doit := func(arg int) interface{} {
		var result *struct{} = nil

		if arg > 0 {
			result = &struct{}{}
		}

		return result
	}

	if res := doit(-1); res != nil {
		fmt.Println("good result:", res) //prints: good result: <nil>
		//'res' is not 'nil', but its value is 'nil'
	}

	doit2 := func(arg int) interface{} {
		var result *struct{} = nil

		if arg > 0 {
			result = &struct{}{}
		} else {
			return nil //return an explicit 'nil'
		}

		return result
	}

	if res := doit2(-1); res != nil {
		fmt.Println("good result:", res)
	} else {
		fmt.Println("bad result (res is nil)") //here as expected
	}

}

func main() {
	x := 1
	fmt.Println(x) //prints 1
	{
		fmt.Println(x) //prints 1
		x := 2
		fmt.Println(x) //prints 2
	}
	fmt.Println(x) //prints 1 (bad if you need 2)

	var s []int
	fmt.Println(append(s, 1))
	fmt.Println(s)

	arraysAreValues()
	slicesArePointers()

	m := map[string]string{"one": "a", "two": "", "three": "c"}
	if v := m["two"]; v == "" { //incorrect
		fmt.Println("no entry")
	}
	if _, ok := m["nonexist"]; !ok {
		fmt.Println("no entry")
	}

	data := "ě"
	fmt.Println(len(data))                    //prints: 2
	fmt.Println(len([]rune(data)))            //prints: 1
	fmt.Println(utf8.RuneCountInString(data)) //prints: 1
	data = "é"                               // weird char, different from é
	fmt.Println(len(data))                    //prints: 3
	fmt.Println(len([]rune(data)))            //prints: 2
	fmt.Println(utf8.RuneCountInString(data)) //prints: 2

	in := MyData{1, "two"}
	fmt.Printf("%#v\n", in) //prints main.MyData{One:1, two:"two"}
	encoded, _ := json.Marshal(in)
	fmt.Println(string(encoded)) //prints {"One":1}
	var out MyData
	json.Unmarshal(encoded, &out)
	fmt.Printf("%#v\n", out) //prints main.MyData{One:1, two:""}

	waitgroup()
	value_pointer()

	json_gotchas()

	type datacompare struct {
		num     int
		fp      float32
		complex complex64
		str     string
		char    rune
		yes     bool
		events  <-chan string
		handler interface{}
		ref     *byte
		raw     [10]byte
		// checks  [10]func() bool   //not comparable
		// doit    func() bool       //not comparable
		// m       map[string]string //not comparable
		// bytes   []byte            //not comparable
	}

	type dataNOTcompare struct {
		checks [10]func() bool   //not comparable
		doit   func() bool       //not comparable
		m      map[string]string //not comparable
		bytes  []byte            //not comparable
	}

	v1 := datacompare{}
	v2 := datacompare{}
	fmt.Println("v1 == v2:", v1 == v2) //prints: v1 == v2: true

	v3 := dataNOTcompare{}
	v4 := dataNOTcompare{}
	fmt.Println("v3 == v4:", reflect.DeepEqual(v3, v4)) //prints: v3 == v4: true

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2)) //prints: m1 == m2: true

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2)) //prints: s1 == s2: true

	var b1 []byte = nil
	b2 := []byte{}
	fmt.Println("b1 == b2:", reflect.DeepEqual(b1, b2)) //prints: b1 == b2: false
	fmt.Println("b1 == b2:", bytes.Equal(b1, b2))       //prints: b1 == b2: true

	var str string = "one"
	var inter interface{} = "one"
	fmt.Println("str == inter:", str == inter, reflect.DeepEqual(str, inter))
	//prints: str == inter: true true

	vv1 := []string{"one", "two"}
	vv2 := []interface{}{"one", "two"}
	fmt.Println("vv1 == vv2:", reflect.DeepEqual(vv1, vv2))
	//prints: vv1 == vv2: false (not ok)

	ces1 := "ěščřß"
	ces2 := "ĚŠČŘß"
	fmt.Println("ces1 == ces2:", ces1 == ces2)
	fmt.Println("ces1 == ces2:", strings.ToUpper(ces1) == strings.ToUpper(ces2))
	fmt.Println("ces1 == ces2:", strings.ToLower(ces1) == strings.ToLower(ces2))
	fmt.Println("ces1 == ces2:", strings.EqualFold(ces1, ces2))

	rangedata := []int{1, 2, 3}
	for _, v := range rangedata {
		v *= 10 //original item is not changed, value passed by range is copy
	}
	fmt.Println("rangedata:", rangedata) //prints rangedata: [1 2 3]
	for i, _ := range rangedata {
		rangedata[i] *= 10
	}
	fmt.Println("rangedata:", rangedata) //prints rangedata: [10 20 30]
	rangedatapointer := []*struct{ num int }{{1}, {2}, {3}}
	for _, v := range rangedatapointer {
		v.num *= 10
	}
	fmt.Println(rangedatapointer[0], rangedatapointer[1], rangedatapointer[2])    //prints &{10} &{20} &{30}
	fmt.Println(*rangedatapointer[0], *rangedatapointer[1], *rangedatapointer[2]) //prints {10} {20} {30}

	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) //prints: 10000 10000 <byte_addr_x>
	resliced := raw[:3]
	fmt.Println(len(resliced), cap(resliced), &resliced[0]) //prints: 3 10000 <byte_addr_x>
	copied := make([]byte, 3)
	copy(copied, raw[:3])
	fmt.Println(len(copied), cap(copied), &copied[0]) //prints: 3 3 <byte_addr_y>

	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')
	dir1 := path[:sepIndex]
	dir3 := path[:sepIndex:sepIndex] // full slice expression
	dir2 := path[sepIndex+1:]
	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAA
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => BBBBBBBBB
	fmt.Println("dir3 =>", string(dir3)) //prints: dir3 => AAAA
	dir3 = append(dir3, "suffix"...)
	path = bytes.Join([][]byte{dir3, dir2}, []byte{'/'})
	fmt.Println("new path =>", string(path))
	dir1 = append(dir1, "suffix"...)
	path = bytes.Join([][]byte{dir1, dir2}, []byte{'/'})
	fmt.Println("new path =>", string(path))

	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => uffixBBBB (not ok)
	fmt.Println("dir3 =>", string(dir3)) //prints: dir3 => AAAAsuffix

	slice1 := []int{1, 2, 3}
	fmt.Println(len(slice1), cap(slice1), slice1) //prints 3 3 [1 2 3]
	slice2 := slice1[1:]
	fmt.Println(len(slice2), cap(slice2), slice2) //prints 2 2 [2 3]
	for i := range slice2 {
		slice2[i] += 20
	}
	//still referencing the same array
	fmt.Println(slice1) //prints [1 22 23]
	fmt.Println(slice2) //prints [22 23]
	slice2 = append(slice2, 4)
	for i := range slice2 {
		slice2[i] += 10
	}
	//slice1 is now "stale"
	fmt.Println(slice1) //prints [1 22 23]
	fmt.Println(slice2) //prints [32 33 14]

	dataThree := []string{"one", "two", "three"}
	for _, v := range dataThree {
		go func() {
			fmt.Println(v)
		}()
	}
	//goroutines print: three, three, three
	time.Sleep(1 * time.Second)
	for _, v := range dataThree {
		vcopy := v //
		go func() {
			fmt.Println(vcopy)
		}()
	}
	//goroutines print: one, two, three
	time.Sleep(1 * time.Second)
	for _, v := range dataThree {
		go func(in string) {
			fmt.Println(in)
		}(v)
	}
	//goroutines print: one, two, three
	time.Sleep(1 * time.Second)

	deferf()

	type mapslice struct {
		name string
	}
	mymap := map[string]mapslice{"x": {"one"}}
	// mymap["x"].name = "two" //error
	fmt.Println(mymap)
	r := mymap["x"]
	r.name = "two"
	mymap["x"] = r
	fmt.Printf("%v\n", mymap) //prints: map[x:{two}]
	mymap2 := map[string]*mapslice{"x": {"one"}}
	mymap2["x"].name = "three" //ok
	fmt.Println(mymap2["x"])   //prints: &{three}
	myslice := []mapslice{{"one"}}
	myslice[0].name = "two" //ok
	fmt.Println(myslice)    //prints: [{two}]

	interfacesnil()
}
