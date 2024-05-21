# QP - Query Parser from URL Query to SQL

Query Parser (QP) is a Go library designed to streamline the creation of dynamic SQL queries for web applications. It provides an easy-to-use API for handling filtering through GET queries and manages translations and validations for user inputs.

## Installation

To install the QP library, run the following command:

```bash
go get -u github.com/nex-gen-tech/qp
```

## Idea

The inspiration for this library came from the article "[REST API Design: Filtering, Sorting, and Pagination](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)." The concepts in the article are useful for projects that involve lists with various filtering options.

## Quick Start

Refer to the `examples` folder and the test files for more comprehensive examples.

```go
package main

import (
	"errors"
	"fmt"
	"net/url"

	rqp "github.com/nex-gen-tech/qp"
)

func main() {
	url, _ := url.Parse("http://localhost/?sort=name,-id&limit=10&id=1&i[eq]=5&s[eq]=one&email[like]=*tim*|name[like]=*tim*")
	q, _ := rqp.NewParse(url.Query(), rqp.Validations{
		"limit:required": rqp.MinMax(10, 100), // limit must be present and between 10 and 100
		"sort":           rqp.In("id", "name"),   // sort could be "id" or "name"
		"s":              rqp.In("one", "two"),   // filter: s - string with equal comparison
		"id:int":         nil,                   // filter: id is an integer without additional validation
		"i:int": func(value interface{}) error { // custom validation function for "i"
			if value.(int) > 1 && value.(int) < 10 {
				return nil
			}
			return errors.New("i: must be greater than 1 and less than 10")
		},
		"email": nil,
		"name":  nil,
	})

	fmt.Println(q.SQL("table")) // SELECT * FROM table WHERE id = ? AND i = ? AND s = ? AND (email LIKE ? OR name LIKE ?) ORDER BY name, id DESC LIMIT 10
	fmt.Println(q.Where())      // id = ? AND i = ? AND s = ? AND (email LIKE ? OR name LIKE ?)
	fmt.Println(q.Args())       // [1 5 one %tim% %tim%]

	q.AddValidation("fields", rqp.In("id", "name"))
	q.SetUrlString("http://localhost/?fields=id,name&limit=10")
	q.Parse()

	fmt.Println(q.SQL("table")) // SELECT id, name FROM table ORDER BY id LIMIT 10
	fmt.Println(q.Select())     // id, name
	fmt.Println(q.Args())       // []
}
```

## Top-Level Fields

- `fields`: Specifies the fields for the SELECT clause, separated by commas (","). For example, `&fields=id,name`. If not provided, "*" is used by default. To use this filter, define a validation function. Use `rqp.In("id", "name")` to limit the fields in your query.
- `sort`: Specifies the sorting fields list, separated by commas (","). This field must be validated. A +/- prefix indicates ASC/DESC sorting. For example, `&sort=+id,-name` results in `ORDER BY id, name DESC`. Validate fields using `rqp.In("id", "name")`.
- `limit`: Specifies the limit for the LIMIT clause. The value should be greater than 0 by default. Validation for `limit` is optional. Use `rqp.Max(100)` to set an upper threshold.
- `offset`: Specifies the offset for the OFFSET clause. The value should be greater than or equal to 0 by default. Validation for `offset` is optional.

### Example URL Queries
- `?fields=id,name&sort=-created_at&limit=20&offset=5`
- `?sort=name,-id&limit=50&status[eq]=active&category[like]=tech`
- `?fields=product_id,product_name&sort=price&price[gt]=100&in_stock[eq]=true`

## Validation Modifiers

- `:required`: Indicates that the parameter is required and must be present in the query string. Raises an error if not found.
- `:int`: Specifies that the parameter must be convertible to an integer type. Raises an error if not.
- `:bool`: Specifies that the parameter must be convertible to a boolean type. Raises an error if not.

## Supported Types

- `string`: The default type for all filters unless specified otherwise. Can be compared using `eq, ne, gt, lt, gte, lte, like, ilike, nlike, nilike, in, nin, is, not` methods. (`nlike, nilike` represent `NOT LIKE, NOT ILIKE`; `in, nin` represent `IN, NOT IN`; and `is, not` are for NULL comparisons: `IS NULL, IS NOT NULL`.)
- `int`: Integer type. Must be specified with the tag ":int". Can be compared using `eq, ne, gt, lt, gte, lte, in, nin` methods.
- `bool`: Boolean type. Must be specified with the tag ":bool". Can be compared using the `eq` method.

## Date Usage

Here's an example demonstrating how to handle date values:

```go
import (
	"fmt"
	"net/url"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func main() {
	url, _ := url.Parse("http://localhost/?created_at[eq]=2020-10-02")
	q, _ := rqp.NewParse(url.Query(), rqp.Validations{
		"created_at": func(v interface{}) error {
			s, ok := v.(string)
			if !ok {
				return rqp.ErrBadFormat
			}
			return validation.Validate(s, validation.Date("2006-01-02"))
		},
	})

	q.ReplaceNames(rqp.Replacer{"created_at": "DATE(created_at)"})

	fmt.Println(q.SQL("table")) // SELECT * FROM table WHERE DATE(created_at) = ?
}
```

## Advanced Usage

QP offers advanced features like custom validation functions, handling multiple OR conditions, and more.

### Custom Validation Functions

You can define custom validation functions to perform specific checks on query parameters.

```go
func validateAge(value interface{}) error {
	age, ok := value.(int)
	if !ok {
		return rqp.ErrBadFormat
	}
	if age < 18 {
		return errors.New("age must be at least 18")
	}
	return nil
}
```

### Handling Multiple OR Conditions

QP supports handling multiple OR conditions within filters.

```go
q.AddORFilters(func(query *rqp.Query) {
	query.AddFilter("name", rqp.LIKE, "John")
	query.AddFilter("email", rqp.LIKE, "john@example.com")
})
```

### Replacing Field Names

You can replace field names in filters, fields, and sorting using the `ReplaceNames` method.

```go
q.ReplaceNames(rqp.Replacer{"created_at": "DATE(created_at)"})
```
## **[Guide for Frontend Developers](./docs/frontend-guide.md)**


## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open a pull request or issue on the GitHub repository to suggest improvements or report any issues.

## Acknowledgments

Special thanks to the contributors and supporters of this project.