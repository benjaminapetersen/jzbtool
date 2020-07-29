# unjzb

A tiny go program for decoding a jzb string into JSON, or encoding JSON into a jzb.

## Option flags:

```bash
  -color
    	if jzb provided, color the JSON output
  -jbz string
    	oops! its jzb, but we got your back. pass in jbz, we will assume you meant jzb and will also decode it to JSON
  -json string
    	pass in JSON as a string to encode it to jzb
  -jzb string
    	pass in jzb as a string to decode it to JSON
  -pretty
    	if jzb provided, pretty print the JSON output
```

Note that flags of any of the following format are valid:

```bash
jbztool -json {name: "Jane" }
jbztool --json {name: "Jane" }
jbztool -json={name: "Jane" }
jbztool --json={name: "Jane" }
// except for the 2 boolean flags, the = is required:
jbztool -pretty=true -jbz eJxSqlYCBAAA__8BgQDA
jbztool -pretty=true -color=true -jbz eJxSqlYCBAAA__8BgQDA
```

## Examples:

```bash
jbztool -json {"name": "Jane"} // eJxSqlYCBAAA__8BgQDA
jbztool -jzb eJxSqlYCBAAA__8BgQDA // { "name": "Jane" }
jbztool -pretty=true -jzb eJxSqlYCBAAA__8BgQDA // { "name": "Jane" }
```

## Gotchas

Remember that `{ "name": "Jane" }` is valid JSON but `{ name: "Jane" }` is not.  The parser expects strict JSON. 