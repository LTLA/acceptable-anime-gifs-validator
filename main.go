package main

import (
    "os"
    "encoding/json"
    "path/filepath"
    "strings"
    "errors"
    "fmt"
    "flag"
    "io/ioutil"
)

// Definitions of structs and constants. 

type GifInfo struct {
    Path string `json:"path"`
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

func LoadGifMetadata(dir, show, base string) (*GifInfo, error) {
    handle, err := os.Open(filepath.Join(dir, show, base))
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

    if (len(base) == len(suffix)) {
        return nil, errors.New("'base' should contain a non-empty name before '" + suffix + "'")
    }
    prefix := base[:len(base) - len(suffix)]
    gif_name := prefix + ".gif"
    output.Path = show + "/" + gif_name // don't use filepath.Join, we're making URL components here.

    checkpath := filepath.Join(dir, show, gif_name)
    if _, err := os.Stat(checkpath); os.IsNotExist(err) {
        return nil, errors.New("GIF should exist at '" + checkpath + "'")
    }

    return &output, nil
}

func LoadShowMetadata(dir, show string) (*ShowInfo, error) {
    path := filepath.Join(dir, filepath.Clean(show) + ".json")

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

    return &output, nil
}

// Overall collator function.

func CollateMetadata(dir string) ([]GifInfo, []ShowInfo, error) {
    var gifs []GifInfo
    var shows []ShowInfo

    entries, err := ioutil.ReadDir(dir)
    if (err != nil) {
        return nil, nil, err
    }

    for _, info := range entries {
        if (info.IsDir()) {
            show_info, err := LoadShowMetadata(dir, info.Name())
            if (err != nil) {
                err = errors.New("failed to load metadata for '" + info.Name() + "':\n" + err.Error())
                return nil, nil, err
            }
            shows = append(shows, *show_info)

            subentries, err := ioutil.ReadDir(filepath.Join(dir, info.Name()))
            if (err != nil) {
                return nil, nil, err
            }

            for _, subinfo := range subentries {
                if (!strings.HasSuffix(subinfo.Name(), ".json")) {
                    continue
                }

                gif_info, err := LoadGifMetadata(dir, info.Name(), subinfo.Name())
                if (err != nil) {
                    err = errors.New("failed to load metadata for '" + info.Name() + "/" + subinfo.Name() + "':\n" + err.Error())
                    return nil, nil, err
                }

                // Updating the show ID to use the actual ID, not our path.
                gif_info.ShowId = show_info.Id

                // Checking that each GIF has valid characters.
                for _, y := range gif_info.Characters {
                    _, found := show_info.Characters[y]
                    if (!found) {
                        return nil, nil, errors.New("did not find listing for '" + y + "' in '" + gif_info.Path + "'")
                    }
                }

                gifs = append(gifs, *gif_info)
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
