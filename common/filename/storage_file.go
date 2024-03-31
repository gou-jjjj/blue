package filename

import "fmt"

func StorageName(base string, idx int) string {
	return fmt.Sprintf("%s/storage_db_%d", base, idx)
}
