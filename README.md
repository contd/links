# Links REST Service

This is a rest api service for managing saved links for personal use and to replace the storage currently used by the front end web app.

## Development

Besure to add to and run tests for development.  This will create a test database `saved_test.sqlite` if one does not exist.  The file is not deleted after tests run so if changes are made to an existing table's structure, the database file should be removed before running the tests. Run the following commands to get verbose testing results and then get basic code coverage:

```bash
go test -v .
go test -cover .
```


### Code Coverage

For getting code coverage information and a nice html output the following can and was run to generate both the `coverage.out` file and opens a browser showing the nice html coverage report:

```bash
go test -coverprofile coverage.out .
go tool cover  -html=coverage.out
```

Also, to add the frequency of lines executed add the `-covermode` flag to the test command:

```bash
go test -covermode=count -coverprofile coverage.out .
go tool cover  -html=coverage.out
```

### Docker Build

Once the tests all pass you can build your own docker image with the included docker file like so:

```bash
docker build -t contd/links .
```
If you use `go install` the binary will expect the saved.sqlite file to be in the same directory.  You can override this by passing an environment variable like so:

```bash
SQLITE_PATH=/some/other/path/saved.sqlite links
```

This assumes your `PATH` includes `$GOPATH/bin` and you must include the file name of the sqlite database file you want to use.  The tables are not created by the application and it expects them to be there so you can use the one created from running the tests and just rename it.

## Docker Running

To run this in docker and persist the sqlite database, use the following once you've created an image from the `Dockerfile`:

```bash
docker run --name golinks -d -p 5555:5555 -v $GOPATH/src/githib.com/contd/links:/data contd/links
```

Then use the `docker-links.service` file to make the service autostart on boot:

```bash
sudo cp $GOPATH/src/github.com/contd/links/docker-links.service /etc/systemd/system/
sudo systemctl enable docker-links.service
sudo systemctl start docker-links.service
```

