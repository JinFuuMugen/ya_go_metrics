// !build prod

package storage

func Reset() {
	defaultStorage = NewStorage()
}
