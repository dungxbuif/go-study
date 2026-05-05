package singlylinkedlist

type Element struct {
	next  *Element
	list  *List
	Value any
}

type List struct {
	head, tail *Element
	len        int
}


func (l *List) Init() *List {
	l.head = nil
	l.tail = nil
	l.len = 0
	return l
}

func New() *List {
	return new(List).Init()
}

func (e *Element) Next() *Element {
	if e.list == nil {
		return nil
	}
	return e.next
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Front() *Element {
	return l.head
}

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
	if mark.list != l {
		return nil
	}
	return l.insert(&Element{Value: v}, mark)
}

// Remove xóa phần tử e khỏi list.
func (l *List) Remove(e *Element) any {
	if e.list != l {
		return nil
	}

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
		}
	}

	e.next = nil
	e.list = nil
	l.len--
	return e.Value
}

func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.head == e {
		return
	}
	l.Remove(e)
	l.PushFront(e.Value)
}

func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.tail == e {
		return
	}
	l.Remove(e)
	l.PushBack(e.Value)
}

