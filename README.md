# mockdata

[![Go](https://github.com/kumarvv/mockdata/actions/workflows/go.yml/badge.svg)](https://github.com/kumarvv/mockdata/actions/workflows/go.yml)

A command-line mock data generator written in Go. Define your data schema in a YAML config file and generate realistic test data as JSON, CSV, or SQL INSERT statements.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration File](#configuration-file)
- [Target Types](#target-types)
- [Function Types](#function-types)
  - [Static / Literal](#static--literal)
  - [Sequence](#sequence)
  - [Names & Identity](#names--identity)
  - [Contact & Location](#contact--location)
  - [Numbers](#numbers)
  - [Dates & Times](#dates--times)
  - [Random Choice](#random-choice)
  - [Text](#text)
  - [Advanced](#advanced)
- [Sample Files](#sample-files)

---

## Installation

```bash
git clone https://github.com/kumarvv/mockdata.git
cd mockdata
go build -o mockdata .
```

---

## Usage

```bash
./mockdata <config-file>
```

**Example:**

```bash
./mockdata sample.yml
```

The tool reads the YAML config file, generates data according to the schema, and writes output files to the directory specified in `to_path`. Output files are named `<table_name>.<target_type>` (e.g., `users.json`, `tags.csv`, `orders.sql`).

---

## Configuration File

Config files are written in YAML with two top-level keys: `target` and `tables`.

```yaml
target:
  type: json              # Output format: json | csv | sql (required)
  to_path: /output/dir   # Directory to write output files (required)

tables:
  - name: users           # Table/file name (required)
    row_count: 100        # Number of rows to generate (optional, default: 1)
    columns:
      - column_name: function_name(param1=value1,param2=value2)
      - id: uuid()
      - name: random_full_name()
```

### Target Fields

| Field        | Required | Description |
|--------------|----------|-------------|
| `type`       | Yes      | Output format — `json`, `csv`, or `sql` |
| `to_path`    | Yes      | Directory path where output files are written |

### Table Fields

| Field       | Required | Description |
|-------------|----------|-------------|
| `name`      | Yes      | Table name; used as the output file name |
| `row_count` | No       | Number of rows to generate (default: `1`) |
| `columns`   | Yes      | List of column definitions as `column_name: function_expr` pairs |

### Column Syntax

Each column is a key-value pair where the key is the column name and the value is a function expression:

```yaml
- column_name: function_name(param1=value1,param2=value2)
```

Parameters are comma-separated `key=value` pairs inside parentheses. For functions that take a single `value` parameter, the shorthand positional form works too:

```yaml
- status: string(active)         # shorthand
- status: string(value=active)   # explicit
```

---

## Target Types

### `json`

Generates a JSON array of objects written to `<table_name>.json`.

```yaml
target:
  type: json
  to_path: /output
```

Output example (`users.json`):
```json
[
  {
    "id": "a1b2c3d4-...",
    "name": "Jane Smith",
    "email": "jane.smith@example.com"
  }
]
```

### `csv`

Generates a CSV file with a header row written to `<table_name>.csv`. Values are quoted and escaped automatically.

```yaml
target:
  type: csv
  to_path: /output
```

Output example (`users.csv`):
```
id,name,email
a1b2c3d4-...,Jane Smith,jane.smith@example.com
```

### `sql`

Generates SQL `INSERT` statements written to `<table_name>.sql`. Table names are lowercased. String values are quoted.

```yaml
target:
  type: sql
  to_path: /output
```

Output example (`users.sql`):
```sql
INSERT INTO users (id, name, email) VALUES ('a1b2c3d4-...', 'Jane Smith', 'jane.smith@example.com');
```

---

## Function Types

### Static / Literal

These functions output a fixed value on every row.

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `string` | `value` (required) | Static string | `string(active)` |
| `integer` | `value` (required) | Static integer | `integer(42)` |
| `float` | `value` (required) | Static float | `float(3.14)` |
| `boolean` | `value` (required) | Static boolean — `true`, `false`, `y`, `yes` | `boolean(true)` |
| `date` | `value` (required) | Static date in `YYYY-MM-DD` format | `date(2024-01-15)` |
| `datetime` | `value` (required) | Static datetime in `YYYY-MM-DD HH:MM:SS` format | `datetime(2024-01-15 09:30:00)` |
| `sql` | `value` (required) | SQL expression (requires DB connection) | `sql(NOW())` |

**Usage:**

```yaml
columns:
  - source: string(IMPORT)
  - version: integer(1)
  - price: float(9.99)
  - is_active: boolean(true)
  - start_date: date(2024-01-01)
  - created_at: datetime(2024-01-01 00:00:00)
```

---

### Sequence

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `uuid` | none | Generates a UUID v4 | `uuid()` |
| `serial` | `start` (optional, default: `1`) | Auto-incrementing integer starting from `start` | `serial(start=100)` |

**Usage:**

```yaml
columns:
  - id: uuid()
  - seq_num: serial()
  - order_num: serial(start=1000)
```

---

### Names & Identity

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_title` | `gender`, `case` | Random title (Mr., Mrs., Dr., etc.) | `random_title(gender=female)` |
| `random_first_name` | `gender`, `case` | Random first name | `random_first_name(gender=male)` |
| `random_last_name` | `case` | Random last name | `random_last_name()` |
| `random_full_name` | `gender`, `case` | Random full name | `random_full_name()` |
| `random_gender` | `case` | Random gender (`male` or `female`) | `random_gender()` |

**Parameters:**

- `gender` — `male` or `female`; if omitted, randomly selected per row
- `case` — `upper` or `lower`; if omitted, mixed case is used

**Usage:**

```yaml
columns:
  - title: random_title(gender=female)
  - first_name: random_first_name()
  - last_name: random_last_name(case=upper)
  - full_name: random_full_name(gender=male)
  - gender: random_gender()
```

---

### Contact & Location

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_email` | `case`, `domain` | Random email address | `random_email(domain=example.com)` |
| `random_phone` | none | Random phone number | `random_phone()` |
| `random_address` | `case`, `country` | Random full street address | `random_address()` |
| `random_street` | `case` | Random street name | `random_street()` |
| `random_city` | `case`, `country` | Random city name | `random_city(country=US)` |
| `random_state` | `case`, `country` | Random state (full name) | `random_state(country=US)` |
| `random_state2` | `case`, `country` | Random state abbreviation (2-char) | `random_state2()` |
| `random_country` | `case` | Random full country name | `random_country()` |
| `random_country2` | `case` | Random 2-char country code | `random_country2()` |
| `random_country3` | `case` | Random 3-char country code | `random_country3()` |
| `random_currency` | `case` | Random currency code | `random_currency()` |

**Parameters:**

- `case` — `upper` or `lower`
- `country` — country filter (e.g., `US`)
- `domain` — custom email domain (e.g., `mycompany.com`)

**Usage:**

```yaml
columns:
  - email: random_email()
  - work_email: random_email(domain=acme.com)
  - phone: random_phone()
  - address: random_address()
  - street: random_street()
  - city: random_city(country=US)
  - state: random_state()
  - state_code: random_state2()
  - country: random_country()
  - country_code: random_country2()
  - currency: random_currency(case=upper)
```

---

### Numbers

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_number` | `min`, `max` | Random integer in range | `random_number(min=1,max=100)` |
| `random_decimal` | `min`, `max` | Random decimal (float) in range | `random_decimal(min=0,max=9.99)` |
| `random_range` | `min`, `max` | Random number within range | `random_range(min=10,max=50)` |
| `random_boolean` | none | Random `true` or `false` | `random_boolean()` |

**Usage:**

```yaml
columns:
  - age: random_number(min=18,max=99)
  - score: random_decimal(min=0,max=10)
  - quantity: random_range(min=1,max=500)
  - is_verified: random_boolean()
```

---

### Dates & Times

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_date` | `format` | Random date; default format `2006-01-02` | `random_date(format=2006-01-02)` |
| `random_day` | none | Random day of month (1–31) | `random_day()` |
| `random_month` | none | Random month name | `random_month()` |
| `random_year` | none | Random year (1900–2999) | `random_year()` |

> **Note:** The `format` parameter for `random_date` uses Go's reference time: `2006-01-02 15:04:05`.

**Usage:**

```yaml
columns:
  - birth_date: random_date()
  - event_date: random_date(format=01/02/2006)
  - day: random_day()
  - month: random_month()
  - year: random_year()
```

---

### Random Choice

Pick a random value from a pipe-separated list of options.

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_in_string` | `value` (required) | Random string from `\|`-separated list | `random_in_string(admin\|user\|viewer)` |
| `random_in_integer` | `value` (required) | Random integer from `\|`-separated list | `random_in_integer(1\|2\|3\|4\|5)` |
| `random_in_float` | `value` (required) | Random float from `\|`-separated list | `random_in_float(0.5\|1.0\|1.5\|2.0)` |

**Usage:**

```yaml
columns:
  - role: random_in_string(admin|editor|viewer)
  - status: random_in_string(active|inactive|pending)
  - priority: random_in_integer(1|2|3)
  - discount: random_in_float(0.05|0.10|0.15|0.20)
```

---

### Text

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_string` | `len`, `min`, `max`, `case` | Random alphanumeric string | `random_string(len=10)` |
| `random_paragraph` | `len`, `min`, `max`, `case` | Random paragraph of text | `random_paragraph(min=50,max=200)` |

**Parameters:**

- `len` — exact character length
- `min` / `max` — random length between min and max characters
- `case` — `upper` or `lower`

**Usage:**

```yaml
columns:
  - token: random_string(len=32,case=upper)
  - username: random_string(min=5,max=15)
  - bio: random_paragraph(min=100,max=500)
  - notes: random_paragraph(len=200,case=lower)
```

---

### Advanced

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `random_format` | `numOfPairs`, `separator` | Formatted random string pairs joined by separator | `random_format(numOfPairs=3,separator=-)` |
| `random_from_sql` | `value` (required) | Random value from a SQL query result (requires DB connection) | `random_from_sql(SELECT id FROM products)` |

**Usage:**

```yaml
columns:
  - tracking_code: random_format(numOfPairs=4,separator=-)
  - product_id: random_from_sql(SELECT id FROM products ORDER BY RANDOM() LIMIT 1)
```

---

## Sample Files

### JSON Output — User Records

Generates 10 user records as a JSON file.

```yaml
target:
  type: json
  to_path: /tmp/mockdata

tables:
  - name: users
    row_count: 10
    columns:
      - id: uuid()
      - seq: serial()
      - title: random_title(gender=male)
      - first_name: random_first_name()
      - last_name: random_last_name()
      - full_name: random_full_name()
      - gender: random_gender()
      - email: random_email()
      - phone: random_phone()
      - address: random_address()
      - city: random_city()
      - state: random_state2()
      - country: random_country2()
      - username: random_string(len=12,case=lower)
      - role: random_in_string(admin|editor|viewer)
      - rating: random_in_integer(1|2|3|4|5)
      - score: random_decimal(min=0,max=100)
      - is_active: boolean(true)
      - source: string(WEB)
      - created_at: random_date()
```

**Sample output (`users.json`):**

```json
[
  {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "seq": 1,
    "title": "Mr.",
    "first_name": "James",
    "last_name": "Peterson",
    "full_name": "James Peterson",
    "gender": "male",
    "email": "james.peterson@randommail.com",
    "phone": "555-867-5309",
    "address": "742 Evergreen Terrace",
    "city": "Springfield",
    "state": "IL",
    "country": "US",
    "username": "jamesp2024ab",
    "role": "editor",
    "rating": 4,
    "score": 73.42,
    "is_active": true,
    "source": "WEB",
    "created_at": "2023-07-14"
  }
]
```

---

### CSV Output — Tag Records

Generates 5 tag records as a CSV file.

```yaml
target:
  type: csv
  to_path: /tmp/mockdata

tables:
  - name: tags
    row_count: 5
    columns:
      - uuid: uuid()
      - username: random_string(len=25,case=upper)
      - full_name: random_full_name(gender=male)
      - email: random_email()
      - created_at: random_date()
      - sys_key: string(ABCD)
      - is_active: boolean(true)
      - version: integer(123)
      - rating: random_in_integer(1|2|3|4|5)
```

**Sample output (`tags.csv`):**

```
uuid,username,full_name,email,created_at,sys_key,is_active,version,rating
f47ac10b-58cc-4372-a567-0e02b2c3d479,XKQZRMJWDNLFPVHBCEAYTOU,James Peterson,james.peterson@randommail.com,2023-07-14,ABCD,true,123,3
```

---

### SQL Output — Order Records

Generates 20 order records as SQL INSERT statements.

```yaml
target:
  type: sql
  to_path: /tmp/mockdata

tables:
  - name: orders
    row_count: 20
    columns:
      - order_id: uuid()
      - customer_id: uuid()
      - status: random_in_string(pending|processing|shipped|delivered|cancelled)
      - total_amount: random_decimal(min=10,max=9999)
      - item_count: random_number(min=1,max=20)
      - currency: random_currency(case=upper)
      - ship_country: random_country2()
      - notes: random_paragraph(min=10,max=80)
      - created_at: random_date()
      - is_paid: random_boolean()
```

**Sample output (`orders.sql`):**

```sql
INSERT INTO orders (order_id, customer_id, status, total_amount, item_count, currency, ship_country, notes, created_at, is_paid) VALUES ('f47ac10b-...', 'a1b2c3d4-...', 'shipped', 349.75, 3, 'USD', 'US', 'Please leave at front door.', '2024-02-10', true);
```

---

### Multi-Table Config

Generate multiple related tables in a single run.

```yaml
target:
  type: json
  to_path: /tmp/mockdata

tables:
  - name: categories
    row_count: 5
    columns:
      - id: serial()
      - name: random_in_string(Electronics|Clothing|Books|Home|Sports)
      - code: random_string(len=4,case=upper)

  - name: products
    row_count: 50
    columns:
      - id: uuid()
      - sku: random_string(len=8,case=upper)
      - name: random_string(min=5,max=30)
      - price: random_decimal(min=1,max=999)
      - stock: random_number(min=0,max=500)
      - is_available: random_boolean()
      - created_at: random_date()

  - name: customers
    row_count: 25
    columns:
      - id: uuid()
      - first_name: random_first_name()
      - last_name: random_last_name()
      - email: random_email()
      - phone: random_phone()
      - city: random_city()
      - country: random_country2()
      - joined_at: random_date()
```
