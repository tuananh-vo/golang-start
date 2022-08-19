\*\*\*Let’s take a moment to talk through these files and folders and explain the purpose that
they’ll serve in our finished project.

The bin directory will contain our compiled application binaries, ready for deployment
to a production server.
The cmd/api directory will contain the application-specific code for our Greenlight API
application. This will include the code for running the server, reading and writing HTTP
requests, and managing authentication.
The internal directory will contain various ancillary packages used by our API. It will
contain the code for interacting with our database, doing data validation, sending emails
and so on. Basically, any code which isn’t application-specific and can potentially be
reused will live in here. Our Go code under cmd/api will import the packages in the
internal directory (but never the other way around).
The migrations directory will contain the SQL migration files for our database.
The remote directory will contain the configuration files and setup scripts for our
production server.
The go.mod file will declare our project dependencies, versions and module path.
The Makefile will contain recipes for automating common administrative tasks — like
auditing our Go code, building binaries, and executing database migrations.
===========================================================

1. Run example flag: go run ./cmd/api -port=3030 -env=production
2. go test -run=^$ -bench=. -benchmem -count=3 -benchtime=5s
3. customize output string : fmt.Sprintf("%d mins", r)
4. response: fmt.Fprintf()
5. Download file large: curl -d @/tmp/largefile.json localhost:4000/v1/movies
6. Install postgres
   1. psql --version
   2. On Unix-based systems you can check your /etc/passwd file to
      confirm this, like so: cat /etc/passwd | grep 'postgres'
   3. sudo -u postgres psql
   4. SELECT current_user;
   5. CREATE DATABASE greenlight;
   6. Connect: \c greenlight
   7. Create user: CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
   8. Install extensions: CREATE EXTENSION IF NOT EXISTS citext;
   9. Connect data postgres: psql --host=localhost --dbname=greenlight --username=greenlight
7. Migrate postgres
   macos: brew install golang-migrate
   window, linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
   mv migrate.linux-amd64 $GOPATH/bin/migrate
8. Create folder/file migrate table:
   migrate create -seq -ext=.sql -dir=./migrations create_movies_table
   migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints
9. Run migrate
   migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
10. limit connection database: go run ./cmd/api -db-max-open-conns=1
11. index: migrate create -seq -ext .sql -dir ./migrations add_movies_indexes
12. migrate create -seq -ext=.sql -dir=./migrations create_users_table
13. run : migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
14. go get golang.org/x/crypto/bcrypt@latest
15. migrate create -seq -ext .sql -dir ./migrations create_tokens_table
16. migrate create -seq -ext .sql -dir ./migrations add_permissions
17. Support add permission to database
    psql $GREENLIGHT_DB_DSN

    -- Set the activated field for alice@example.com to true.
    UPDATE users SET activated = true WHERE email = 'alice@example.com';
    -- Give all users the 'movies:read' permission
    INSERT INTO users_permissions
    SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM users;
    -- Give faith@example.com the 'movies:write' permission
    INSERT INTO users_permissions
    VALUES (
    (SELECT id FROM users WHERE email = 'anhvt@gmail.com'),
    (SELECT id FROM permissions WHERE code = 'movies:write')
    );
    -- List all activated users and their permissions.
    SELECT email, array_agg(permissions.code) as permissions
    FROM permissions
    INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
    INNER JOIN users ON users_permissions.user_id = users.id
    WHERE users.activated = true
    GROUP BY email;

18. Supporting multiple dynamic origins: go run ./cmd/api -cors-trusted-origins="https://www.example.com https://staging.example.com"
19. Run example cors
    1. go run ./cmd/examples/cors/simple
    2. run project main: go run ./cmd/api -cors-trusted-origins="http://localhost:9000"
20. Run test display magic
    1. go run ./cmd/api -limiter-enabled=false -db-max-open-conns=50 -db-max-idle-conns=50 -db-max-idle-time=20s -port=4000
21. Load test with hey tool
    1. BODY='{"email": "alice@example.com", "password": "pa55word"}'
    2. hey -d "$BODY" -m "POST" http://localhost:4000/v1/tokens/authentication
22. Recording HTTP Status Codes: go get github.com/felixge/httpsnoop@v1.0.1
23. make run(main), make up(update database), make migration name=create_example_table (migrate new table)
24. go env GOCACHE
25. Force all packages to be rebuilt: go build -a -o=/bin/foo ./cmd/foo
26. go clean -cache
