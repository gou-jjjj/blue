package filename

import (
	"blue/common/rand"
	"fmt"
	"time"
)

func StorageName(base string, idx int) string {
	return fmt.Sprintf("%s/storage_db_%d", base, idx)
}

func LogName() string {
	r := rand.RandString(8)
	return fmt.Sprintf("%s-%s", time.Now().Format("2006:01:02-15:04:05"), r)
}
