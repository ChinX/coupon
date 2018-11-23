package model

func (a *Activity) List(from, count int) (int64, []*Activity) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(a)
	offset := (from - 1)*count
	if n == 0 || offset >= int(n){
		return 0, list
	}

	engine.Desc("id").Limit(count, offset).Find(&list)
	return n, list
}

func (a *Bargain) List(from, count int) (int64, []*Bargain) {
	list := make([]*Bargain, 0)
	n, _ := engine.Count(a)
	offset := (from - 1)*count
	if n == 0 || offset >= int(n){
		return 0, list
	}

	engine.Desc("id").Limit(count, offset).Find(&list)
	return n, list
}

