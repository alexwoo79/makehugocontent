package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // 使用 modernc.org 的 SQLite 驱动
)

var (
	DataDB *sql.DB // 文章列表数据库
	UserDB *sql.DB // 用户数据库
)

func Init() {
	var err error
	UserDB, err = sql.Open("sqlite", "./database/users.db")
	if err != nil {
		log.Fatal("用户数据库打开失败:", err)
	}

	// 用户表 来创建用户密码和部门数据
	_, _ = UserDB.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
		);`)

	// 创建部门表
	_, _ = UserDB.Exec(`
		CREATE TABLE IF NOT EXISTS departments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
			);`)

	// 创建用户-部门关联表（多对多关系）
	_, _ = UserDB.Exec(`
			CREATE TABLE IF NOT EXISTS user_departments (
				user_id INTEGER,
				dept_id INTEGER,
				FOREIGN KEY(user_id) REFERENCES users(id),
				FOREIGN KEY(dept_id) REFERENCES departments(id),
				UNIQUE(user_id, dept_id)
				);`)

	// 创建一些部门
	_, _ = UserDB.Exec(`
				INSERT OR IGNORE INTO departments (name) VALUES
				('general'),
				('pm'),
				('cz'),
				('bim'),
				('zixun')
				;`)

	// 3. 插入管理员账户（用户名 admin，密码 123456，角色 admin）
	password := "123456"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("管理员密码加密失败:", err)
	}

	_, err = UserDB.Exec(`
		INSERT OR IGNORE INTO users (username, password, role)
		VALUES (?, ?, 'admin');
	`, "admin", string(hashed))
	if err != nil {
		log.Fatal("管理员插入失败:", err)
	}

	DataDB, err = sql.Open("sqlite", "./database/data.db")
	if err != nil {
		log.Fatal("文章列表数据库打开失败:", err)
	}

	// 上传文章列表记录
	_, _ = DataDB.Exec(`
	CREATE TABLE IF NOT EXISTS uploads(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		filepath TEXT,
		user_id INTEGER,
		author TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
}

func InsertTestUsers() {
	users := []struct {
		Username    string
		Password    string
		Role        string
		Departments []string
	}{
		{"alice", "123456", "viewer", []string{"pm"}},
		{"bob", "123456", "editor", []string{"cz", "bim"}},
		{"carol", "123456", "admin", []string{"general", "cz", "pm"}},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		res, err := UserDB.Exec("INSERT OR IGNORE INTO users(username, password, role) VALUES(?, ?, ?)", u.Username, hash, u.Role)
		if err != nil {
			log.Println("插入用户失败:", u.Username, err)
			continue
		}
		uid, _ := res.LastInsertId()

		for _, dept := range u.Departments {
			var deptID int
			err := UserDB.QueryRow("SELECT id FROM departments WHERE name=?", dept).Scan(&deptID)
			if err == nil {
				_, _ = UserDB.Exec("INSERT OR IGNORE INTO user_departments(user_id, dept_id) VALUES(?, ?)", uid, deptID)
			}
		}
	}
}
