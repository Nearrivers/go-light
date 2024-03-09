package app_errors

import (
	"fmt"
	"os"
)

type RuntimeError struct {
	Err error
}

func (re RuntimeError) Error() string {
	errorByte := []byte(fmt.Sprintf("Une erreur est survenue pendant l'exécution de l'applicaiton : %s", re.Err.Error()))
	err := os.WriteFile("../runtime_errors.txt", errorByte, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("Une erreur est survenue pendant l'exécution de l'application : %s", re.Err.Error())
}
