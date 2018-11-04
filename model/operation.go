package model

func Insert(tb interface{}) bool {
	n, err := engine.Insert(tb)
	if err != nil || n==0 {
		return false
	}
	return true
}

func Update(tb interface{}) bool {
	n, err := engine.Update(tb)
	if err != nil || n==0 {
		return false
	}
	return true
}

func Find(tb interface{}) bool {
	err := engine.Find(tb)
	if err != nil {
		return false
	}
	return true
}

func Get(tb interface{}) bool {
	ok, err := engine.Get(tb)
	if err != nil || !ok {
		return false
	}
	return true
}

func Delete(tb interface{}) bool {
	n, err := engine.Delete(tb)
	if err != nil || n==0 {
		return false
	}
	return true
}
