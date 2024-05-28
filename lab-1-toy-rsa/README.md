# Toy RSA Algorithm Lab: Understanding Asymmetric Encryption

## Introduction

Welcome to the Toy RSA Algorithm Lab! In this hands-on exercise, you'll explore the principles of asymmetric encryption by implementing a simple or 'toy' RSA algorithm using Go. You'll learn how to encrypt, decrypt, sign, and verify messages, gaining a deeper understanding of cryptography.

> Disclaimer - this code is to demonstrate the theory of RSA and should not be used in production. Golang has secure cryptographic libraries which can be imported and used for performing these operations. In addition, Elliptic curve algorithms are now considered more secure in practice for both key exchange and digital signatures (e.g. see this [Trail of Bits blog post](https://blog.trailofbits.com/2019/07/08/fuck-rsa/)). We are simply using RSA as an example to demonstrate the principles of asymmetric cryptography, given that the mathematical details are easier to understand!

## Background

Asymmetric encryption is achieved by algorithms which are easy to compute in one direction, but prohibitively difficult to compute in reverse. Algorithms achieve this through the creation of a key pair consisting of a public and private key. The public part can be safely shared with others, but the private part must be kept secret. Asymmetric algorithms can be used for:

- Encryption: a person or process with access to a third party's public key can encrypt data, which can only be decrypted by an entity in possession of the corresponding private key;
- Digital Signatures: any person or process with access to the private key can create a digital signature (using an appropriate hashing algorithm in combination with asymmetric cryptography), which can be verified by a process with access to the public key;
- Key Exchange: two parties could use an algorithm such as the Diffie-Hellman key exchange scheme to agree on a shared secret by separately performing an offline computation involving their private key, and an authentic copy of the other party's public key.

RSA is based on the principle that if we take two large, prime integers `p` and `q` and multiply them to form a 'modulus' `n`, the multiplication operation is computationally easy, but factoring `n` to recover `p` and `q` is computationally infeasible given the amount of computing power / time it would take. The modulus `n` is used by both the public and private keys and provides the link between them. The length of `n` in bits is the 'key length'.

Along with the computation of `n`, the RSA algorithm requires the computation of the totient $t=(p-1)(q-1)$. We now have all the information we need to derive the public and private key:

- The 'public exponent' `e` is chosen to be relatively prime to `t`. The public key is then made up of `n` and `e`
- The 'private exponent' `d` is computed such that $de \mod t = 1$. `mod` stands for the modulo operation, which gives the remainder when two numbers are divided, e.g. $8 \mod 3 = 2$. The private key is then made up of `n` and `d`.

If Alice wants to send a secret integer `M` to Bob, she can encrypt the integer by performing the following operation using Bob's public key $(n,e)$: $C = M^e \mod n$. Bob can now decrypt the ciphertext `C` by computing $M = C^d \mod n$

## Code Overview

We have provided a toy implementation of the RSA algorithm (using some very small prime numbers!) in [toy-rsa.go](toy-rsa.go). Take a look at the code, which can be broken down as follows.

### Step 1: Setting Up the RSA Struct

The `RSA` struct is the core of our implementation. Here's a breakdown of its components:

- **Initialization (`NewRSA` function)**: Takes two prime numbers, `p` and `q`, and calculates `n = p * q` and `t = (p - 1) * (q - 1)`.

- **Public Key (`PubKey` method)**: Finds a number `e` that is relatively prime to `t`.

- **Private Key (`PrivKey` method)**: Calculates the private key `d` such that `(d * e) % t = 1`.

### Step 2: Encryption and Decryption

- **Encrypting an Integer (`encryptInteger` function)**: Takes a message as an integer and encrypts it using the public key.

- **Decrypting an Integer (`decryptInteger` function)**: Takes an encrypted integer and decrypts it using the private key.

### Step 3: Signing and Verifying Messages

- **Signing an Integer (`signInteger` function)**: Signs an integer (e.g., a hash) using the private key.

- **Verifying a Signed Integer (`verifySignedInteger` function)**: Verifies a signed integer using the public key.

### Step 4: Putting It All Together

The provided code includes examples of:

- Encrypting and decrypting a user-input message.
- Signing a message and verifying the signature.
- Demonstrating how a one-character change (e.g., appending a character) in the message leads to a vastly different hash and signature.

## Running the Lab

1. Ensure you have [Go](https://go.dev/) installed on your system.
2. Open a terminal or command prompt.
3. Clone this repo and navigate to the directory containing the `toy-rsa.go` file.
4. Run the command `go run toy-rsa.go`.
5. Enter a plain-text message when prompted and hit Enter when prompted to progress through each step.
6. Observe the output, including the encrypted message, decrypted message, and message signature.
7. To exit the user-input loop, press `Ctrl + C` or hit Enter to continue.

### In-line References

This lab provides rich in-line comments to help de-mystify cryptography, specifically how the RSA algorithm works to encrypt, decrypt, sign, and verify messages. These in-line comments help break down the code into bite-size bits. As such, it is recommended that students pause when given terminal prompts and reference the specific functions / code blocks that the prompt calls out to complete the next step of the cryptographic process.

## Conclusion

This lab provides a hands-on introduction to the RSA algorithm and asymmetric encryption using Go. By working through the code and running the provided examples, you'll gain valuable insights into encryption, decryption, signing, and verification. The interactive user-input loop allows for experimentation with different messages and showcases the impact of even minor changes to the text.

Feel free to experiment with different messages, prime numbers, or even extend the code with additional features. Happy coding!
