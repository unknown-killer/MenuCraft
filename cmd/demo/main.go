// Command demo exports real public-menu pages into a static GitHub Pages demo.
// It seeds a new temporary database and never opens the operator's database.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"ResturantMenue/internal/config"
	"ResturantMenue/internal/database"
	"ResturantMenue/internal/repository"
	"ResturantMenue/internal/server"
	"ResturantMenue/internal/theme"
	"ResturantMenue/web"
)

type demoTheme struct {
	Key     string   `json:"key"`
	Fa      string   `json:"fa"`
	En      string   `json:"en"`
	Colors  []string `json:"colors"`
	Mode    string   `json:"mode"`
	Layout  string   `json:"layout"`
	Columns int      `json:"columns"`
	Radius  int      `json:"radius"`
}

func main() {
	output := flag.String("output", "_site", "directory for generated static demo")
	flag.Parse()
	if err := generate(*output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate(output string) error {
	temp, err := os.MkdirTemp("", "menucraft-demo-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	dbPath := filepath.Join(temp, "demo.db")
	uploads := filepath.Join(temp, "uploads")
	seed := exec.Command("go", "run", "./cmd/seed")
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(variable), "MENU_") {
			seed.Env = append(seed.Env, variable)
		}
	}
	seed.Env = append(seed.Env, "MENU_DB="+dbPath, "MENU_UPLOADS="+uploads, "MENU_ENV=production")
	if result, err := seed.CombinedOutput(); err != nil {
		return fmt.Errorf("seed isolated demo: %w: %s", err, result)
	}
	db, err := database.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	store := repository.New(db)
	venue, err := store.Restaurants.BySlug("demo")
	if err != nil {
		return err
	}
	// Do not publish working links to fictitious contact details from the seed.
	venue.Phone = ""
	venue.Instagram = ""
	venue.Telegram = ""
	venue.WhatsApp = ""
	venue.Website = ""
	if err := store.Restaurants.Update(venue); err != nil {
		return err
	}
	cfg := &config.Config{DatabasePath: dbPath, UploadDir: uploads, SecretKey: []byte("isolated-static-demo-only-no-live-sessions"), SessionDuration: time.Hour, DefaultLocale: "fa", Env: "production", MaxUploadBytes: 8 << 20}
	handler := server.New(cfg, store).Handler()
	var catalogue []demoTheme
	for _, preset := range theme.Presets() {
		value := theme.ApplyPreset(preset.Key)
		if err := store.Themes.Save(venue.ID, preset.Key, value); err != nil {
			return err
		}
		palette := value.Colors
		if value.Mode == "dark" {
			palette = value.DarkColors
		}
		catalogue = append(catalogue, demoTheme{preset.Key, preset.LabelFa, preset.LabelEn, []string{palette.Primary, palette.Background, palette.Surface, palette.Text}, value.Mode, value.Layout.Mode, value.Layout.Columns, value.Cards.Radius})
		for _, language := range []string{"fa", "en"} {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest("GET", "http://demo.local/m/demo?lang="+language, nil))
			if response.Code != 200 {
				return fmt.Errorf("render %s/%s: HTTP %d", preset.Key, language, response.Code)
			}
			page := strings.ReplaceAll(response.Body.String(), "/static/", "../../static/")
			page = strings.ReplaceAll(page, "?lang=fa", "fa.html")
			page = strings.ReplaceAll(page, "?lang=en", "en.html")
			// Explain the static snapshot without pretending there is a public admin server.
			label := "🧪 Sample menu · Open the interactive studio"
			if language == "fa" {
				label = "🧪 منوی نمونه · ورود به استودیوی آزمایشی"
			}
			page = strings.Replace(page, "<body ", "<body data-static-demo=\"true\" ", 1)
			page = strings.Replace(page, "</body>", "<p class=\"demo-home\"><a href=\"../../?lang="+language+"\" target=\"_top\">"+label+"</a></p><style>.demo-home{text-align:center;padding:20px;font-size:14px;border-top:1px solid var(--c-border)}@media print{.demo-home{display:none}}</style></body>", 1)
			if err := writeFile(filepath.Join(output, "menus", preset.Key, language+".html"), []byte(page)); err != nil {
				return err
			}
		}
	}
	if err := copyTree(web.Static(), filepath.Join(output, "static")); err != nil {
		return err
	}
	if err := copyTree(os.DirFS("site"), output); err != nil {
		return err
	}
	data, err := json.Marshal(catalogue)
	if err != nil {
		return err
	}
	if err := writeFile(filepath.Join(output, "themes.json"), data); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(output, ".nojekyll"), nil); err != nil {
		return err
	}
	fmt.Printf("Static demo: %d themes × 2 languages in %s\n", len(catalogue), output)
	return nil
}

func copyTree(source fs.FS, destination string) error {
	return fs.WalkDir(source, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(source, path)
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(destination, filepath.FromSlash(path)), data)
	})
}
func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
