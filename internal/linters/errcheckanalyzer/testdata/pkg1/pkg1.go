package pkg1

//lint:file-ignore U1000 игнорируем неиспользуемый код, так как он нужен для тестирования multi checker

import "fmt"

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func errCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	mulfunc(5)           // want "expression returns unchecked error"
	res, _ := mulfunc(5) // want "assignment with unchecked error"
	fmt.Println(res)     // want "expression returns unchecked error"
	go mulfunc(5)        // want "go statement with unchecked error"
	defer mulfunc(5)     // want "defer with unchecked error"
}
