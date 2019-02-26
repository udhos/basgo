100 screen 0
110 for fg=0 to 15
115 color 7, 0: print fg;": ";
120 for bg=15 to 0 step -1
130 color fg, bg
140 print "x";
150 next
160 print
170 next
180 print "hit any key"
190 print input$(1)
