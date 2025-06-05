package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	defaultIterations = 4096
	saltLength       = 16
	keyLength        = 32
)

type Config struct {
	UseStdin   bool
	ShowHelp   bool
	Iterations int
}

func main() {
	config := parseFlags()

	if config.ShowHelp {
		showHelp()
		os.Exit(0)
	}

	var password string
	var err error

	if config.UseStdin {
		password, err = readPasswordFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password from stdin: %v\n", err)
			os.Exit(1)
		}
	} else {
		password, err = promptPassword()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
	}

	if err := validatePassword(password); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid password: %v\n", err)
		os.Exit(1)
	}

	hash, err := generateSCRAMSHA256(password, config.Iterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating SCRAM-SHA-256: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(hash)
}

func parseFlags() Config {
	config := Config{}
	
	flag.BoolVar(&config.UseStdin, "stdin", false, "Read password from stdin instead of prompting")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help message")
	flag.IntVar(&config.Iterations, "iterations", defaultIterations, "Number of PBKDF2 iterations")
	flag.IntVar(&config.Iterations, "i", defaultIterations, "Number of PBKDF2 iterations")
	
	flag.Parse()
	
	return config
}

func showHelp() {
	fmt.Println("SCRAM-SHA-256 Password Generator")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Printf("  %s [OPTIONS]\n", os.Args[0])
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -stdin           Read password from stdin instead of prompting")
	fmt.Println("  -h, -help        Show this help message")
	fmt.Println("  -i, -iterations  Number of PBKDF2 iterations (default: 4096)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Printf("  %s                    # Prompt for password\n", os.Args[0])
	fmt.Printf("  echo 'mypass' | %s -stdin  # Read from stdin\n", os.Args[0])
	fmt.Printf("  %s -i 8192               # Custom iterations\n", os.Args[0])
	fmt.Println()
	fmt.Println("INSTALLATION:")
	fmt.Println("  go install github.com/SonOfBytes/scram-sha-256@latest")
}

func promptPassword() (string, error) {
	fmt.Print("Password: ")
	
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	
	fmt.Println()
	return string(passwordBytes), nil
}

func readPasswordFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}
	
	return strings.TrimRight(password, "\r\n"), nil
}

func validatePassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	
	if !utf8.ValidString(password) {
		return fmt.Errorf("password must be valid UTF-8")
	}
	
	return nil
}

func generateSCRAMSHA256(password string, iterations int) (string, error) {
	if iterations < 1 {
		return "", fmt.Errorf("iterations must be at least 1")
	}
	
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	
	saltedPassword := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
	
	clientKey := hmac.New(sha256.New, saltedPassword)
	clientKey.Write([]byte("Client Key"))
	clientKeyBytes := clientKey.Sum(nil)
	
	storedKey := sha256.Sum256(clientKeyBytes)
	
	serverKey := hmac.New(sha256.New, saltedPassword)
	serverKey.Write([]byte("Server Key"))
	serverKeyBytes := serverKey.Sum(nil)
	
	saltB64 := base64.StdEncoding.EncodeToString(salt)
	storedKeyB64 := base64.StdEncoding.EncodeToString(storedKey[:])
	serverKeyB64 := base64.StdEncoding.EncodeToString(serverKeyBytes)
	
	result := fmt.Sprintf("SCRAM-SHA-256$%d:%s$%s:%s", iterations, saltB64, storedKeyB64, serverKeyB64)
	
	return result, nil
}