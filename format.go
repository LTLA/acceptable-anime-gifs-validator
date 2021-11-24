package main

import (
    "os"
    "encoding/json"
    "path/filepath"
    "strings"
    "errors"
    "fmt"
    "flag"
    "io/fs"
)

// Definitions of structs and constants. 

type GifInfo struct {
    Id string `json:"id"`
    ShowId string `json:"show_id"`
    Characters []string `json:"characters"`
    Sentiments []string `json:"sentiments"`
    Url string `json:"url"`
}

type ShowInfo struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Characters map[string]string `json:"characters"`
}

const suffix = ".json"

// Loading functions.

func LoadGifMetadata(path string) (*GifInfo, error) {
    handle, err := os.Open(path)
    if (err != nil) {
        return nil, err
    }
    defer handle.Close()

    dec := json.NewDecoder(handle)

    var output GifInfo
    err = dec.Decode(&output)
    if (err != nil) {
        return nil, err
    }

    parts := strings.Split(path, string(os.PathSeparator))
    if (len(parts) < 2) {
        return nil, errors.New("'path' should contain a subdirectory with the show ID")
    }
    base := parts[len(parts) - 1]
    show := parts[len(parts) - 2]
    output.ShowId = show

    if (len(base) == len(suffix)) {
        return nil, errors.New("'path' should contain a non-empty name before '" + suffix + "'")
    }
    prefix := base[:len(base) - len(suffix)]
    gif_name := prefix + ".gif"
    output.Id = show + "/" + gif_name // don't use filepath.Join, we're making URL components here.

    parts[len(parts) - 1] = gif_name
    checkpath := filepath.Join(parts...)
    if _, err := os.Stat(checkpath); os.IsNotExist(err) {
        return nil, errors.New("GIF should exist at '" + path + "'")
    }

    return &output, nil
}

func LoadShowMetadata(dir string) (*ShowInfo, error) {
    path := filepath.Clean(dir) + ".json"

    handle, err := os.Open(path)
    if (err != nil) {
        return nil, err
    }
    defer handle.Close()

    dec := json.NewDecoder(handle)

    var output ShowInfo
    err = dec.Decode(&output)
    if (err != nil) {
        return nil, err
    }

    _, show := filepath.Split(dir)
    output.Id = show
    return &output, nil
}

// Overall collator function.

func CollateMetadata(dir string) ([]GifInfo, []ShowInfo, error) {
    var gifs []GifInfo
    var shows []ShowInfo
    show_ptrs := make(map[string]*ShowInfo)

    err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if (info.IsDir()) {
            if (dir != path) {
                show_info, err := LoadShowMetadata(path)
                if (err != nil) {
                    return err
                }
                shows = append(shows, *show_info)
                show_ptrs[show_info.Id] = show_info
            }
            return nil
        }

        // Avoid loading in show metadata, if this occurs inside 'dir' without further nesting.
        stub, err := filepath.Rel(dir, path)
        if (err != nil) {
            return err
        }

        subdir, _ := filepath.Split(stub)
        if (subdir == "" || !strings.HasSuffix(path, suffix)) {
            return nil
        }

        gif_info, err := LoadGifMetadata(path)
        if (err != nil) {
            return err
        }
        gifs = append(gifs, *gif_info)

        return nil
    })

    if (err != nil) {
        return nil, nil, err
    }

    // Looping through and checking that each GIF has valid characters.
    for _, x := range gifs {
        curshow, found := show_ptrs[x.ShowId]
        if (!found) {
            return nil, nil, errors.New("did not find show-level metadata for '" + x.Id + "'")
        }

        for _, y := range x.Characters {
            _, found := curshow.Characters[y]
            if (!found) {
                return nil, nil, errors.New("did not find listing for '" + y + "' in '" + x.Id + "'")
            }
        }
    }

    return gifs, shows, nil
}

func DumpToJSON(path string, manifest interface{}) error {
    stuff, err := json.MarshalIndent(manifest, "", "    ")
    if (err != nil) {
        return err
    }

    handle, err := os.Create(path)
    if (err != nil) {
        return err
    }

    _, err = handle.Write(stuff)
    if (err != nil) {
        handle.Close()
        return err
    }

    err = handle.Close()
    if (err != nil) {
        return err
    }

    return nil
}

// Finally, the main loop.

func main() {
    dir := flag.String("dir", "", "Directory containing the GIF and show metadata")
    out := flag.String("out", ".", "Directory in which to store the output manifests")
    flag.Parse()

    if (*dir == "") {
        fmt.Fprintln(os.Stderr, "need to specify directory containing the metadata")
        os.Exit(1)
    }

    if (*out == "") {
        fmt.Fprintln(os.Stderr, "need to specify directory for output")
        os.Exit(1)
    }

    gifs, shows, err := CollateMetadata(*dir)
    if (err != nil) {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    // Dumping the output to file.
    err = os.MkdirAll(*out, 0755)
    if (err != nil) {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    err = DumpToJSON(filepath.Join(*out, "shows.json"), shows)
    if (err != nil) {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    err = DumpToJSON(filepath.Join(*out, "gifs.json"), gifs)
    if (err != nil) {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    os.Exit(0)
}
