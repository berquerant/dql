# dql

Find files or directories by sql.

## Syntax

```
SELECT select_expr [, select_expr ...]
[WHERE where_condition]
[GROUP BY col_name]
[HAVING having_condition]
[ORDER BY order_by_expr [ASC | DESC]]
[LIMIT row_count [OFFSET offset]]
```

### SELECT

`select_expr` is a column that you want.

```
select name, size;
```

Select all columns.

```
select all;
```

is equivalent to

```
select name, size, mode, mod_time, is_dir;
```

Give a temporary name:

```
select name as N;
```

The temporary name is available in `WHERE`, `HAVING` and `ORDER BY`.

Aggregations are available with conditions below:

a. without `GROUP BY` and select aggregations only.
b. with `GROUP BY` then except `GROUP BY` column.

### WHERE

`where_condition` is a condition expr, if the evaluated value of a row is true then the row is selected.

```
select all where is_dir;
```

### GROUP BY

`GROUP BY` aggregates rows by `col_name`.
`SELECT` must contain `col_name` and not contain raw columns but `col_name`.

```
select is_dir, count(name) group by is_dir;
```

### HAVING

`having_condition` is a condition expr, if the evaluated value of a row is true then the row is selected.
`HAVING` must be written with `GROUP BY`.

```
select mode, count(name) group by mode having count(name) > 5;
```

### ORDER BY

`order_by_expr` is a expr on which to sort rows.

```
select len(name) as nlen, name order by nlen desc;
```

### LIMIT

`row_count` constrains the number of the result rows.

```
select all limit 3;
```

If `offset` is exist, ignore first `offset` rows.

```
select all limit 3 offset 5;
```

## Columns

- `name` is the path.
- `size` is the size, number of bytes in the file.
- `mode` is the file mode, the entry type and permissions.
- `mod_time` is the last modification time.
- `is_dir` is the limited entry type, if true, the row comes from a file.

## Data types

| Name   | Description    | Example       |
|--------|----------------|---------------|
| int    | integer        | 10, -1        |
| float  | floating point | 1.2, -0.5     |
| string | string         | "str"         |
| bool   | bool           | (no literals) |

Hereafter, int or float are referred to as number,
and a string literal matched with `[01]+` is referred to as bits.

## Operators

The operators in the more lower row has the higher precedence.

| Format                | Description           | Argument Types            | Result Type | Example                     |
|-----------------------|-----------------------|---------------------------|-------------|-----------------------------|
| or                    | or                    | bool, bool                | bool        | is_dir or size < 100        |
| and                   | and                   | bool, bool                | bool        | is_dir and name = "x"       |
| xor                   | xor                   | bool, bool                | bool        | size < 100 xor name = "x"   |
| =                     | equal                 | any, any (same type)      | bool        | name = "x"                  |
| <>                    | not equal             | any, any (same type)      | bool        | name <> "x"                 |
| <                     | less than             | any, any (same type)      | bool        | size < 100                  |
| <=                    | less than or equal    | any, any (same type)      | bool        | size <= 100                 |
| >                     | greater than          | any, any (same type)      | bool        | size > 100                  |
| >=                    | greater than or equal | any, any (same type)      | bool        | size >= 100                 |
| . in (...)            | within                | any, list of any          | bool        | name in ("x", "y")          |
| . not in (...)        | not within            | any, list of any          | bool        | name not in ("x", "y")      |
| . between . and .     | between               | any, any, any (same type) | bool        | size between 10 and 100     |
| . not between . and . | not between           | any, any, any (same type) | bool        | size not between 10 and 100 |
| . like REGEX          | matched string        | string, string            | bool        | name like "log$"            |
| . not like REGEX      | not matched string    | string, string            | bool        | name not like "log$"        |
| +                     | add                   | number, number            | number      | size + 1                    |
| \-                    | subtract              | number, number            | number      | size - 1                    |
| \*                    | multiply              | number, number            | number      | size * 2                    |
| /                     | division              | number, number            | number      | size / 2                    |
| \|                    | bit or                | int or bits, int or bits  | int         | 3 \| 4                      |
| &                     | bit and               | int or bits, int or bits  | int         | 3 & 4                       |
| ^                     | bit xor               | int or bits, int or bits  | int         | 3 ^ 4                       |
| +                     | noop                  | any                       | any         | +1                          |
| \-                    | unary minus           | number                    | number      | -1                          |
| ~                     | bit not               | int or bits               | int         | ~15                         |
| not                   | not                   | bool                      | bool        | not is_dir                  |

## Functions

The function converts an expr or a column into some value.

| Format     | Description                      | Argument Types | Result Type | Example                  |
|------------|----------------------------------|----------------|-------------|--------------------------|
| pow(x, y)  | x to the power of y              | number, number | number      | pow(2, 3)                |
| ceil(x)    | ceiling                          | number         | int         | ceil(2.3)                |
| floor(x)   | floor                            | number         | int         | floor(2.3)               |
| len(x)     | length of string                 | string         | int         | len("length")            |
| base(x)    | the last element of path         | string         | string      | base("dir/file")         |
| dir(x)     | all but the last element of path | string         | string      | dir("dir/file")          |
| ext(x)     | the file name extension          | string         | string      | ext("dired.elc")         |
| bin2int(x) | bits to int                      | bits           | int         | bin2int("1010")          |
| int2bin(x) | int to bits                      | int            | bits        | int2bin(10)              |
| cast(x, y) | cast x to y                      | any            | string      | cast(10, "string")       |
| now()      | the current local time           |                | int         | now()                    |
| depth(x)   | the depth of the path            | name           | int         | depth("/home/user")      |
| grep(x, y) | `grep x y`                       | string         | string      | grep("lambda", "map.py") |

### Cast

`cast(value, "destination type")` cast value to destination type.

If the type of the value equals the destination type, then the result is the value itself.
More conversion rules are below:

| Value Type | Destination Type | Description                                                                  |
|------------|------------------|------------------------------------------------------------------------------|
| float      | int              | floor                                                                        |
| string     | int              | parse                                                                        |
| bool       | int              | true into 1, false into 0                                                    |
| int        | float            | type conversion only                                                         |
| string     | float            | parse                                                                        |
| bool       | float            | true into 1.0, false into 0.0                                                |
| any        | string           | format as string                                                             |
| number     | bool             | value != 0                                                                   |
| string     | bool             | value != ""                                                                  |
| string     | timestamp        | parse string like "2006-01-02T15:04:05Z07:00" (RFC3339) into timestamp (int) |
| number     | time             | timestamp into time (string)                                                 |
| number     | duration         | seconds into duration (string)                                               |

Undefined rules cause conversion errors.

### Aggregations

The aggregation converts the rows into some value.

| Format     | Description           | Argument Types | Result Type     | Example       |
|------------|-----------------------|----------------|-----------------|---------------|
| count(x)   | number of the rows    | any            | int             | count(name)   |
| min(x)     | minimum of the rows   | any            | any (same type) | min(size)     |
| max(x)     | maximum of the rows   | any            | any (same type) | max(name)     |
| product(x) | product of the rows   | number         | number          | product(size) |
| sum(x)     | summation of the rows | number         | number          | sum(size)     |
| avg(x)     | average of the rows   | number         | number          | avg(size)     |

## Reserved words

The reserved words are case insensitive.

```
select where having group by order limit as asc desc like in not and or xor between offset
```
