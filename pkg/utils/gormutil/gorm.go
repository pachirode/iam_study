package gormutil

const DefaultLimit = 1000

type LimitAndOffset struct {
	Offset int
	Limit  int
}

func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	o, l := 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
