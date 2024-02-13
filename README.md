#cuesheet

`cuesheet` is a go library, that builds a (parse) tree out of a cue file. 
It doesn't support all cue files, but the most common ones. 
### Usage
```go
    input := "TITLE S.O.S."
    parse := NewParser(input)
	tree, err := parse.Start()
	for _, r := tree.Traverse() {
        ...
    }
```




