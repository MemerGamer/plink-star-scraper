# Plink Star Scraper

I made this project to help the [plink](https://github.com/darksinge/plink.nvim) neovim plugins development. This project is a simple web scraper that gets the star count of a plugin and saves it to a file.

## Limitations

Since I didn't wanted to overload the API I only tested for the keyword: "telescope".
I curled the request to the `search-4-telescope.json` file.

```console
curl -o search-4-telescope.json API/search\?q\=telescope
```

## How the Script Works

1. The script reads the JSON file containing the search results.
2. For each result, it checks if the URL starts with `https://github.com/`. If not, it attempts to concatenate `https://github.com/` with the plugin's name since the name often matches the GitHub repository.
3. If the URL concatenation fails, the star count is set to -1.
4. If successful, the script fetches the star count from the GitHub repository and updates it in the search result.
5. Finally, the updated star counts are saved to a new JSON file named search-4-telescope.json-updated.json.

## Running the script

```console
    go run main.go
```

or you can compile it and run the binary (it will be faster)

```console
    go build main.go
    ./main
```

| Before compilation                                           | After compilation                                    |
| ------------------------------------------------------------ | ---------------------------------------------------- |
| `go run main.go 0,68s user 0,15s system 6% cpu 13,167 total` | `./main 0,43s user 0,06s system 3% cpu 12,849 total` |
