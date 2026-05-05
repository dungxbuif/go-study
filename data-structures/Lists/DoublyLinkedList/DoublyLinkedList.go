package doublylinkedlist

// Element là một node trong danh sách liên kết đôi.
type Element struct {
	next, prev *Element
	list       *List
	Value      any
}

// Next trả về phần tử kế tiếp.
func (e *Element) Next() *Element {
	if e.list == nil {
		return nil
	}
	return e.next
}

// Prev trả về phần tử đứng trước.
func (e *Element) Prev() *Element {
	if e.list == nil {
		return nil
	}
	return e.prev
}

// List đại diện cho danh sách liên kết đôi.
type List struct {
	head, tail *Element
	len        int
}

// Init khởi tạo hoặc làm sạch danh sách.
func (l *List) Init() *List {
	l.head = nil
	l.tail = nil
	l.len = 0
	return l
}

// New tạo một danh sách mới.
func New() *List {
	return new(List).Init()
}

// Len trả về số lượng phần tử.
func (l *List) Len() int {
	return l.len
}

// Front trả về phần tử đầu tiên.
func (l *List) Front() *Element {
	return l.head
}

// Back trả về phần tử cuối cùng.
func (l *List) Back() *Element {
	return l.tail
}

func (l *List) lazyInit() {
	if l.head == nil && l.len == 0 {
		l.Init()
	}
}

// insert chèn phần tử e vào sau phần tử at.
func (l *List) insert(e, at *Element) *Element {
	e.list = l
	if at == nil {
		// Chèn vào đầu
		e.next = l.head
		e.prev = nil
		if l.head != nil {
			l.head.prev = e
		}
		l.head = e
		if l.tail == nil {
			l.tail = e
		}
	} else {
		// Chèn vào sau at
		e.next = at.next
		e.prev = at
		if at.next != nil {
			at.next.prev = e
		}
		at.next = e
		if at == l.tail {
			l.tail = e
		}
	}
	l.len++
	return e
}

// PushFront thêm vào đầu. O(1)
func (l *List) PushFront(v any) *Element {
	l.lazyInit()
	return l.insert(&Element{Value: v}, nil)
}

// PushBack thêm vào cuối. O(1)
func (l *List) PushBack(v any) *Element {
	l.lazyInit()
	return l.insert(&Element{Value: v}, l.tail)
}

// InsertAfter chèn v vào sau mark. O(1)
func (l *List) InsertAfter(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insert(&Element{Value: v}, mark)
}

// InsertBefore chèn v vào trước mark. O(1)
func (l *List) InsertBefore(v any, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	return l.insert(&Element{Value: v}, mark.prev)
}

// Remove xóa phần tử e khỏi list. O(1)
func (l *List) Remove(e *Element) any {
	if e.list != l {
		return nil
	}

	if e.prev != nil {
		e.prev.next = e.next
	} else {
		l.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		l.tail = e.prev
	}

	val := e.Value
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--
	return val
}

// MoveToFront di chuyển e lên đầu. O(1)
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.head == e {
		return
	}
	l.Remove(e)
	l.PushFront(e.Value)
}

// MoveToBack di chuyển e xuống cuối. O(1)
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.tail == e {
		return
	}
	l.Remove(e)
	l.PushBack(e.Value)
}
