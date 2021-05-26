# skyham

I am building an amateur radio satellite tracking system.  It is both hardware
and software.  The core computing hardware is Raspberry Pi.  It is up and
running with go compiler and gobot package installed, tested and working.

I will write up my progress here as I go along.  My second step is work out
the orbital mechanics.  My starting point is this paper.

https://apps.dtic.mil/dtic/tr/fulltext/u2/1027338.pdf

I am not sure of the applicable IP laws, so I will just like it here.

I have carried out his Adelaide calculations and I have included it in the
directory "background".

To use the planets program, here is what you need to do:
1. Build a planet file like the two that I have put in the data directory
2. You can get this data from the reference in either paper mentioned
3. Build and run the program
4. Use the command line interface to change the program parameters
5. There is some error checking in the program, but be careful
6. The built in data in the interface is Don Koks data, but you can change it

Here is the test data and its comparison to the data from the site:
timeanddate.com (web for short).  All of this data is for May 22 UTC.  
Time column is also UTC.

Planet Mars:

| Location | Time | Azimuth | Elevation | Azimuth (web) | Elevation (web) |
---------------------------------------------------------------------------
| Berlin   | 2:55 |   6.9   |    13.5    |      7.7     |   13.5          |
| Tokyo    | 3:01 |  94.5   |    49.3    |     95.1     |   50.1          |
| San Fran | 3:04 | 268.6   |    42.5    |    269.4     |   41.9          |
| NYC      | 3:08 | 296.0   |     6.0    |    296.8     |    5.5          |
| London   | 3:12 | 357.9   |   -14.9    |    358.8     |  -14.8          |
| Jo'Burg  | 3:15 |  90.3   |   -65.6    |     89.6     |  -64.7          |

I will add more data in the future.
