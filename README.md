# rkn-bgp

Fetches list of resources banned in Russia.

Deduplicates and summarizes records into set of subnets.

Outputs them in BIRD static route format.

# Usage

This tool is supposed to fill list of static routes that are about to be redistributed into BGP. So BIRD daemon should be set up with minimal configuration that looks like this:

```
log syslog all;
router id 1.1.1.1;

protocol static static_bgp {
        import all;
        include "prefixlist.txt";
}

protocol bgp OurRouter {
        description "Our Router";
        neighbor 2.2.2.2 as 22222;
        import none;
        export where proto = "static_bgp";
        local as 11111;
        passive off;
        multihop;
}

```

After that periodic cron job should be set up:

`*/20 * * * * /opt/rkn-bgp > /etc/bird/prefixlist.txt && service bird reload`

It will populate prefix list and reload daemon every 20 minutes, keeping BGP routes up to date.

If you prefer using Docker change cron job to this:

`*/20 * * * * docker run --rm zarbis/rkn-bgp > /etc/bird/prefixlist.txt && service bird reload`

## Options

`-listURL` - custom URL to fetch `.csv` dump from (Defaults to official GitHub repo: https://github.com/zapret-info/z-i).

`-hostPrefix` - prefix to assume when summarizing single host records, setting lower than default value may significantly lower number of resulting subnets at cost of some false-positives (Default: 32).
