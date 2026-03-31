package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thirteen37/setapp/internal/model"
	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func Open() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}
	dbPath := filepath.Join(home, "Library", "Application Support", "Setapp", "Default", "Databases", "Apps.sqlite")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Setapp database not found at %s. Is Setapp installed?", dbPath)
	}
	conn, err := sql.Open("sqlite", dbPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}
	return &DB{db: conn}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

const appColumns = `Z_PK, ZIDENTIFIER, ZNAME, ZBUNDLEIDENTIFIER, ZVENDORNAME, ZTAGLINE,
	ZMARKETINGDESCRIPTION, ZMARKETINGURLSTRING, ZSHARINGURLSTRING, ZSIZE,
	ZMARKETINGVERSIONSTRING, ZMINOSVERSION, ZURLSCHEME, ZJOINEDKEYWORDS,
	ZFIRSTRELEASEDATE, ZLASTRELEASEDATE`

// qualifiedAppColumns prefixes each column in appColumns with the given table alias.
func qualifiedAppColumns(alias string) string {
	cols := strings.Split(appColumns, ",")
	for i, c := range cols {
		cols[i] = alias + "." + strings.TrimSpace(c)
	}
	return strings.Join(cols, ", ")
}

func scanApp(row interface{ Scan(...any) error }) (model.App, error) {
	var a model.App
	var bundleID, vendor, tagline, desc, marketingURL, sharingURL sql.NullString
	var version, minOS, urlScheme, keywords sql.NullString
	var size sql.NullInt64
	var firstRelease, lastRelease sql.NullFloat64
	var identifier sql.NullInt64

	err := row.Scan(
		&a.PK, &identifier, &a.Name, &bundleID, &vendor, &tagline,
		&desc, &marketingURL, &sharingURL, &size,
		&version, &minOS, &urlScheme, &keywords,
		&firstRelease, &lastRelease,
	)
	if err != nil {
		return a, err
	}

	if identifier.Valid {
		a.Identifier = int(identifier.Int64)
	}
	a.BundleIdentifier = bundleID.String
	a.Vendor = vendor.String
	a.Tagline = tagline.String
	a.Description = desc.String
	a.MarketingURL = marketingURL.String
	a.SharingURL = sharingURL.String
	if size.Valid {
		a.Size = size.Int64
	}
	a.Version = version.String
	a.MinOS = minOS.String
	a.URLScheme = urlScheme.String
	a.Keywords = keywords.String
	if firstRelease.Valid {
		a.FirstRelease = &firstRelease.Float64
	}
	if lastRelease.Valid {
		a.LastRelease = &lastRelease.Float64
	}

	return a, nil
}

func (d *DB) AllApps() ([]model.App, error) {
	rows, err := d.db.Query("SELECT " + appColumns + " FROM ZAPP ORDER BY ZNAME COLLATE NOCASE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.App
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}

func (d *DB) SearchApps(query string) ([]model.App, error) {
	like := "%" + query + "%"
	rows, err := d.db.Query(
		"SELECT "+appColumns+" FROM ZAPP WHERE ZNAME LIKE ? COLLATE NOCASE OR ZJOINEDKEYWORDS LIKE ? COLLATE NOCASE OR ZTAGLINE LIKE ? COLLATE NOCASE OR ZVENDORNAME LIKE ? COLLATE NOCASE ORDER BY ZNAME COLLATE NOCASE",
		like, like, like, like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.App
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}

func (d *DB) FindApp(name string) (*model.App, error) {
	// Exact match on name or search identifier
	row := d.db.QueryRow(
		"SELECT "+appColumns+" FROM ZAPP WHERE ZNAME = ? COLLATE NOCASE OR ZSEARCHIDENTIFIER = ? COLLATE NOCASE",
		name, name,
	)
	a, err := scanApp(row)
	if err == nil {
		return &a, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Substring match
	like := "%" + name + "%"
	rows, err := d.db.Query(
		"SELECT "+appColumns+" FROM ZAPP WHERE ZNAME LIKE ? COLLATE NOCASE ORDER BY ZNAME COLLATE NOCASE",
		like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.App
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		matches = append(matches, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no app found matching %q", name)
	}
	if len(matches) == 1 {
		return &matches[0], nil
	}

	names := make([]string, len(matches))
	for i, m := range matches {
		names[i] = m.Name
	}
	return nil, fmt.Errorf("ambiguous app name %q, matches: %s", name, strings.Join(names, ", "))
}

func (d *DB) LoadCategories(apps []model.App) error {
	rows, err := d.db.Query(
		"SELECT j.Z_1APPLICATIONS, c.ZNAME FROM Z_1SETAPPCATEGORIES j JOIN ZSETAPPCATEGORY c ON j.Z_20SETAPPCATEGORIES = c.Z_PK",
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	catMap := make(map[int][]string)
	for rows.Next() {
		var pk int
		var name string
		if err := rows.Scan(&pk, &name); err != nil {
			return err
		}
		catMap[pk] = append(catMap[pk], name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i := range apps {
		if cats, ok := catMap[apps[i].PK]; ok {
			apps[i].Categories = cats
		}
	}
	return nil
}

func (d *DB) AppCategories(pk int) ([]string, error) {
	rows, err := d.db.Query(
		"SELECT c.ZNAME FROM ZSETAPPCATEGORY c JOIN Z_1SETAPPCATEGORIES j ON j.Z_20SETAPPCATEGORIES = c.Z_PK WHERE j.Z_1APPLICATIONS = ?",
		pk,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		cats = append(cats, name)
	}
	return cats, rows.Err()
}

func (d *DB) AllCategories() ([]model.Category, error) {
	rows, err := d.db.Query(
		"SELECT Z_PK, ZIDENTIFIER, ZNAME, ZCATEGORYDESCRIPTION, ZPOSITION FROM ZSETAPPCATEGORY ORDER BY ZPOSITION",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.PK, &c.Identifier, &c.Name, &c.Description, &c.Position); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (d *DB) AppsByCategory(categoryName string) ([]model.App, error) {
	rows, err := d.db.Query(
		"SELECT "+qualifiedAppColumns("a")+" FROM ZAPP a JOIN Z_1SETAPPCATEGORIES j ON j.Z_1APPLICATIONS = a.Z_PK JOIN ZSETAPPCATEGORY c ON j.Z_20SETAPPCATEGORIES = c.Z_PK WHERE c.ZNAME = ? COLLATE NOCASE ORDER BY a.ZNAME COLLATE NOCASE",
		categoryName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.App
	for rows.Next() {
		a, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}
