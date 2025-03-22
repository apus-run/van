/*
    // 初始化一个大小为 2 的环形缓冲区
    ringBuffer := buffer.NewRingGrowing(2)

    // 向缓冲区写入数据
    ringBuffer.WriteOne("test01")
    ringBuffer.WriteOne("test02")

    // 缓冲区满时，继续写入会扩展缓冲区
    ringBuffer.WriteOne("test03") // 此时缓存区将自动扩展为4个位置

    // 读取并消费缓冲区中的数据
    data, ok := ringBuffer.ReadOne()
    if ok {
        fmt.Println("Read:", data) // 输出: Read: test01
    }

    data, ok = ringBuffer.ReadOne()
    if ok {
        fmt.Println("Read:", data) // 输出: Read: test02
    }

    // 读取并消费缓冲区中的数据
    data, ok = ringBuffer.ReadOne()
    if ok {
        fmt.Println("Read:", data) // 输出: Read: test03
    }

    data, ok = ringBuffer.ReadOne()
    if ok {
        fmt.Println("Read:", data) // 无输出
    }
*/

package buffer

// RingGrowing is a growing ring buffer.
// Not thread safe.
type RingGrowing struct {
	data     []interface{}
	n        int // Size of Data
	beg      int // First available element
	readable int // Number of data items available
}

// NewRingGrowing constructs a new RingGrowing instance with provided parameters.
func NewRingGrowing(initialSize int) *RingGrowing {
	return &RingGrowing{
		data: make([]interface{}, initialSize),
		n:    initialSize,
	}
}

// ReadOne reads (consumes) first item from the buffer if it is available, otherwise returns false.
func (r *RingGrowing) ReadOne() (data interface{}, ok bool) {
	if r.readable == 0 {
		return nil, false
	}
	r.readable--
	element := r.data[r.beg]
	r.data[r.beg] = nil // Remove reference to the object to help GC
	if r.beg == r.n-1 {
		// Was the last element
		r.beg = 0
	} else {
		r.beg++
	}
	return element, true
}

// WriteOne adds an item to the end of the buffer, growing it if it is full.
func (r *RingGrowing) WriteOne(data interface{}) {
	if r.readable == r.n {
		// Time to grow
		newN := r.n * 2
		newData := make([]interface{}, newN)
		to := r.beg + r.readable
		if to <= r.n {
			copy(newData, r.data[r.beg:to])
		} else {
			copied := copy(newData, r.data[r.beg:])
			copy(newData[copied:], r.data[:(to%r.n)])
		}
		r.beg = 0
		r.data = newData
		r.n = newN
	}
	r.data[(r.readable+r.beg)%r.n] = data
	r.readable++
}

// Len returns the number of items in the buffer.
func (r *RingGrowing) Len() int {
	return r.readable
}

// Cap returns the capacity of the buffer.
func (r *RingGrowing) Cap() int {
	return r.n
}
