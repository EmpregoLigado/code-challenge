//Package null implements the cipher interface for testing purposes only
package null

//Cipher implements the crypt.Cipher interface
type Cipher struct{}

//Encrypt implements the encrypt func of crypt.Cipher
func (pk Cipher) Encrypt(message string) (string, error) {
	return message, nil
}

//Decrypt implements the decrypt func of crypt.Cipher
func (pk Cipher) Decrypt(message string) (string, error) {
	return message, nil
}
