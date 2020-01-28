package canopen

func DicMustParse(a *DicObjectDic, err error) *DicObjectDic {
	if err != nil {
		panic(err)
	}

	return a
}
