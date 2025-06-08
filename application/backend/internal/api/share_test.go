package api

import (
	"testing"
	"time"
)

func TestValidateSharingExpiryDate(t *testing.T) {
	// Zero time should be invalid
	if validateSharingExpiryDate(time.Time{}) {
		t.Error("Expected zero time to be invalid")
	}

	// Future date should be valid
	futureDate := time.Now().UTC().Add(24 * time.Hour)
	if !validateSharingExpiryDate(futureDate) {
		t.Error("Expected future date to be valid")
	}

	// Past date should be invalid
	pastDate := time.Now().UTC().Add(-24 * time.Hour)
	if validateSharingExpiryDate(pastDate) {
		t.Error("Expected past date to be invalid")
	}
}

func TestSharingExpiryDatePastDue(t *testing.T) {
	now := time.Now().UTC()

	// Same day should not be past due
	if sharingExpiryDatePastDue(now) {
		t.Error("Expected same day to not be past due")
	}

	// Yesterday should be past due
	yesterday := now.AddDate(0, 0, -1)
	if !sharingExpiryDatePastDue(yesterday) {
		t.Error("Expected yesterday to be past due")
	}

	// Tomorrow should not be past due
	tomorrow := now.AddDate(0, 0, 1)
	if sharingExpiryDatePastDue(tomorrow) {
		t.Error("Expected tomorrow to not be past due")
	}
}
