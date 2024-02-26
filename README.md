#cuesheet

`cuesheet` is a go library, that builds a (parse) tree out of a cue file. 

### Usage
```go
  
	c, err := cuesheet.New(string(source))

	if err != nil {
		log.Fatal(err)
	}

	for i, v := range c.Files() {
		fmt.Printf("#[%d] %s\n", i+1, v.Path)
		for j, track := range v.Tracks {
			fmt.Printf("  |%d| %d\n", j+1, track.ID)
			for k, id := range track.Indices {
				fmt.Printf("  {%d} %d %d %d\n", k+1, id.Minutes, id.Seconds, id.Frames)
			}
		}
	}
```




