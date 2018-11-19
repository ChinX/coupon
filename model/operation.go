package model

func (a *Activity) List(from, count int) (int64, []*Activity) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(a)
	if n == 0 {
		return 0, list
	}
	if from == 0 {
		engine.Desc("id").Limit(count, 0).Find(&list)
	} else {
		engine.Where("id < ?", from).Desc("id").Limit(count, 0).Find(&list)
	}
	return n, list
}

