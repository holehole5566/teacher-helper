package services

import (
	"teacher-wails/internal/models"
	"testing"
)

func TestMidtermPersistenceAndValidation(t *testing.T) {
	dir := t.TempDir()
	dh, err := NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	expected, actual := 30, 28
	record := models.MidtermExam{Date: "2026-11-10", Subjects: [8]string{" 國語 ", "自習"}, Expected: &expected, Actual: &actual}
	if err := dh.SaveMidtermExam(record); err != nil {
		t.Fatal(err)
	}
	if err := dh.SaveMidtermExam(models.MidtermExam{Date: "2026-11-11", Subjects: [8]string{"數學"}}); err != nil {
		t.Fatal(err)
	}
	reopened, err := NewDataHandlerWithPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.GetMidtermExam(record.Date)
	if err != nil {
		t.Fatal(err)
	}
	if got.Subjects[0] != "國語" || got.Subjects[1] != "自習" || got.Expected == nil || *got.Expected != 30 || got.Actual == nil || *got.Actual != 28 {
		t.Fatalf("unexpected record: %+v", got)
	}
	empty, err := reopened.GetMidtermExam("2026-11-12")
	if err != nil || empty.Expected != nil || empty.Actual != nil || empty.Date != "2026-11-12" {
		t.Fatalf("unexpected empty record: %+v, %v", empty, err)
	}
	negative, excess := -1, 31
	for _, invalid := range []models.MidtermExam{
		{Date: "2026-02-30"},
		{Date: record.Date, Expected: &negative},
		{Date: record.Date, Expected: &expected, Actual: &excess},
	} {
		if err := reopened.SaveMidtermExam(invalid); err == nil {
			t.Fatalf("accepted invalid record: %+v", invalid)
		}
	}
	got, _ = reopened.GetMidtermExam(record.Date)
	if got.Actual == nil || *got.Actual != 28 {
		t.Fatal("invalid save changed persisted attendance")
	}
	zero := 0
	if err := reopened.SaveMidtermExam(models.MidtermExam{Date: record.Date, Expected: &zero, Actual: &zero}); err != nil {
		t.Fatal(err)
	}
	got, _ = reopened.GetMidtermExam(record.Date)
	if got.Actual == nil || *got.Actual != 0 {
		t.Fatal("zero attendance was lost")
	}
}
