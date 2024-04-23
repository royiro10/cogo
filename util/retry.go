package util

func MakeRetryable[T any](fn func() (T, error), retryCount int) func() (T, error) {
	return func() (T, error) {
		var err error
		var result T

		for tries := 0; tries < retryCount; tries++ {
			result, err = fn()
			if err == nil {
				return result, nil
			}
		}

		return result, err
	}
}
