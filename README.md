# Fetch transaction log from Privat24 for GnuCash

## Summary

Command line tool to fetch transaction logs using
[Privat24 API](https://api.privatbank.ua/#p24/orders) and
convert them to
[Quicken Interchange Format](https://en.wikipedia.org/wiki/Quicken_Interchange_Format),
which can be easily imported to [GnuCash](https://www.gnucash.org/).

## License

It uses a [FreeBSD License](http://www.freebsd.org/copyright/freebsd-license.html).
You can obtain the license online or in the file LICENSE on
the top of the sources tree.

## Build

You need [Golang](https://golang.org/) 1.14 to build it:

```
make
```

On success, `p24fetch` executable will be generated.

## Run the tool

```
./p24fetch etc/merchants.json
```

## Run unit tests

```
make test
```

## How it works

Possible usage workflow is:

1. Launch p24fetch weekly (for instance, with [cron](https://en.wikipedia.org/wiki/Cron));
2. Run GnuCash, import generated QIF files, then remove the files;
3. Add unsorted transactions to GnuCash manually,
 update sorting rules in the `rules.json` config file.

On every successful merchant processing the date and time of the last
fetched transaction will be stored under `run/dedup` directory (see
`dedup_dir` setting in `merchants.json` configuration file), so
the next time `p24fetch` will pull transactions from the Privat24 API
only new transactions will be processed.

## Configuration

### Main configuration file -- `merchants.json`

An example can be found in [etc/merchants.json.example](etc/merchants.json.example).
It has two main sections:

* _merchants_ -- list of all known Privat24 Merchants;
* _defaults_ -- convenience dict which will fill up missing
 fields in merchants dicts.

Schema and meaning of the fields of every _merchants_ entry and
_defaults_ dict are documented in
[config/config.go](config/config.go).

### Account mapping rules -- `rules.json`

An example can be found in [etc/rules.json.example](etc/rules.json.example).
It has three sections:

* _accounts_ -- a mapping from account shorthand IDs to account IDs in
 your GnuCash Ledger;
* _ignore_ -- regexp patterns of transactions which should be ignored;
* _rules_ -- array of ShorthandAccountID to Regexp patterns mappings.

How transactions are matched against regexps:

Every fetched transaction has two mandatory fields: beneficiary name and
transaction note. Transaction considered matching particular regexp pattern
when beneficiary name OR transaction note match the regexp.

How transactions are processed:

After being fetched from Privat24 API, every transaction matched against
_ignore_ patterns. On match, it will not be processed further. Then the
_rules_ array will be traversed to find a match between transaction and
one of configured accounts.

All matched ransactions will be exported to `results` directory (see
`results_dir` setting in `merchans.json` config) as a QIF file using
GnuCash account name as a beneficiary name. Transaction note will contain
both origin beneficiary name and transaction note.

All ignored transactions will be stored as JSON files under `results/ignored`
directory (see `results_dir` setting in `merchans.json` config).

All transactions not matched to any known account will be:

* stored as JSON files under `results/unsorted` directory (see
 `results_dir` setting in `merchans.json` config);
* (optional) messaged to configured Slack Channel.

## Links

* https://api.privatbank.ua/#p24/orders
* https://en.wikipedia.org/wiki/Quicken_Interchange_Format
* https://www.gnucash.org/
* http://www.freebsd.org/copyright/freebsd-license.html
