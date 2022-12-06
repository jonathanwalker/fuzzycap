<h1 align="center">
  <br>
    <img align="center" width=500 height=200 src="assets/fuzzycap.png">
  <br>
  <br>
</h1>

FuzzyCap is a utility that allows you to take screenshots from a list of urls and returns the perceptual hash of each within a json array. This can help you determine assets that are similar in nature, group assets, or identify the application behind the urls. 

## Usage

Usage of fuzzycap is simple and all you need to do is specify a file with a list of urls. `fuzzycap -i urls.txt`. This will then output the url, screenshot location, and fuzzy hash of the site. 

```
Usage: fuzzycap [--input INPUT]

Options:
  --input INPUT, -i INPUT
                         input file with list of urls
  --help, -h             display this help and exit
```

## License

This is licensed as [BSD Clause 3](https://tldrlegal.com/license/bsd-3-clause-license-(revised)) in which if you wish to use it outside of the use cases, feel free to let me know. Feel free to contribute, distribute, and use!