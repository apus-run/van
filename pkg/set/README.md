在 Kubernetes 的 k8s.io/utils 库中，set 包含了一些有用的集合数据结构和方法。这些方法的设计目的主要是为了解决 Go 中常见的集合操作（例如去重、查找等）的痛点

Set 结构体

Set 是一个封装了 map 的数据结构，用于实现集合的功能。通常集合会保证元素是唯一的，而 map 天生就能提供去重的功能。因此，Set 实际上就是基于 Go 内置的 map 类型来实现的，map[Type]struct{}

元素操作
1. Insert(items ...T)
• 功能：向集合中添加一个或多个元素
package main

import (
    "fmt"
    "k8s.io/utils/set"
)

func main() {
    // 创建一个
    set1 := set.New[string]("name", "age", "address")
    // 向集合中添加一个或者多个元素
    set1.Insert("apple", "banana", "name")
    
    // 打印集合内容
    fmt.Println(set1)
}
go run main.go
map[address:{} age:{} apple:{} banana:{} name:{}]
2. Delete(items ...T)
• 功能：从集合中删除一个或多个元素
func main() {
    // 创建一个
    set1 := set.New[string]("name", "age", "address")
    
    // 向集合中添加一个或者多个元素
    set1.Insert("apple", "banana", "name")
    
    // 集合中删除name元素
    set1.Delete("name")
    
    // 打印集合内容
    fmt.Println(set1)
}
go run main.go
map[address:{} age:{} apple:{} banana:{}]
3. Has(item T) bool
• 功能：检查元素是否存在于集合中
func main() {
    // 创建一个
    set1 := set.New[string]("name", "age", "address")
    
    // 判断指定元素是否存在
    exists := set1.Has("name")
    fmt.Println(exists)
}
go run main.go
true
4. HasAll(items ...E) bool
• 功能：方法用于检查集合是否包含所有指定的元素
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    fmt.Println(set1.HasAll("name", "age"))
}
go run main.go
true
5. HasAny(items ...string) bool
• 功能：方法用于检查集合是否至少包含一个指定的元素
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    fmt.Println(set1.HasAny("name", "huawei", "aws"))
}
go run main.go
true
6. PopAny() (string, bool)
• 功能：方法用于从集合中随机弹出一个元素，并返回它
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    set2 := set1.Clone()
    item, ok := set1.PopAny()
    if ok {
        fmt.Println("Popped item:", item)
    }
    fmt.Println(set1.SortedList())
    fmt.Println(set2.SortedList())
}
go run main.go
Popped item: name
[address age apple]
[address age apple name]
7. Len() int
• 功能：返回集合中元素的数量
func main() {
    // 创建一个
    set1 := set.New[string]("name", "age", "address")
    size := set1.Len()
    fmt.Println(size)
}
go run main.go
3
8. Union(s Set[T]) Set[T]
• 功能：返回当前集合与另一个集合的并集（新集合）
func main() {
    // 创建集合
    set1 := set.New[string]("name", "age", "address")
    set2 := set.New[string]("apple", "huawei", "aws")
    // 合并两个集合
    unionSet := set1.Union(set2)
    fmt.Println(unionSet)
}
go run main.go
map[address:{} age:{} apple:{} aws:{} huawei:{} name:{}]
9. Intersection(s Set[T]) Set[T]
• 功能：返回当前集合与另一个集合的交集（新集合）
func main() {
    // 创建集合
    set1 := set.New[string]("name", "age", "address")
    set2 := set.New[string]("apple", "huawei", "aws", "age")
    intersectSet := set1.Intersection(set2)
    fmt.Println(intersectSet)
}
go run main.go
map[age:{}]
10. Difference(s Set[T]) Set[T]
• 功能：返回当前集合与另一个集合的差集（当前集合有而另一个集合没有的元素）
func main() {
    set1 := set.New[string]("name", "age", "address")
    set2 := set.New[string]("apple", "huawei", "aws", "age")
    diffSet := set1.Difference(set2)
    fmt.Println(diffSet)
}
go run main.go
map[address:{} name:{}]
11. SymmetricDifference(s Set[T]) Set[T]
• 功能：返回两个集合的对称差集（仅在其中一个集合中存在的元素）
func main() {
    set1 := set.New[string]("name", "age", "address")
    set2 := set.New[string]("apple", "huawei", "aws", "age")
    symDiffSet := set1.SymmetricDifference(set2)
    fmt.Println(symDiffSet)
}
go run main.go
map[address:{} apple:{} aws:{} huawei:{} name:{}]
12. IsSuperset(s Set[T]) bool
• 功能：判断当前集合是否是另一个集合的超集
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    set2 := set.New[string]("age")
    isSuperset := set1.IsSuperset(set2)
    fmt.Println(isSuperset)
}
go run main.go
true
13. Equal(s Set[T]) bool
• 功能：判断两个集合是否相等（元素完全相同）
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    set2 := set.New[string]("huawei", "age")
    fmt.Println(set1.Equal(set2))
}
go run main.go
false
14. Clone() Set[T]
• 功能：深拷贝当前集合，返回一个独立的新集合
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    // 克隆后的集合与原集合无引用关系
    set2 := set1.Clone()
    set1.Clear()
    fmt.Println("set1 值: ", set1)
    fmt.Println("set2 值: ", set2)
}
go run main.go
set1 值:  map[]
set2 值:  map[address:{} age:{} apple:{} name:{}]
集合转换
1. UnsortedList() []T
• 功能：将集合转换为无序的切片
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    slice := set1.UnsortedList()
    fmt.Println(slice)
}
go run main.go
[name age address apple]
2. SortedList() []T
• 功能：将集合转换为按自然顺序排序的切片（仅支持可排序类型，如 string、int）
当调用 SortedList() []T 方法且集合的泛型类型 T 为 string 时，排序规则是区分大小写的字典序（lexicographical order），即基于字符的 Unicode 码点（ASCII 值）进行排序

排序规则详解

• 区分大小写
• 大写字母（如 A、B）会排在小写字母（如 a、b）之前，因为大写字母的 Unicode 码点更小;例如："Apple" < "apple" < "Banana" < "banana"
• 字典序（逐字符比较）
• 按字符串的每个字符依次比较 Unicode 码点，直到找到差异;例如"abc" < "abd"（第三个字符 c < d）; "a" < "aa"（长度短者优先）
• 自然顺序
• 等同于 Go 标准库 sort.Strings() 的默认行为
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    slice := set1.SortedList()
    fmt.Println(slice)
}
go run main.go
[address age apple name]
其他工具方法
1. Clear()
• 功能：清空集合中的所有元素
func main() {
    set1 := set.New[string]("name", "age", "address", "apple")
    fmt.Println(set1.Clear())
}
go run main.go
map[]