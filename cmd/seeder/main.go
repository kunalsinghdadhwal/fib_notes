package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "seed":
		handleSeed()
	case "clear":
		handleClear()
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		printVersion()
	default:
		fmt.Printf("Error: Unknown command '%s'\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func handleSeed() {
	fmt.Printf("Starting database seeding...\n\n")

	userCount := 10
	if len(os.Args) > 2 {
		if count, err := parseUserCount(os.Args[2]); err == nil {
			userCount = count
		} else {
			fmt.Printf("Warning: Invalid user count '%s', using default: %d\n", os.Args[2], userCount)
		}
	}

	if err := seeder.SeedDatabase(userCount); err != nil {
		fmt.Printf("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database seeded successfully with %d users!\n", userCount)
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
	fmt.Printf("  seeder <command> [options]\n\n")
	fmt.Printf("COMMANDS:\n")
	fmt.Printf("  seed [count]    Seed the database with dummy users and notes\n")
	fmt.Printf("                  count: Number of users to create (default: 10, max: 100)\n")
	fmt.Printf("  clear           Clear all users and notes from the database\n")
	fmt.Printf("  help            Show this help message\n")
	fmt.Printf("  version         Show version information\n\n")
	fmt.Printf("EXAMPLES:\n")
	fmt.Printf("  seeder seed           # Create 10 users with 5+ notes each\n")
	fmt.Printf("  seeder seed 25        # Create 25 users with 5+ notes each\n")
	fmt.Printf("  seeder clear          # Clear all data\n")
	fmt.Printf("  seeder help           # Show this help\n\n")
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

func parseUserCount(countStr string) (int, error) {
	var count int
	if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil {
		return 0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("user count must be at least 1")
	}
	if count > 100 {
		return 0, fmt.Errorf("user count cannot exceed 100")
	}
	return count, nil
}
