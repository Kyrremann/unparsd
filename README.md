Unparsd
=======

A CLI tool that parses Untappd check-in data and generates JSON statistics for use with a static Jekyll site (e.g. https://github.com/Kyrremann/beers).

## Commands

### `generate` — build statistics from a local file

```
unparsd generate [options]

  -u, --untappd=<path>      Path to untappd.json export or a directory of
                            per-year JSON files (default: untappd.json)
  -o, --output=<dir>        Output directory for generated files (default: ./)
  -s, --all-styles=<path>   Path to all-styles.json; omit to scrape Untappd live
```

Reads check-in data and writes statistics to `<output>/_data/`:

| File                     | Contents                                                                    |
|--------------------------|-----------------------------------------------------------------------------|
| `allmy.json`             | Global stats + per-year breakdown (streak, ABV distribution, weekly counts) |
| `beers.json`             | Most-checked-in beers                                                       |
| `breweries.json`         | Brewery stats                                                               |
| `countries.json`         | Check-ins by brewery country                                                |
| `styles.json`            | Beer style breakdown                                                        |
| `venues.json`            | Top venues                                                                  |
| `serving_types.json`     | Serving type distribution                                                   |
| `flavors.json`           | Flavor profile stats                                                        |
| `rating_deltas.json`     | Personal vs. global rating deltas                                           |
| `missing_styles.json`    | Styles not found in Untappd's style list                                    |
| `missing_countries.json` | Countries with no check-ins                                                 |

Monthly Jekyll pages are written to `<output>/_monthly/`.

### `fetch` — incrementally pull new check-ins from the Untappd API

```
unparsd fetch [options]

  -u  --username=<name>   Untappd username to fetch check-ins for (required)
  -o, --output=<dir>      Directory to write per-year JSON files into
                          (default: ./checkins)
```

Credentials are read from environment variables:

```
UNTAPPD_CLIENT_ID=<your client id>
UNTAPPD_CLIENT_SECRET=<your client secret>
```

**How it works:** on each run the tool reads the highest stored `checkin_id` from the `checkins/` directory, then pages backwards from the newest check-in on the API until it reaches that ID.
New check-ins are merged into per-year files (`checkins/2024.json`, etc.) sorted newest-first.  
On the first run (no existing files) the full history is fetched.

Graceful shutdown: sending `SIGINT` or `SIGTERM` (Ctrl+C) saves any check-ins fetched so far before exiting.

## One-time historical import

If you have an Untappd data export (`untappd.json`), split it into per-year files before running `generate` on the `checkins/` directory:

```
python3 split_checkins.py untappd.json checkins/
unparsd generate --untappd checkins/
```
