package hw04lrucache

/*
List is an interface of doubly linked list.
*/
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

/*
ListItem is an item of doubly linked list.
*/
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length    int
	listStart *ListItem
	listEnd   *ListItem
}

// NewList creates new list.
func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.listStart
}

func (l *list) Back() *ListItem {
	return l.listEnd
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.length == 0 {
		l.listStart = item
		l.listEnd = item
	} else {
		item.Next = l.listStart
		l.listStart.Prev = item
		l.listStart = item
	}
	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.length == 0 {
		l.listStart = item
		l.listEnd = item
	} else {
		item.Prev = l.listEnd
		l.listEnd.Next = item
		l.listEnd = item
	}
	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.listStart = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.listEnd = i.Prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
