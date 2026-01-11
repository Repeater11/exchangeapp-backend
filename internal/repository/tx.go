package repository

import "gorm.io/gorm"

type Transactioner interface {
	Transaction(fn func(tx *gorm.DB) error) error
}

type ThreadRepoWithTx interface {
	WithTx(tx *gorm.DB) ThreadRepository
}

type ThreadLikeRepoWithTx interface {
	WithTx(tx *gorm.DB) ThreadLikeRepository
}

type ThreadLikeCounterWithTx interface {
	WithTx(tx *gorm.DB) ThreadLikeCounter
}
