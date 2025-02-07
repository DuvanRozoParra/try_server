package pkg

type Queue[T any] struct {
	missions []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		missions: make([]T, 0),
	}
}

func (q *Queue[T]) Enqueue(mission T) {
	q.missions = append(q.missions, mission)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.missions) == 0 {
		return zero, false
	}
	mission := q.missions[0]
	q.missions = q.missions[1:]
	return mission, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if len(q.missions) == 0 {
		return zero, false
	}
	return q.missions[0], true
}

func (q *Queue[T]) Count() int {
	return len(q.missions)
}
