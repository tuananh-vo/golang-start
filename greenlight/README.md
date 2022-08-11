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
