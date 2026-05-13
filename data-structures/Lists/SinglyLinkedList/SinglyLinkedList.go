package singlylinkedlist

// Element là một node trong danh sách liên kết.
type Element struct {
	next  *Element
	Value any
}

// Next trả về phần tử kế tiếp.
func (e *Element) Next() *Element {
	return e.next
}

// List đại diện cho danh sách liên kết đơn.
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
	if at == nil {
		// Chèn vào đầu
		e.next = l.head
		l.head = e
		if l.tail == nil {
			l.tail = e
		}
	} else {
		// Chèn vào sau at
		e.next = at.next
		at.next = e
		if at == l.tail {
			l.tail = e
		}
	}
	l.len++
	return e
}

// PushFront thêm vào đầu.
func (l *List) PushFront(v any) *Element {
	l.lazyInit()
	return l.insert(&Element{Value: v}, nil)
}

// PushBack thêm vào cuối.
func (l *List) PushBack(v any) *Element {
	l.lazyInit()
	return l.insert(&Element{Value: v}, l.tail)
}

// InsertAfter chèn v vào sau mark.
func (l *List) InsertAfter(v any, mark *Element) *Element {
	return l.insert(&Element{Value: v}, mark)
}

// Remove xóa phần tử e khỏi list.
func (l *List) Remove(e *Element) any {
	if l.head == e {
		l.head = e.next
		if l.head == nil {
			l.tail = nil
		}
	} else {
		prev := l.head
		for prev != nil && prev.next != e {
			prev = prev.next
		}
		if prev != nil {
			prev.next = e.next
			if e == l.tail {
				l.tail = prev
			}
		} else {
			// Không tìm thấy e trong list
			return nil
		}
	}

	val := e.Value
	e.next = nil
	l.len--
	return val
}

// MoveToFront di chuyển e lên đầu.
func (l *List) MoveToFront(e *Element) {
	if l.head == e {
		return
	}
	val := l.Remove(e)
	if val != nil {
		l.PushFront(val)
	}
}

// MoveToBack di chuyển e xuống cuối.
func (l *List) MoveToBack(e *Element) {
	if l.tail == e {
		return
	}
	val := l.Remove(e)
	if val != nil {
		l.PushBack(val)
	}
}
