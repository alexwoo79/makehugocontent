-- 创建部门表
CREATE TABLE IF NOT EXISTS departments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- 创建用户-部门关联表（多对多关系）
CREATE TABLE IF NOT EXISTS user_departments (
    user_id INTEGER,
    dept_id INTEGER,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(dept_id) REFERENCES departments(id),
    UNIQUE(user_id, dept_id) -- 防止重复插入
);

--创建一些部门--
INSERT OR IGNORE INTO departments (name) VALUES
  ('general'),
  ('pm'),
  ('cz'),
  ('bim'),
  ('zixun')
;

