package main

// Employee interface
type Employee interface {
	sendFeedback()
}

// Operator interface
type Operator interface {
	AssignEmployee()
	AddReview()
	UpdateReview()
	AddEmployee(Employee)
	UpdateEmployee(Employee)
	RemoveEmployee(Employee)
}
