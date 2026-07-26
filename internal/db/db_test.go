package db

import (
	"os"
	"testing"
)

func TestGuildSettings(t *testing.T) {
	path := "test_settings.db"
	defer os.Remove(path)

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	gs, err := db.GetGuildSettings("123")
	if err != nil {
		t.Fatalf("get guild settings: %v", err)
	}
	if gs.GuildID != "123" || gs.Language != "en" || gs.TideStation != "QUB" {
		t.Errorf("unexpected defaults: %+v", gs)
	}

	if err := db.SetLanguage("123", "tc"); err != nil {
		t.Fatalf("set language: %v", err)
	}
	if err := db.SetTideStation("123", "SPW"); err != nil {
		t.Fatalf("set tide station: %v", err)
	}

	gs, err = db.GetGuildSettings("123")
	if err != nil {
		t.Fatalf("get guild settings after update: %v", err)
	}
	if gs.Language != "tc" || gs.TideStation != "SPW" {
		t.Errorf("unexpected updated values: %+v", gs)
	}
}

func TestWarningState(t *testing.T) {
	path := "test_warnings.db"
	defer os.Remove(path)

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	state, err := db.GetWarningState("WRAIN")
	if err != nil {
		t.Fatalf("get warning state: %v", err)
	}
	if state != nil {
		t.Errorf("expected nil state, got %+v", state)
	}

	if err := db.SaveWarningState("WRAIN", "WRAINA", "ISSUE", "t1", "t2"); err != nil {
		t.Fatalf("save warning state: %v", err)
	}

	state, err = db.GetWarningState("WRAIN")
	if err != nil {
		t.Fatalf("get warning state after save: %v", err)
	}
	if state == nil || state.Code != "WRAIN" || state.ActionCode != "ISSUE" {
		t.Errorf("unexpected state: %+v", state)
	}

	if err := db.SaveWarningState("WRAIN", "", "UPDATE", "t1", "t3"); err != nil {
		t.Fatalf("save updated warning state: %v", err)
	}

	latest, err := db.LatestWarningStates()
	if err != nil {
		t.Fatalf("latest warning states: %v", err)
	}
	if len(latest) != 1 {
		t.Errorf("expected 1 latest state, got %d", len(latest))
	}
	if latest[0].ActionCode != "UPDATE" {
		t.Errorf("expected latest action UPDATE, got %s", latest[0].ActionCode)
	}
}

func TestTipsUpdateTime(t *testing.T) {
	path := "test_tips.db"
	defer os.Remove(path)

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if _, err := db.GetLatestTipsUpdateTime(); err != nil {
		t.Fatalf("get latest tips update time: %v", err)
	}

	if err := db.SaveTipsUpdateTime("2026-07-26T12:00:00+08:00"); err != nil {
		t.Fatalf("save tips update time: %v", err)
	}

	latest, err := db.GetLatestTipsUpdateTime()
	if err != nil {
		t.Fatalf("get latest tips update time after save: %v", err)
	}
	if latest != "2026-07-26T12:00:00+08:00" {
		t.Errorf("unexpected latest time: %q", latest)
	}
}
