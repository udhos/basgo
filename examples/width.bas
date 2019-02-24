
100 screen 0
110 w=25
120 print "setting WIDTH=";w
130 print
140 width w
150 for i=1 to 200:print chr$((i mod 10)+asc("0"));:next
160 print
170 print
180 print "done. hit any key";input$(1)

