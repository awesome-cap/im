package async

func Async(fun func()){
	go func() {
		fun()
	}()
}
