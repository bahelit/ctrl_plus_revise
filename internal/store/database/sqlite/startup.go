package sqlite

import (
	"database/sql"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

func GetDatabase() (*DB, error) {
	sqlFile := filepath.Join(createFolder(), "ctrl_plus_revise.db")

	SQLiteConn, err := sql.Open("sqlite3", sqlFile)
	if err != nil {
		slog.Error("failed to open database connection", "path", sqlFile, "error", err)
		return nil, err
	}
	return &DB{Conn: SQLiteConn}, nil
}

func createFolder() string {
	u, err := user.Current()
	if err != nil {
		slog.Error("Failed to get user", "err", err)
		return ""
	}

	folderName := "CtrlPlusRevise"
	folderPath := filepath.Join(u.HomeDir, folderName, "DB")

	if _, err = os.Stat(folderPath); !os.IsNotExist(err) {
		slog.Debug("Folder already exists", "folder", folderName)
	} else {
		err = os.Mkdir(folderPath, 0755)
		if err != nil && !os.IsExist(err) {
			slog.Error("Failed to create folder", "err", err)
			return ""
		}
		slog.Info("Folder created successfully", "folder", folderPath)
	}

	return folderPath
}
