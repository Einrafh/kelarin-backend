package database

import (
	"fmt"
	"log"
	"os"

	"kelarin-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase loads environment variables, connects to Postgres,
// auto‑migrates models, and ensures cascade foreign key constraints.
func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found, using system environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "db"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.WorkspaceUser{},
		&models.BoardList{},
		&models.Card{},
		&models.Subtask{},
		&models.CardAssignee{},
		&models.CardAttachment{},
		&models.CardLabel{},
		&models.CardComment{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	ensureCascadeFK(db)

	DB = db
	log.Println("Database connected, migrated, and cascade constraints ensured")
}

// ensureCascadeConstraints is a universal helper that iterates over the given association names
// for a model and ensures that the foreign key constraints are recreated with ON DELETE CASCADE
// based on the model's tags.
func ensureCascadeConstraints(db *gorm.DB, model interface{}, associations []string) {
	// List any legacy constraint names that might exist for this model.
	legacyConstraints := []string{"fk_workspaces_collaborators"}

	// First, drop any legacy constraints if they exist.
	for _, legacyName := range legacyConstraints {
		if db.Migrator().HasConstraint(model, legacyName) {
			if err := db.Migrator().DropConstraint(model, legacyName); err != nil {
				log.Fatalf("Failed to drop legacy constraint %q on model %T: %v", legacyName, model, err)
			}
			log.Printf("Dropped legacy constraint %q on model %T", legacyName, model)
		}
	}

	// Now, process the associations provided.
	for _, association := range associations {
		if db.Migrator().HasConstraint(model, association) {
			if err := db.Migrator().DropConstraint(model, association); err != nil {
				log.Fatalf("Failed to drop constraint %q on model %T: %v", association, model, err)
			}
			log.Printf("Dropped existing constraint %q on model %T", association, model)
		}

		if err := db.Migrator().CreateConstraint(model, association); err != nil {
			log.Fatalf("Failed to create cascade constraint %q on model %T: %v", association, model, err)
		}
		log.Printf("Ensured cascade constraint on association %q for model %T", association, model)
	}
}

// ensureCascadeFK applies cascade constraints for models that require them.
// When adding a new model that requires ON DELETE CASCADE, add an appropriate call below.
func ensureCascadeFK(db *gorm.DB) {
	// For WorkspaceUser, ensure both the "Workspace" and "User" constraints have ON DELETE CASCADE.
	ensureCascadeConstraints(db, &models.WorkspaceUser{}, []string{"Workspace", "User"})
	// For other models, add similar calls. Example:
	// ensureCascadeConstraints(db, &models.YourNewModel{}, []string{"AssociationName1", "AssociationName2"})
}

// getEnv returns the environment variable or fallback if not set.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
