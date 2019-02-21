110 print "hit ENTER to initialize screen 0"
120 print input$(1)
130 screen 0
140 print "screen 0 initialized. hit key to see term width test "
150 print input$(1)
160 for i=1 to 200:print chr$((i mod 10)+asc("0"));:next
170 print:print "width test done. hit key to see term scroll test"
180 print input$(1)
190 for i=1 to 30:print i:next
200 print "scroll test done. hit key to test keyboard input"
210 i$=input$(1)
220 while i$<>"q"
230 print "key:";i$;" asc:";asc(i$)
240 print "hit any other key, or q to exit"
250 i$=input$(1)
260 wend
