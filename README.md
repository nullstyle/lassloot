# lassloot - because I can't justify the license cost for lastools - LIDAR point cloud code and tools for dorks

So, hypothetically, you find yourself with a piece of land you'd like to develop.   You're a Person of the Present Age(tm) and so you hire some dronefolk to LIDAR scan the plot.  The day of scanning comes and goes, buzzings were heard, five to ten thousand dollars have left your bank account and you're left with a series of classified and unclassified ASPRS LAS 1.4 files.  These point cloud files represent the data read from the LIDAR sensors and you'd like to do **all sorts of fun stuff** with the data.

You'd like to 3d print a copy of the plot to share with your parents. They'd appreciate the print more than a 3d view, regardless your webgl skills. You'd like to import a mesh of the plot into blender so you can sculpt out ideas for earthworks.  It turns out working the earth virtually is a lot easier than hiring some hardhats to run a dozerloadervator to push actual dirt, rocks, and roots around.  Most importantly, the architects you've hired would love some topographical contours that they may sacrifice to their Autodesk-ian gods. You encounter a road block:  How do I get this primo-expensive-data into these various formats?  The dronefolk want an upcharge for each transmogrification: That won't do.  Surely their cpu time and licensing fees are not that expensive?  Perhaps there is some existing code out there to do the work for you? Why does my finger smell like this? All of these questions hit you at once.

Googles are googles, links are clinked, thunder booms and XKCD re-asserts its prescience.  [They might not have been around since 2003 and it looks like they're German rather than Nebraskan, but still -- All LIDAR code out there flows through Martin Isenberg][0].  Holy fuck 1500 euros just license to las2iso?  Time to write some code I guess.

## Capabilities 

- Reads LAS 1.4 files, somewhat

## Discapabilites

- Spatial referencing of any kind.
- Compressed point records

## Usage

```go

import (
	"github.com/nullstyle/lassloot"
)

err, pc := lassloot.NewPointCloudFromPath("somedots.las")

... more code forthcoming ...
	

```

## Documentation

https://pkg.dev.go/github.com/nullstyle/lasloot

## Disclaimer

First of all, please see LICENSE.txt for the definitive word.  Informally, this code is open and available to use but don't expect much out of me from it.  If you need support, I advise you to use LASTools and go pay Martin, if you want your own features added plan on forking the code, and if you need clarification please open a github issue.  This code is developed for my own personal use but made available should you wish to learn from it or hack it up for yourself.

## Contributions

If we're talking realtalk here, the best chances you have to contribute to this repository is to provide some constructive criticism about this library for your use case.  Bonus points if the feedback comes in a wonderful pull request that reads as if it sprung from my own mind.  Points will be deducted if it's a sloppy pull request and you'll probably just get the PR closed without response.

I'm not interested in running a Real Open Source Project as you might find described in a book published by Stripe publishing.  Please fork this repo if you're going to add friction to my life, ain't nobody got time for that.

[0]: https://xkcd.com/2347/