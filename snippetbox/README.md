# golang-start

1. go mod init name
2. go run . || go run name
3. go run ./cmd/web -addr=":9999"
4. go get link package: ex: go get github.com/go-sql-driver/mysql@v1
5. go mod verify
6. go mod download
7. upgrade version: go get -u github.com/foo/bar@v2.0.0
8. upgrade minor or patch: go get -u github.com/foo/bar
9. Remove unused packages: go get github.com/foo/bar@none
10. Remove all unused package in go.sum, go.mod: go mod tidy -v
11. Create data mysql
12. Flow middleware: secureHeaders → servemux → application handler → servemux → secureHeaders
13. Flow log: logRequest ↔ secureHeaders ↔ servemux ↔ application handler
14. TLS: go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
15. access mysql in terminal: mysql -D snippetbox -u root -p

    -- Create a new UTF-8 `snippetbox` database. CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE
    utf8mb4_unicode_ci; -- Switch to using the `snippetbox` database. USE snippetbox; -- Create a `snippets` table.
    CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, title VARCHAR(100) NOT NULL, content TEXT NOT NULL, created DATETIME
    NOT NULL, expires DATETIME NOT NULL
    ); -- Add an index on the created column. CREATE INDEX idx_snippets_created ON snippets(created); -- Add some dummy
    records (which we'll use in the next couple of chapters). INSERT INTO snippets (title, content, created, expires)
    VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō', UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    ); INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki', UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    ); INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo', UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
    ); ========Create user for mysql============= CREATE USER 'web'@'localhost'; GRANT SELECT, INSERT, UPDATE ON
    snippetbox.\* TO 'web'@'localhost'; -- Important: Make sure to swap 'pass' with a password of your own choosing.
    ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

mysql -D snippetbox -u web -p ====================Command Database==================== DB.Query() is used for SELECT
queries which return multiple rows. DB.QueryRow() is used for SELECT queries which return a single row. DB.Exec() is
used for statements which don’t return rows (like INSERT and DELETE).
