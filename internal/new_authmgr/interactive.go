package new_authmgr

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/amadeusitgroup/cds/internal/cerr"
	cg "github.com/amadeusitgroup/cds/internal/global"
	nc "github.com/amadeusitgroup/cds/internal/new_bo"
)

var isInteractive bool

func init() {
	stdinStat, _ := os.Stdin.Stat()
	isInteractive = (stdinStat.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

// DefaultPrompt returns the default password prompt message.
func DefaultPrompt() string {
	return "Enter your office LDAP user password (Hint: windows account password):"
}

// PromptPassword reads a password from stdin (interactive or piped)
// and returns it as a string.
func PromptPassword(message string) (string, error) {
	if isInteractive {
		return askForPasswordInteractive(message)
	}
	return askForPasswordFromStdin(message)
}

// PromptAndSavePassword interactively prompts for a password, creates a
// password credential, and saves it under the given key in the store.
func PromptAndSavePassword(store *Store, credKey, login, prompt string) (nc.Credential, error) {
	pwd, err := PromptPassword(prompt)
	if err != nil {
		return nc.Credential{}, fmt.Errorf("failed to read password for %q: %w", credKey, err)
	}

	cred := nc.NewPasswordCredential(login, pwd)
	if err := store.Set(credKey, cred); err != nil {
		return nc.Credential{}, fmt.Errorf("failed to save credential %q: %w", credKey, err)
	}

	return cred, nil
}

func askForPasswordInteractive(message string) (string, error) {
	fmt.Print(message)
	byteSecret, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println(cg.EmptyStr) // newline after password entry
	if err != nil {
		return cg.EmptyStr, cerr.AppendError("Unable to read password", err)
	}
	return string(byteSecret), nil
}

func askForPasswordFromStdin(message string) (string, error) {
	fmt.Print(message)
	return readLineFromStdin()
}

func readLineFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return cg.EmptyStr, cerr.NewError("Failed to acquire a new line from stdin")
}
