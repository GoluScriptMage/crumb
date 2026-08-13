package tests

import (
	"crumb/cmd"
	"crumb/helpers"
	"crumb/store"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// setupTestDB redirects store I/O to a temp file and returns a cleanup func.
func setupTestDB(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data.json")
	store.SetDbPathOverride(dbPath)
	return func() {
		store.SetDbPathOverride("")
		os.Remove(dbPath)
	}
}

func TestTask_AddAndList(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"write tests", "fix bug"}); err != nil {
		t.Fatalf("add tasks failed: %v", err)
	}

	data, err := store.ReadData()
	if err != nil {
		t.Fatalf("read data failed: %v", err)
	}
	if len(data.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(data.Tasks))
	}
	if data.Tasks[0].Text != "write tests" || data.Tasks[1].Text != "fix bug" {
		t.Fatalf("unexpected task texts: %+v", data.Tasks)
	}
	for _, task := range data.Tasks {
		if task.Status != "pending" {
			t.Fatalf("expected pending status, got %s", task.Status)
		}
		if len(task.ID) != 3 {
			t.Fatalf("expected 3-char id, got %q", task.ID)
		}
	}

	if err := cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{}); err != nil {
		t.Fatalf("list tasks failed: %v", err)
	}
}

func TestTask_ClearRequiresDouble(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"one"})
	if err := cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"clear"}); err != nil {
		t.Fatalf("single clear returned error: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Tasks) != 1 {
		t.Fatalf("single 'clear' should not clear tasks, got %d", len(data.Tasks))
	}

	if err := cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"clear", "clear"}); err != nil {
		t.Fatalf("double clear failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Tasks) != 0 {
		t.Fatalf("double 'clear' should clear tasks, got %d", len(data.Tasks))
	}
}

func TestDone_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"task one"})
	data, _ := store.ReadData()
	id := data.Tasks[0].ID

	if err := cmd.DoneCmd().RunE(cmd.DoneCmd(), []string{id}); err != nil {
		t.Fatalf("done failed: %v", err)
	}
	data, _ = store.ReadData()
	if data.Tasks[0].Status != "done" {
		t.Fatalf("expected done status, got %s", data.Tasks[0].Status)
	}

	if err := cmd.DoneCmd().RunE(cmd.DoneCmd(), []string{"nope"}); err != nil {
		t.Fatalf("done missing id returned error: %v", err)
	}
}

func TestCancel_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"task one"})
	data, _ := store.ReadData()
	id := data.Tasks[0].ID

	if err := cmd.CancelCmd().RunE(cmd.CancelCmd(), []string{id}); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	data, _ = store.ReadData()
	if data.Tasks[0].Status != "canceled" {
		t.Fatalf("expected canceled status, got %s", data.Tasks[0].Status)
	}

	if err := cmd.CancelCmd().RunE(cmd.CancelCmd(), []string{}); err != nil {
		t.Fatalf("cancel no args returned error: %v", err)
	}
}

func TestDelete_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"keep", "remove"})
	data, _ := store.ReadData()
	removeID := data.Tasks[1].ID

	if err := cmd.DeleteCmd().RunE(cmd.DeleteCmd(), []string{removeID}); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Tasks) != 1 {
		t.Fatalf("expected 1 task after delete, got %d", len(data.Tasks))
	}
	if data.Tasks[0].Text != "keep" {
		t.Fatalf("wrong task removed: %+v", data.Tasks)
	}

	if err := cmd.DeleteCmd().RunE(cmd.DeleteCmd(), []string{}); err != nil {
		t.Fatalf("delete no args returned error: %v", err)
	}
}

func TestClear_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"a", "b", "c"})
	if err := cmd.ClearCmd().RunE(cmd.ClearCmd(), []string{"clear"}); err != nil {
		t.Fatalf("clear failed: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Tasks) != 0 {
		t.Fatalf("expected 0 tasks after clear, got %d", len(data.Tasks))
	}

	if err := cmd.ClearCmd().RunE(cmd.ClearCmd(), []string{}); err != nil {
		t.Fatalf("clear no args returned error: %v", err)
	}
}

func TestIdea_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.IdeaCmd().RunE(cmd.IdeaCmd(), []string{"idea one", "idea two"}); err != nil {
		t.Fatalf("add ideas failed: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Ideas) != 2 {
		t.Fatalf("expected 2 ideas, got %d", len(data.Ideas))
	}

	if err := cmd.IdeaCmd().RunE(cmd.IdeaCmd(), []string{}); err != nil {
		t.Fatalf("list ideas failed: %v", err)
	}
}

func TestNote_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{"a note"}); err != nil {
		t.Fatalf("add note failed: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Notes) != 1 || data.Notes[0] != "a note" {
		t.Fatalf("unexpected notes: %+v", data.Notes)
	}

	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{}); err != nil {
		t.Fatalf("list notes failed: %v", err)
	}

	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{"too", "many"}); err != nil {
		t.Fatalf("note multiple args returned error: %v", err)
	}
	fmt.Println("All done`")
}

func TestNext_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.NextCmd().RunE(cmd.NextCmd(), []string{"ship it"}); err != nil {
		t.Fatalf("set next failed: %v", err)
	}
	data, _ := store.ReadData()
	if data.Next != "ship it" {
		t.Fatalf("expected 'ship it', got %q", data.Next)
	}

	if err := cmd.NextCmd().RunE(cmd.NextCmd(), []string{"clear"}); err != nil {
		t.Fatalf("clear next failed: %v", err)
	}
	data, _ = store.ReadData()
	if data.Next != "" {
		t.Fatalf("expected empty next, got %q", data.Next)
	}

	if err := cmd.NextCmd().RunE(cmd.NextCmd(), []string{}); err != nil {
		t.Fatalf("show next failed: %v", err)
	}

	if err := cmd.NextCmd().RunE(cmd.NextCmd(), []string{"a", "b"}); err != nil {
		t.Fatalf("next multiple args returned error: %v", err)
	}
}

func TestGenerateShortId(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		id := cmd.GenerateShortId()
		if len(id) != 3 {
			t.Fatalf("expected 3-char id, got %q", id)
		}
		for _, c := range id {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Fatalf("non-hex char in id %q", id)
			}
		}
		seen[id] = true
	}
	if len(seen) < 40 {
		t.Fatalf("ids not random enough, only %d distinct of 50", len(seen))
	}
}

func TestFormatStatus(t *testing.T) {
	for _, s := range []string{"done", "canceled", "failed", "pending", "weird"} {
		if helpers.FormatStatus(s) == "" {
			t.Fatalf("FormatStatus(%q) returned empty", s)
		}
	}
}

func TestStore_ReadDataMissingFile(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	data, err := store.ReadData()
	if err != nil {
		t.Fatalf("ReadData on missing file should not error: %v", err)
	}
	if data.Notes == nil {
		t.Fatalf("expected non-nil Notes slice on empty data")
	}
	if len(data.Tasks) != 0 || len(data.Ideas) != 0 {
		t.Fatalf("expected empty collections, got %+v", data)
	}
}

func TestStore_WriteReadRoundTrip(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	in := store.CrumbData{
		Next:  "focus",
		Tasks: []store.Task{{ID: "abc", Text: "t", Status: "pending"}},
		Ideas: []string{"idea"},
		Notes: []string{"note"},
	}
	if err := store.WriteData(in); err != nil {
		t.Fatalf("WriteData failed: %v", err)
	}

	out, err := store.ReadData()
	if err != nil {
		t.Fatalf("ReadData failed: %v", err)
	}
	if out.Next != "focus" || len(out.Tasks) != 1 || out.Tasks[0].ID != "abc" ||
		len(out.Ideas) != 1 || len(out.Notes) != 1 {
		t.Fatalf("round-trip mismatch: %+v", out)
	}
}

func TestStore_UpdateAppliesAndPersists(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := store.Update(func(d *store.CrumbData) error {
		d.Notes = append(d.Notes, "via update")
		return nil
	}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	data, _ := store.ReadData()
	if len(data.Notes) != 1 || data.Notes[0] != "via update" {
		t.Fatalf("Update did not persist: %+v", data.Notes)
	}
}

func TestStore_GetDbPathOverride(t *testing.T) {
	dir := t.TempDir()
	custom := dir + "/custom.json"
	store.SetDbPathOverride(custom)
	defer store.SetDbPathOverride("")

	got, err := store.GetDbPath()
	if err != nil {
		t.Fatalf("GetDbPath failed: %v", err)
	}
	if got != custom {
		t.Fatalf("expected overridden path %q, got %q", custom, got)
	}
}
