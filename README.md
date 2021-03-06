# GoRethink - RethinkDB Driver for Go 

[![GitHub tag](https://img.shields.io/github/tag/dancannon/gorethink.svg?style=flat)]()
[![GoDoc](https://godoc.org/github.com/dancannon/gorethink?status.png)](https://godoc.org/github.com/dancannon/gorethink)
[![wercker status](https://app.wercker.com/status/e315e764041af8e80f0c68280d4b4de2/s/master "wercker status")](https://app.wercker.com/project/bykey/e315e764041af8e80f0c68280d4b4de2) 

[Go](http://golang.org/) driver for [RethinkDB](http://www.rethinkdb.com/) 


Current version: v0.6.3 (RethinkDB v1.16) 

**Version 0.6 introduced some small API changes and some significant internal changes, for more information check the [change log](CHANGELOG.md) and please be aware the driver is not yet stable**

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dancannon/gorethink?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Installation

```sh
go get -u github.com/dancannon/gorethink
```

## Connection

### Basic Connection

Setting up a basic connection with RethinkDB is simple:

```go
import (
    r "github.com/dancannon/gorethink"
)

var session *r.Session

session, err := r.Connect(r.ConnectOpts{
    Address:  "localhost:28015",
    Database: "test",
})

if err != nil {
    log.Fatalln(err.Error())
}

```
See the [documentation](http://godoc.org/github.com/dancannon/gorethink#Connect) for a list of supported arguments to Connect().

### Connection Pool

The driver uses a connection pool at all times, by default it creates and frees connections automatically. It's safe for concurrent use by multiple goroutines.

To configure the connection pool `MaxIdle`, `MaxOpen` and `IdleTimeout` can be specified during connection. If you wish to change the value of `MaxIdle` or `MaxOpen` during runtime then the functions `SetMaxIdleConns` and `SetMaxOpenConns` can be used.

```go
import (
    r "github.com/dancannon/gorethink"
)

var session *r.Session

session, err := r.Connect(r.ConnectOpts{
    Address:  "localhost:28015",
    Database: "test",
    MaxIdle: 10,
    MaxOpen: 10,
})
if err != nil {
    log.Fatalln(err.Error())
}

session.SetMaxOpenConns(5)
```

A pre-configured [Pool](http://godoc.org/github.com/dancannon/gorethink#Pool) instance can also be passed to Connect().


## Query Functions

This library is based on the official drivers so the code on the [API](http://www.rethinkdb.com/api/) page should require very few changes to work.

To view full documentation for the query functions check the [GoDoc](http://godoc.org/github.com/dancannon/gorethink#Term)

Slice Expr Example
```go
r.Expr([]interface{}{1, 2, 3, 4, 5}).Run(session)
```
Map Expr Example
```go
r.Expr(map[string]interface{}{"a": 1, "b": 2, "c": 3}).Run(session)
```
Get Example
```go
r.Db("database").Table("table").Get("GUID").Run(session)
```
Map Example (Func)
```go
r.Expr([]interface{}{1, 2, 3, 4, 5}).Map(func (row Term) interface{} {
    return row.Add(1)
}).Run(session)
```
Map Example (Implicit)
```go
r.Expr([]interface{}{1, 2, 3, 4, 5}).Map(r.Row.Add(1)).Run(session)
```
Between (Optional Args) Example
```go
r.Db("database").Table("table").Between(1, 10, r.BetweenOpts{
    Index: "num",
    RightBound: "closed",
}).Run(session)
```


### Optional Arguments

As shown above in the Between example optional arguments are passed to the function as a struct. Each function that has optional arguments as a related struct. This structs are named in the format FunctionNameOpts, for example BetweenOpts is the related struct for Between.

## Results

Different result types are returned depending on what function is used to execute the query.

- `Run` returns a cursor which can be used to view all rows returned.
- `RunWrite` returns a WriteResponse and should be used for queries such as Insert, Update, etc...
- `Exec` sends a query to the server and closes the connection immediately after reading the response from the database. If you do not wish to wait for the response then you can set the `NoReply` flag.

Example:

```go
res, err := r.Db("database").Table("tablename").Get(key).Run(session)
if err != nil {
    // error
}
```

Cursors have a number of methods available for accessing the query results

- `Next` retrieves the next document from the result set, blocking if necessary.
- `All` retrieves all documents from the result set into the provided slice.
- `One` retrieves the first document from the result se.

Examples:

```go
var row interface{}
for res.Next(&row) {
    // Do something with row
}
if res.Err() != nil {
    // error
}
```

```go
var rows []interface{}
err := res.All(&rows)
if err != nil {
    // error
}
```

```go
var row interface{}
err := res.One(&row)
if err == r.ErrEmptyResult {
    // row not found
}
if err != nil {
    // error
}
```

## Encoding/Decoding Structs
When passing structs to Expr(And functions that use Expr such as Insert, Update) the structs are encoded into a map before being sent to the server. Each exported field is added to the map unless

  - the field's tag is "-", or
  - the field is empty and its tag specifies the "omitempty" option.

Each fields default name in the map is the field name but can be specified in the struct field's tag value. The "gorethink" key in
the struct field's tag value is the key name, followed by an optional comma
and options. Examples:

```go
// Field is ignored by this package.
Field int `gorethink:"-"`
// Field appears as key "myName".
Field int `gorethink:"myName"`
// Field appears as key "myName" and
// the field is omitted from the object if its value is empty,
// as defined above.
Field int `gorethink:"myName,omitempty"`
// Field appears as key "Field" (the default), but
// the field is skipped if empty.
// Note the leading comma.
Field int `gorethink:",omitempty"`
```

Alternatively you can implement the FieldMapper interface  by providing the FieldMap function which returns a map of strings in the form of `"FieldName": "NewName"`. For example:

```go
type A struct {
    Field int
}

func (a A) FieldMap() map[string]string {
    return map[string]string{
        "Field": "myName",
    }
}
```

## Examples

View other examples on the [wiki](https://github.com/dancannon/gorethink/wiki/Examples).

## License

Copyright 2013 Daniel Cannon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
