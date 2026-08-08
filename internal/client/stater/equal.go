package stater

import "fmt"

func Equal(a *State, b *State) bool {
	if a == nil || b == nil {
		fmt.Println("nil")
		return false
	}

	if a.SnapshotID != b.SnapshotID {
		fmt.Printf("snapshot not equal, a: %d, b: %d\n", a.SnapshotID, b.SnapshotID)
		return false
	}

	if len(a.FileStates) != len(b.FileStates) {
		fmt.Println("file counts not equal")
		return false
	}

	for k, fsa := range a.FileStates {
		fsb, ok := b.FileStates[k]
		if !ok {
			fmt.Println("file in a does not exist in b")
			return false
		}

		if fsa.Size != fsb.Size {
			fmt.Println("sizes dont match")
			return false
		}

		if fsa.ModTime != fsb.ModTime {
			fmt.Println("modtime dont match")
			return false
		}
	}

	return true
}
