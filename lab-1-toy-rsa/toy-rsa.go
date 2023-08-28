package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// RSA struct to hold the prime numbers p and q, their product n, and the totient t
type RSA struct {
	p, q, n, t *big.Int
}

// NewRSA constructs an RSA instance with given prime numbers p and q
func NewRSA(p, q int64) *RSA {
	pBig := big.NewInt(p)             // Convert p to big.Int
	qBig := big.NewInt(q)             // Convert q to big.Int
	n := new(big.Int).Mul(pBig, qBig) // Compute n = p * q
	t := new(big.Int).Mul(            // Compute t = (p-1) * (q-1)
		new(big.Int).Sub(pBig, big.NewInt(1)),
		new(big.Int).Sub(qBig, big.NewInt(1)))
	return &RSA{pBig, qBig, n, t}
}

// PubKey computes and returns the public key for the RSA instance
func (rsa *RSA) PubKey() *big.Int {
	for i := int64(2); i < rsa.t.Int64(); i++ {
		if new(big.Int).GCD(nil, nil, big.NewInt(i), rsa.t).Cmp(big.NewInt(1)) == 0 {
			return big.NewInt(i) // Return i as public key if GCD(i, t) is 1
		}
	}
	return big.NewInt(0) // Should not reach here
}

// PrivKey computes and returns the private key for the RSA instance
func (rsa *RSA) PrivKey() *big.Int {
	e := rsa.PubKey()    // Get public key
	j := big.NewInt(0)   // Initialize j to 0
	one := big.NewInt(1) // Define one as big.Int value 1
	for {
		if new(big.Int).Mod(new(big.Int).Mul(j, e), rsa.t).Cmp(one) == 0 {
			return j // Return j as private key if (j * e) mod t is 1
		}
		j.Add(j, one) // Increment j by 1
	}
}

// encryptInteger encrypts an integer using the RSA public key
func encryptInteger(rsa *RSA, mes int) *big.Int {
	e := rsa.PubKey()                                        // Get public key
	ct := new(big.Int).Exp(big.NewInt(int64(mes)), e, rsa.n) // Compute ciphertext as (mes^e) mod n
	return ct                                                // Return ciphertext
}

// decryptInteger decrypts an integer using the RSA private key
func decryptInteger(rsa *RSA, ct *big.Int) *big.Int {
	d := rsa.PrivKey()                    // Get private key
	mes := new(big.Int).Exp(ct, d, rsa.n) // Compute message as (ct^d) mod n
	return mes                            // Return decrypted message
}

// signInteger signs an integer (hash) using the RSA private key
func signInteger(rsa *RSA, hash int) *big.Int {
	d := rsa.PrivKey()
	sig := new(big.Int).Exp(big.NewInt(int64(hash)), d, rsa.n) // Compute signature as (hash^d) mod n
	return sig                                                 // Return signature
}

// verifySignedInteger verifies a signed integer using the RSA public key
func verifySignedInteger(rsa *RSA, sig *big.Int) *big.Int {
	e := rsa.PubKey()
	hash := new(big.Int).Exp(sig, e, rsa.n) // Compute hash as (sig^e) mod n
	return hash                             // Return hash
}

// Main function to setup signal handling and run the exercise function
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go exercise()

	<-ctx.Done()
}

// Exercise function to handle the logic for interacting with the user,
// encrypting and decrypting messages, and signing / verifying signatures
// using the defined RSA functions above.
func exercise() {
	r := bufio.NewReader(os.Stdin)
	toyRSA := NewRSA(53, 59)

	for {
		fmt.Print("Enter a plain-text message: ")

		message, _ := r.ReadString('\n')
		message = strings.TrimSpace(message)
		fmt.Println()

		fmt.Println(
			`🔐 The program will now encrypt the provided message using the encryptInteger function
   which takes each character of the message, converts it to an integer,and encrypts it
   using the RSA public key.`)
		waitUserInput(r)

		// Encrypting each character of the message
		var encryptedMessage []*big.Int
		for _, c := range message {
			num := int(c)
			encryptedInt := encryptInteger(toyRSA, num)
			encryptedMessage = append(encryptedMessage, encryptedInt)
		}

		fmt.Printf("The encrypted message is: %v\n\n", encryptedMessage)
		fmt.Println(
			`🔓 The program will now decrypt your provided message using the decryptInteger function,
   which takes each encrypted integer and decrypts it using the RSA private key, converting
   it back to the original character.`)
		waitUserInput(r)

		// Decrypting each encrypted integer
		var decryptedMessage string
		for _, i := range encryptedMessage {
			num := decryptInteger(toyRSA, i)
			character := string(rune(num.Int64()))
			decryptedMessage += character
		}

		fmt.Printf("The decrypted message is: %v\n\n", decryptedMessage)
		fmt.Println(
			`🔏 The program will now generate the message signature using the signInteger function
   which it takes the SHA 256 hash of the message and signs it using the RSA private key.`)
		waitUserInput(r)

		// Signing the message
		hasher := sha256.New()
		hasher.Write([]byte(message))
		messageHash := fmt.Sprintf("%x", hasher.Sum(nil))

		// Signing each character of the message hash
		var messageSignature []*big.Int
		for _, c := range messageHash {
			num := int(c)
			sig := signInteger(toyRSA, num)
			messageSignature = append(messageSignature, sig)
		}

		fmt.Printf("The message signature is: %v\n\n", messageSignature)
		fmt.Println(
			`🔏🔍 The program will now generate the message hash using the verifySignedInteger function
     which takes the signature and verifies it using the RSA public key, deriving the original hash.`,
		)
		waitUserInput(r)

		// Verifying the message signature
		var hashComparison string
		for _, i := range messageSignature {
			num := verifySignedInteger(toyRSA, i)
			character := string(rune(num.Int64()))
			hashComparison += character
		}

		fmt.Printf("The message hash is: %v\n", messageHash)
		fmt.Printf("The hash derived from the message signature is: %v\n\n", hashComparison)

		fmt.Println(
			`💨 To illustrate the effect of changing a single character to the encrypted output,
   the program will now append a random character to your provided message.`)
		waitUserInput(r)

		// Appending a random uppercase letter to create message2
		message2 := message + string(rune(rand.Intn(26)+65))

		fmt.Printf("The plain-text message with one character appended is: %v\n\n", message2)
		fmt.Println(
			`🔏 The program will now generate the message signature for the edited message, it takes
   the SHA256 hash of the edited message and signs it using the RSA private key.`,
		)
		waitUserInput(r)

		// Signing the edited message
		hasher2 := sha256.New()
		hasher2.Write([]byte(message2))
		messageHash2 := fmt.Sprintf("%x", hasher2.Sum(nil))

		// Signing each character of the edited message hash
		var messageSignature2 []*big.Int
		for _, c := range messageHash2 {
			num := int(c)
			sig := signInteger(toyRSA, num)
			messageSignature2 = append(messageSignature2, sig)
		}

		fmt.Printf("The message signature with only one character appended is: %v\n\n\n", messageSignature2)

		fmt.Println(
			`🎉🔒🔮 Congratulations! You've gone through the basics of the RSA algorithm
       and are one step closer to demystifying cryptography! 🔮🔒🎉

> Would you like to go through the process again?
  Press Enter to continue or Ctrl+C to exit.`)
		_, _ = r.ReadString('\n')
		fmt.Print("\n\n")
	}
}

// waitUserInput wait until the user press Enter.
func waitUserInput(r *bufio.Reader) {
	fmt.Println()
	fmt.Println("> Press Enter to continue.")
	_, _ = r.ReadString('\n')
}
