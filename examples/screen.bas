10 print "hit ENTER to initialize screen 0"
20 print input$(1)
30 screen 0
40 print "screen 0 initialized. hit key to see term scroll test"
50 print input$(1)
60 for i=1 to 30:print i:next
70 print "scroll test done. hit key to test keyboard input"
80 while i$<>"q"
90 print "key:";i$
100 print "hit any other key, or q to exit"
110 i$=input$(1)
120 wend
