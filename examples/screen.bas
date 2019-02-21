10 print "hit ENTER to initialize screen 0"
20 print input$(1)
30 screen 0
40 print "screen 0 initialized. hit key to see term scroll test"
50 print input$(1)
60 for i=1 to 30:print i:next
70 print "scroll test done. hit key to test keyboard input"
80 i$=input$(1)
90 while i$<>"q"
100 print "key:";i$;" asc:";asc(i$)
110 print "hit any other key, or q to exit"
120 i$=input$(1)
130 wend
