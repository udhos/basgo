110 print "hit ENTER to initialize screen 0"
120 print input$(1)
130 screen 0
140 print "screen 0 initialized": print "hit key to see CLS test"
150 print input$(1)
160 cls
170 print "CLS tested. hit key to see term width test"
180 print input$(1)
190 for i=1 to 200:print chr$((i mod 10)+asc("0"));:next
200 print:print "width test done. hit key to see term scroll test"
210 print input$(1)
220 for i=1 to 30:print i:next
230 print "scroll test done. hit key to test keyboard input"
240 i$=input$(1)
250 while i$<>"q"
260 print "key:";i$;" asc:";asc(i$)
270 print "hit any other key, or q to exit"
280 i$=input$(1)
290 wend
