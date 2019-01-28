10 read n$, birth$
20 gosub 100
30 print string$(20, "-"); " "; string$(20, "-")
40 for i=1 to 3
50 read n$, birth$
60 gosub 100
70 next
80 end
100 print n$; space$(20-len(n$)); " "; birth$
110 return
1000 data "NAME", "BIRTHDAY"
1010 data "smith, john", "January 27, 2019"
1020 data "denvers, carol", "February 10, 2018"
1030 data "wayne, bruce", "March 30, 2017"
