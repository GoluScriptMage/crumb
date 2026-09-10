package tests

import (
	"crumb/cmd"
	"crumb/store"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		if len(task.ID) != 4 {
			t.Fatalf("expected 4-char id, got %q", task.ID)
		}
	}

	if err := cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{}); err != nil {
		t.Fatalf("list tasks failed: %v", err)
	}
}

func TestFail_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"task fail me"})
	data, _ := store.ReadData()
	id := data.Tasks[0].ID

	if err := cmd.FailCmd().RunE(cmd.FailCmd(), []string{id}); err != nil {
		t.Fatalf("fail failed: %v", err)
	}
	data, _ = store.ReadData()
	if data.Tasks[0].Status != "failed" {
		t.Fatalf("expected failed status, got %s", data.Tasks[0].Status)
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

func TestDel_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cmd.TaskCmd().RunE(cmd.TaskCmd(), []string{"keep", "remove"})
	data, _ := store.ReadData()
	removeID := data.Tasks[1].ID

	if err := cmd.DelCmd().RunE(cmd.DelCmd(), []string{removeID}); err != nil {
		t.Fatalf("del failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Tasks) != 1 {
		t.Fatalf("expected 1 task after del, got %d", len(data.Tasks))
	}
	if data.Tasks[0].Text != "keep" {
		t.Fatalf("wrong task removed: %+v", data.Tasks)
	}

	// Del all
	if err := cmd.DelCmd().RunE(cmd.DelCmd(), []string{"all"}); err != nil {
		t.Fatalf("del all failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Tasks) != 0 {
		t.Fatalf("expected 0 tasks after del all, got %d", len(data.Tasks))
	}
}

func TestIdea_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.IdeaCmd().RunE(cmd.IdeaCmd(), []string{"idea one", "idea two"}); err != nil {
		t.Fatalf("add ideas failed: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Ideas) != 1 {
		t.Fatalf("expected 1 idea joined, got %d", len(data.Ideas))
	}

	if err := cmd.IdeaCmd().RunE(cmd.IdeaCmd(), []string{}); err != nil {
		t.Fatalf("list ideas failed: %v", err)
	}

	// Delete idea
	if err := cmd.IdeaCmd().RunE(cmd.IdeaCmd(), []string{"del", "1"}); err != nil {
		t.Fatalf("del idea failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Ideas) != 0 {
		t.Fatalf("expected 0 ideas after delete, got %d", len(data.Ideas))
	}
}

func TestNote_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{"a note with spaces"}); err != nil {
		t.Fatalf("add note failed: %v", err)
	}
	data, _ := store.ReadData()
	if len(data.Notes) != 1 || data.Notes[0] != "a note with spaces" {
		t.Fatalf("unexpected notes: %+v", data.Notes)
	}

	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{}); err != nil {
		t.Fatalf("list notes failed: %v", err)
	}

	// Delete note
	if err := cmd.NoteCmd().RunE(cmd.NoteCmd(), []string{"del", "1"}); err != nil {
		t.Fatalf("del note failed: %v", err)
	}
	data, _ = store.ReadData()
	if len(data.Notes) != 0 {
		t.Fatalf("expected 0 notes after delete, got %d", len(data.Notes))
	}
}

func TestTimer_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	if err := cmd.TimerCmd().RunE(cmd.TimerCmd(), []string{"50", "Doing DSA"}); err != nil {
		t.Fatalf("start timer failed: %v", err)
	}
	data, _ := store.ReadData()
	if data.Timer == nil || data.Timer.Minutes != 50 || data.Timer.Task != "Doing DSA" {
		t.Fatalf("unexpected timer state: %+v", data.Timer)
	}

	// View timer
	if err := cmd.TimerCmd().RunE(cmd.TimerCmd(), []string{}); err != nil {
		t.Fatalf("view timer failed: %v", err)
	}

	// Stop timer
	if err := cmd.TimerCmd().RunE(cmd.TimerCmd(), []string{"stop"}); err != nil {
		t.Fatalf("stop timer failed: %v", err)
	}
	data, _ = store.ReadData()
	if data.Timer != nil {
		t.Fatalf("expected timer to be nil after stop, got %+v", data.Timer)
	}
}

func TestVersion_Command(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.VersionCmd().RunE(cmd.VersionCmd(), []string{})

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	output := string(out)

	want := "crumb version " + cmd.Version()
	if !strings.Contains(output, want) {
		t.Fatalf("expected version output to contain %q, got %q", want, output)
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
		Tasks: []store.Task{{ID: "abc", Text: "t", Status: "pending"}},
		Ideas: []string{"idea"},
		Notes: []string{"note"},
		Timer: &store.TimerState{Task: "dsa", Minutes: 25, StartedAt: 100, Duration: 1500},
	}
	if err := store.WriteData(in); err != nil {
		t.Fatalf("WriteData failed: %v", err)
	}

	out, err := store.ReadData()
	if err != nil {
		t.Fatalf("ReadData failed: %v", err)
	}
	if len(out.Tasks) != 1 || out.Tasks[0].ID != "abc" ||
		len(out.Ideas) != 1 || len(out.Notes) != 1 ||
		out.Timer == nil || out.Timer.Task != "dsa" {
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
