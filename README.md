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

## Configuration

### Main configuration file -- `merchants.json`

An example

### Account mapping rules -- `rules.json`

An example

## Links

* https://api.privatbank.ua/#p24/orders
* https://en.wikipedia.org/wiki/Quicken_Interchange_Format
* https://www.gnucash.org/
* http://www.freebsd.org/copyright/freebsd-license.html
