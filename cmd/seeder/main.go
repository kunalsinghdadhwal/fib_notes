package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kunalsinghdadhwal/fib_notes/db"
	"github.com/kunalsinghdadhwal/fib_notes/seeder"
)

const (
	version = "1.0.0"
	appName = "FibNotes Database Seeder"
)

func main() {
	if err := db.Connect(); err != nil {
		fmt.Printf("Error: Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Check for specific commands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "clear":
			handleClear()
			return
		case "help", "-h", "--help":
			printHelp()
			return
		case "version", "-v", "--version":
			printVersion()
			return
		case "-n":
			// Handle -n flag directly
			handleSeed()
			return
		}
	}

	// If no arguments or unrecognized command, default to seed with parsing all args
	handleSeed()
}

func handleSeed() {
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)
	userCount := seedCmd.Int("n", 10, "Number of users to create (default: 10, max: 100)")

	seedCmd.Parse(os.Args[1:])

	if *userCount < 1 || *userCount > 100 {
		fmt.Printf("Error: User count must be between 1 and 100\n")
		os.Exit(1)
	}

	fmt.Printf("Starting database seeding...\n\n")

	if err := seeder.SeedDatabase(*userCount); err != nil {
		fmt.Printf("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database seeded successfully with %d users!\n", *userCount)
}

func handleClear() {
	fmt.Printf("Clearing database...\n")

	if err := seeder.ClearDatabase(); err != nil {
		fmt.Printf("Failed to clear database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database cleared successfully!\n")
}

func printHelp() {
	fmt.Printf("%s v%s\n\n", appName, version)
	fmt.Printf("USAGE:\n")
	fmt.Printf("  seed [options]      Default command - seed the database\n")
	fmt.Printf("  seed <command>\n\n")
	fmt.Printf("OPTIONS:\n")
	fmt.Printf("  -n <count>          Number of users to create (default: 10, max: 100)\n\n")
	fmt.Printf("COMMANDS:\n")
	fmt.Printf("  clear               Clear all users and notes from the database\n")
	fmt.Printf("  help                Show this help message\n")
	fmt.Printf("  version             Show version information\n\n")
	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  seed                # Create 10 users with 5+ notes each\n")
	fmt.Printf("  seed -n 25          # Create 25 users with 5+ notes each\n")
	fmt.Printf("  seed clear          # Clear all data\n")
	fmt.Printf("  seed help           # Show this help\n\n")
	fmt.Printf("NOTES:\n")
	fmt.Printf("  - Each user will have 5-10 random notes\n")
	fmt.Printf("  - User names are authentic Indian names\n")
	fmt.Printf("  - Email addresses are generated based on names\n")
	fmt.Printf("  - Default password for all seeded users is 'password123'\n")
}

func printVersion() {
	fmt.Printf("%s v%s\n", appName, version)
	fmt.Printf("Built for FibNotes API\n")
}
