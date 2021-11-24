# Format the acceptable anime GIF database

This repository contains code to format the [acceptable anime GIF](https://github.com/LTLA/acceptable-anime-gifs) database.
The aim is to easily collate the GIF metadata into a single set of manifests for easy consumption by REST services.
To use, simply download the prebuilt releases binaries (or run `go build .` yourself) and then supply a directory containing the GIF registry:

```sh
./anime-gif-formatter -dir registry
```

The registry should contain subdirectories for each show, which in turn contain the GIFs relating to that show:

```
registry/
- <SHOW_1>.json
- <SHOW_1>/
  - <GIF_1>.gif
  - <GIF_1>.json
  - <GIF_2>.gif
  - <GIF_2>.json
  - ..
- <SHOW_2>.json
- <SHOW_2>/
  - <GIF_1>.gif
  - <GIF_1>.json
  - ...
- ...
```

Note that the `SHOW_*` and `GIF_*` are arbitrary - any value can be used as long as they are unique within their respective directories.
That is, each show has its own `SHOW_*` name, while each GIF within a given show has a different `GIF_*` name (which does not need to be unique across shows).

The JSON files contain the metadata for each show and GIF.
For each show, the metadata should contain:

- `id`: the [MyAnimeList](https://myanimelist.net) identifier for the show.
- `name`: the name of the show.
- `characters`: an object containing the name and MyAnimeList identifier for relevant characters.

For example, we might have:

```js
{
    "id":"10165",
    "name":"Nichijou",
    "characters": {
        "Nano Shinonome": "10422",
        "Hakase Shinonome": "41055",
        "Mio Naganohara": "40081"
    }
}
```

For each GIF, the metadata should contain:

- `characters`: an array of strings naming the characters involved in the GIF.
  Each character named in this manner should also be listed in the show's metadata.
- `sentiments`: an array of strings listing the sentiments expressed in the GIF.
  (Controlled vocabulary coming soon.)
- `url`: string containing the original source of the GIF.

So, for example:

```js
{
    "characters": [
        "Mio Naganohara"
    ],
    "sentiments": [
        "attack"
    ],
    "url": "https://25.media.tumblr.com/tumblr_m6r2fnOPGO1qzvtljo1_500.gif"
}
```

The formatter will then produce two JSON manifests containing arrays of objects.
The first is `gifs.json`, where each object describes a GIF and has the following fields:

- `path`: string containing the path to the GIF file relative to the input directory.
- `show_id`: string containing the MyAnimeList identifier for the show in which the GIF occurs.
- `characters`: an array of strings as described for the `<GIF_*>.json` file.
- `sentiments`: an array of strings as described for the `<GIF_*>.json` file.
- `url`: string as described for the `<GIF_*>.json` file.

The second file is `shows.json`, which also contains an array of objects describing the individual shows.
Each object has exactly the same contents as the `<SHOW_*>.json` JSON file for each show.
