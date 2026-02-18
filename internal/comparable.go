package internal

type Comparable interface {
	Compare(other any) (int, bool)
}
