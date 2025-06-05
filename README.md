# SCRAM-SHA-256 Password Generator

A simple command-line tool for generating SCRAM-SHA-256 password hashes. This tool is useful for creating secure password hashes for authentication systems that use the SCRAM-SHA-256 mechanism.

## Installation

```bash
go install github.com/SonOfBytes/scram-sha-256@latest
```

## Usage

### Interactive Mode (Default)
Prompt for password input:
```bash
scram-sha-256
```

### Stdin Mode
Read password from stdin:
```bash
echo 'mypassword' | scram-sha-256 -stdin
```

### Custom Iterations
Specify the number of PBKDF2 iterations:
```bash
scram-sha-256 -i 8192
```

### Help
Display usage information:
```bash
scram-sha-256 -help
```

## Options

| Flag | Description |
|------|-------------|
| `-stdin` | Read password from stdin instead of prompting |
| `-h`, `-help` | Show help message |
| `-i`, `-iterations` | Number of PBKDF2 iterations (default: 4096) |

## Output Format

The tool outputs a SCRAM-SHA-256 hash in the following format:
```
SCRAM-SHA-256$<iterations>:<salt>$<stored_key>:<server_key>
```

Example output:
```
SCRAM-SHA-256$4096:2JOi/fCI9fLQbBIvOGYjZg==$xT/AX0HenOPJXGjpyhuSPpXwBlrRtO3NglO9pJs8ICg=:4Eu09+DaYd3rb27kkrCKQloHGoR6X9pqVgzSxDakeGU=
```

## Examples

### Basic usage with password prompt
```bash
$ scram-sha-256
Password: [hidden input]
SCRAM-SHA-256$4096:abc123...=$def456...:ghi789...
```

### Using with pipes
```bash
$ echo 'mypassword' | scram-sha-256 -stdin
SCRAM-SHA-256$4096:xyz789...=$uvw012...:rst345...
```

### Custom iterations for higher security
```bash
$ echo 'mypassword' | scram-sha-256 -stdin -i 10000
SCRAM-SHA-256$10000:mno678...=$pqr901...:stu234...
```

## Security Features

- **Secure password input**: Interactive mode uses terminal password masking
- **Random salt generation**: Each hash uses a cryptographically secure random salt
- **Configurable iterations**: Adjustable PBKDF2 iteration count for computational hardness
- **Input validation**: Validates UTF-8 encoding and non-empty passwords
- **Memory safety**: Uses Go's built-in security features

## Technical Details

- **Algorithm**: SCRAM-SHA-256 as defined in RFC 7677
- **Key derivation**: PBKDF2 with SHA-256
- **Default iterations**: 4096 (configurable)
- **Salt length**: 16 bytes
- **Key length**: 32 bytes
- **Dependencies**: Only uses Go standard library and golang.org/x packages

## Error Handling

The tool provides clear error messages and appropriate exit codes:

- **Exit code 0**: Success
- **Exit code 1**: Error (invalid input, file I/O error, etc.)

Common error scenarios:
- Empty password input
- Invalid UTF-8 in password
- I/O errors when reading from stdin
- Invalid iteration count (< 1)

## Building from Source

```bash
git clone https://github.com/SonOfBytes/scram-sha-256.git
cd scram-sha-256
go build
```

## License

This project is open source. See the repository for license details.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests on GitHub.