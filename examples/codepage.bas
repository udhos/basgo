100 screen 0
110 cls
120 print "printing all characters in current codepage"
130 for i=0 to 255
140 rem locate 23,1: print i;"  "
150 locate i\16+3, 2*(i MOD 16)+1
160 if i<>12 then print chr$(i);
170 rem locate 22,1: print "hit any key -- or q to exit":i$=input$(1): locate 22,1: print space$(20);
180 rem if i$="q" then end
190 next
300 print:print
310 print "hit any key"
320 print input$(1)
