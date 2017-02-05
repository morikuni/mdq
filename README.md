# mdq
[![Build Status](https://travis-ci.org/morikuni/mdq.svg?branch=master)](https://travis-ci.org/morikuni/mdq)

mdq queries multiple databases in parallel and output the results in JSON.

## Example

```sh
$ mysql -u root --database example -e "select * from hello";

+------+---------+
| id   | message |
+------+---------+
|    1 | hello   |
|    2 | world   |
+------+---------+


$ cat ~/.config/mdq/config.yaml
dbs:
  - name: "example_db"
    driver: "mysql"
    dsn: "root@/example"
    tags: ["example_tag"]


$ mdq -q "select * from hello" --tag example_tag | jq .
[
  {
    "Database": "example_db",
    "Columns": [
      "id",
      "message"
    ],
    "Rows": [
      {
        "id": "1",
        "message": "hello"
      },
      {
        "id": "2",
        "message": "world"
      }
    ]
  }
]
```

## Supported Databases and DSN

- MySQL
    - `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`
    - [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#dsn-data-source-name)
- PostgreSQL
    - `host=host dbname=dbname user=user password=password ssl=require`
    - [github.com/lib/pq](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters)
