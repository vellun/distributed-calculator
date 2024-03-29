package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func Connect() *pgx.Conn {
	params := GetDBParams()
	db_url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", params.Username, params.Password,
		params.Host, params.Port, params.DBName)
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Printf("Unable to acquire a database connection: %v\n", err)
		return nil
	}
	return conn
}

func InitRepository() { // Создание всех таблиц в бд и заполнение таблиц операций и агентов
	var create_tables_stmt = `
		CREATE TABLE IF NOT EXISTS expressions(
			id INT generated always AS IDENTITY PRIMARY KEY,
			expression VARCHAR NOT NULL,
			status VARCHAR NOT NULL,
			result VARCHAR,
			started_at INT,
			ended_at INT);
		CREATE TABLE IF NOT EXISTS operations(
			id INT generated always AS IDENTITY PRIMARY KEY,
			name VARCHAR NOT NULL,
			duration INT);
		CREATE TABLE IF NOT EXISTS computing_resources(
			id INT generated always AS IDENTITY PRIMARY KEY,
			last_active INT,
			ind INT,
			status VARCHAR);
		CREATE TABLE IF NOT EXISTS tasks(
			id INT generated always AS IDENTITY PRIMARY KEY,
			operand1 VARCHAR,
			operand2 VARCHAR,
			task_id1 INT,
			task_id2 INT,
			operation_id INT,
			status VARCHAR,
			expression_id INT,
			seq_number INT,
			FOREIGN KEY (expression_id) REFERENCES expressions (id),
			FOREIGN KEY (operation_id) REFERENCES operations (id),
			FOREIGN KEY (task_id1) REFERENCES tasks (id),
			FOREIGN KEY (task_id2) REFERENCES tasks (id))`

	conn := Connect()
	defer conn.Close(context.Background())
	_, err := conn.Exec(context.Background(), create_tables_stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec for create tables failed: %v\n", err)
		return
	}
	fmt.Println("All tables succesfully created ;)")

	conn = Connect()
	defer conn.Close(context.Background())
	var n int
	// Проверяем пустая ли таблица с вычислителями
	num, _ := conn.Query(context.Background(), "SELECT COUNT(*) FROM computing_resources;")
	for num.Next() {
		num.Scan(&n)
		// Если не пустая, ничего не делаем и выходим
		if n == 0 {
			// Добавляем вычислители
			for i := 0; i < 3; i++ {
				conn := Connect()
				defer conn.Close(context.Background())
				stmt := "INSERT INTO computing_resources(ind, status, last_active) VALUES (%d, 'running', %d)"
				_, err := conn.Exec(context.Background(), fmt.Sprintf(stmt, i+1, 0))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Exec for set default computing resources failed: %v\n", err)
					return
				}
			}
			fmt.Println("Succesfully inserted default computing resources")
		}
	}

	conn = Connect()
	defer conn.Close(context.Background())

	// Проверяем пустая ли таблица с операциями
	num, _ = conn.Query(context.Background(), "SELECT COUNT(*) FROM operations;")
	for num.Next() {
		num.Scan(&n)
		// Если не пустая, ничего не делаем и выходим
		if n != 0 {
			return
		}
	}

	conn = Connect()
	defer conn.Close(context.Background())

	// Добавляем доступные операции
	var insertStmt string = `INSERT INTO operations(name, duration) VALUES 
							('+', 10),
							('-', 10), 
							('*', 10), 
							('/', 10)`
	_, err = conn.Exec(context.Background(), insertStmt)
	if err != nil {
		fmt.Printf("Exec for insert operations into table failed: %v\n", err)
	}
	fmt.Println("Succesfully inserted default operations)")
}
