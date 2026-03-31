package db

import (
	"database/sql"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	for _, ddl := range []string{
		`CREATE TABLE ZAPP (
			Z_PK INTEGER PRIMARY KEY,
			ZIDENTIFIER INTEGER,
			ZNAME TEXT,
			ZBUNDLEIDENTIFIER TEXT,
			ZVENDORNAME TEXT,
			ZTAGLINE TEXT,
			ZMARKETINGDESCRIPTION TEXT,
			ZMARKETINGURLSTRING TEXT,
			ZSHARINGURLSTRING TEXT,
			ZSIZE INTEGER,
			ZMARKETINGVERSIONSTRING TEXT,
			ZMINOSVERSION TEXT,
			ZURLSCHEME TEXT,
			ZJOINEDKEYWORDS TEXT,
			ZFIRSTRELEASEDATE REAL,
			ZLASTRELEASEDATE REAL,
			ZSEARCHIDENTIFIER TEXT
		)`,
		`CREATE TABLE ZSETAPPCATEGORY (
			Z_PK INTEGER PRIMARY KEY,
			ZIDENTIFIER INTEGER,
			ZNAME TEXT,
			ZCATEGORYDESCRIPTION TEXT,
			ZPOSITION INTEGER
		)`,
		`CREATE TABLE Z_1SETAPPCATEGORIES (
			Z_1APPLICATIONS INTEGER,
			Z_20SETAPPCATEGORIES INTEGER
		)`,
	} {
		if _, err := conn.Exec(ddl); err != nil {
			t.Fatal(err)
		}
	}

	return &DB{db: conn}
}

func seedTestData(t *testing.T, d *DB) {
	t.Helper()
	for _, q := range []string{
		`INSERT INTO ZAPP (Z_PK, ZIDENTIFIER, ZNAME, ZVENDORNAME, ZTAGLINE, ZJOINEDKEYWORDS, ZSEARCHIDENTIFIER)
		 VALUES (1, 100, 'CleanMyMac', 'MacPaw', 'Clean your Mac', 'cleaner,utility', 'cleanmymac')`,
		`INSERT INTO ZAPP (Z_PK, ZIDENTIFIER, ZNAME, ZVENDORNAME, ZTAGLINE, ZJOINEDKEYWORDS, ZSEARCHIDENTIFIER)
		 VALUES (2, 200, 'Bartender', 'Surtees Studios', 'Organize menu bar', 'menu,bar', 'bartender')`,
		`INSERT INTO ZAPP (Z_PK, ZIDENTIFIER, ZNAME, ZVENDORNAME, ZTAGLINE, ZJOINEDKEYWORDS, ZSEARCHIDENTIFIER)
		 VALUES (3, 300, 'CleanMyMac X', 'MacPaw', 'Advanced cleaning', 'cleaner', 'cleanmymacx')`,
		`INSERT INTO ZSETAPPCATEGORY (Z_PK, ZIDENTIFIER, ZNAME, ZCATEGORYDESCRIPTION, ZPOSITION)
		 VALUES (1, 10, 'Utilities', 'Utility apps', 2)`,
		`INSERT INTO ZSETAPPCATEGORY (Z_PK, ZIDENTIFIER, ZNAME, ZCATEGORYDESCRIPTION, ZPOSITION)
		 VALUES (2, 20, 'Productivity', 'Productivity apps', 1)`,
		`INSERT INTO Z_1SETAPPCATEGORIES (Z_1APPLICATIONS, Z_20SETAPPCATEGORIES) VALUES (1, 1)`,
		`INSERT INTO Z_1SETAPPCATEGORIES (Z_1APPLICATIONS, Z_20SETAPPCATEGORIES) VALUES (2, 2)`,
		`INSERT INTO Z_1SETAPPCATEGORIES (Z_1APPLICATIONS, Z_20SETAPPCATEGORIES) VALUES (1, 2)`,
	} {
		if _, err := d.db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllApps(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	apps, err := d.AllApps()
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 3 {
		t.Fatalf("got %d apps, want 3", len(apps))
	}
	// Should be sorted case-insensitively by name
	if apps[0].Name != "Bartender" {
		t.Errorf("first app = %q, want %q", apps[0].Name, "Bartender")
	}
	if apps[1].Name != "CleanMyMac" {
		t.Errorf("second app = %q, want %q", apps[1].Name, "CleanMyMac")
	}
}

func TestSearchApps(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	t.Run("by name", func(t *testing.T) {
		apps, err := d.SearchApps("bartender")
		if err != nil {
			t.Fatal(err)
		}
		if len(apps) != 1 || apps[0].Name != "Bartender" {
			t.Errorf("got %v, want [Bartender]", apps)
		}
	})

	t.Run("by keyword", func(t *testing.T) {
		apps, err := d.SearchApps("cleaner")
		if err != nil {
			t.Fatal(err)
		}
		if len(apps) != 2 {
			t.Errorf("got %d apps, want 2", len(apps))
		}
	})

	t.Run("by vendor", func(t *testing.T) {
		apps, err := d.SearchApps("MacPaw")
		if err != nil {
			t.Fatal(err)
		}
		if len(apps) != 2 {
			t.Errorf("got %d apps, want 2", len(apps))
		}
	})

	t.Run("no match", func(t *testing.T) {
		apps, err := d.SearchApps("nonexistent")
		if err != nil {
			t.Fatal(err)
		}
		if len(apps) != 0 {
			t.Errorf("got %d apps, want 0", len(apps))
		}
	})
}

func TestFindApp(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	t.Run("exact name", func(t *testing.T) {
		app, err := d.FindApp("Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if app.Name != "Bartender" {
			t.Errorf("got %q, want %q", app.Name, "Bartender")
		}
	})

	t.Run("search identifier", func(t *testing.T) {
		app, err := d.FindApp("bartender")
		if err != nil {
			t.Fatal(err)
		}
		if app.Name != "Bartender" {
			t.Errorf("got %q, want %q", app.Name, "Bartender")
		}
	})

	t.Run("substring single match", func(t *testing.T) {
		app, err := d.FindApp("artend")
		if err != nil {
			t.Fatal(err)
		}
		if app.Name != "Bartender" {
			t.Errorf("got %q, want %q", app.Name, "Bartender")
		}
	})

	t.Run("ambiguous", func(t *testing.T) {
		// "Clean" matches both "CleanMyMac" and "CleanMyMac X" via substring
		_, err := d.FindApp("Clean")
		if err == nil {
			t.Fatal("expected error for ambiguous match")
		}
	})

	t.Run("no match", func(t *testing.T) {
		_, err := d.FindApp("nonexistent")
		if err == nil {
			t.Fatal("expected error for no match")
		}
	})
}

func TestAppCategories(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	cats, err := d.AppCategories(1) // CleanMyMac -> Utilities, Productivity
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 2 {
		t.Fatalf("got %d categories, want 2", len(cats))
	}
}

func TestLoadCategories(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	apps := []model.App{
		{PK: 1, Name: "CleanMyMac"},
		{PK: 2, Name: "Bartender"},
		{PK: 3, Name: "CleanMyMac X"},
	}
	if err := d.LoadCategories(apps); err != nil {
		t.Fatal(err)
	}
	if len(apps[0].Categories) != 2 {
		t.Errorf("CleanMyMac categories = %d, want 2", len(apps[0].Categories))
	}
	if len(apps[1].Categories) != 1 {
		t.Errorf("Bartender categories = %d, want 1", len(apps[1].Categories))
	}
	if len(apps[2].Categories) != 0 {
		t.Errorf("CleanMyMac X categories = %d, want 0", len(apps[2].Categories))
	}
}

func TestAllCategories(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	cats, err := d.AllCategories()
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 2 {
		t.Fatalf("got %d categories, want 2", len(cats))
	}
	// Should be sorted by position: Productivity (1), Utilities (2)
	if cats[0].Name != "Productivity" {
		t.Errorf("first category = %q, want %q", cats[0].Name, "Productivity")
	}
	if cats[1].Name != "Utilities" {
		t.Errorf("second category = %q, want %q", cats[1].Name, "Utilities")
	}
}

func TestAppsByCategory(t *testing.T) {
	d := newTestDB(t)
	seedTestData(t, d)

	apps, err := d.AppsByCategory("Utilities")
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 1 || apps[0].Name != "CleanMyMac" {
		t.Errorf("got %v, want [CleanMyMac]", apps)
	}
}
