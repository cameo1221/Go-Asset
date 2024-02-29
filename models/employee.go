package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee in the system
type Employee struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	ArchivedAt *time.Time `json:"archive_at,omitempty"`
}

// EmployeeModel represents the model for employee operations
type EmployeeModel struct {
	DB *sql.DB
}

// CreateEmployee creates a new employee in the database
func (em *EmployeeModel) CreateEmployee(employee *Employee) error {
	query := `
		INSERT INTO employee (id, name, email, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	employee.ID = uuid.New()
	err := em.DB.QueryRow(query, employee.ID,employee.Name, employee.Email, employee.Role, employee.CreatedAt).Scan(&employee.ID)
	if err != nil {
		return err
	}

	
	return nil
}

// UpdateEmployee updates an existing employee in the database
func (em *EmployeeModel) UpdateEmployee(employee *Employee) error {
	query := `
		UPDATE employee
		SET name = $1, email = $2, role = $3, archive_at = $4
		WHERE id = $5
	`

	_, err := em.DB.Exec(query, employee.Name, employee.Email, employee.Role, employee.ArchivedAt, employee.ID)
	return err
}

// ArchiveEmployee archives an existing employee in the database
func (em *EmployeeModel) ArchiveEmployee(id uuid.UUID) error {
	query := `
		UPDATE employee
		SET archive_at = $1
		WHERE id = $2
	`

	_, err := em.DB.Exec(query, time.Now(), id)
	return err
}

// GetEmployeeByID retrieves an employee from the database by its ID
func (em *EmployeeModel) GetEmployeeByID(id uuid.UUID) (*Employee, error) {
	query := `
		SELECT id, name, email, role, created_at, archive_at
		FROM employee
		WHERE id = $1
	`

	employee := &Employee{}
	err := em.DB.QueryRow(query, id).Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Role, &employee.CreatedAt, &employee.ArchivedAt)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

// GetAllEmployees retrieves all employees from the database
func (em *EmployeeModel) GetAllEmployees() ([]*Employee, error) {
	query := `
		SELECT id, name, email, role, created_at, archive_at
		FROM employee
	`

	rows, err := em.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		employee := &Employee{}
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Role, &employee.CreatedAt, &employee.ArchivedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
