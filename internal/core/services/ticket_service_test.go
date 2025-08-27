package services

import (
	"testing"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestTicketService_CalculateFinalFields(t *testing.T) {
	service := NewTicketService(nil, nil)
	
	parent := domain.Ticket{
		CustomFields: map[string]string{
			"Priority": "High",
			"Sprint": "10",
		},
	}
	
	task := domain.Task{
		CustomFields: map[string]string{
			"Priority": "Low",
		},
	}
	
	result := service.calculateFinalFields(parent, task)
	
	// Assert: Priority should be overridden to "Low", Sprint should be inherited as "10"
	if result["Priority"] != "Low" {
		t.Errorf("Expected Priority to be 'Low', got '%s'", result["Priority"])
	}
	
	if result["Sprint"] != "10" {
		t.Errorf("Expected Sprint to be '10', got '%s'", result["Sprint"])
	}
}