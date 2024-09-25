package database

import (
	"errors"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bahelit/ctrl_plus_revise/internal/store/database/sqlite"
	"github.com/bahelit/ctrl_plus_revise/internal/store/models/chat"
)

type Chat interface {
	GetAllChats(user string) ([]*chat.Chat, error)
	SaveChat(user string, yakity *chat.Chat) (string, error)
}

type ChatBot struct {
	SQL *sqlite.DB
}

func NewSQLiteDB() (*ChatBot, error) {
	db, err := sqlite.GetDatabase()
	if err != nil {
		return nil, err
	}
	cb := &ChatBot{SQL: db}
	err = cb.CreateTable()
	if err != nil {
		return nil, err
	}
	return cb, nil
}

func (db *ChatBot) CreateTable() error {
	sqlStmt := `
	create table if not exists chat (id integer not null primary key, model integer, context blob, owner text, title text, questions text, responses text);
	`
	_, err := db.SQL.Conn.Exec(sqlStmt)
	if err != nil {
		slog.Error("Failed to crate table", "error", err, "SQL", sqlStmt)
		return err
	}
	return nil
}

func (db *ChatBot) GetAllChats(user string) ([]*chat.Chat, error) {
	rows, err := db.SQL.Conn.Query("select id, model, context, title, questions, responses from chat where owner = ?", user)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		chats     []*chat.Chat
		questions string
		responses string
	)
	for rows.Next() {
		var chatEntry chat.Chat
		chatEntry.ID = new(int64)
		chatContext := new([]byte)
		err = rows.Scan(chatEntry.ID, &chatEntry.Model, chatContext, &chatEntry.Title, &questions, &responses)
		if err != nil {
			slog.Error("Failed to scan row", "error", err, "row", rows)
			return nil, err
		}
		chatEntry.ContextFromDB(*chatContext)
		chatEntry.SetQuestions(questions)
		chatEntry.SetResponses(responses)
		chats = append(chats, &chatEntry)
	}
	err = rows.Err()
	if err != nil {
		slog.Error("Failed to scan rows", "error", err, "rows", rows)
		return nil, err
	}
	slog.Debug("Getting all chats", "found", len(chats), "user", user)
	return chats, nil
}

func (db *ChatBot) SaveChat(chatEntry *chat.Chat) error {
	tx, err := db.SQL.Conn.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO chat(model, context, owner, title, questions, responses) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		slog.Error("Failed to prepare statement", "error", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(chatEntry.Model, chatEntry.ContextToDB(), chatEntry.Owner, chatEntry.Title, chatEntry.QuestionsToString(), chatEntry.ResponsesToString())
	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		return nil
	}
	chatID, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert id", "error", err)
		return err
	}
	chatEntry.ID = &chatID
	return nil
}

func (db *ChatBot) UpdateChat(chatEntry *chat.Chat) error {
	tx, err := db.SQL.Conn.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}
	chatCtx := chatEntry.ContextToDB()
	stmt, err := tx.Prepare("UPDATE chat SET context = $1, questions = $2, responses = $3 WHERE id=$4")
	if err != nil {
		slog.Error("Failed to prepare statement", "error", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(chatCtx, chatEntry.QuestionsToString(), chatEntry.ResponsesToString(), chatEntry.ID)
	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		slog.Warn("Chat message not updated", "id", chatEntry.ID)
		return errors.New("no rows updated")
	}
	return nil
}

func (db *ChatBot) DeleteChat(id int64) error {
	tx, err := db.SQL.Conn.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}
	stmt, err := tx.Prepare("delete from chat where id=?")
	if err != nil {
		slog.Error("Failed to prepare statement", "error", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		slog.Warn("Chat message not deleted", "id", id)
		return errors.New("no rows updated")
	}
	return nil
}
