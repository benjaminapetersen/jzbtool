# jzbtool

A tiny go program for decoding a jzb string into JSON, or encoding JSON into a jzb.

![Example Usage](https://i.imgur.com/v0N9c7m.png)

## Example Usage:

```bash
jbztool -jzb eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO // { "name": "Jane" }
jbztool -json '{"name": "Jane"}' // eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO
```


## Instructions

1. Download the binary for your platform (OSX, Linux)
2. Install it in your ~/bin dir (or somewhere that lets you easily use it)
3. Run `jzbtool -jbz eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO` to see the json decoded (and pretty printed!) output
4. Run `jzbtool -jbz eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO --color=false --pretty=false` to see the json blob without fancy features
5. Run `jsbtool -json  '{"name": "Jane"}'` to see the jzb encoded output string

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
jbztool -json '{"name": "Jane" }'
jbztool --json '{"name": "Jane" }'
jbztool -json='{"name": "Jane" }'
jbztool --json='{"name": "Jane" }'
// except for the 2 boolean flags, the = is required:
jbztool -pretty=true -jbz eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO
jbztool -pretty=true -color=true -jbz eJxSqo5RykvMTY1RsopR8krMS41RqlUCBAAA__9I_AaO
```

## Gotchas

- Remember that `{ "name": "Jane" }` is valid JSON but `{ name: "Jane" }` and `{ 'name': 'Jane' }` are not.  The parser expects strict JSON. Use double quotes!
- The cli tool expects the JSON arg passed within single quotes.  For example, `jbztool -json '{"name": "Jane" }'`.

## Credits

[Allan McNaughton](https://github.com/amcnaughton) did all of the real work, I just added the pretty print options. 