package repository

type ThreadLikeCounter interface {
	IncrementLikeCount(threadID uint, delta int) error
	GetLikeCount(threadID uint) (int64, error)
}
