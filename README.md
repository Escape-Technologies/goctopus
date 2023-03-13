# goctopus

Blazing fast graphql fingerprinting toolbox.

> ⚠️ Goctopus is still in very early development. Breaking changes are expected.

`````

                    .-'   `'.
                   /         \
                   |         ;
                   |         |           ___.--,
          _.._     |0) ~ (0) |    _.---'`__.-( (_.
   __.--'`_.. '.__.\    '--. \_.-' ,.--'`     `""`
  ( ,.--'`   ',__ /./;   ;, '.__.'`    __
  _`) )  .---.__.' / |   |\   \__..--""  """--.,_
 `---' .'.''-._.-'`_./  /\ '.  \ _.-~~~````~~~-._`-.__.'
       | |  .' _.-' |  |  \  \  '.               `~---`
        \ \/ .'     \  \   '. '-._)
         \/ /        \  \    `=.__`~-.
     jgs / /\         `) )    / / `"".`\
   , _.-'.'\ \        / /    ( (     / /
    `--~`   ) )    .-'.'      '.'.  | (
           (/`    ( (`          ) )  '-;
            `      '-;         (-'
                  _
  __ _  ___   ___| |_ ___  _ __  _   _ ___
 / _` |/ _ \ / __| __/ _ \| '_ \| | | / __|
| (_| | (_) | (__| || (_) | |_) | |_| \__ \
 \__, |\___/ \___|\__\___/| .__/ \__,_|___/
 |___/                    |_|
INFO[0000] Starting 100 workers
INFO[0000] Found: {"domain":"gontoz.escape.tech","type":"IS_GRAPHQL","url":"https://gontoz.escape.tech"}
INFO[0002] Done. Found 1 graphql endpoints
`````

## Usage

Using go:

```BASH
go install -v github.com/escape-technologies/goctopus@latest
goctopus -i input.txt -o output.jsonl
```

Using docker:

```BASH
docker run --rm -it -v $(pwd):/data escape/goctopus -i /data/input.txt -o /data/output.jsonl
```

## Main options & features

### Input

For now, goctopus only supports input from a file.  
The file path should be specified using the `-i` flag.  
The input file should contain a list of endpoints and/or urls separated by newlines.  
This is an example of a valid input file:

```
example.com
https://example.com/graphql
escape.tech
https://example.com/api
```

### Introspection fingerprinting

The `-introspect` flag enables introspection fingerprinting.  
If enabled, goctopus will detect if the introspection of graphql endpoints is enabled.

### Subdomain enumeration

The `-subdomain` flag enables subdomain enumeration.  
If enabled, goctopus will try to find graphql endpoints on subdomains of the given domains.  
The enumeration is done using [subfinder](https://github.com/projectdiscovery/subfinder).

### Field suggestion fingerprinting

The `-suggest` flag enables field suggestion fingerprinting.  
This option needs the introspection fingerprinting (`-introspect`) to be enabled.
When enabled, goctopus will try to detect if the graphql endpoint has field suggestion enabled, if the introspection is closed.  
This is useful to bruteforce fields and/or types when introspection is disabled, with tools such as [ClairvoyaceNext](https://github.com/Escape-Technologies/ClairvoyanceNext).

## Aditionnal options

```BASH
  -i string
    	Input file
  -introspect
    	Enable introspection fingerprinting
  -o string
    	Output file (json-lines format) (default "output.jsonl")
  -s	Silent
  -subdomain
    	Enable subdomain enumeration
  -suggest
    	Enable fields suggestion fingerprinting.
    	Needs "introspection" to be enabled.
  -t int
    	Request timeout (seconds) (default 30)
  -v	Verbose
  -w int
    	Max workers (default 100)
  -webhook string
    	Webhook URL
```

## Roadmap

- [ ] Better wordlist for field suggestion fingerprinting, to improve the detection performance and detection rate.
- [ ] Engine fingerprinting.
- [ ] Refactor to make goctopus usable as a go package.
- [ ] Better flags.
- [ ] Direct cli input.
- [ ] Improve performance further.
- [ ] Resume from output file.
- [ ] Custom ascii art.
