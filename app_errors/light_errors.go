package app_errors

import (
	"fmt"
	"os"
	"time"
)

type LightError struct {
	Operation   string
	Description string
}

func (le LightError) Error() string {
	errorByte := []byte(fmt.Sprintf("L'opération '%s' à échouée à %v", le.Operation, time.Now()))
	err := os.WriteFile("../light_errors.txt", errorByte, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
	}

	if len(le.Description) > 0 {
		return fmt.Sprintf("L'opération '%s' à échouée à %v. Détail: %s ", le.Operation, time.Now(), le.Description)
	}

	return fmt.Sprintf("L'opération '%s' à échouée à %v. ", le.Operation, time.Now())
}
