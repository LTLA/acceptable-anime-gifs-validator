# Format the acceptable anime GIF database

This repository contains code to format the [acceptable anime GIF](https://github.com/LTLA/acceptable-anime-gifs) database.
The aim is to easily collate the GIF metadata into a single set of manifests for easy consumption by REST services.
To use, simply download the prebuilt binaries (or run `go build .`) and then pass the path to a directory containing the GIF registry:

```sh
./anime-gif-formatter -dir registry
```

The registry should contain subdirectories for each show, which in turn contain the GIFs relating to that show:

```
registry/
- <SHOW_ID_1>.json
- <SHOW_ID_1>/
  - <GIF_ID_1>.gif
  - <GIF_ID_1>.json
  - <GIF_ID_2>.gif
  - <GIF_ID_2>.json
  - ..
- <SHOW_ID_2>.json
- <SHOW_ID_2>/
  - <GIF_ID_1>.gif
  - <GIF_ID_1>.json
  - ...
- ...
```

The JSON files contain the metadata for each show and GIF.
For each show, the metadata should contain:

- `id`: the [MyAnimeList](https://myanimelist.net) identifier for the show.
- `name`: the name of the show.
- `characters`: an object containing the name and identifier for relevant characters.

For example, we might have:

```js
{
    "id":"10165",
    "name":"Nichijou",
    "characters": {
        "Nano Shinonome": "10422",
        "Hakase Shinonome": "41055"
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
