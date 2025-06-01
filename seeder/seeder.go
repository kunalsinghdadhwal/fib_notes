package seeder

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kunalsinghdadhwal/fib_notes/db"
	"github.com/kunalsinghdadhwal/fib_notes/models"
)

// Names contains authentic names for seeding
var Names = []string{
	"Arjun Kumar", "Priya Sharma", "Rahul Singh", "Anjali Patel", "Vikram Gupta",
	"Sneha Reddy", "Aditya Verma", "Kavya Nair", "Rohan Shah", "Meera Joshi",
	"Karan Malhotra", "Pooja Agarwal", "Siddharth Rao", "Riya Kapoor", "Abhishek Dubey",
	"Shreya Iyer", "Nikhil Bhatt", "Divya Saxena", "Harsh Pandey", "Ananya Singh",
	"Dev Khanna", "Ishita Bansal", "Yash Tiwari", "Nisha Chopra", "Akash Srivastava",
	"Tanvi Mehta", "Aryan Jain", "Kritika Mishra", "Varun Desai", "Shweta Kulkarni",
	"Mohit Goyal", "Sakshi Aggarwal", "Rishab Bhatia", "Aditi Singhania", "Kartik Sethi",
	"Manasi Rao", "Ayush Chandra", "Rhea Mittal", "Pranav Awasthi", "Simran Kaur",
	"Sameer Khan", "Gayatri Pillai", "Ankit Singhal", "Radhika Sinha", "Vivek Thakur",
	"Aparna Kothari", "Shubham Goel", "Megha Bajaj", "Rohit Kumar", "Swati Bhardwaj",
}

// NoteTitles contains professional note titles
var NoteTitles = []string{
	"Meeting Notes", "Project Requirements", "Weekly Goals", "Code Review Comments",
	"Bug Report Analysis", "Feature Implementation Plan", "Database Schema Design",
	"API Documentation", "Testing Strategy", "Deployment Checklist",
	"Performance Optimization", "Security Audit Notes", "Client Feedback Summary",
	"Sprint Planning", "Technical Specifications", "Architecture Overview",
	"User Story Analysis", "Error Handling Guidelines", "Configuration Management",
	"Monitoring Setup", "Backup Procedures", "Release Planning",
	"Team Standup Notes", "Research Findings", "Best Practices Documentation",
	"Troubleshooting Guide", "Integration Testing", "Code Standards",
	"Review Action Items", "System Requirements", "Database Migration Notes",
	"Performance Metrics", "Security Compliance", "DevOps Pipeline",
	"Quality Assurance", "Technical Debt Analysis", "Refactoring Plan",
	"Infrastructure Setup", "Monitoring Dashboard", "Log Analysis",
	"Capacity Planning", "Service Level Agreements", "Disaster Recovery Plan",
	"Version Control Guidelines", "CI/CD Implementation", "Documentation Standards",
	"Knowledge Transfer", "Training Materials", "Process Improvement",
	"Risk Assessment", "Technical Interview Questions", "Code Architecture Review",
}

// NoteContents contains professional note content templates
var NoteContents = []string{
	"Reviewed the current implementation and identified several areas for improvement. Need to focus on optimizing database queries and implementing proper caching mechanisms.",
	"Discussed project timeline with stakeholders. Delivery date confirmed for end of month. All features must be thoroughly tested before deployment.",
	"Analyzed user feedback from the latest release. Users are requesting better search functionality and improved mobile responsiveness.",
	"Conducted code review session with the team. Found several potential security vulnerabilities that need immediate attention.",
	"Documented the new API endpoints and their respective request/response formats. Updated the API documentation accordingly.",
	"Investigated performance bottlenecks in the application. Database connections seem to be the primary issue affecting response times.",
	"Planned the database migration strategy for the upcoming release. Need to ensure zero downtime during the migration process.",
	"Reviewed error logs from production environment. Identified recurring issues that need to be addressed in the next patch release.",
	"Outlined the testing strategy for the new feature implementation. Unit tests, integration tests, and end-to-end tests are required.",
	"Documented the deployment process and created automated scripts for streamlined releases across different environments.",
	"Analyzed system metrics and identified areas where resource utilization can be optimized for better cost efficiency.",
	"Reviewed security protocols and updated authentication mechanisms to comply with latest industry standards.",
	"Gathered requirements from business stakeholders for the upcoming quarterly release. Priority features have been identified.",
	"Conducted team retrospective meeting. Identified process improvements and action items for better collaboration.",
	"Documented technical specifications for the new microservice architecture. Defined service boundaries and communication protocols.",
	"Analyzed user journey through the application and identified pain points that need UX improvements.",
	"Reviewed third-party integrations and updated API keys and configurations for better security practices.",
	"Planned capacity scaling strategy for handling increased traffic during peak usage periods.",
	"Documented disaster recovery procedures and tested backup restoration processes to ensure data integrity.",
	"Analyzed competitor features and market trends to identify opportunities for product enhancement.",
	"Conducted technical interview sessions and documented evaluation criteria for new team members.",
	"Reviewed monitoring alerts and fine-tuned thresholds to reduce false positives while maintaining system observability.",
	"Documented lessons learned from the recent production incident and updated incident response procedures.",
	"Planned knowledge transfer sessions for critical system components to ensure team knowledge redundancy.",
	"Analyzed technical debt accumulated in the codebase and prioritized refactoring tasks for the next sprint.",
}

func SeedDatabase(userCount int) error {
	fmt.Printf("Creating %d users with random notes...\n", userCount)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	users := make([]models.User, 0, userCount)
	allNotes := make([]models.Note, 0, userCount*8)

	for i := 0; i < userCount; i++ {
		name := Names[rand.Intn(len(Names))]
		email := generateEmail(name, i)

		user := models.User{
			Name:     name,
			Email:    email,
			Password: "password123",
		}

		if err := db.DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", name, err)
		}

		users = append(users, user)

		noteCount := rand.Intn(6) + 5 // 5 to 10 notes
		for j := 0; j < noteCount; j++ {
			note := models.Note{
				UserID:  user.ID,
				Title:   NoteTitles[rand.Intn(len(NoteTitles))],
				Content: NoteContents[rand.Intn(len(NoteContents))],
			}
			allNotes = append(allNotes, note)
		}

		if (i+1)%10 == 0 || i+1 == userCount {
			fmt.Printf("  Created %d/%d users\n", i+1, userCount)
		}
	}

	fmt.Printf("Creating %d notes...\n", len(allNotes))
	if err := db.DB.CreateInBatches(allNotes, 100).Error; err != nil {
		return fmt.Errorf("failed to create notes: %w", err)
	}

	fmt.Printf("  Created %d notes successfully\n", len(allNotes))
	return nil
}

func ClearDatabase() error {
	if err := db.DB.Exec("DELETE FROM notes").Error; err != nil {
		return fmt.Errorf("failed to clear notes: %w", err)
	}

	if err := db.DB.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	fmt.Printf("  Cleared all users and notes from database\n")
	return nil
}

func generateEmail(name string, index int) string {
	emailName := ""
	for _, char := range name {
		if char >= 'A' && char <= 'Z' {
			emailName += string(char + 32)
		} else if char >= 'a' && char <= 'z' {
			emailName += string(char)
		} else if char == ' ' {
			emailName += "."
		}
	}

	// Add index to ensure uniqueness
	domains := []string{"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "company.com"}
	domain := domains[rand.Intn(len(domains))]

	return fmt.Sprintf("%s%d@%s", emailName, index+1, domain)
}

// GetStats returns database statistics
func GetStats() (int64, int64, error) {
	var userCount, noteCount int64

	if err := db.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if err := db.DB.Model(&models.Note{}).Count(&noteCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count notes: %w", err)
	}

	return userCount, noteCount, nil
}
