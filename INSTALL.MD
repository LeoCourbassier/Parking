# Install Go

# Test
```bash
$ make test
```

# Compile
```bash
$ make
```

# Run
Create file named `configs.sh` with the following:
```bash
export DATABASE_USER=postgres
export DATABASE_PASS=postgres
export DATABASE_URL=localhost
export DATABASE_PORT=5432
export DATABASE_SCHEMA=mlabs
```

Then, source it:
```bash
$ source configs.sh
```

Run the application:
```bash
$ ./cmd/br.com.mlabs
```

# Test image recognition
You can use postman.

* Set url to: `localhost:4000/parking/in`
* Go to body, set: `form-data`
* In keys, set: `plate`, and change type to `file` (right side of the text box)
* In values, click upload
* Make the request :)