package login

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func HandleConfigCommand() {
	name := flag.String("name", "User", "Please type your name.")
	email := flag.String("email", "user@example.com", "Please type your email.")

	// Parse the flags
	flag.CommandLine.Parse(os.Args[2:])

	// Get the set of defined flags
	validFlags := make(map[string]struct{})
	flag.VisitAll(func(f *flag.Flag) {
		validFlags[f.Name] = struct{}{}
	})

	// Check for unknown flags
	unknownFlag := false
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flagName := strings.TrimPrefix(arg, "-")
			if strings.Contains(flagName, "=") {
				flagName = strings.SplitN(flagName, "=", 2)[0]
			}
			if _, valid := validFlags[flagName]; !valid {
				unknownFlag = true
				fmt.Fprintf(os.Stderr, "Unknown option: %s\n", flagName)
			}
		}
	}

	// If there's an unknown flag, print the valid options and exit
	if unknownFlag {
		fmt.Fprintln(os.Stderr, "Valid options are:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s: %s\n", f.Name, f.Usage)
		})
		os.Exit(1)
	}

	// Check if at least one valid flag was set
	atLeastOneFlagSet := false
	flag.Visit(func(f *flag.Flag) {
		atLeastOneFlagSet = true
	})

	if !atLeastOneFlagSet {
		fmt.Fprintln(os.Stderr, "Error: At least one valid option must be provided.")
		fmt.Fprintln(os.Stderr, "Valid options are:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s: %s\n", f.Name, f.Usage)
		})
		os.Exit(1)
	}

	emailValue := *email
	if !strings.Contains(emailValue, "@") || strings.HasPrefix(emailValue, "@") || strings.HasSuffix(emailValue, "@") {
		fmt.Fprintln(os.Stderr, "Error: Invalid email address. Email must contain '@' and it cannot be the first or last character.")
		os.Exit(1)
	}

	// Final parsing to ensure all flags are processed correctly
	flag.Parse()

	// Print the configured login
	fmt.Println("Login configured:")
	fmt.Printf("- Name: %s\n", *name)
	fmt.Printf("- Email: %s\n", *email)
}
