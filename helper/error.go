package helper

func PanicfIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
