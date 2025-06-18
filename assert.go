package main

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func assertResultError[T any](result T, err error) T {
	assertError(err)
	return result
}

func assertCondition[T any](condition bool, exception func() T) {
	if !condition {
		panic(exception())
	}
}
