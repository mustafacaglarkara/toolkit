package toolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnoprsştuvyzABCDEFGHIJKLMNOPRSŞTUVYZ0123456789_+"

// Tools is the type used to instantiate this module. Any variable of this will access
// to all the method with the reciever *Tools
type Tools struct{}

// RandomString returns n kadar karma karakter üretip döner
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
