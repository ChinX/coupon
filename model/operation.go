package model

func (a *Activity) List(from, count int) (int64, []*Activity) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(a)
	offset := (from - 1) * count
	if n == 0 || offset >= int(n) {
		return 0, list
	}

	engine.Desc("id").Limit(count, offset).Find(&list)
	return n, list
}

func (t *Task) Last() *Task {
	list := make([]*Task, 0)
	engine.Desc("id").Where("user_id=?", t.UserID).Limit(1).Find(&list)
	if len(list) == 0 {
		return t
	}
	return list[0]
}

func (b *Bargain) List(from, count int) (int64, []*Bargain) {
	list := make([]*Bargain, 0)
	n, _ := engine.Count(b)
	offset := (from - 1) * count
	if n == 0 || offset >= int(n) {
		return 0, list
	}

	engine.Desc("id").Where("task_id=?", b.TaskID).Limit(count, offset).Find(&list)
	return n, list
}
